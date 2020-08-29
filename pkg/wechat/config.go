package wechat

import (
	"fmt"
	"gin-web/pkg/global"
	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	"github.com/silenceper/wechat/v2/officialaccount"
	"github.com/silenceper/wechat/v2/officialaccount/config"
)

// 获取公众号配置
func GetOfficialAccount() *officialaccount.OfficialAccount {
	wc := wechat.NewWechat()

	var cacheIns cache.Cache
	if global.Conf.System.UseRedis {
		redisOpts := &cache.RedisOpts{
			Host:     fmt.Sprintf("%s:%d", global.Conf.Redis.Host, global.Conf.Redis.Port),
			Database: global.Conf.Redis.Database,
			Password: global.Conf.Redis.Password,
		}
		cacheIns = cache.NewRedis(redisOpts)
	} else {
		cacheIns = cache.NewMemory()
	}
	cfg := &config.Config{
		AppID:          global.Conf.WeChat.Official.AppId,
		AppSecret:      global.Conf.WeChat.Official.AppSecret,
		EncodingAESKey: global.Conf.WeChat.Official.Encoding,
		Cache:          cacheIns,
	}
	return wc.GetOfficialAccount(cfg)
}
