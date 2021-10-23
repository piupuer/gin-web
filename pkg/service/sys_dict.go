package service

import (
	"fmt"
	"gin-web/pkg/request"
	"github.com/piupuer/go-helper/ms"
	"gorm.io/gorm"
	"strings"
)

func (my MysqlService) GetDictData(dictName, dictDataKey string) ms.SysDictData {
	dict, err := my.GetDictDataWithErr(dictName, dictDataKey)
	if err != nil || dict == nil {
		return ms.SysDictData{}
	}
	return *dict
}

func (my MysqlService) GetDictDataWithErr(dictName, dictDataKey string) (*ms.SysDictData, error) {
	oldCache, ok := CacheGetDictNameAndKey(my.Q.Ctx, dictName, dictDataKey)
	if ok {
		return oldCache, nil
	}
	list := make([]ms.SysDictData, 0)
	my.Q.Tx.
		Model(&ms.SysDictData{}).
		Preload("Dict").
		Order("created_at DESC").
		Find(&list)
	for _, data := range list {
		if data.Dict.Name == dictName && data.Key == dictDataKey {
			CacheSetDictNameAndKey(my.Q.Ctx, dictName, dictDataKey, data)
			return &data, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (my MysqlService) FindDictDataByName(name string) ([]ms.SysDictData, error) {
	oldCache, ok := CacheGetDictName(my.Q.Ctx, name)
	if ok {
		return oldCache, nil
	}
	list := make([]ms.SysDictData, 0)
	my.Q.Tx.
		Model(&ms.SysDictData{}).
		Preload("Dict").
		Order("sort").
		Find(&list)
	newList := make([]ms.SysDictData, 0)
	for _, data := range list {
		if data.Dict.Name == name {
			newList = append(newList, data)
		}
	}
	CacheSetDictName(my.Q.Ctx, name, newList)
	return newList, nil
}

func (my MysqlService) FindDict(req *request.DictReq) []ms.SysDict {
	list := make([]ms.SysDict, 0)
	query := my.Q.Tx.
		Model(&ms.SysDict{}).
		Preload("DictDatas").
		Order("created_at DESC")
	name := strings.TrimSpace(req.Name)
	if name != "" {
		query = query.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	desc := strings.TrimSpace(req.Desc)
	if desc != "" {
		query = query.Where("desc = ?", desc)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	my.Q.FindWithPage(query, &req.Page, &list)
	return list
}

func (my MysqlService) FindDictData(req *request.DictDataReq) []ms.SysDictData {
	list := make([]ms.SysDictData, 0)
	query := my.Q.Tx.
		Model(&ms.SysDictData{}).
		Preload("Dict").
		Order("created_at DESC")
	key := strings.TrimSpace(req.Key)
	if key != "" {
		query = query.Where("key LIKE ?", fmt.Sprintf("%%%s%%", key))
	}
	val := strings.TrimSpace(req.Val)
	if val != "" {
		query = query.Where("val LIKE ?", fmt.Sprintf("%%%s%%", val))
	}
	attr := strings.TrimSpace(req.Attr)
	if attr != "" {
		query = query.Where("attr = ?", attr)
	}
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}
	if req.DictId != nil {
		query = query.Where("dict_id = ?", *req.DictId)
	}
	my.Q.FindWithPage(query, &req.Page, &list)
	return list
}

func (my MysqlService) CreateDict(req *request.CreateDictReq) (err error) {
	err = my.Q.Create(req, new(ms.SysDict))
	CacheFlushDictName(my.Q.Ctx)
	CacheFlushDictNameAndKey(my.Q.Ctx)
	return
}

func (my MysqlService) UpdateDictById(id uint, req request.UpdateDictReq) (err error) {
	err = my.Q.UpdateById(id, req, new(ms.SysDict))
	CacheFlushDictName(my.Q.Ctx)
	CacheFlushDictNameAndKey(my.Q.Ctx)
	return
}

func (my MysqlService) DeleteDictByIds(ids []uint) (err error) {
	err = my.Q.DeleteByIds(ids, new(ms.SysDict))
	CacheFlushDictName(my.Q.Ctx)
	CacheFlushDictNameAndKey(my.Q.Ctx)
	return
}

func (my MysqlService) CreateDictData(req *request.CreateDictDataReq) (err error) {
	err = my.Q.Create(req, new(ms.SysDictData))
	CacheFlushDictName(my.Q.Ctx)
	CacheFlushDictNameAndKey(my.Q.Ctx)
	return
}

func (my MysqlService) UpdateDictDataById(id uint, req request.UpdateDictDataReq) (err error) {
	err = my.Q.UpdateById(id, req, new(ms.SysDictData))
	CacheFlushDictName(my.Q.Ctx)
	CacheFlushDictNameAndKey(my.Q.Ctx)
	return
}

func (my MysqlService) DeleteDictDataByIds(ids []uint) (err error) {
	err = my.Q.DeleteByIds(ids, new(ms.SysDictData))
	CacheFlushDictName(my.Q.Ctx)
	CacheFlushDictNameAndKey(my.Q.Ctx)
	return
}
