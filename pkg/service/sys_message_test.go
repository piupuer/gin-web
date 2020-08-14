package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/tests"
	"testing"
)

func TestMysqlService_BatchCreateOneToOneMessage(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)

	message := models.SysMessage{
		FromUserId: 1,
		Content:    "你好, 欢迎使用",
	}

	err := s.BatchCreateOneToOneMessage(message, []uint{1, 2, 3, 4})
	if err != nil {
		panic(err)
	}
}

func TestMysqlService_BatchCreateOneToManyMessage(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)

	message := models.SysMessage{
		FromUserId: 1,
		Content:    "你们好, 欢迎使用",
	}

	err := s.BatchCreateOneToManyMessage(message, []uint{1, 2})
	if err != nil {
		panic(err)
	}
}

func TestMysqlService_CreateSystemMessage(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)

	message := models.SysMessage{
		FromUserId: 1,
		Content:    "大家好, 出一则通知",
	}

	err := s.CreateSystemMessage(message)
	if err != nil {
		panic(err)
	}
}

func TestMysqlService_UpdateMessageByUserId(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)
	err := s.SyncMessageByUserIds([]uint{1})
	if err != nil {
		panic(err)
	}
}

func TestMysqlService_GetUnReadMessages(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)
	list, err := s.GetUnReadMessages(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
}

func TestMysqlService_GetReadMessages(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)
	list, err := s.GetReadMessages(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
}

func TestMysqlService_GetDeletedMessages(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)
	list, err := s.GetDeletedMessages(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(list)
}

func TestMysqlService_UpdateMessageRead(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)
	err := s.UpdateMessageRead(1)
	if err != nil {
		panic(err)
	}
}

func TestMysqlService_UpdateMessageDeleted(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)
	err := s.UpdateMessageDeleted(1)
	if err != nil {
		panic(err)
	}
}

func TestMysqlService_UpdateAllMessageRead(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)
	err := s.UpdateAllMessageRead(1)
	if err != nil {
		panic(err)
	}
}

func TestMysqlService_UpdateAllMessageDeleted(t *testing.T) {
	tests.InitTestEnv()
	s := New(nil)
	err := s.UpdateAllMessageDeleted(1)
	if err != nil {
		panic(err)
	}
}
