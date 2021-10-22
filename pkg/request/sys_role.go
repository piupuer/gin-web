package request

import (
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

type RoleReq struct {
	Name            string `json:"name" form:"name"`
	Keyword         string `json:"keyword" form:"keyword"`
	CurrentRoleSort uint   `json:"currentRoleSort"`
	Status          *uint  `json:"status" form:"status"`
	resp.Page
}

type CreateRoleReq struct {
	CurrentRoleSort uint          `json:"currentRoleSort"`
	Name            string        `json:"name" validate:"required"`
	Keyword         string        `json:"keyword" validate:"required"`
	Sort            *req.NullUint `json:"sort" validate:"required"`
	Desc            string        `json:"desc"`
	Status          *req.NullUint `json:"status"`
}

func (s CreateRoleReq) FieldTrans() map[string]string {
	m := make(map[string]string, 0)
	m["Name"] = "name"
	m["Keyword"] = "keyword"
	m["Sort"] = "sort"
	return m
}

type UpdateRoleReq struct {
	Name    *string       `json:"name"`
	Keyword *string       `json:"keyword"`
	Sort    *req.NullUint `json:"sort"`
	Desc    *string       `json:"desc"`
	Status  *req.NullUint `json:"status"`
}
