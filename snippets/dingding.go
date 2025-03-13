package snippets

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-snippets-test/config"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

type DingTalkBot struct {
	webhook      string        // 机器人 Webhook 地址
	secret       string        // 安全加签密钥
	messageCache sync.Map      // 消息去重缓存
	cacheTTL     time.Duration // 消息有效时间
	logger       *log.Logger   // 日志记录
}

func NewDingTalkBot() *DingTalkBot {
	// 读取配置文件
	_ = config.InitConfig()
	webhook := config.GetConfig().DingBot.Webhook
	secret := config.GetConfig().DingBot.Secret

	bot := &DingTalkBot{
		webhook:  webhook,
		secret:   secret,
		cacheTTL: time.Minute, // 消息缓存 1 分钟
		logger:   log.Default(),
	}
	go bot.cleanupMessageCache() // 启动清理任务
	return bot
}

func (bot *DingTalkBot) cleanupMessageCache() {
	for {
		time.Sleep(bot.cacheTTL) // 每隔 1 分钟清理一次
		bot.messageCache.Range(func(key, value interface{}) bool {
			if time.Since(value.(time.Time)) > bot.cacheTTL {
				bot.messageCache.Delete(key)
			}
			return true
		})
	}
}

func (bot *DingTalkBot) sign() string {
	if bot.secret == "" {
		return ""
	}
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10) // 毫秒级时间戳
	stringToSign := timestamp + "\n" + bot.secret
	h := hmac.New(sha256.New, []byte(bot.secret))
	h.Write([]byte(stringToSign))
	sign := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf("&timestamp=%s&sign=%s", timestamp, url.QueryEscape(sign))
}

func (bot *DingTalkBot) sendMessage(msgType string, message map[string]interface{}) {
	bot.logger.Println(fmt.Sprintf("开始发送%s类型的消息...", msgType))
	requestURL := bot.webhook + bot.sign()
	// 组装消息内容
	payload, err := json.Marshal(message)
	if err != nil {
		bot.logger.Println("JSON 消息解析失败:", err)
		return
	}

	msgKey := string(payload) // 直接用 JSON 字符串做 Key
	// **1. 检查是否 1 分钟内发送过**
	if lastTime, exists := bot.messageCache.Load(msgKey); exists {
		if time.Since(lastTime.(time.Time)) < bot.cacheTTL {
			bot.logger.Println("⚠️ 1 分钟内已发送过相同消息，跳过发送")
			return
		}
	}
	// **2. 记录当前消息的时间**
	bot.messageCache.Store(msgKey, time.Now())

	// 发送 HTTP 请求
	resp, err := http.Post(requestURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		bot.logger.Println("发送失败:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode != 200 {
		bot.logger.Println(fmt.Sprintf("钉钉 API 返回错误: %d", resp.StatusCode))
	} else {
		bot.logger.Println("[DingDing] 告警发送成功:", message)
	}
}

// SendTextMessage 发送 Text 消息
func (bot *DingTalkBot) SendTextMessage(content string, atMobiles []string, isAtAll bool) {
	message := map[string]interface{}{
		"at": map[string]interface{}{
			"atMobiles": atMobiles, // @指定手机号的用户
			"isAtAll":   isAtAll,   // 是否 @所有人
		},
		"msgtype": "text",
		"text": map[string]string{
			"content": content,
		},
	}
	bot.sendMessage("Text", message)
}

// SendMarkdownMessage 发送 Markdown 消息
func (bot *DingTalkBot) SendMarkdownMessage(title, text string, atMobiles []string, isAtAll bool) {
	message := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"title": title,
			"text":  text,
		},
		"at": map[string]interface{}{
			"atMobiles": atMobiles,
			"isAtAll":   isAtAll,
		},
	}
	bot.sendMessage("Markdown", message)
}

// SendLinkMessage 发送链接消息
func (bot *DingTalkBot) SendLinkMessage(title, text, messageURL, picURL string) {
	message := map[string]interface{}{
		"msgtype": "link",
		"link": map[string]string{
			"title":      title,
			"text":       text,
			"messageUrl": messageURL,
			"picUrl":     picURL,
		},
	}
	bot.sendMessage("Link", message)
}

// SendActionCardMessage 发送 ActionCard 消息
func (bot *DingTalkBot) SendActionCardMessage(title, text, singleTitle, singleURL string) {
	message := map[string]interface{}{
		"msgtype": "actionCard",
		"actionCard": map[string]string{
			"title":          title,
			"text":           text,
			"singleTitle":    singleTitle,
			"singleURL":      singleURL,
			"btnOrientation": "0",
		},
	}
	bot.sendMessage("ActionCard", message)
}

// SendFeedCardMessage 发送 FeedCard 消息 (多链接消息)
func (bot *DingTalkBot) SendFeedCardMessage() {
	message := map[string]interface{}{
		"msgtype": "feedCard",
		"feedCard": map[string]interface{}{
			"links": []map[string]string{
				{
					"title":      "示例文章 1",
					"messageURL": "https://www.example.com/article1",
					"picURL":     "https://www.example.com/image1.jpg",
				},
				{
					"title":      "示例文章 2",
					"messageURL": "https://www.example.com/article2",
					"picURL":     "https://www.example.com/image2.jpg",
				},
			},
		},
	}
	bot.sendMessage("FeedCard", message)
}
