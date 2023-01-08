package router

import (
	"HappyShopTogether/handler"
    "HappyShopTogether/middleware"
    "github.com/gin-gonic/gin"
)

func setOrderRouter(router *gin.RouterGroup) {
	shareBillGroup := router.Group("/share_bill")
    shareBillGroup.Use(middleware.CustomerOnly())
	{
		// 创建拼单订单
        shareBillGroup.POST("/create", handler.ShareBillCreateHandler)
		// 参与拼团
        shareBillGroup.PUT("/join", handler.ShareBillJoinHandler)
		// 获取拼单链接
//		shareBillGroup.GET("/link", handler.DefaultHandler)
		// 获取拼单列表
        shareBillGroup.GET("/list", handler.ShareBillListHandler)
		// 获取拼单订单细节
        shareBillGroup.GET("/detail/:ID", handler.ShareBillDetailHandler)
	}

    // 商家发货
    router.PUT("/order/merchant/confirm/:ID", handler.OrderMerchantConfirmHandler)

	orderGroup := router.Group("/order")
    orderGroup.Use(middleware.CustomerOnly())
	{
		// 订单列表获取
		orderGroup.GET("/list/customer", handler.OrderCustomerListHandler)
		// 获取订单细节
		orderGroup.GET("/detail/:ID", handler.OrderDetailHandler)
		// 顾客收货
        orderGroup.PUT("/customer/confirm/:ID", handler.OrderCustomerConfirmHandler)
	}
}
