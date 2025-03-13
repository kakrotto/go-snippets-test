package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Proxy   Proxy   `yaml:"proxy"`
	DingBot DingBot `yaml:"dingBot"`
	TgBot   TgBot   `yaml:"tgBot"`
}

type Proxy struct {
	Url string `yaml:"url"`
}

type DingBot struct {
	Webhook string `yaml:"webhook"`
	Secret  string `yaml:"secret"`
}

type TgBot struct {
	Token  string `yaml:"token"`
	UserId string `yaml:"user_id"`
}

var Conf Config

func InitConfig() error {
	dir, err := filepath.Abs("..")
	if err != nil {
		log.Fatal(err)
	}

	// 拼接配置文件路径
	configFilePath := filepath.Join(dir, "config.yaml")
	fmt.Println("配置文件路径:", configFilePath)

	// 打开配置文件
	file, err := os.Open(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 使用 yaml 库解析配置文件
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&Conf); err != nil {
		return err
	}
	return nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return &Conf
}
