package router

//import (
//    "HappyShopTogether/handler"
//    "github.com/gin-gonic/gin"
//)
//
//func setChatRouter(router *gin.RouterGroup) {
//    group := router.Group("/chat")
//    {
//        // 获取联系人列表
//        group.GET("/contact/list", handler.DefaultHandler)
//        // 获取聊天记录
//        group.GET("/record/get", handler.DefaultHandler)
//        // 发送消息
//        group.POST("/send", handler.DefaultHandler)
//        // 标记已读
//        group.PUT("/read/status/:contactWith", handler.DefaultHandler)
//        // 及时获取聊天记录 websocket 待学习
//        group.GET("/quick/get", handler.DefaultHandler)
//    }
//}
