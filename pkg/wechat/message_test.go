package wechat

import (
	"gin-web/models"
	"gin-web/tests"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"testing"
	"time"
)

func TestSendTplMessage(t *testing.T) {
	tests.InitTestEnv()
	msg := message.TemplateMessage{
		ToUser:     "xxx",
		TemplateID: "xxx",
		Data: map[string]*message.TemplateDataItem{
			"first": {
				Value: "日常事项定时提醒",
			},
			"keyword1": {
				Value: "每日购买",
			},
			"keyword2": {
				Value: "请到商城下单支付一单(杨博士店有一分钱的单)",
			},
			"keyword3": {
				Value: models.LocalTime{
					Time: time.Now(),
				}.String(),
			},
			"remark": {
				Value: "下单完成记得将截图发到群里哦~",
			},
		},
	}
	msg.MiniProgram.AppID = "xxx"
	msg.MiniProgram.PagePath = "pages/index/index"
	err := SendTplMessage(&msg)
	if err != nil {
		panic(err)
	}
}
