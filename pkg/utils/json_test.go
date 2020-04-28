package utils

import (
	"go-shipment-api/models"
	"reflect"
	"testing"
)

func TestStruct2Json(t *testing.T) {
	type args struct {
		obj interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "测试user转换为json case1",
			args: args{obj: &models.SysUser{
				Model: models.Model{
					Id: 1,
				},
				Username: "张三",
				Password: "123456",
			}},
			want: `{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"username":"张三","password":"123456"}`,
		},
		{
			name: "测试user转换为json case2",
			args: args{obj: &models.SysUser{
				Model: models.Model{
					Id: 2,
				},
				Username: "李四",
				Password: "654321",
			}},
			want: `{"ID":2,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"username":"李四","password":"654321"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Struct2Json(tt.args.obj); got != tt.want {
				t.Errorf("Struct2Json() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJsonI2Struct(t *testing.T) {
	type args struct {
		str string
		obj interface{}
	}
	var obj models.SysUser
	tests := []struct {
		name string
		args args
	}{
		{
			name: "测试json转换为user case1",
			args: args{
				str: `{"ID":1,"CreatedAt":"0001-01-01T00:00:00Z","UpdatedAt":"0001-01-01T00:00:00Z","DeletedAt":null,"username":"张三","password":"123456"}`,
				obj: obj,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if JsonI2Struct(tt.args.str, &tt.args.obj); reflect.ValueOf(&tt.args.obj).IsNil() {
				t.Errorf("转换失败, struct为空")
			}
		})
	}
}
