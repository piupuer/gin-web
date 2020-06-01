package service

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/request"
	"gin-web/pkg/response"
	"gin-web/tests"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"testing"
	"time"
)

var (
	approval = models.SysWorkflowLogStateApproval
	deny     = models.SysWorkflowLogStateDeny
	cancel   = models.SysWorkflowLogStateCancel
	restart  = models.SysWorkflowLogStateRestart
	end      = models.SysWorkflowLogStateEnd
)

func TestMysqlService_CreateWorkflow(t *testing.T) {
	tests.InitTestEnv()

	s := New(nil)
	var req request.CreateWorkflowRequestStruct
	req.Name = "请假流程"
	req.Creator = "系统"
	req.Desc = "用于员工请假审批"
	s.CreateWorkflow(&req)
}

func TestMysqlService_CreateWorkflowNode(t *testing.T) {
	tests.InitTestEnv()

	s := New(nil)
	var req1 request.CreateWorkflowNodeRequestStruct
	req1.FlowId = 1
	req1.RoleId = 1
	req1.Name = "主管审批"
	req1.Creator = "系统"
	req1.Desc = "主管[李四]审批请假单"
	s.CreateWorkflowNode(&req1)
	var req2 request.CreateWorkflowNodeRequestStruct
	req2.FlowId = 1
	req2.RoleId = 2
	req2.Name = "总经理审批"
	req2.Creator = "系统"
	req2.Desc = "总经理[张三]审批请假单"
	s.CreateWorkflowNode(&req2)
	var req3 request.CreateWorkflowNodeRequestStruct
	req3.FlowId = 1
	req3.RoleId = 3
	req3.Name = "董事长审批"
	req3.Creator = "系统"
	req3.Desc = "董事长[管理员]审批请假单"
	s.CreateWorkflowNode(&req3)
}

func TestMysqlService_CreateWorkflowLine(t *testing.T) {
	tests.InitTestEnv()

	s := New(nil)
	s.CreateWorkflowLine([][]uint{
		{1, 3},
		{2, 3},
		{3},
	})
}

func TestMysqlService_WorkflowTransition(t *testing.T) {
	tests.InitTestEnv()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := New(nil)
	// 测试1:
	// 测试目标: 配置角色, 每一层级单人审核
	testName := "测试1"
	// 1.创建审核角色1, 角色2, 角色3, 普通角色4
	var role1 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role1)
	var role2 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role2)
	var role3 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role3)
	var role4 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role4)
	// 2.创建用户1,2,3,4属于角色1
	var user1 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user1)
	var user2 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user2)
	var user3 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user3)
	var user4 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user4)
	// 3.创建用户5,6,7属于角色2
	var user5 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user5)
	var user6 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user6)
	var user7 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user7)
	// 4.创建用户8,9属于角色3
	var user8 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user8)
	var user9 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user9)
	// 5.创建用户10用于提交审批单
	var user10 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role4.Id,
	}
	s.tx.Create(&user10)

	falsePtr := false
	// 6.构建请假流程流水线 角色1审批=>角色2审批=>角色3审批
	var flow = models.SysWorkflow{
		Uuid:              uuid.NewV4().String(),
		Category:          models.SysWorkflowCategoryOnlyOneApproval, // 只需要有1个人通过
		SubmitUserConfirm: &falsePtr,                                 // 不需要提交人确认
		Self:              &falsePtr,                                 // 不能自我审批
		Name:              fmt.Sprintf("%s工作流程%d", testName, r.Intn(1000000)),
		Creator:           "系统",
		Desc:              "用于员工请假审批",
	}
	s.tx.Create(&flow)
	var node1 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管1审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node1)
	var node2 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管2审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node2)
	var node3 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理1审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node3)
	var node4 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理2审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node4)
	var node5 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理3审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node5)
	var node6 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role3.Id,
		Name:    testName + "董事长审批",
		Creator: "系统",
		Desc:    "3级审批",
	}
	s.tx.Create(&node6)
	// 2个主管/3个总经理/1个董事长
	s.CreateWorkflowLine([][]uint{
		{
			node1.Id,
			node2.Id,
		},
		{
			node3.Id,
			node4.Id,
			node5.Id,
		},
		{
			node6.Id,
		},
	})

	// 7.构建审批单
	// 提交一次
	err := s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		SubmitUserId:   user10.Id,                             // 用户10提交
	})

	fmt.Println("提交人初次提交", err)

	// 再提交一次
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		SubmitUserId:   user10.Id,                             // 用户10提交
	})

	fmt.Println("提交人重复提交", err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user2.Id,                              // 用户2审批
		ApprovalOpinion: "1级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user2.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user7.Id,                              // 用户7审批
		ApprovalOpinion: "2级拒绝",
		ApprovalStatus:  &deny,
	})

	fmt.Println("2级审批, 用户id", user7.Id, err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user3.Id,                              // 用户3审批
		ApprovalOpinion: "1级再次通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级再次审批, 用户id", user3.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user6.Id,                              // 用户6审批
		ApprovalOpinion: "2级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("2级审批, 用户id", user6.Id, err)

	// user8,9其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user9.Id,                              // 用户9审批
		ApprovalOpinion: "3级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("3级审批, 用户id", user9.Id, err)
}

