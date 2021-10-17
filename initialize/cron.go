package initialize

import (
	"context"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/piupuer/go-helper/pkg/job"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// 初始化定时任务
func Cron() {
	j, err := job.New(
		job.Config{
			RedisClient: global.Redis,
		},
		job.WithLogger(global.Log),
		job.WithContext(ctx),
		job.WithAutoRequestId,
	)
	if err != nil {
		panic(err)
	}
	if global.Conf.Upload.CompressImageCronTask != "" {
		j.AddTask(job.GoodTask{
			Name: "compress",
			Expr: global.Conf.Upload.CompressImageCronTask,
			Func: compress,
		}).Start()
	}
	global.Log.Debug(ctx, "初始化定时任务完成")
}

var dirs []string

func compress(c context.Context) error {
	requestId, _ := c.Value(global.RequestIdContextKey).(string)
	ctx := global.RequestIdContext(requestId)
	global.Log.Info(ctx, "[定时任务][图片压缩]准备开始...")
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
			if utils.Contains(dirs, currentDir) {
				global.Log.Debug(ctx, "[定时任务][图片压缩]目录%s已扫描, 跳过", currentDir)
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
						global.Log.Debug(ctx, "[定时任务][图片压缩]目录%s为源文件保存目录, 跳过", path)
						return nil
					}
					// 保存源文件
					err = utils.CompressImageSaveOriginal(path, global.Conf.Upload.CompressImageOriginalSaveDir)
				} else {
					// 不保存源文件
					err = utils.CompressImage(path)
				}
				if err != nil {
					global.Log.Error(ctx, "[定时任务][图片压缩]压缩失败, 当前文件%s, %v", path, err)
				} else {
					global.Log.Info(ctx, "[定时任务][图片压缩]压缩成功, 当前文件%s", path)
				}
				return nil
			})
			if !utils.Contains(dirs, currentDir) {
				dirs = append(dirs, currentDir)
			}
		}
	}
	global.Log.Info(ctx, "[定时任务][图片压缩]任务结束")
	return nil
}
