package controllers

import (
	"app/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// FindAll - query-database-all
func (p *Product) _FindAll(ctx *gin.Context) {
	query1CacheKey := "items::product"
	query2CacheKey := "items::page"

	var serializedProduct []productRespons
	var paging *pagingResult

	cacheItems, err := p.Cache.MGet([]string{query1CacheKey, query2CacheKey})
	if err != nil {
		log.Println(err.Error())
	}

	productJS := cacheItems[0]
	pageJS := cacheItems[1]

	// fmt.Println("productJS: ", len(productJS.(string)))
	// fmt.Println("pageJS: ", len(pageJS.(string)))

	// Found query #1 cache
	if productJS != nil && len(productJS.(string)) > 0 {
		// ctx.Log("cache hit")

		err := json.Unmarshal([]byte(productJS.(string)), &serializedProduct)
		if err != nil {
			p.Cache.Del(query1CacheKey)
			log.Println(err.Error())
		}

	}

	itemToCaches := map[string]interface{}{}

	// var paginationItem pagination
	var paginationItem *pagingResult
	if productJS == nil {
		fmt.Println("productJS == nil : ", productJS)
		var products []models.Product
		query := p.DB.Preload("Category").Order("id desc")
		if category := ctx.Query("category"); category != "" {
			c, _ := strconv.Atoi(category)
			query = query.Where("category_id = ?", c)
		}

		// paginationItem = pagination{ctx: ctx, query: query, records: &products}
		pagination := pagination{ctx: ctx, query: query, records: &products}
		paginationItem = pagination.paginate()
		copier.Copy(&serializedProduct, &products)

		fmt.Println("serializedProduct : ", len(products))

		itemToCaches[query1CacheKey] = serializedProduct
	}

	// Found query #2 cache
	if pageJS != nil && len(pageJS.(string)) > 0 {
		// ctx.Log("cache hit")
		// counter, err = strconv.Atoi(counterJS.(string))

		err := json.Unmarshal([]byte(pageJS.(string)), &paging)
		if err != nil {
			p.Cache.Del(query2CacheKey)
			log.Println(err.Error())
		}
	}

	if paging == nil {
		paging = paginationItem
		itemToCaches[query2CacheKey] = paging
	}

	// fmt.Println("Found query #2 cache: ", itemToCaches[query2CacheKey])

	// fmt.Println("check...", itemToCaches[query1CacheKey])
	if len(itemToCaches) > 0 {
		timeToExpire := 10 * time.Second // m
		fmt.Println("MSET")

		// Set cache using MSET
		err := p.Cache.MSet(itemToCaches)
		if err != nil {
			log.Println(err.Error())
		}

		// Set time to expire
		keys := []string{}
		for k := range itemToCaches {
			keys = append(keys, k)
		}
		err = p.Cache.Expires(keys, timeToExpire)
		if err != nil {
			log.Println(err.Error())
		}
	}

	// serializedProduct := []productRespons{}
	// copier.Copy(&serializedProduct, &products)

	ctx.JSON(http.StatusOK, gin.H{"products": producsPaging{Items: serializedProduct, Paging: paging}})
}
