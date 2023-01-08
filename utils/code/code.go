package code

import "github.com/gin-gonic/gin"

// common Code :10
var (
	OK = 20000

	// commodity codes

	NotOnShelf = 20010
	NotEnough  = 20011

	// sharebill

	ShareBillDone = 20020
	UserInTeam    = 20021

	BadRequest           = 40000
	InvalidIdentity      = 40001
	InvalidPhone         = 40002
	InvalidPassword      = 40003
	UserExist            = 40010
	UserNotExist         = 40011
	ShopNameExist        = 40012
	PhoneExist           = 40013
	MissingItems         = 40021
	MissingCart          = 40022
	MissingShareBill     = 40023
	MissingOrder         = 40024
	MissingMerchant      = 40025
	MissingAddress       = 40026
	OrderNotInDue        = 40031
	OrderNotInCommodity  = 40032
	UnAuthorized         = 40300
	TokenInvalid         = 40301
	TokenExpired         = 40302
	PhoneORPasswordError = 40310
	UnMatchedID          = 40320

	ServerError = 50000
	InsertError = 50001
	DropError   = 50002
	CheckError  = 50003
	UpdateError = 50004
	DBEmpty     = 60001
) // auth Code: 20

type MsgCode struct {
	Msg  string
	Code int
}

func ginFunctionFactor(msgCode *MsgCode) func(c *gin.Context) {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": msgCode.Code,
			"msg":  msgCode.Msg,
			"data": gin.H{},
		})
	}
}

var (
	GinOKEmpty    = ginFunctionFactor(&MsgCode{Msg: "OK", Code: OK})
	GinNotOnShelf = ginFunctionFactor(&MsgCode{Msg: "NotOnShelf", Code: NotOnShelf})
	GinNotEnough  = ginFunctionFactor(&MsgCode{Msg: "NotEnough", Code: NotEnough})

	GinShareBillDone = ginFunctionFactor(&MsgCode{Msg: "ShareBillDone", Code: ShareBillDone})
	GinUserInTeam    = ginFunctionFactor(&MsgCode{Msg: "UserInTeam", Code: UserInTeam})

	GinBadRequest          = ginFunctionFactor(&MsgCode{Msg: "BadRequest", Code: BadRequest})
	GinMissingItems        = ginFunctionFactor(&MsgCode{Msg: "MissingItems", Code: MissingItems})
	GinMissingCart         = ginFunctionFactor(&MsgCode{Msg: "MissingCart", Code: MissingCart})
	GinMissingShareBill    = ginFunctionFactor(&MsgCode{Msg: "MissingShareBill", Code: MissingShareBill})
	GinMissingOrder        = ginFunctionFactor(&MsgCode{Msg: "MissingOrder", Code: MissingOrder})
	GinMissingMerchant     = ginFunctionFactor(&MsgCode{Msg: "MissingMerchant", Code: MissingMerchant})
	GinMissingAddress      = ginFunctionFactor(&MsgCode{Msg: "MissingAddress", Code: MissingAddress})
	GinOrderNotInDue       = ginFunctionFactor(&MsgCode{Msg: "OrderNotInDue", Code: OrderNotInDue})
	GinOrderNotInCommodity = ginFunctionFactor(&MsgCode{Msg: "OrderNotInCommodity", Code: OrderNotInCommodity})

	GinUnAuthorized         = ginFunctionFactor(&MsgCode{Msg: "UnAuthorized", Code: UnAuthorized})
	GinUserNotExist         = ginFunctionFactor(&MsgCode{Msg: "UserNotExist", Code: UserNotExist})
	GinPhoneORPasswordError = ginFunctionFactor(&MsgCode{Msg: "PhoneORPasswordError", Code: PhoneORPasswordError})
	GinUnMatchedID          = ginFunctionFactor(&MsgCode{Msg: "UnMatchedID", Code: UnMatchedID})

	GinServerError = ginFunctionFactor(&MsgCode{Msg: "ServerError", Code: ServerError})
)

func GinEmptyMsgCode(c *gin.Context, msgCode *MsgCode) {
	c.JSON(200, gin.H{
		"code": msgCode.Code,
		"msg":  msgCode.Msg,
		"data": gin.H{},
	})
}

func GinOKPayload(c *gin.Context, payload *gin.H) {
	c.JSON(200, gin.H{
		"code": OK,
		"msg":  "OK",
		"data": *payload,
	})
}

func GinOKPayloadAny(c *gin.Context, payload any) {
	c.JSON(200, gin.H{
		"code": OK,
		"msg":  "OK",
		"data": payload,
	})
}

/*
400	Bad Request	表示其他错误，就是4xx都无法描述的前端发生的错误
401	Authentication	表示认证类型的错误
403	Authorization	表示授权的错误（认证和授权的区别在于：认证表示“识别前来访问的是谁”，而授权则是“赋予特定用户执行特定操作的权限”）
404	Not Found	表示访问的数据不存在
405	Method Not Allowed	表示可以访问接口，但是使用的HTTP方法不允许
408	Request Timeout	表示前端发送的请求到服务器所需的时间太长
409	Conflict	表示资源发生了冲突，比如使用已被注册邮箱地址注册时，就引起冲突
410	Gone	表示访问的资源不存在。不单表示资源不存在，还进一步告知该资源该资源曾经存在但目前已消失
415	Unsupported Media Type	表示服务器端不支持客户端请求首部Content-Type里指定的数据格式
416	Range Not Satisfiable	表示无法提供Range请求中的指定的那段包体
429	Too Many Requests	表示客户端发送请求的速率过快


状态码	名称	说明
500	Internal Server Error	表示服务器内部错误，且不属于以下错误类型
502	Bad Gateway	代理服务器无法获取到合法资源
504	Gateway Timeout	表示代理服务器无法及时的从上游获得响应
507	Insufficient Storage	表示服务器没有足够的空间处理请求
*/