func TestMysqlService_WorkflowTransition2(t *testing.T) {
	tests.InitTestEnv()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := New(nil)
	// 测试1
	// 测试目标: 配置具体审批人, 每一层级全部人审核
	testName := "测试2"
	// 1.创建审核角色1, 角色2, 角色3, 普通角色4
	var role1 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role1)
	var role2 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role2)
	var role3 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role3)
	var role4 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role4)
	// 2.创建用户1,2,3,4属于角色1
	var user1 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user1)
	var user2 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user2)
	var user3 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user3)
	var user4 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user4)
	// 3.创建用户5,6,7属于角色2
	var user5 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user5)
	var user6 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user6)
	var user7 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user7)
	// 4.创建用户8,9属于角色3
	var user8 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user8)
	var user9 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user9)
	// 5.创建用户10用于提交审批单
	var user10 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role4.Id,
	}
	s.tx.Create(&user10)

	falsePtr := false
	// 6.构建请假流程流水线 角色1审批=>角色2审批=>角色3审批
	var flow = models.SysWorkflow{
		Uuid:              uuid.NewV4().String(),
		Category:          models.SysWorkflowCategoryAllApproval, // 全部人通过
		SubmitUserConfirm: &falsePtr,                             // 不需要提交人确认
		Self:              &falsePtr,                             // 不能自我审批
		Name:              fmt.Sprintf("%s工作流程%d", testName, r.Intn(1000000)),
		Creator:           "系统",
		Desc:              "用于员工请假审批",
	}
	s.tx.Create(&flow)
	var node1 = models.SysWorkflowNode{
		FlowId: flow.Id,
		Users: []models.SysUser{
			user1,
			user2,
		},
		Name:    testName + "主管1审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node1)
	var node2 = models.SysWorkflowNode{
		FlowId: flow.Id,
		Users: []models.SysUser{
			user3,
			user4,
		},
		Name:    testName + "主管2审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node2)
	var node3 = models.SysWorkflowNode{
		FlowId: flow.Id,
		Users: []models.SysUser{
			user6,
		},
		Name:    testName + "总经理1审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node3)
	var node4 = models.SysWorkflowNode{
		FlowId: flow.Id,
		Users: []models.SysUser{
			user5,
		},
		Name:    testName + "总经理2审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node4)
	var node5 = models.SysWorkflowNode{
		FlowId: flow.Id,
		Users: []models.SysUser{
			user7,
		},
		Name:    testName + "总经理3审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node5)
	var node6 = models.SysWorkflowNode{
		FlowId: flow.Id,
		Users: []models.SysUser{
			user8,
			user9,
		},
		Name:    testName + "董事长审批",
		Creator: "系统",
		Desc:    "3级审批",
	}
	s.tx.Create(&node6)
	// 2个主管/3个总经理/1个董事长
	s.CreateWorkflowLine([][]uint{
		{
			node1.Id,
			node2.Id,
		},
		{
			node3.Id,
			node4.Id,
			node5.Id,
		},
		{
			node6.Id,
		},
	})

	// 7.构建审批单
	// 提交一次
	err := s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		SubmitUserId:   user10.Id,                             // 用户10提交
	})

	fmt.Println("提交人初次提交", err)

	// 再提交一次
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		SubmitUserId:   user10.Id,                             // 用户10提交
	})

	fmt.Println("提交人重复提交", err)

	fmt.Println(fmt.Sprintf("用户%d待审批列表", user1.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user1.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user2.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user2.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user3.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user3.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user4.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user4.Id))

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	// user1,2,3,4挨个审批(不分顺序)
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user2.Id,                              // 用户2审批
		ApprovalOpinion: "1级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user2.Id, err)

	fmt.Println(fmt.Sprintf("用户%d待审批列表", user1.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user1.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user2.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user2.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user3.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user3.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user4.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user4.Id))

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user4.Id,                              // 用户4审批
		ApprovalOpinion: "1级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user4.Id, err)

	fmt.Println(fmt.Sprintf("用户%d待审批列表", user1.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user1.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user2.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user2.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user3.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user3.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user4.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user4.Id))

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user1.Id,                              // 用户1审批
		ApprovalOpinion: "1级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user1.Id, err)

	fmt.Println(fmt.Sprintf("用户%d待审批列表", user1.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user1.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user2.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user2.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user3.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user3.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user4.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user4.Id))

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user3.Id,                              // 用户3审批
		ApprovalOpinion: "1级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user3.Id, err)

	fmt.Println(fmt.Sprintf("用户%d待审批列表", user1.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user1.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user2.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user2.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user3.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user3.Id))
	fmt.Println(fmt.Sprintf("用户%d待审批列表", user4.Id))
	fmt.Println(s.GetWorkflowApprovingList(flow.Id, user10.Id, user4.Id))

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user7.Id,                              // 用户7审批
		ApprovalOpinion: "2级拒绝",
		ApprovalStatus:  &deny,
	})

	fmt.Println("2级审批, 用户id", user7.Id, err)

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	// user1,2,3,4再次全部审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user3.Id,                              // 用户3审批
		ApprovalOpinion: "1级再次通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级再次审批, 用户id", user3.Id, err)

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user2.Id,                              // 用户4审批
		ApprovalOpinion: "1级再次通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级再次审批, 用户id", user2.Id, err)

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user1.Id,                              // 用户4审批
		ApprovalOpinion: "1级再次通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级再次审批, 用户id", user1.Id, err)

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user4.Id,                              // 用户4审批
		ApprovalOpinion: "1级再次通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级再次审批, 用户id", user4.Id, err)

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	// user5,6,7全部审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user5.Id,                              // 用户5审批
		ApprovalOpinion: "2级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("2级审批, 用户id", user5.Id, err)

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user6.Id,                              // 用户6审批
		ApprovalOpinion: "2级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("2级审批, 用户id", user6.Id, err)

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user7.Id,                              // 用户7审批
		ApprovalOpinion: "2级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("2级审批, 用户id", user7.Id, err)

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	// user8,9全部审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user9.Id,                              // 用户9审批
		ApprovalOpinion: "3级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("3级审批, 用户id", user9.Id, err)

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user8.Id,                              // 用户8审批
		ApprovalOpinion: "3级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("3级审批, 用户id", user8.Id, err)

	fmt.Println("下一审批人")
	fmt.Println(s.GetWorkflowNextApprovingUsers(flow.Id, user10.Id))

}

