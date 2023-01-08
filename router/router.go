package router

import (
    "HappyShopTogether/middleware"
    "github.com/gin-gonic/gin"
)

var Router *gin.Engine

func Init() {
	Router = gin.Default()

    // 不需要token
    donotNeedAuthorize := Router.Group("/api/v1")
    {
        setAuthorRouterWithoutAuthorize(donotNeedAuthorize)
        setSearchRouter(donotNeedAuthorize)
    }

	needAuthorize := Router.Group("/api/v1")
    needAuthorize.Use(middleware.TokenAuthorize())
	{
        setAuthorRouter(needAuthorize)
        
        setUserInfoRouter(needAuthorize)
        setCommodityRouter(needAuthorize)
        setOrderRouter(needAuthorize)
//        setChatRouter(needAuthorize)
//        setStatisticsRouter(needAuthorize)
	}
}
