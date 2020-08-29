package wechat

import (
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func SendTplMessage(msg *message.TemplateMessage) error {
	o := GetOfficialAccount()
	_, err := o.GetTemplate().Send(msg)
	return err
}
