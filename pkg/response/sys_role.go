package response

type RoleResp struct {
	BaseData
	Name    string `json:"name"`
	Keyword string `json:"keyword"`
	Sort    uint   `json:"sort"`
	Desc    string `json:"desc"`
	Status  *uint  `json:"status"`
	Creator string `json:"creator"`
}
