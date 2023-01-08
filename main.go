package main

import (
	"HappyShopTogether/model"
	"HappyShopTogether/router"
    "HappyShopTogether/utils"
    "fmt"
    "log"
)

// basic setting

func main() {
	model.Db.Init() 
	defer model.Db.Close()

	router.Init()
    ginConfig := fmt.Sprintf(":%d", utils.GlobalConfig.GinConfig.Port)
	if err := router.Router.Run(ginConfig); err != nil {
		log.Print("路由错误")
		log.Fatal(err)
	}

}
