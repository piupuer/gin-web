package service

import (
	"github.com/piupuer/go-helper/pkg/fsm"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

// find finite state machine
func (my MysqlService) FindFsm(r req.FsmMachine) ([]resp.FsmMachine, error) {
	f := fsm.New(my.Q.Tx)
	return f.FindMachine(r)
}

// create finite state machine
func (my MysqlService) CreateFsm(r req.FsmCreateMachine) error {
	f := fsm.New(my.Q.Tx)
	_, err := f.CreateMachine(r)
	return err
}

// find waiting approve log
func (my MysqlService) FindFsmApprovingLog(r req.FsmPendingLog) ([]fsm.Log, error) {
	f := fsm.New(my.Q.Tx)
	return f.FindPendingLogByApprover(r)
}
