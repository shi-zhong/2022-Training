package utils

import (
	"HappyShopTogether/utils/code"
	"crypto/md5"
	"encoding/hex"
	"github.com/golang-jwt/jwt/v4"
	mRand "math/rand"
	"time"
)

type TokenPayload struct {
	Type  uint8
	ID    uint
	Phone string
}
type TokenPayloadClaims struct {
	TokenPayload
	jwt.RegisteredClaims // 注册当前结构体为 Claims
}

func TokenEecode(payload *TokenPayload) (string, error) {
	// 创建秘钥
	key := []byte(GlobalConfig.SecretKey.Private)

	// 创建Token结构体
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenPayloadClaims{
		TokenPayload: *payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(GlobalConfig.GinConfig.TokenExpires) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	})
	// 调用加密方法，发挥Token字符串
	signedString, err := token.SignedString(key)
	return signedString, err
}

/*
newWithClaims
type Token struct {
    Raw       string        //原始令牌
    Method    SigningMethod   // 加密方法 比如sha256加密
    Header    map[string]interface{} // token头信息
    Claims    Claims  // 加密配置，比如超时时间等
    Signature string  // 加密后的字符串
    Valid     bool   // 是否校验
}
*/

func TokenDecode(signedString string) (*TokenPayloadClaims, *code.MsgCode, error) {
	// 根据Token字符串解析成Claims结构体
	token, err := jwt.ParseWithClaims(signedString, &TokenPayloadClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(GlobalConfig.SecretKey.Private), nil
	})

	// 简要错误处理
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 || ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, &code.MsgCode{
					Msg: "TokenInvalid", Code: code.TokenInvalid,
				}, err
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, &code.MsgCode{Msg: "TokenExpired", Code: code.TokenExpired}, err
			} else {
				return nil, &code.MsgCode{Msg: "ServerError", Code: code.ServerError}, err
			}
		}
	}

	if claims, ok := token.Claims.(*TokenPayloadClaims); ok && token.Valid {
		return claims, &code.MsgCode{Msg: "OK", Code: code.OK}, nil
	}
	return nil, &code.MsgCode{Msg: "ServerError", Code: code.ServerError}, err
}

func Salt(len int) string {
	var salt string
	var j int
	mRand.Seed(time.Now().UnixNano())
	for i := 0; i < len; i++ {
		j = mRand.Intn(52)
		if j < 26 {
			salt += string(rune(j + 65))
		} else {
			salt += string(rune(j + 71))
		}
	}
	return salt
}

func PasswordMD5(pwd, salt string) (string, error) {

	hash := md5.New()

	_, err := hash.Write([]byte(pwd + salt))

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
