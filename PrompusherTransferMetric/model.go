/**
 * Package main
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/11/21 16:13
 */

package main

import (
	"bytes"
	"errors"
	"fmt"
	"goStudy/PrompusherTransferMetric/promexport"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var (
	// ErrNoFileInfo means fileinfo is nil
	ErrNoFileInfo = errors.New("no-fileinfo")
	// ErrWrongFilename means wrong filename
	ErrWrongFilename = errors.New("wrong-filename")
)

type (
	Plugin struct {
		File         string
		LogFile      string
		FileModTime  int64
		LastExecTime int64
		Cycle        int64
		ReloadTime   int64
		timeout      time.Duration
	}
	ExecResult struct {
		Metrics []*promexport.Metric
		Error   error
		IsZero  bool
	}
)

func parseFilename(file string) (int64, error) {
	idx := strings.Index(file, "_")
	if idx < 0 {
		return 0, ErrWrongFilename
	}
	return strconv.ParseInt(file[:idx], 10, 64)
}

func calTimeout(cycle int64) time.Duration {
	duration := cycle - 1
	if duration < 30 {
		duration = 30
	}
	return time.Duration(duration) * time.Second
}

func (p *Plugin) ShouldExec(t int64) bool {
	return t-p.LastExecTime >= p.Cycle
}

func (p *Plugin) tryCleanLogFile() {
	info, _ := os.Stat(p.LogFile)
	if info == nil {
		return
	}
	// > 200M，删除文件
	if info.Size() > 1024*1024*200 {
		os.RemoveAll(p.LogFile)
	}
}

func (p *Plugin) writeLog(logData []byte) {
	if p.LogFile == "" {
		return
	}
	if rand.Intn(5) == 1 { // 随机处理清理文件
		p.tryCleanLogFile()
	}
	os.MkdirAll(filepath.Dir(p.LogFile), os.ModePerm)
	f, err := os.OpenFile(p.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return
	}
	f.WriteString(time.Now().Format(time.RFC3339Nano))
	f.WriteString("\t")
	f.Write(logData)
	f.Close()
}

func (p *Plugin) Exec() ([]byte, error) {
	p.LastExecTime = time.Now().Unix()

	var (
		stdout = bytes.NewBuffer(nil)
		stderr = bytes.NewBuffer(nil)
	)
	cmd := exec.Command(p.File)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	time.AfterFunc(p.timeout, func() {
		if cmd.Process != nil {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
	})

	err := cmd.Run()

	if stderr.Len() > 0 {
		errBytes := stderr.Bytes()
		if len(errBytes) > 256 { // 限制长度
			errBytes = errBytes[:255]
		}
		p.writeLog(errBytes)
	}

	if err != nil {
		if cmd.Process != nil {
			syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		return nil, fmt.Errorf("%v:%s", err, stderr.String())
	}
	return stdout.Bytes(), nil
}

func NewPlugin(file string, logFile string, info os.FileInfo) (*Plugin, error) {
	cycle, err := parseFilename(filepath.Base(file))
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, ErrNoFileInfo
	}
	p := &Plugin{
		File:         file,
		LogFile:      logFile,
		FileModTime:  info.ModTime().Unix(),
		Cycle:        cycle,
		LastExecTime: time.Now().Unix() - cycle + rand.Int63n(30) + 1, // 随机 30s 内开始执行
		timeout:      calTimeout(cycle),
	}
	return p, nil
}
