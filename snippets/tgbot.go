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

func QTelegramBot() {
	_ = config.InitConfig()
	token := config.GetConfig().TgBot.Token
	proxy := config.GetConfig().Proxy.Url

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
	bot.Debug = true // 开启 Debug 模式，观察日志

	// 获取 Bot 更新
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60 // 设置超时时间
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		// 处理消息
		if update.Message != nil {
			msgText := update.Message.Text   // 用户输入的文本
			chatID := update.Message.Chat.ID // 聊天 ID

			var replyText string

			// 简单的自动回复逻辑
			switch msgText {
			case "你好":
				replyText = "你好，有什么可以帮助你的？😊"
			case "帮助":
				replyText = "这是一个 FAQ Bot，你可以问我一些常见问题！📌"
			case "你是谁":
				replyText = "我是你的 Telegram 机器人 🤖"
			default:
				replyText = "抱歉，我不明白你的意思。可以输入 '帮助' 查看支持的指令。"
			}

			// 发送消息
			msg := tgbotapi.NewMessage(chatID, replyText)
			bot.Send(msg)
		}
	}
}
