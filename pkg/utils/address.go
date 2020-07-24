package utils

import (
	"fmt"
	"gin-web/pkg/global"
	"io/ioutil"
	"net/http"
)

type IpResp struct {
	Status   string `json:"status"`
	Province string `json:"province"`
	City     string `json:"city"`
}

// 获取IP真实地址
func GetIpRealLocation(ip string) string {
	resp, err := http.Get(fmt.Sprintf("https://restapi.amap.com/v3/ip?ip=%s&key=%s", ip, "0d1c332f93207b553ad59d6bc72e6e88"))
	address := "未知地址"
	if err != nil {
		global.Log.Error(fmt.Sprintf("[GetIpRealLocation]IP地址查询失败: %v", err))
		return address
	}
	defer resp.Body.Close()
	// 读取响应数据
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		global.Log.Error(fmt.Sprintf("[GetIpRealLocation]IP地址查询失败: %v", err))
		return address
	}
	// json数据转结构体
	var json IpResp
	Json2Struct(string(data), &json)
	if json.Status == "1" {
		address = json.Province
		if json.City != "" {
			address += json.City
		}
	}
	return address
}
