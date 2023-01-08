package router

import (
	"HappyShopTogether/handler"
    "HappyShopTogether/middleware"
    "github.com/gin-gonic/gin"
)

func setUserInfoRouter(router *gin.RouterGroup) {
	group := router.Group("/user")
	{
		// 用户信息查询
		userInfo := group.Group("/info")
		{
			// 顾客基本信息获取
			userInfo.GET("/get/customer", handler.GetCustomerHandler)
			userInfo.GET("/get/merchant", handler.GetMerchantHandler)
			// 顾客信息修改
			userInfo.PUT("/update/customer", handler.UpdateCustomerHandler)
			userInfo.PUT("/update/merchant", handler.UpdateMerhantHandler)
		}

		userAddress := group.Group("/address")
        userAddress.Use(middleware.CustomerOnly())
		{
            // 用户only
			userAddress.POST("/add", handler.AddressAddHandler)
			userAddress.PUT("/update", handler.AddressUpdateHandler)
			userAddress.PUT("/default/:addressID", handler.AddressDefaultHandler)
			userAddress.DELETE("/delete/:addressID", handler.AddressDeleteHandler)
		}
		// 购物车列表
		userCart := group.Group("/cart")
        userCart.Use(middleware.CustomerOnly())
		{
			userCart.GET("/list", handler.CartListHandler)
			userCart.POST("/add/:ID", handler.CartListAddHandler)         // commodity_id
			userCart.DELETE("/remove/:ID", handler.CartListRemoveHandler) // cart_id
		}

//		group.GET("/history", handler.DefaultHandler)
	}
}
