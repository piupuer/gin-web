package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/utils"
	"gorm.io/gorm"
	"strings"
)

func (my MysqlService) FindMachine(req *request.MachineReq) ([]models.SysMachine, error) {
	var err error
	list := make([]models.SysMachine, 0)
	query := my.Q.Tx.
		Model(&models.SysMachine{}).
		Order("created_at DESC")
	host := strings.TrimSpace(req.Host)
	if host != "" {
		query = query.Where("host LIKE ?", fmt.Sprintf("%%%s%%", host))
	}
	loginName := strings.TrimSpace(req.LoginName)
	if loginName != "" {
		query = query.Where("login_name LIKE ?", fmt.Sprintf("%%%s%%", loginName))
	}
	if req.Status != nil {
		if *req.Status > 0 {
			query = query.Where("status = ?", 1)
		} else {
			query = query.Where("status = ?", 0)
		}
	}
	err = my.Q.FindWithPage(query, &req.Page, &list)
	return list, err
}

// connect machine
func (my MysqlService) ConnectMachine(id uint) error {
	var oldMachine models.SysMachine
	query := my.Q.Tx.Model(&oldMachine).Where("id = ?", id).First(&oldMachine)
	if query.Error == gorm.ErrRecordNotFound {
		return gorm.ErrRecordNotFound
	}

	err := initRemoteMachine(&oldMachine)
	var newMachine models.SysMachine
	unConnectedStatus := models.SysMachineStatusUnhealthy
	normalStatus := models.SysMachineStatusHealthy
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

// init machine
func initRemoteMachine(machine *models.SysMachine) error {
	config := utils.SshConfig{
		LoginName: machine.LoginName,
		LoginPwd:  machine.LoginPwd,
		Port:      machine.SshPort,
		Host:      machine.Host,
		Timeout:   2,
	}
	cmds := []string{
		// system version
		"lsb_release -d | cut -f 2 -d : | awk '$1=$1'",
		// system arch
		"arch",
		// system username
		"uname -n",
		// cpu info
		"cat /proc/cpuinfo | grep name | cut -f 2 -d : | uniq | awk '$1=$1'",
		// cpu cores
		"cat /proc/cpuinfo| grep 'cpu cores' | uniq | awk '{print $4}'",
		// cpu processor
		"cat /proc/cpuinfo | grep 'processor' | wc -l",
		// memory(GB)
		"cat /proc/meminfo | grep MemTotal | awk '{printf (\"%.2fG\\n\", $2 / 1024 / 1024)}'",
		// disk(GB)
		"df -h / | head -n 2 | tail -n 1 | awk '{print $2}'",
	}
	res := utils.ExecRemoteShell(config, cmds)
	if res.Err != nil {
		return res.Err
	}

	info := strings.Split(strings.TrimSuffix(res.Result, "\n"), "\n")
	if len(info) != len(cmds) {
		return fmt.Errorf("read machine info failed")
	}

	normalStatus := models.SysMachineStatusHealthy

	machine.Status = &normalStatus
	machine.Version = info[0]
	machine.Arch = info[1]
	machine.Name = info[2]
	machine.Cpu = fmt.Sprintf("%s cores %s processor | %s", info[4], info[5], info[3])
	machine.Memory = info[6]
	machine.Disk = info[7]

	return nil
}