func TestMysqlService_WorkflowTransition3(t *testing.T) {
	tests.InitTestEnv()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := New(nil)
	// 测试3:
	// 测试目标: 配置角色, 每一层级单人审核, 审核到最后一人然后反向全部拒绝, 是否能结束
	testName := "测试3"
	// 1.创建审核角色1, 角色2, 角色3, 普通角色4
	var role1 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role1)
	var role2 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role2)
	var role3 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role3)
	var role4 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role4)
	// 2.创建用户1,2,3,4属于角色1
	var user1 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user1)
	var user2 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user2)
	var user3 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user3)
	var user4 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user4)
	// 3.创建用户5,6,7属于角色2
	var user5 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user5)
	var user6 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user6)
	var user7 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user7)
	// 4.创建用户8,9属于角色3
	var user8 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user8)
	var user9 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user9)
	// 5.创建用户10用于提交审批单
	var user10 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role4.Id,
	}
	s.tx.Create(&user10)

	falsePtr := false
	// 6.构建请假流程流水线 角色1审批=>角色2审批=>角色3审批
	var flow = models.SysWorkflow{
		Uuid:              uuid.NewV4().String(),
		Category:          models.SysWorkflowCategoryOnlyOneApproval, // 只需要有1个人通过
		SubmitUserConfirm: &falsePtr,                                 // 不需要提交人确认
		Self:              &falsePtr,                                 // 不能自我审批
		Name:              fmt.Sprintf("%s工作流程%d", testName, r.Intn(1000000)),
		Creator:           "系统",
		Desc:              "用于员工请假审批",
	}
	s.tx.Create(&flow)
	var node1 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管1审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node1)
	var node2 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管2审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node2)
	var node3 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理1审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node3)
	var node4 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理2审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node4)
	var node5 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理3审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node5)
	var node6 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role3.Id,
		Name:    testName + "董事长审批",
		Creator: "系统",
		Desc:    "3级审批",
	}
	s.tx.Create(&node6)
	// 2个主管/3个总经理/1个董事长
	s.CreateWorkflowLine([][]uint{
		{
			node1.Id,
			node2.Id,
		},
		{
			node3.Id,
			node4.Id,
			node5.Id,
		},
		{
			node6.Id,
		},
	})

	// 7.构建审批单
	// 提交一次
	err := s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		SubmitUserId:   user10.Id,                             // 用户10提交
	})

	fmt.Println("提交人初次提交", err)

	// 再提交一次
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		SubmitUserId:   user10.Id,                             // 用户10提交
	})

	fmt.Println("提交人重复提交", err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user2.Id,                              // 用户2审批
		ApprovalOpinion: "1级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user2.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user6.Id,                              // 用户6审批
		ApprovalOpinion: "2级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("2级审批, 用户id", user6.Id, err)

	// user8,9其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user9.Id,                              // 用户9审批
		ApprovalOpinion: "3级拒绝",
		ApprovalStatus:  &deny,
	})

	fmt.Println("3级审批, 用户id", user9.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user5.Id,                              // 用户5审批
		ApprovalOpinion: "2级拒绝",
		ApprovalStatus:  &deny,
	})

	fmt.Println("2级审批, 用户id", user5.Id, err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user1.Id,                              // 用户1审批
		ApprovalOpinion: "1级拒绝",
		ApprovalStatus:  &deny,
	})

	fmt.Println("1级审批, 用户id", user1.Id, err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user3.Id,                              // 用户3审批
		ApprovalOpinion: "1级再次拒绝",
		ApprovalStatus:  &deny,
	})

	fmt.Println("1级审批, 用户id", user3.Id, err)

}

