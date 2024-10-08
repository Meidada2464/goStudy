/**
 * Package util
 * @Author fuqiang.li <fuqiang.li@baishan.com>
 * @Date 2024/1/18
 */

package util

import (
	"bytes"
	"os/exec"
	"syscall"
	"time"
)

// RunCommand 执行命令行
func RunCommand(name string, arg []string, timeout time.Duration) ([]byte, []byte, error) {
	var (
		stdout = bytes.NewBuffer(nil)
		stderr = bytes.NewBuffer(nil)
	)
	cmd := exec.Command(name, arg...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// 强制 kill
	time.AfterFunc(timeout, func() {
		if cmd.Process != nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
	})

	err := cmd.Run()
	return stdout.Bytes(), stderr.Bytes(), err
}
