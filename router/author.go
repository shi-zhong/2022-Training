package router

import (
	"HappyShopTogether/handler"
	"github.com/gin-gonic/gin"
)

// 需要token鉴权
func setAuthorRouter(router *gin.RouterGroup) {
	group := router.Group("/author")
	{
		// 上传私钥
//		group.PUT("/private_key", handler.DefaultHandler)
		group.PUT("/mobile/update", handler.MobileUpdateHandler)
		group.PUT("/password/update", handler.PasswordUpdateHandler)
	}
}

func setAuthorRouterWithoutAuthorize(router *gin.RouterGroup) {
	group := router.Group("/author")
	{
		// 获取公钥
//		group.GET("/public_key", handler.DefaultHandler)

		group.POST("/login", handler.LoginHandler)
		group.POST("/register/customer", handler.RegisterCustomerHandler)
		group.POST("/register/merchant", handler.RegisterMerchantHandler)
		// 店铺名是否存在
		group.POST("/shopname/check", handler.ShopnameCheckHandler)
		group.POST("/mobile/check", handler.MobileCheckHandler)

	}
}
