package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	DbConfig  DbConfig  `json:"dbConfig"`
	GinConfig GinConfig `json:"ginConfig"`
	SecretKey SecretKey `json:"secretKey"`
	Global    global    `json:"global"`
}

type DbConfig struct {
	Host         string `json:"host"`
	Password     string `json:"password"`
	Port         uint   `json:"port"`
	DatabaseName string `json:"db"`
}

type GinConfig struct {
	Port         uint `json:"port"`
	TokenExpires int  `json:"tokenExpires"`
}

type SecretKey struct {
	Public  string
	Private string
}

type global struct {
	GroupMemberCount   uint `json:"groupMemberCount"`
	GroupForMemberTime uint `josn:"groupForMemberTime"`
}

func ConfigRead() (Config, error) {
	jsonFile, err := os.Open("config.json")

	if err != nil {
		fmt.Println("文件打开失败")
		return Config{}, err
	}
	defer jsonFile.Close()

	jsonData, err := ioutil.ReadAll(jsonFile)

	var config Config
	json.Unmarshal(jsonData, &config)

	if err != nil {
		fmt.Println("文件读取失败")
		return Config{}, err
	}
	return config, nil
}

var GlobalConfig, _ = ConfigRead()
