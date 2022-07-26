package v1

import (
	"context"
	"gin-web/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/piupuer/go-helper/pkg/req"
	"github.com/piupuer/go-helper/pkg/resp"
)

func FsmTransition(ctx context.Context, logs ...resp.FsmApprovalLog) error {
	my := service.New(ctx)
	return my.FsmTransition(logs...)
}

func GetFsmDetail(c *gin.Context, detail req.FsmSubmitterDetail) []resp.FsmSubmitterDetail {
	my := service.New(c)
	return my.GetFsmDetail(detail)
}

func UpdateFsmDetail(c *gin.Context, detail req.UpdateFsmSubmitterDetail) error {
	my := service.New(c)
	return my.UpdateFsmDetail(detail)
}
