package initialize

import (
	"context"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/robfig/cron/v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 初始化定时任务
func Cron() {
	go func() {
		// 新建cron实例
		c := cron.New()
		// 自定义日志
		log := new(CronCustomLogger)
		log.ctx = global.RequestIdContext("")
		// 图片压缩任务
		if global.Conf.Upload.CompressImageCronTask != "" {
			// SkipIfStillRunning作用是如果前一个任务未执行完成将跳过新任务
			c.AddJob(
				global.Conf.Upload.CompressImageCronTask,
				cron.NewChain(cron.SkipIfStillRunning(log)).Then(
					&CompressImageJob{
						ctx: log.ctx,
					},
				),
			)
		}
		// 启动调度
		c.Start()
	}()
	global.Log.Debug(ctx, "初始化定时任务完成")
}

// 自定义日志输出
type CronCustomLogger struct {
	ctx context.Context
	cron.Logger
}

func (s CronCustomLogger) Info(msg string, keysAndValues ...interface{}) {
	global.Log.Info(s.ctx, "[定时任务]%s", msg)
}

func (s CronCustomLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	global.Log.Error(s.ctx, "[定时任务]%s, err: %v", msg, err)
}

// 图片压缩定时job
type CompressImageJob struct {
	ctx context.Context
	// 已经压缩过的目录
	Dirs []string
}

func (s *CompressImageJob) Run() {
	global.Log.Info(s.ctx, "[定时任务][图片压缩]准备开始...")
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
				global.Log.Debug(s.ctx, "[定时任务][图片压缩]目录%s已扫描, 跳过", currentDir)
				continue
			}
			filepath.Walk(currentDir, func(path string, fi os.FileInfo, errBack error) error {
				if errBack != nil {
					return errBack
				}
				var err error
				// 压缩图片
				if global.Conf.Upload.CompressImageOriginalSaveDir != "" {
					if strings.Contains(path, global.Conf.Upload.CompressImageOriginalSaveDir) {
						global.Log.Debug(s.ctx, "[定时任务][图片压缩]目录%s为源文件保存目录, 跳过", path)
						return nil
					}
					// 保存源文件
					err = utils.CompressImageSaveOriginal(path, global.Conf.Upload.CompressImageOriginalSaveDir)
				} else {
					// 不保存源文件
					err = utils.CompressImage(path)
				}
				if err != nil {
					global.Log.Error(s.ctx, "[定时任务][图片压缩]压缩失败, 当前文件%s, %v", path, err)
				} else {
					global.Log.Info(s.ctx, "[定时任务][图片压缩]压缩成功, 当前文件%s", path)
				}
				return nil
			})
			if !utils.Contains(s.Dirs, currentDir) {
				s.Dirs = append(s.Dirs, currentDir)
			}
		}
	}
	global.Log.Info(s.ctx, "[定时任务][图片压缩]任务结束")
}
