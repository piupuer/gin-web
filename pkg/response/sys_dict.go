package response

import "github.com/piupuer/go-helper/pkg/resp"

type DictResp struct {
	resp.Base
	Name      string         `json:"name"`
	Desc      string         `json:"desc"`
	Status    uint           `json:"status"`
	Remark    string         `json:"remark"`
	DictDatas []DictDataResp `json:"dictDatas"`
}

type DictDataResp struct {
	resp.Base
	Key      string   `json:"key"`
	Val      string   `json:"val"`
	Attr     string   `json:"attr"`
	Addition string   `json:"addition"`
	Sort     uint     `json:"sort"`
	Status   uint     `json:"status"`
	Remark   string   `json:"remark"`
	DictId   uint     `json:"dictId"`
	Dict     DictResp `json:"dict"`
}
