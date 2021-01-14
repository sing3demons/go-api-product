package controllers

import (

	"mime/multipart"

	"github.com/gin-gonic/gin"
)

//Product - struct
type Product struct{}

type createProductForm struct {
	Title string                `form:"title" binding:"required"`
	Body  string                `form:"body" binding:"required"`
	Image *multipart.FileHeader `form:"image" binding:"required"`
}


func (p *Product) FindAll(ctx *gin.Context) {
	
}

func (p *Product) FindOne(ctx *gin.Context) {

}

func (p *Product) Update(ctx *gin.Context) {

}