func TestMysqlService_WorkflowTransition4(t *testing.T) {
	tests.InitTestEnv()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := New(nil)
	// 测试4:
	// 测试目标: 配置角色, 每一层级单人审核, 开启提交人确认, 审批是否能结束
	testName := "测试4"
	// 1.创建审核角色1, 角色2, 角色3, 普通角色4
	var role1 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role1)
	var role2 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role2)
	var role3 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role3)
	var role4 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role4)
	// 2.创建用户1,2,3,4属于角色1
	var user1 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user1)
	var user2 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user2)
	var user3 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user3)
	var user4 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user4)
	// 3.创建用户5,6,7属于角色2
	var user5 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user5)
	var user6 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user6)
	var user7 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user7)
	// 4.创建用户8,9属于角色3
	var user8 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user8)
	var user9 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user9)
	// 5.创建用户10用于提交审批单
	var user10 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role4.Id,
	}
	s.tx.Create(&user10)

	truePtr := true
	falsePtr := false
	// 6.构建请假流程流水线 角色1审批=>角色2审批=>角色3审批
	var flow = models.SysWorkflow{
		Uuid:              uuid.NewV4().String(),
		Category:          models.SysWorkflowCategoryOnlyOneApproval, // 只需要有1个人通过
		SubmitUserConfirm: &truePtr,                                  // 需要提交人确认
		Self:              &falsePtr,                                 // 不能自我审批
		Name:              fmt.Sprintf("%s工作流程%d", testName, r.Intn(1000000)),
		Creator:           "系统",
		Desc:              "用于员工请假审批",
	}
	s.tx.Create(&flow)
	var node1 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管1审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node1)
	var node2 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管2审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node2)
	var node3 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理1审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node3)
	var node4 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理2审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node4)
	var node5 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理3审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node5)
	var node6 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role3.Id,
		Name:    testName + "董事长审批",
		Creator: "系统",
		Desc:    "3级审批",
	}
	s.tx.Create(&node6)
	// 2个主管/3个总经理/1个董事长
	s.CreateWorkflowLine([][]uint{
		{
			node1.Id,
			node2.Id,
		},
		{
			node3.Id,
			node4.Id,
			node5.Id,
		},
		{
			node6.Id,
		},
	})

	// 7.构建审批单
	// 提交一次
	err := s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		SubmitUserId:   user10.Id,                             // 用户10提交
	})

	fmt.Println("提交人初次提交", err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user2.Id,                              // 用户2审批
		ApprovalOpinion: "1级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user2.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user7.Id,                              // 用户7审批
		ApprovalOpinion: "2级拒绝",
		ApprovalStatus:  &deny,
	})

	fmt.Println("2级审批, 用户id", user7.Id, err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user3.Id,                              // 用户3审批
		ApprovalOpinion: "1级再次通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级再次审批, 用户id", user3.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user6.Id,                              // 用户6审批
		ApprovalOpinion: "2级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("2级审批, 用户id", user6.Id, err)

	// user8,9其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user9.Id,                              // 用户9审批
		ApprovalOpinion: "3级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("3级审批, 用户id", user9.Id, err)

	// 提交人确认
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		ApprovalId:     user10.Id,                             // 用户10确认
		ApprovalStatus: &deny,                                 // 拒绝: 将会失败
	})

	fmt.Println("提交人确认", err)

	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		ApprovalId:     user10.Id,                             // 用户10确认
		ApprovalStatus: &approval,                             // 同意
	})

	fmt.Println("提交人确认", err)
}

