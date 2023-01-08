package router

import (
	"HappyShopTogether/handler"
    "HappyShopTogether/middleware"
    "github.com/gin-gonic/gin"
)

func setCommodityRouter(router *gin.RouterGroup) {

    // 店铺内商品列表
    router.GET("/commodity/list/shop/:ID", handler.CommodityShopListHandler)

	group := router.Group("/commodity")
    group.Use(middleware.MerchantOnly())
    {
		group.POST("/create", handler.CommodityCreateHandler)
		group.POST("/update", handler.CommodityUpdateHandler)
		// 上架下架
		group.PUT("/status", handler.CommodityStatusHandler)
		group.PUT("/delete/:ID", handler.CommodityDeleteHandler)
		group.GET("/detail/:ID", handler.CommodityDetailHandler)
	}
}
