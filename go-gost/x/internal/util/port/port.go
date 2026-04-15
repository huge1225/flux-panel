package port

import (
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runConntrackDelete(port int) error {
	commands := [][]string{
		{"-D", "-p", "tcp", "--dport", fmt.Sprintf("%d", port)},
		{"-D", "-p", "tcp", "--sport", fmt.Sprintf("%d", port)},
		{"-D", "-p", "udp", "--dport", fmt.Sprintf("%d", port)},
		{"-D", "-p", "udp", "--sport", fmt.Sprintf("%d", port)},
	}

	var runErrs []string
	for _, args := range commands {
		cmd := exec.Command("conntrack", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			// conntrack 没匹配连接时会返回非零，不作为失败处理
			if !strings.Contains(string(out), "0 flow entries have been deleted") &&
				!strings.Contains(string(out), "Operation failed") {
				runErrs = append(runErrs, fmt.Sprintf("%v: %s", err, strings.TrimSpace(string(out))))
			}
		}
	}

	if len(runErrs) > 0 {
		return fmt.Errorf(strings.Join(runErrs, "; "))
	}
	return nil
}

func ForceClosePortConnections(addr string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("⚠️ ForceClosePortConnections panic recovered: %v\n", r)
			err = nil // 永远返回 nil
		}
	}()

	if addr == "" {
		fmt.Println("⚠️ 地址为空")
		return nil
	}

	_, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		fmt.Printf("⚠️ 地址解析失败: %v\n", err)
		return nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Printf("⚠️ 端口非法: %v\n", err)
		return nil
	}

	cmd := exec.Command("tcpkill", "-i", "any", "port", fmt.Sprintf("%d", port))
	if err := cmd.Start(); err != nil {
		fmt.Printf("⚠️ 启动 tcpkill 失败，尝试 conntrack(nft): %v\n", err)
		if conntrackErr := runConntrackDelete(port); conntrackErr != nil {
			fmt.Printf("⚠️ conntrack 清理失败: %v\n", conntrackErr)
		} else {
			fmt.Printf("✅ 已通过 conntrack(nft) 清理端口 %d 连接跟踪条目\n", port)
		}
		return nil
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("⚠️ tcpkill goroutine panic recovered: %v\n", r)
			}
		}()
		time.Sleep(2 * time.Second)
		if cmd.Process != nil {
			if err := cmd.Process.Kill(); err != nil {
				fmt.Printf("⚠️ 终止 tcpkill 失败: %v\n", err)
			}
		}
	}()

	fmt.Printf("✅ 正在断开端口 %d 上的所有连接...\n", port)
	return nil
}
