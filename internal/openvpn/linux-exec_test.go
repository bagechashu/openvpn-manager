package openvpn

import (
	"fmt"
	"testing"
)

func TestExecCommand(t *testing.T) {
	command := "echo"
	args := []string{"hello", "world"} // 你的命令参数，这里使用 ls 的 -l 和 -a 参数作为示例

	// 创建新的进程
	cmd := execCommand(command, args)

	// 执行命令并获取输出
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("命令执行出错:", err)
		return
	}

	// 输出命令的输出结果
	fmt.Printf("Cmd Output: %s", string(output))

	got := string(output)
	want := "hello world\n"
	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}
