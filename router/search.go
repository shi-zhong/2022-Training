package router

import (
	"HappyShopTogether/handler"
	"github.com/gin-gonic/gin"
)

func setSearchRouter(router *gin.RouterGroup) {
	//    group := router.Group("")
	router.GET("/search", handler.SearchHandler)
}
