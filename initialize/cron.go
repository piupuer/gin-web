package initialize

import (
	"fmt"
	"gin-web/models"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"gin-web/pkg/wechat"
	"github.com/robfig/cron/v3"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 初始化定时任务
func InitCron() {
	go func() {
		// 新建cron实例
		c := cron.New()
		// 自定义日志
		log := new(CronCustomLogger)
		// 图片压缩任务
		if global.Conf.Upload.CompressImageCronTask != "" {
			// SkipIfStillRunning作用是如果前一个任务未执行完成将跳过新任务
			c.AddJob(global.Conf.Upload.CompressImageCronTask, cron.NewChain(cron.SkipIfStillRunning(log)).Then(&CompressImageJob{}))
		}
		// 微信消息任务
		if global.Conf.WeChat.Official.TplMessageCronTask.Expr != "" {
			job := new(WeChatTplMessageJob)
			job.Users = strings.Split(global.Conf.WeChat.Official.TplMessageCronTask.Users, ",")
			// SkipIfStillRunning作用是如果前一个任务未执行完成将跳过新任务
			c.AddJob(global.Conf.WeChat.Official.TplMessageCronTask.Expr, cron.NewChain(cron.SkipIfStillRunning(log)).Then(job))
		}
		// 启动调度
		c.Start()
	}()
}

// 自定义日志输出
type CronCustomLogger struct {
	cron.Logger
}

func (s CronCustomLogger) Info(msg string, keysAndValues ...interface{}) {
	global.Log.Info(fmt.Sprintf("[定时任务]%s", msg))
}

func (s CronCustomLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	global.Log.Error(fmt.Sprintf("[定时任务]%s, err: %v", msg, err))
}

// 图片压缩定时job
type CompressImageJob struct {
	// 已经压缩过的目录
	Dirs []string
}

func (s *CompressImageJob) Run() {
	global.Log.Info("[定时任务][图片压缩]准备开始...")
	// 获取工作目录
	pwd := utils.GetWorkDir()
	// 默认目录为文件上传目录
	compressDir := pwd + "/" + global.Conf.Upload.SaveDir
	// 配置了自定义压缩根目录
	if global.Conf.Upload.CompressImageRootDir != "" {
		compressDir = global.Conf.Upload.CompressImageRootDir
	}
	// 获取全部子目录
	childDirList, _ := ioutil.ReadDir(compressDir)

	for _, info := range childDirList {
		if info.IsDir() {
			currentDir := compressDir + "/" + info.Name()
			if utils.Contains(s.Dirs, currentDir) {
				global.Log.Debug(fmt.Sprintf("[定时任务][图片压缩]目录%s已扫描, 跳过", currentDir))
				continue
			}
			filepath.Walk(currentDir, func(path string, fi os.FileInfo, errBack error) error {
				if errBack != nil {
					return errBack
				}
				var err error
				// 压缩图片
				if global.Conf.Upload.CompressImageOriginalSaveDir != "" {
					// 保存源文件
					err = utils.CompressImageSaveOriginal(path, global.Conf.Upload.CompressImageOriginalSaveDir)
				} else {
					// 不保存源文件
					err = utils.CompressImage(path)
				}
				if err != nil {
					global.Log.Error(fmt.Sprintf("[定时任务][图片压缩]%v", err))
				}
				return nil
			})
			if !utils.Contains(s.Dirs, currentDir) {
				s.Dirs = append(s.Dirs, currentDir)
			}
		}
	}
}

const CurrentIndexKey = "we_chat_tpl_message_job_current_index"

// 微信模板消息通知定时job
type WeChatTplMessageJob struct {
	// 当前index
	Current int
	// 用户微信号列表
	Users []string
}

func (s *WeChatTplMessageJob) Run() {
	global.Log.Info("[定时任务][微信模板消息]准备开始...")
	l := len(s.Users)
	if l == 0 {
		global.Log.Warn("[定时任务][微信模板消息]用户列表未配置")
		return
	}
	// 从redis中读取index
	if global.Conf.System.UseRedis {
		current, _ := global.Redis.Get(CurrentIndexKey).Int64()
		s.Current = int(current)
	}
	// 不得超过最大长度
	if l <= s.Current {
		s.Current = 0
	}
	currentUser := s.Users[s.Current]
	msg := message.TemplateMessage{
		ToUser:     currentUser,
		TemplateID: global.Conf.WeChat.Official.TplMessageCronTask.TemplateId,
		Data: map[string]*message.TemplateDataItem{
			"first": {
				Value: "日常事项定时提醒",
			},
			"keyword1": {
				Value: "每日购买",
			},
			"keyword2": {
				Value: "请到商城下单支付一单(杨博士店有一分钱的单)",
			},
			"keyword3": {
				Value: models.LocalTime{
					Time: time.Now(),
				}.String(),
			},
			"remark": {
				Value: "下单完成记得将截图发到群里哦~",
			},
		},
	}
	msg.MiniProgram.AppID = global.Conf.WeChat.Official.TplMessageCronTask.MiniProgramAppId
	msg.MiniProgram.PagePath = global.Conf.WeChat.Official.TplMessageCronTask.MiniProgramPagePath
	err := wechat.SendTplMessage(&msg)
	if err == nil {
		s.Current++
		// 保存到redis
		if global.Conf.System.UseRedis {
			global.Redis.Set(CurrentIndexKey, s.Current, 0)
		}
	}
}
