package initialize

import (
	"context"
	"gin-web/pkg/global"
	"gin-web/pkg/utils"
	"github.com/piupuer/go-helper/pkg/constant"
	"github.com/piupuer/go-helper/pkg/job"
	"github.com/piupuer/go-helper/pkg/query"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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
	global.Log.Debug(ctx, "initialize cron job success")
}

var dirs []string

func compress(c context.Context) error {
	ctx := query.NewRequestId(c, constant.MiddlewareRequestIdCtxKey)
	global.Log.Info(ctx, "[cron job][image compress]starting...")
	pwd := utils.GetWorkDir()
	// compress dir is default upload save dir
	compressDir := pwd + "/" + global.Conf.Upload.SaveDir
	if global.Conf.Upload.CompressImageRootDir != "" {
		compressDir = global.Conf.Upload.CompressImageRootDir
	}
	childDirList, _ := ioutil.ReadDir(compressDir)

	for _, info := range childDirList {
		if info.IsDir() {
			currentDir := compressDir + "/" + info.Name()
			if utils.Contains(dirs, currentDir) {
				global.Log.Debug(ctx, "[cron job][image compress]dir %s scanned, skip", currentDir)
				continue
			}
			filepath.Walk(currentDir, func(path string, fi os.FileInfo, errBack error) error {
				if errBack != nil {
					return errBack
				}
				var err error
				if global.Conf.Upload.CompressImageOriginalSaveDir != "" {
					if strings.Contains(path, global.Conf.Upload.CompressImageOriginalSaveDir) {
						global.Log.Debug(ctx, "[cron job][image compress]dir %s is original dir, skip", path)
						return nil
					}
					// save original file
					err = utils.CompressImageSaveOriginal(path, global.Conf.Upload.CompressImageOriginalSaveDir)
				} else {
					// direct compression
					err = utils.CompressImage(path)
				}
				if err != nil {
					global.Log.Error(ctx, "[cron job][image compress]compress filename %s failed: err", path, err)
				} else {
					global.Log.Info(ctx, "[cron job][image compress]compress filename %s success", path)
				}
				return nil
			})
			if !utils.Contains(dirs, currentDir) {
				dirs = append(dirs, currentDir)
			}
		}
	}
	global.Log.Info(ctx, "[cron job][image compress]ended")
	return nil
}
