package controllers

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/sing3demons/app/cache"
	"github.com/sing3demons/app/models"

	"github.com/gin-gonic/gin"
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

type productResponse struct {
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
	Items  []productResponse `json:"items"`
	Paging *pagingResult     `json:"paging"`
}

func setEnvPath(name string, id uint) string {
	path := "uploads/" + name + "/" + strconv.Itoa(int(id))
	os.MkdirAll(path, os.ModePerm)
	return path
}

func DestroyImage(product *models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	subPath := strings.Split(product.Image, "/")

	fileName := subPath[10] + "/" + subPath[11] + "/" + subPath[12]
	fmt.Printf("%s \n", fileName)

	cld, public_id, err := NewCloudinary(fileName)
	if err != nil {
		return err
	}

	cld.Upload.Destroy(ctx, uploader.DestroyParams{PublicID: public_id})
	return nil
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

	path := setEnvPath("products", products.ID)

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
