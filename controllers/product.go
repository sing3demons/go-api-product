package controllers

import (
	"net/http"
	"strconv"

	"github.com/sing3demons/app/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// FindAll godoc
// @Summary Show an products
// @Tags products
// @Accept  json
// @Produce  json
// @Param page query uint false "page"
// @Param limit query uint false "limit"
// @Success 200 {object} productsPaging
// @Router /api/v1/products [get]
func (p *Product) FindAll(ctx *gin.Context) {
	var products []models.Product

	query := p.DB.Preload("Category").Order("id desc")

	if category := ctx.Query("category"); category != "" {
		c, _ := strconv.Atoi(category)
		query = query.Where("category_id = ?", c)
	}

	pagination := pagination{ctx: ctx, query: p.DB, records: &products}
	paging := pagination.paginate()

	serializedProduct := []productResponse{}
	copier.Copy(&serializedProduct, &products)

	ctx.JSON(http.StatusOK, gin.H{"products": productsPaging{Items: serializedProduct, Paging: paging}})
}

// FindOne godoc
// @Summary FindOne - /:id
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path int true "ID"
// @Success 200 {object} productResponse
// @Router /api/v1/products/{id} [get]
func (p *Product) FindOne(ctx *gin.Context) {
	product, err := p.findProductByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	serializedProduct := productResponse{}
	copier.Copy(&serializedProduct, &product)
	ctx.JSON(http.StatusOK, gin.H{"product": serializedProduct})

}

// Create godoc
// @Summary add an product
// @Description add by form product
// @Tags products
// @Accept  mpfd
// @Produce  json
// @Security BearerAuth
// @Param name formData string true "name"
// @Param desc formData string true "desc"
// @Param price formData int true "price"
// @Param image formData file true "image"
// @Param categoryId formData uint true "categoryId"
// @Success 201 {object} productResponse
// @Router /api/v1/products [post]
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

	serializedProduct := productResponse{}
	copier.Copy(&serializedProduct, &product)

	ctx.JSON(http.StatusCreated, gin.H{"product": serializedProduct})

}

//Update - update --> patch
// UpdateAll godoc
// @Summary update an products
// @Description update by form product
// @Tags products
// @Accept  mpfd
// @Produce  json
// @Security BearerAuth
// @Param id path string true "id"
// @Param name formData string false "name"
// @Param desc formData string false "desc"
// @Param price formData int false "price"
// @Param image formData file false "image"
// @Param categoryId formData uint false "categoryId"
// @Success 200 {object} productResponse
// @Router /api/v1/products/{id} [patch]
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

	p.setProductImage(ctx, product)
	if err := p.DB.Save(&product).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error})
		return
	}

	// p.setProductImage(ctx, product)

	var serializedProduct productResponse
	copier.Copy(&serializedProduct, &product)
	ctx.JSON(http.StatusOK, gin.H{"product": serializedProduct})
}

// Delete godoc
// @Summary	delete an product
// @Description	delete by json product
// @Tags	products
// @Accept	json
// @Produce	json
// @Security BearerAuth
// @Param id path string true "id"
// @Success 200 {object} productResponse
// @Failure	422  {object} string "Bad Request"
// @Failure	404  {object}  map[string]any	"{"error": "not found"}"
// @Router /api/v1/products/{id} [delete]
func (p *Product) Delete(ctx *gin.Context) {
	product, err := p.findProductByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// DestroyImage(product)
	p.DB.Unscoped().Delete(&product)

	ctx.Status(http.StatusNoContent)
}
