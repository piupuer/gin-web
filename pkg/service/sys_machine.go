package service

import (
	"errors"
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"gorm.io/gorm"
	"strings"
)

// 获取机器
func (s *MysqlService) GetMachines(req *request.MachineListRequestStruct) ([]models.SysMachine, error) {
	var err error
	list := make([]models.SysMachine, 0)
	db := global.Mysql.
		Table(new(models.SysMachine).TableName()).
		Order("created_at DESC")
	host := strings.TrimSpace(req.Host)
	if host != "" {
		db = db.Where("host LIKE ?", fmt.Sprintf("%%%s%%", host))
	}
	loginName := strings.TrimSpace(req.LoginName)
	if loginName != "" {
		db = db.Where("login_name LIKE ?", fmt.Sprintf("%%%s%%", loginName))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}
	if req.Status != nil {
		if *req.Status > 0 {
			db = db.Where("status = ?", 1)
		} else {
			db = db.Where("status = ?", 0)
		}
	}
	// 查询条数
	err = db.Count(&req.PageInfo.Total).Error
	if err == nil {
		if req.PageInfo.NoPagination {
			// 不使用分页
			err = db.Find(&list).Error
		} else {
			// 获取分页参数
			limit, offset := req.GetLimit()
			err = db.Limit(limit).Offset(offset).Find(&list).Error
		}
	}
	return list, err
}

// 验证机器状态
func (s *MysqlService) ConnectMachine(id uint) error {
	var oldMachine models.SysMachine
	query := s.tx.Model(oldMachine).Where("id = ?", id).First(&oldMachine)
	if query.Error == gorm.ErrRecordNotFound {
		return errors.New("记录不存在")
	}

	// 初始化机器
	err := initRemoteMachine(&oldMachine)
	var newMachine models.SysMachine
	unConnectedStatus := models.SysMachineStatusUnConnected
	normalStatus := models.SysMachineStatusNormal
	if err != nil {
		newMachine.Status = &unConnectedStatus
		query.Updates(newMachine)
		return err
	}
	newMachine.Status = &normalStatus
	newMachine.Version = oldMachine.Version
	newMachine.Name = oldMachine.Name
	newMachine.Arch = oldMachine.Arch
	newMachine.Cpu = oldMachine.Cpu
	newMachine.Memory = oldMachine.Memory
	newMachine.Disk = oldMachine.Disk
	query.Updates(newMachine)
	return nil
}

// 初始化机器信息
func initRemoteMachine(machine *models.SysMachine) error {
	config := machine.GetSshConfig(2)
	cmds := []string{
		// 查看系统版本
		"lsb_release -d | cut -f 2 -d : | awk '$1=$1'",
		// 查看系统架构
		"arch",
		// 查看机器名字
		"uname -n",
		// 查看cpu型号
		"cat /proc/cpuinfo | grep name | cut -f 2 -d : | uniq | awk '$1=$1'",
		// 查看cpu核数
		"cat /proc/cpuinfo| grep 'cpu cores' | uniq | awk '{print $4}'",
		// 查看cpu线程数
		"cat /proc/cpuinfo | grep 'processor' | wc -l",
		// 查看内存大小（G）
		"cat /proc/meminfo | grep MemTotal | awk '{printf (\"%.2fG\\n\", $2 / 1024 / 1024)}'",
		// 查看硬盘容量（G）
		"df -h / | head -n 2 | tail -n 1 | awk '{print $2}'",
	}
	res := utils.ExecRemoteShell(config, cmds)
	if res.Err != nil {
		return res.Err
	}

	info := strings.Split(strings.TrimSuffix(res.Result, "\n"), "\n")
	if len(info) != len(cmds) {
		return fmt.Errorf("读取机器信息有误")
	}

	normalStatus := models.SysMachineStatusNormal

	machine.Status = &normalStatus
	machine.Version = info[0]
	machine.Arch = info[1]
	machine.Name = info[2]
	machine.Cpu = fmt.Sprintf("%s核%s线程 | %s", info[4], info[5], info[3])
	machine.Memory = info[6]
	machine.Disk = info[7]

	return nil
}
