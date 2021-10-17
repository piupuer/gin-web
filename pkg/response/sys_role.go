package response

import "github.com/piupuer/go-helper/pkg/resp"

type RoleResp struct {
	resp.Base
	Name    string `json:"name"`
	Keyword string `json:"keyword"`
	Sort    uint   `json:"sort"`
	Desc    string `json:"desc"`
	Status  *uint  `json:"status"`
	Creator string `json:"creator"`
}
