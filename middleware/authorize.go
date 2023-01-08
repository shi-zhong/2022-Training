package middleware

import (
    "HappyShopTogether/model"
    "HappyShopTogether/utils"
	"HappyShopTogether/utils/code"
	_ "fmt"
	"github.com/gin-gonic/gin"
)

func TokenAuthorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		// prehandle
		token := c.GetHeader("Token")

		tokenPayloadClaims, msgCode, _ := utils.TokenDecode(token)

		if msgCode.Code == code.TokenInvalid || msgCode.Code == code.ServerError {
			code.GinEmptyMsgCode(c, msgCode)
			c.Abort()
			return
		}

		if msgCode.Code == code.TokenExpired {
			// 后期改成更新模式
			code.GinEmptyMsgCode(c, msgCode)
			c.Abort()
			return
		}

		c.Set("ID", tokenPayloadClaims.TokenPayload.ID)
		c.Set("Type", tokenPayloadClaims.TokenPayload.Type)
		c.Set("Phone", tokenPayloadClaims.TokenPayload.Phone)

		c.Next()

		// afterhandle

	}
}

func CustomerOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, Type, _ := utils.GetTokenInfo(c)

        if Type != model.UserTypeCustomer {
            code.GinUnAuthorized(c)
            c.Abort()
            return
        }

		c.Next()
	}
}

func MerchantOnly() gin.HandlerFunc {
    return func(c *gin.Context) {
        _, Type, _ := utils.GetTokenInfo(c)

        if Type != model.UserTypeMerchant {
            code.GinUnAuthorized(c)
            c.Abort()
            return
        }
        c.Next()
    }
}