func TestMysqlService_WorkflowTransition5(t *testing.T) {
	tests.InitTestEnv()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := New(nil)
	// 测试5:
	// 测试目标: 配置角色, 每一层级单人审核, 开启提交人确认, 审批是否能回退
	testName := "测试5"
	// 1.创建审核角色1, 角色2, 角色3, 普通角色4
	var role1 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role1)
	var role2 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role2)
	var role3 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role3)
	var role4 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role4)
	// 2.创建用户1,2,3,4属于角色1
	var user1 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user1)
	var user2 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user2)
	var user3 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user3)
	var user4 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user4)
	// 3.创建用户5,6,7属于角色2
	var user5 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user5)
	var user6 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user6)
	var user7 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user7)
	// 4.创建用户8,9属于角色3
	var user8 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user8)
	var user9 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user9)
	// 5.创建用户10用于提交审批单
	var user10 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role4.Id,
	}
	s.tx.Create(&user10)

	truePtr := true
	falsePtr := false
	// 6.构建请假流程流水线 角色1审批=>角色2审批=>角色3审批
	var flow = models.SysWorkflow{
		Uuid:              uuid.NewV4().String(),
		Category:          models.SysWorkflowCategoryOnlyOneApproval, // 只需要有1个人通过
		SubmitUserConfirm: &truePtr,                                  // 需要提交人确认
		Self:              &falsePtr,                                 // 不能自我审批
		Name:              fmt.Sprintf("%s工作流程%d", testName, r.Intn(1000000)),
		Creator:           "系统",
		Desc:              "用于员工请假审批",
	}
	s.tx.Create(&flow)
	var node1 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管1审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node1)
	var node2 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管2审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node2)
	var node3 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理1审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node3)
	var node4 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理2审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node4)
	var node5 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理3审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node5)
	var node6 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role3.Id,
		Name:    testName + "董事长审批",
		Creator: "系统",
		Desc:    "3级审批",
	}
	s.tx.Create(&node6)
	// 2个主管/3个总经理/1个董事长
	s.CreateWorkflowLine([][]uint{
		{
			node1.Id,
			node2.Id,
		},
		{
			node3.Id,
			node4.Id,
			node5.Id,
		},
		{
			node6.Id,
		},
	})

	// 7.构建审批单
	// 提交一次
	err := s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		SubmitUserId:   user10.Id,                             // 用户10提交
	})

	fmt.Println("提交人初次提交", err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user2.Id,                              // 用户2审批
		ApprovalOpinion: "1级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user2.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user7.Id,                              // 用户7审批
		ApprovalOpinion: "2级拒绝",
		ApprovalStatus:  &deny,
	})

	fmt.Println("2级审批, 用户id", user7.Id, err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user3.Id,                              // 用户3审批
		ApprovalOpinion: "1级再次通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级再次审批, 用户id", user3.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user6.Id,                              // 用户6审批
		ApprovalOpinion: "2级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("2级审批, 用户id", user6.Id, err)

	// user8,9其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user9.Id,                              // 用户9审批
		ApprovalOpinion: "3级拒绝",
		ApprovalStatus:  &deny,
	})

	fmt.Println("3级审批, 用户id", user9.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user6.Id,                              // 用户6审批
		ApprovalOpinion: "2级拒绝",
		ApprovalStatus:  &deny,
	})

	fmt.Println("2级审批, 用户id", user6.Id, err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user3.Id,                              // 用户3审批
		ApprovalOpinion: "1级拒绝",
		ApprovalStatus:  &deny,
	})

	fmt.Println("1级审批, 用户id", user3.Id, err)

	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		ApprovalId:     user10.Id,                             // 用户10确认
		ApprovalStatus: &approval,                             // 同意
	})

	fmt.Println("提交人确认", err)
}

