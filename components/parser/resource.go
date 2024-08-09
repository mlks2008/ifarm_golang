package parser

import (
	"errors"
	"fmt"
	"github.com/radovskyb/watcher"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type (
	Resource struct {
		FileConfig
		watcher  *watcher.Watcher
		changeFn ConfigChangeFn
	}

	FileConfig struct {
		FilePath    string `json:"filePath"`   // 配置文件路径
		ExtName     string `json:"extName"`    // 文件扩展名
		ReloadTime  int64  `json:"reloadTime"` // 定时重载扫描(秒)
		MonitorPath string `json:"-"`          // 监控路径
	}
)

func NewResource() *Resource {
	return &Resource{
		FileConfig: FileConfig{
			FilePath:    "",
			ExtName:     "",
			ReloadTime:  0,
			MonitorPath: "",
		},
	}
}

func (r *Resource) BuildFilePath(filePath string) *Resource {
	r.FileConfig.FilePath = filePath
	return r
}

func (r *Resource) BuildExtName(extName string) *Resource {
	r.FileConfig.ExtName = extName
	return r
}

func (r *Resource) BuildReloadTime(reloadTime int64) *Resource {
	r.FileConfig.ReloadTime = reloadTime
	return r
}

func (r *Resource) BuildMonitorPath(monitorPath string) *Resource {
	r.FileConfig.MonitorPath = monitorPath
	return r
}

func (r *Resource) Name() string {
	return "file"
}

func (r *Resource) Init(_ IDataConfig) {
	go r.newWatcher()
}

func (r *Resource) ReadBytes(configName string) (data []byte, error error) {
	if configName == "" {
		return nil, errors.New("configName is empty.")
	}

	fullPath, err := joinPath(r.MonitorPath, configName+"."+r.ExtName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("file not found. err = %v", err))
	}

	if isDir(fullPath) {
		return nil, errors.New(fmt.Sprintf("path is dir. fullPath = %s", fullPath))
	}

	data, err = os.ReadFile(fullPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("read file err. err = %v, path = %s", err, fullPath))
	}

	if len(data) < 1 {
		return nil, errors.New(fmt.Sprintf("configName = %s data is err.", configName))
	}

	return data, nil
}

func (r *Resource) OnChange(fn ConfigChangeFn) {
	r.changeFn = fn
}

func (r *Resource) newWatcher() {
	r.watcher = watcher.New()
	r.watcher.SetMaxEvents(20)
	r.watcher.FilterOps(watcher.Write)

	if err := r.watcher.Add(r.MonitorPath); err != nil {
		return
	}

	go func() {
		for {
			select {
			case ev := <-r.watcher.Event:
				{
					if ev.IsDir() {
						continue
					}

					configName := getFileName(ev.FileInfo.Name(), true)

					data, err := r.ReadBytes(configName)
					if err != nil {
						return
					}

					if r.changeFn != nil {
						r.changeFn(configName, data)
					}
				}
			case _ = <-r.watcher.Error:
				{
					return
				}
			case <-r.watcher.Closed:
				return
			}
		}
	}()

	if err := r.watcher.Start(time.Second * time.Duration(r.ReloadTime)); err != nil {
	}
}

func (r *Resource) Stop() {
	if r.watcher == nil {
		return
	}

	_ = r.watcher.Remove(r.MonitorPath)
}

func getFileName(filePath string, removeExt bool) string {
	fileName := path.Base(filePath)
	if removeExt == false {
		return fileName
	}

	var suffix string
	suffix = path.Ext(fileName)

	return strings.TrimSuffix(fileName, suffix)
}

func joinPath(elem ...string) (string, error) {
	filePath := filepath.Join(elem...)

	_, err := os.Stat(filePath)
	if err != nil {
		return filePath, err
	}
	return filePath, nil
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err == nil && info.IsDir() {
		return true
	}
	return false
}
