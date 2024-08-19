package main

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30)*time.Second)
	defer cancel()
	//最大跳数量
	cmd := exec.CommandContext(ctx, "mtr", "104.166.169.130", "-r", "-nz", "-c", "10", "-m", "30", "-p")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return
	}

	if err := cmd.Start(); err != nil {
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("lint", line)
	}
	if err := cmd.Wait(); err != nil {
		return
	}
	//执行超时
	if ctx.Err() == context.DeadlineExceeded {
		return
	}
	return
}