func TestMysqlService_WorkflowTransition6(t *testing.T) {
	tests.InitTestEnv()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := New(nil)
	// 测试6:
	// 测试目标: 配置角色, 每一层级单人审核, 开启自我审批, 看如果提交人是角色中的一个是否能审批
	testName := "测试6"
	// 1.创建审核角色1, 角色2, 角色3
	var role1 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role1)
	var role2 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role2)
	var role3 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role3)
	// 2.创建用户1,2,3,4属于角色1
	var user1 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user1)
	var user2 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user2)
	var user3 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user3)
	var user4 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user4)
	// 3.创建用户5,6,7属于角色2
	var user5 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user5)
	var user6 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user6)
	var user7 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user7)
	// 4.创建用户8,9属于角色3
	var user8 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user8)
	var user9 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user9)
	// 5.创建用户10用于提交审批单=>属于角色1
	var user10 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user10)

	truePtr := true
	falsePtr := false
	// 6.构建请假流程流水线 角色1审批=>角色2审批=>角色3审批
	var flow = models.SysWorkflow{
		Uuid:              uuid.NewV4().String(),
		Category:          models.SysWorkflowCategoryOnlyOneApproval, // 只需要有1个人通过
		SubmitUserConfirm: &falsePtr,                                 // 需要提交人确认
		Self:              &truePtr,                                  // 不能自我审批
		Name:              fmt.Sprintf("%s工作流程%d", testName, r.Intn(1000000)),
		Creator:           "系统",
		Desc:              "用于员工请假审批",
	}
	s.tx.Create(&flow)
	var node1 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管1审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node1)
	var node2 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管2审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node2)
	var node3 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理1审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node3)
	var node4 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理2审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node4)
	var node5 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理3审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node5)
	var node6 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role3.Id,
		Name:    testName + "董事长审批",
		Creator: "系统",
		Desc:    "3级审批",
	}
	s.tx.Create(&node6)
	// 2个主管/3个总经理/1个董事长
	s.CreateWorkflowLine([][]uint{
		{
			node1.Id,
			node2.Id,
		},
		{
			node3.Id,
			node4.Id,
			node5.Id,
		},
		{
			node6.Id,
		},
	})

	// 7.构建审批单
	// 提交一次
	err := s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		SubmitUserId:   user10.Id,                             // 用户10提交
	})

	fmt.Println("提交人初次提交", err)

	// 提交人自我审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user10.Id,                             // 用户10自我审批
		ApprovalOpinion: "1级自我审批通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println(fmt.Sprintf("1级审批, 提交人自我审批, 是否开启自我审批: %v 用户id %d %v", truePtr, user10.Id, err))

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user6.Id,                              // 用户6审批
		ApprovalOpinion: "2级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("2级审批, 用户id", user6.Id, err)

	// user8,9其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user9.Id,                              // 用户9审批
		ApprovalOpinion: "3级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("3级审批, 用户id", user9.Id, err)
}

