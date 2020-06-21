package redis

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/tests"
	"testing"
)

func TestQueryRedis_Count(t *testing.T) {
	tests.InitTestEnv()
	query := New()
	var count uint
	err := global.Mysql.Count(&count).Error
	fmt.Println(err)
	err2 := query.Count(&count).Error
	fmt.Println(err2)
}

func TestQueryRedis_Table(t *testing.T) {
	tests.InitTestEnv()
	query := New()
	var u models.SysUser
	tableName := u.TableName()
	var count uint
	err := global.Mysql.Table(tableName).Count(&count).Error
	fmt.Println(err, count)
	var count2 uint
	err2 := query.Table(tableName).Count(&count2).Error
	fmt.Println(err2, count2)
}

func TestQueryRedis_Find(t *testing.T) {
	tests.InitTestEnv()
	query := New()
	var u models.SysUser
	var us []models.SysUser
	tableName := u.TableName()
	err := global.Mysql.Table(tableName).Where("id > ?", uint(0)).Find(&us).Error
	fmt.Println(err, us)
	var us2 []models.SysUser
	err2 := query.Table(tableName).Where("id", ">", 0).Find(&us2).Error
	fmt.Println(err2, us2)
}

func TestQueryRedis_First(t *testing.T) {
	tests.InitTestEnv()
	query := New()
	var u models.SysUser
	var us models.SysUser
	tableName := u.TableName()
	err := global.Mysql.Table(tableName).Where("id > ?", uint(0)).First(&us).Error
	fmt.Println(err, us)
	var us2 models.SysUser
	err2 := query.Table(tableName).Where("id", ">", 0).First(&us2).Error
	fmt.Println(err2, us2)
}

func TestQueryRedis_Preload(t *testing.T) {
	// 测试preload belongsTo关联
	tests.InitTestEnv()
	query := New()
	var u models.SysUser
	var us []models.SysUser
	tableName := u.TableName()
	err := global.Mysql.Table(tableName).Preload("Role").Find(&us).Error
	fmt.Println(err, us)
	var us2 []models.SysUser
	err2 := query.Table(tableName).Preload("Role").Find(&us2).Error
	fmt.Println(err2, us2)
}

func TestQueryRedis_Preload1(t *testing.T) {
	// 测试preload hasMany关联
	tests.InitTestEnv()
	query := New()
	var u models.SysRole
	var us []models.SysRole
	tableName := u.TableName()
	err := global.Mysql.Table(tableName).Preload("Users").Find(&us).Error
	fmt.Println(err, us)
	var us2 []models.SysRole
	err2 := query.Table(tableName).Preload("Users").Find(&us2).Error
	fmt.Println(err2, us2)
}

func TestQueryRedis_Preload2(t *testing.T) {
	// 测试preload 多个字段
	tests.InitTestEnv()
	query := New()
	var u models.SysUser
	var us []models.SysUser
	tableName := u.TableName()
	err := global.Mysql.Table(tableName).Preload("Role").Preload("Role.Users").Find(&us).Error
	fmt.Println(err, us)
	var us2 []models.SysUser
	err2 := query.Table(tableName).Preload("Role").Preload("Role.Users").Find(&us2).Error
	fmt.Println(err2, us2)
}

func TestQueryRedis_Preload3(t *testing.T) {
	// 测试preload many2many关联
	tests.InitTestEnv()
	query := New()
	var u models.SysRole
	var us models.SysRole
	tableName := u.TableName()
	err := global.Mysql.Table(tableName).Preload("Menus").Where("id = ?", 1).First(&us).Error
	fmt.Println(err, us)
	var us2 models.SysRole
	err2 := query.Table(tableName).Preload("Menus").Where("id", "=", 1).First(&us2).Error
	fmt.Println(err2, us2)
}
