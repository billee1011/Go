// Package logger 提供接口初始化日志的记录， 外部使用 logrus 提供的原生方法来记录日志
package logger

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/Sirupsen/logrus"
)

// 非并发安全的 writer， 用于记录日志
type logWriter struct {
	outputFile bool
	outputErr  bool
	fileName   string
	fileDir    string

	file     *os.File
	fileTime time.Time
}

// 日志文件按小时区分， 如果当前时间和上一次检测的时间相差了一个小时以上， 则打开一个新文件用以记录日志
func (w *logWriter) getFile() {
	now := time.Now().Local()
	if now.Hour() == w.fileTime.Hour() &&
		now.Sub(w.fileTime) < time.Duration(time.Second*3600) {
		return
	}
	w.fileTime = now
	fileName := fmt.Sprintf("%s/%s_%d_%02d_%02d_%02d.log", w.fileDir, w.fileName, now.Year(),
		now.Month(), now.Day(), now.Hour())

	// Umask 是权限的补码,用于设置创建文件和文件夹默认权限
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)

	var err error
	if w.file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0664); err != nil {
		fmt.Println("open log file failed:", err)
	}
}

// 写日志， 根据 logWriter 的状态判断是否需要写入标准错误和文件
func (w *logWriter) Write(buf []byte) (int, error) {
	if w.outputErr {
		os.Stderr.Write(buf)
	}
	if w.outputFile {
		w.getFile()
		if w.file != nil {
			w.file.Write(buf)
		}
	}
	return len(buf), nil
}

// SetupLog 初始化 log
// fileName 为 log 文件名前缀，真实的 log 文件名字为： fileName_YY_MM_DD_HH.log， 如果 fileName 为空或者文件打开失败， 则不会输出到日志文件中。
// dir 为日志文件存放目录， 如果目录打开失败， 将不会输出到文件
// level 为日志级别，高于此日志级别的日志才会被输出。 可取值从低到高分别为 debug, info, warning, error, fatal, panic
// usingErr 是否输出到标准错误， 此值为 true 时不影响输出到文件
func SetupLog(fileName string, dir string, level string, outputErr bool) {
	writer := new(logWriter)
	writer.outputErr = outputErr
	writer.outputFile = false
	if fileName != "" && dir != "" {

		// Umask 是权限的补码,用于设置创建文件和文件夹默认权限
		mask := syscall.Umask(0)
		defer syscall.Umask(mask)

		if err := os.MkdirAll(dir, 0774); err == nil {
			writer.outputFile = true
			writer.fileName = fileName
			writer.fileDir = dir
		}
	}
	if lv, err := logrus.ParseLevel(level); err == nil {
		logrus.SetLevel(lv)
	} else {
		fmt.Println("parse level failed:", level, err)
	}
	logrus.SetOutput(writer)

	logrus.WithFields(logrus.Fields{
		"file_name": fileName,
		"file_dir":  dir,
		"level":     level,
		"outputErr": outputErr,
	}).Debug("setup log")
}
