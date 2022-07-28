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

func GetFsmLogDetail(c *gin.Context, detail req.FsmLogSubmitterDetail) []resp.FsmLogSubmitterDetail {
	my := service.New(c)
	return my.GetFsmLogDetail(detail)
}

func UpdateFsmLogDetail(c *gin.Context, detail req.UpdateFsmLogSubmitterDetail) error {
	my := service.New(c)
	return my.UpdateFsmLogDetail(detail)
}
