### 关于 jwt

安装

``` shell
go install "github.com/golang-jwt/jwt/v4"
```

生成Token
定义claims和serect

``` golang
type MyClaims struct {
    Phone string `json:"phone"`
    jwt.RegisteredClaims  // 注意!这是jwt-go的v4版本新增的，原先是jwt.StandardClaims
}

var MySecret = []byte("手写的从前") // 定义secret，后面会用到

```

iss (issuer)：签发人
exp (expiration time)：过期时间
sub (subject)：主题
aud (audience)：受众
nbf (Not Before)：生效时间
iat (Issued At)：签发时间
jti (JWT ID)：编号
