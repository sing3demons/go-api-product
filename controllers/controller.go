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

	ch := make(chan int)
	go countRecords(query, records, ch)

	offset := (page - 1) * limit
	query.Limit(limit).Offset(offset).Find(records)

	count := <-ch
	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	var nextPage int
	if nextPage == totalPage {
		nextPage = totalPage
	} else {
		nextPage = page + 1
	}

	return &pagingResult{
		Page:      page,
		Limit:     limit,
		Count:     count,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		TotalPage: totalPage,
	}

}

func countRecords(query *gorm.DB, records interface{}, ch chan int) {
	var count int
	query.Model(records).Count(&count)

	ch <- count
}
