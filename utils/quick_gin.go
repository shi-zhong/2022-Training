package utils

import (
	"HappyShopTogether/utils/code"
	"github.com/gin-gonic/gin"
)

// GetTokenInfo ID Type Phone
func GetTokenInfo(c *gin.Context) (uint, uint8, string) {
	ID, exist := c.Get("ID")
	uintID, _ := ID.(uint)

	Type, exist2 := c.Get("Type")
	uintType, _ := Type.(uint)

	Phone, exist3 := c.Get("Phone")
	uintPhone, _ := Phone.(string)

	if !exist || !exist2 || !exist3 {
		code.GinUnAuthorized(c)
		c.Abort()
		return 0, 0, ""
	}
	return uintID, uint8(uintType), uintPhone
}

func QuickBind(c *gin.Context, structPointer any) bool {
	if err := c.BindJSON(structPointer); err != nil {
		code.GinBadRequest(c)
		c.Abort()
		return false
	}
	return true
}

func QuickBindPath(c *gin.Context, structPointer any) bool {
	if err := c.ShouldBindUri(structPointer); err != nil {
		code.GinBadRequest(c)
		c.Abort()
		return false
	}
	return true
}
