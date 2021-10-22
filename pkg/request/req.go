package request

import (
	"gin-web/models"
	"gin-web/pkg/utils"
)

type UpdateMenuIncrementalIdsReq struct {
	Create []uint `json:"create"`
	Delete []uint `json:"delete"`
}

func (in UpdateMenuIncrementalIdsReq) FindIncremental(oldMenuIds []uint, allMenu []models.SysMenu) []uint {
	in.Create = models.FindCheckedMenuId(in.Create, allMenu)
	in.Delete = models.FindCheckedMenuId(in.Delete, allMenu)
	newList := make([]uint, 0)
	for _, oldItem := range oldMenuIds {
		// not in delete
		if !utils.Contains(in.Delete, oldItem) {
			newList = append(newList, oldItem)
		}
	}
	// need create
	return append(newList, in.Create...)
}
