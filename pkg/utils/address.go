package utils

import (
	"fmt"
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
	resp, err := http.Get(fmt.Sprintf("https://restapi.amap.com/v3/ip?ip=%s&key=%s", ip, "9130aaac2b7a920b8bbd5dc9647fbe9e"))
	address := "未知地址"
	if err != nil {
		fmt.Println(fmt.Sprintf("[GetIpRealLocation]IP地址查询失败: %v", err))
		return address
	}
	defer resp.Body.Close()
	// 读取响应数据
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(fmt.Sprintf("[GetIpRealLocation]IP地址查询失败: %v", err))
		return address
	}
	// json数据转结构体
	var result IpResp
	Json2Struct(string(data), &result)
	if result.Status == "1" {
		address = result.Province
		// 城市不为空且城市与省份不重复
		if result.City != "" && result.Province != result.City {
			address += result.City
		}
	}
	return address
}
