package tests

import (
	"go-snippets-test/snippets"
	"testing"
)

func TestDingDing(t *testing.T) {
	bot := snippets.NewDingTalkBot()
	bot.SendTextMessage("Hello, 这是一条测试消息！", []string{""}, false)
	bot.SendMarkdownMessage("Markdown 测试", "### 这是一条 *Markdown* 消息 \n - 点个赞 👍 \n @所有人", []string{}, true)
	bot.SendLinkMessage("查看详情", "点击这里了解更多信息", "", "")
	bot.SendActionCardMessage("欢迎使用钉钉机器人", "点击下方按钮查看详情", "了解更多", "")
	bot.SendFeedCardMessage()
}
