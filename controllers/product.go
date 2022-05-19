package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sing3demons/app/cache"
	"github.com/sing3demons/app/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

//Product - struct
type Product struct {
	DB     *gorm.DB
	Cacher *cache.Cacher
}

type createProductForm struct {
	Name       string                `form:"name" binding:"required"`
	Desc       string                `form:"desc" binding:"required"`
	Price      int                   `form:"price" binding:"required"`
	Image      *multipart.FileHeader `form:"image" binding:"required"`
	CategoryID uint                  `form:"categoryId" binding:"required"`
}

type updateProductForm struct {
	Name       string                `form:"name"`
	Desc       string                `form:"desc"`
	Price      int                   `form:"price"`
	Image      *multipart.FileHeader `form:"image"`
	CategoryID uint                  `form:"categoryId"`
}

type productRespons struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
	Price      int    `json:"price"`
	Image      string `json:"image"`
	CategoryID uint   `json:"categoryId" binding:"required"`
	Category   struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
}

type productsPaging struct {
	Items  []productRespons `json:"items"`
	Paging *pagingResult    `json:"paging"`
}

// FindAll - query-database-all
func (p *Product) FindAll(ctx *gin.Context) {
	query1CacheKey := "items::product"
	query2CacheKey := "items::page"

	serializedProduct := []productRespons{}
	var paging *pagingResult

	cacheItems, err := p.Cacher.MGet([]string{query1CacheKey, query2CacheKey})
	if err != nil {
		log.Println(err.Error())
	}

	productJS := cacheItems[0]
	pageJS := cacheItems[1]

	if productJS != nil && len(productJS.(string)) > 0 {
		err := json.Unmarshal([]byte(productJS.(string)), &serializedProduct)
		if err != nil {
			p.Cacher.Del(query1CacheKey)
			log.Println(err.Error())
		}

	}

	itemToCaches := map[string]interface{}{}

	var paginationItem *pagingResult
	if productJS == nil {
		var products []models.Product
		pagination := pagination{ctx: ctx, query: p.DB, records: &products}
		// pagination := NewPaginationHandler(ctx, p.store, &products)
		paginationItem = pagination.paginate()
		copier.Copy(&serializedProduct, &products)

		itemToCaches[query1CacheKey] = serializedProduct
	}

	if pageJS != nil && len(pageJS.(string)) > 0 {
		err := json.Unmarshal([]byte(pageJS.(string)), &paging)
		if err != nil {
			p.Cacher.Del(query2CacheKey)
			log.Println(err.Error())
		}
	}

	if paging == nil {
		paging = paginationItem
		itemToCaches[query2CacheKey] = paging
	}

	if len(itemToCaches) > 0 {
		timeToExpire := 10 * time.Second // m
		fmt.Println("M_SET")

		// Set cache using MSET
		err := p.Cacher.MSet(itemToCaches)
		if err != nil {
			log.Println(err.Error())
		}

		// Set time to expire
		keys := []string{}
		for k := range itemToCaches {
			keys = append(keys, k)
		}
		err = p.Cacher.Expires(keys, timeToExpire)
		if err != nil {
			log.Println(err.Error())
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"products": productsPaging{Items: serializedProduct, Paging: paging}})
}

// FindOne - first
func (p *Product) FindOne(ctx *gin.Context) {
	product, err := p.findProductByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	serializedProduct := productRespons{}
	copier.Copy(&serializedProduct, &product)
	ctx.JSON(http.StatusOK, gin.H{"product": serializedProduct})

}

// Create - insert data
func (p *Product) Create(ctx *gin.Context) {
	var form createProductForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var product models.Product
	copier.Copy(&product, &form)

	if err := p.DB.Preload("Category").Create(&product).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	p.setProductImage(ctx, &product)

	serializedProduct := productRespons{}
	copier.Copy(&serializedProduct, &product)

	ctx.JSON(http.StatusCreated, gin.H{"product": serializedProduct})

}

//Update - update --> patch
func (p *Product) Update(ctx *gin.Context) {
	var form updateProductForm
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	product, err := p.findProductByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	}

	copier.Copy(&product, &form)

	if err := p.DB.Save(&product).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error})
		return
	}

	p.setProductImage(ctx, product)

	var serializedProduct productRespons
	copier.Copy(&serializedProduct, &product)
	ctx.JSON(http.StatusOK, gin.H{"product": serializedProduct})
}

//Delete - delete
func (p *Product) Delete(ctx *gin.Context) {
	product, err := p.findProductByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	p.DB.Unscoped().Delete(&product)

	ctx.Status(http.StatusNoContent)
}

func (p *Product) setProductImage(ctx *gin.Context, products *models.Product) error {
	file, err := ctx.FormFile("image")
	if err != nil || file == nil {
		return nil
	}

	if products.Image == "" {
		products.Image = strings.Replace(products.Image, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + products.Image)
	}

	path := "uploads/products/" + strconv.Itoa(int(products.ID))
	os.MkdirAll(path, os.ModePerm)

	filename := path + "_" + "product"
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	url, err := cloudinaryUpload(filename)
	if err != nil {
		return err
	}

	fmt.Printf("url %s\n", *url)

	products.Image = *url
	p.DB.Save(products)

	return nil
}

func (p *Product) findProductByID(ctx *gin.Context) (*models.Product, error) {
	var product models.Product
	id := ctx.Param("id")

	if err := p.DB.Preload("Category").First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}
