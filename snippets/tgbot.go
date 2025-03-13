package snippets

import (
	"fmt"
	"go-snippets-test/config"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func RunTelegramBot() {
	// 读取配置文件
	_ = config.InitConfig()
	proxy := config.GetConfig().Proxy.Url
	token := config.GetConfig().TgBot.Token
	userId := config.GetConfig().TgBot.UserId

	proxyURL, err := url.Parse(proxy) // 替换为你的代理地址和端口
	if err != nil {
		log.Fatal("代理地址解析失败:", err)
	}

	// 创建一个 HTTP 客户端，带有代理配置
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}

	// 使用你的 API Token 替换
	bot, err := tgbotapi.NewBotAPIWithClient(token, client)
	if err != nil {
		log.Fatal(err)
	}

	// 设置接收消息的用户 ID（通过 getUpdates 获取）
	userID, _ := strconv.ParseInt(userId, 10, 64)
	// 创建消息对象
	msg := tgbotapi.NewMessage(userID, "这是来自我的服务器的消息！哈哈哈哈")

	// 发送消息
	_, err = bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("消息已发送！")
}
