package strategy

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"github.com/piupuer/go-helper/pkg/fsm"
	"gorm.io/gorm"
)

// 工作流状态流转完成后需要更新目标表, 因此定义下列策略
type AfterTransitionStrategy interface {
	UpdateTarget() error
}

// 请假条审批
type LeaveApproval struct {
	tx       *gorm.DB
	targetId uint
	lastLog  fsm.Log
}

func (s LeaveApproval) UpdateTarget() error {
	var err error
	if *s.lastLog.Status > models.SysWorkflowLogStateSubmit {
		var leave models.SysLeave
		leave.Status = s.lastLog.Status
		// 更新请假审批状态
		err = s.tx.Model(&leave).Where("id = ?", s.targetId).Updates(&leave).Error
	}
	return err
}

// 策略类
type AfterTransitionContext struct {
	Strategy AfterTransitionStrategy
}

// 策略类构造函数
func NewAfterTransitionContext(tx *gorm.DB, targetCategory uint, targetId uint, lastLog fsm.Log) (*AfterTransitionContext, error) {
	ctx := new(AfterTransitionContext)
	switch targetCategory {
	case global.FsmCategoryLeave:
		ctx.Strategy = &LeaveApproval{
			tx:       tx,
			targetId: targetId,
			lastLog:  lastLog,
		}
		break
	default:
		return nil, fmt.Errorf("[NewAfterTransitionContext]策略获取失败, 请检查参数targetCategory: %d", targetCategory)
	}
	return ctx, nil
}