func TestMysqlService_WorkflowTransition7(t *testing.T) {
	tests.InitTestEnv()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	s := New(nil)
	// 测试7:
	// 测试目标: 测试取消与重启
	testName := "测试7"
	// 1.创建审核角色1, 角色2, 角色3, 普通角色4
	var role1 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role1)
	var role2 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role2)
	var role3 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role3)
	var role4 = models.SysRole{
		Name:    fmt.Sprintf("%s角色名称%d", testName, r.Intn(1000000)),
		Keyword: fmt.Sprintf("%s角色关键字%d", testName, r.Intn(1000000)),
		Creator: "系统",
	}
	s.tx.Create(&role4)
	// 2.创建用户1,2,3,4属于角色1
	var user1 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user1)
	var user2 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user2)
	var user3 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user3)
	var user4 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role1.Id,
	}
	s.tx.Create(&user4)
	// 3.创建用户5,6,7属于角色2
	var user5 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user5)
	var user6 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user6)
	var user7 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role2.Id,
	}
	s.tx.Create(&user7)
	// 4.创建用户8,9属于角色3
	var user8 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user8)
	var user9 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role3.Id,
	}
	s.tx.Create(&user9)
	// 5.创建用户10用于提交审批单
	var user10 = models.SysUser{
		Username: fmt.Sprintf("%s用户名%d", testName, r.Intn(1000000)),
		Nickname: fmt.Sprintf("%s用户昵称%d", testName, r.Intn(1000000)),
		Creator:  "系统",
		RoleId:   role4.Id,
	}
	s.tx.Create(&user10)

	falsePtr := false
	// 6.构建请假流程流水线 角色1审批=>角色2审批=>角色3审批
	var flow = models.SysWorkflow{
		Uuid:              uuid.NewV4().String(),
		Category:          models.SysWorkflowCategoryOnlyOneApproval, // 只需要有1个人通过
		SubmitUserConfirm: &falsePtr,                                 // 不需要提交人确认
		Self:              &falsePtr,                                 // 不能自我审批
		Name:              fmt.Sprintf("%s工作流程%d", testName, r.Intn(1000000)),
		Creator:           "系统",
		Desc:              "用于员工请假审批",
	}
	s.tx.Create(&flow)
	var node1 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管1审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node1)
	var node2 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role1.Id,
		Name:    testName + "主管2审批",
		Creator: "系统",
		Desc:    "1级审批",
	}
	s.tx.Create(&node2)
	var node3 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理1审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node3)
	var node4 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理2审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node4)
	var node5 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role2.Id,
		Name:    testName + "总经理3审批",
		Creator: "系统",
		Desc:    "2级审批",
	}
	s.tx.Create(&node5)
	var node6 = models.SysWorkflowNode{
		FlowId:  flow.Id,
		RoleId:  role3.Id,
		Name:    testName + "董事长审批",
		Creator: "系统",
		Desc:    "3级审批",
	}
	s.tx.Create(&node6)
	// 2个主管/3个总经理/1个董事长
	s.CreateWorkflowLine([][]uint{
		{
			node1.Id,
			node2.Id,
		},
		{
			node3.Id,
			node4.Id,
			node5.Id,
		},
		{
			node6.Id,
		},
	})

	// 7.构建审批单
	// 提交一次
	err := s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		SubmitUserId:   user10.Id,                             // 用户10提交
	})

	fmt.Println("提交人初次提交", err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user2.Id,                              // 用户2审批
		ApprovalOpinion: "1级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user2.Id, err)

	// 取消
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		ApprovalId:     user10.Id,                             // 用户10提交
		ApprovalStatus: &cancel,                               // 取消
	})

	fmt.Println("提交人1级审批通过后主动取消", err)

	// 重启
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		ApprovalId:     user10.Id,                             // 用户10提交
		ApprovalStatus: &restart,                              // 重启
	})

	fmt.Println("提交人1级审批取消后主动重启", err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user2.Id,                              // 用户2审批
		ApprovalOpinion: "1级通过(取消后从头开始)",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user2.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user6.Id,                              // 用户6审批
		ApprovalOpinion: "2级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("2级审批, 用户id", user6.Id, err)

	// 取消
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		ApprovalId:     user10.Id,                             // 用户10提交
		ApprovalStatus: &cancel,                               // 取消
	})

	fmt.Println("提交人2级审批通过后主动取消", err)

	// 再次取消
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		ApprovalId:     user10.Id,                             // 用户10提交
		ApprovalStatus: &cancel,                               // 取消
	})

	fmt.Println("提交人2级审批通过后再次主动取消", err)

	// 重启
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		ApprovalId:     user10.Id,                             // 用户10提交
		ApprovalStatus: &restart,                              // 重启
	})

	fmt.Println("提交人2级审批取消后主动重启", err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user2.Id,                              // 用户2审批
		ApprovalOpinion: "1级通过(取消后从头开始)",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user2.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user6.Id,                              // 用户6审批
		ApprovalOpinion: "2级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("2级审批, 用户id", user6.Id, err)

	// user8,9其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user9.Id,                              // 用户9审批
		ApprovalOpinion: "3级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("3级审批, 用户id", user9.Id, err)

	// 取消
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		ApprovalId:     user10.Id,                             // 用户10提交
		ApprovalStatus: &cancel,                               // 取消
	})

	fmt.Println("提交人3级审批通过后主动取消", err)

	// 重启
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:         flow.Id,
		TargetCategory: models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:       user10.Id,                             // 用户10请假
		ApprovalId:     user10.Id,                             // 用户10提交
		ApprovalStatus: &restart,                              // 重启
	})

	fmt.Println("提交人3级审批取消后主动重启", err)

	// user1,2,3,4其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户10请假
		ApprovalId:      user2.Id,                              // 用户2审批
		ApprovalOpinion: "1级通过(取消后从头开始)",
		ApprovalStatus:  &approval,
	})

	fmt.Println("1级审批, 用户id", user2.Id, err)

	// user5,6,7其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user6.Id,                              // 用户6审批
		ApprovalOpinion: "2级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("2级审批, 用户id", user6.Id, err)

	// user8,9其中一人审批
	err = s.WorkflowTransition(&request.WorkflowTransitionRequestStruct{
		FlowId:          flow.Id,
		TargetCategory:  models.SysWorkflowTargetCategoryLeave, // 请假
		TargetId:        user10.Id,                             // 用户2请假
		ApprovalId:      user9.Id,                              // 用户9审批
		ApprovalOpinion: "3级通过",
		ApprovalStatus:  &approval,
	})

	fmt.Println("3级审批, 用户id", user9.Id, err)
}

func TestMysqlService_GetWorkflowLineLogs(t *testing.T) {
	tests.InitTestEnv()

	s := New(nil)
	logs, err := s.GetWorkflowLogs(2, 13)
	res := make([]response.WorkflowLogsListResponseStruct, 0)
	for _, log := range logs {
		res = append(res, response.WorkflowLogsListResponseStruct{
			FlowName:              log.Flow.Name,
			FlowUuid:              log.Flow.Uuid,
			FlowCategoryStr:       models.SysWorkflowCategoryConst[log.Flow.Category],
			FlowTargetCategoryStr: models.SysWorkflowTargetCategoryConst[log.Flow.TargetCategory],
			Status:                log.Status,
			StatusStr:             models.SysWorkflowLogStateConst[*log.Status],
			SubmitUsername:        log.SubmitUser.Username,
			SubmitUserNickname:    log.SubmitUser.Nickname,
			ApprovalUsername:      log.ApprovalUser.Username,
			ApprovalUserNickname:  log.ApprovalUser.Nickname,
			ApprovalOpinion:       log.ApprovalOpinion,
			CreatedAt:             log.Model.CreatedAt,
			UpdatedAt:             log.Model.UpdatedAt,
		})
	}
	fmt.Println(logs, res, err)
}
