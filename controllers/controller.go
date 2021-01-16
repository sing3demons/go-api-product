package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

func pagingResource(ctx *gin.Context, query *gorm.DB, records interface{}) *pagingResult {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "12"))

	var count int
	query.Model(records).Count(&count)

	// 3. Find Records
	// limit, offset
	// limit => 10
	// page => 1, 1 - 10, offset => 0
	// page => 2, 11 - 20, offset => 10
	// page => 3, 21 - 30, offset => 20
	offset := (page - 1) * limit
	query.Offset(offset).Limit(limit).Find(records)

	// 4. total page
	totalPage := int(math.Ceil(float64(count) / float64(limit)))
	// 5. Find nextPage
	var nextPage int
	if nextPage == totalPage {
		nextPage = totalPage
	} else {
		nextPage = totalPage + 1
	}
	// 6. create pagingResult
	return &pagingResult{
		Page:      page,
		Limit:     limit,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		Count:     count,
		TotalPage: totalPage,
	}
}
