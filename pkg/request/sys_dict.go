package request

import (
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

type DictReq struct {
	Name   string        `json:"name" form:"name"`
	Desc   string        `json:"desc" form:"desc"`
	Status *req.NullUint `json:"status" form:"status"`
	resp.Page
}

type CreateDictReq struct {
	Name   string        `json:"name" validate:"required"`
	Desc   string        `json:"desc" validate:"required"`
	Status *req.NullUint `json:"status"`
	Remark string        `json:"remark"`
}

func (s CreateDictReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "字典名称"
	m["Desc"] = "字典描述"
	return m
}

type UpdateDictReq struct {
	Name   *string       `json:"name"`
	Desc   *string       `json:"desc"`
	Status *req.NullUint `json:"status"`
	Remark *string       `json:"remark"`
}

type DictDataReq struct {
	DictId *req.NullUint `json:"dictId" form:"dictId"`
	Key    string        `json:"key" form:"key"`
	Attr   string        `json:"attr" form:"attr"`
	Val    string        `json:"val" form:"val"`
	Status *req.NullUint `json:"status" form:"sort"`
	resp.Page
}

type CreateDictDataReq struct {
	Key      string `json:"key" validate:"required"`
	Val      string `json:"val" validate:"required"`
	Attr     string `json:"attr"`
	Addition string `json:"addition"`
	Sort     *uint  `json:"sort"`
	Status   *uint  `json:"status"`
	Remark   string `json:"remark"`
	DictId   uint   `json:"dictId" validate:"required"`
}

func (s CreateDictDataReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Key"] = "数据键"
	m["Val"] = "数据值"
	m["DictId"] = "字典编号"
	return m
}

type UpdateDictDataReq struct {
	Key      *string       `json:"key"`
	Val      *string       `json:"val"`
	Attr     *string       `json:"attr"`
	Addition *string       `json:"addition"`
	Sort     *req.NullUint `json:"sort"`
	Status   *req.NullUint `json:"status"`
	Remark   *string       `json:"remark"`
	DictId   *req.NullUint `json:"dictId"`
}
