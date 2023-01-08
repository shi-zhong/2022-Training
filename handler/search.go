package handler

import (
	"HappyShopTogether/model/dbop"
	"HappyShopTogether/utils/code"
	"github.com/gin-gonic/gin"
)

func SearchHandler(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")
	searchKeys := c.QueryArray("search")

	commodities, shops := dbop.SearchCommoditiesLimitPage(limit, page, searchKeys)

	code.GinOKPayload(c, &gin.H{
		"commodities": commodities,
		"shops":       shops,
	})
}
