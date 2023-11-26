package openvpn

import (
	"os"
	"os/exec"
	"syscall"
)

func execCommand(command string, args []string) *exec.Cmd {
	cmd := exec.Command(command, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	// 为了通过 cmd.Output() 获取执行结果, 不在此处输出 Stdout
	// cmd.Stdout = os.Stdout

	cmd.Stderr = os.Stderr
	return cmd
}
