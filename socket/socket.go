package socket

import (
	"andtest/logger"
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func Socket(tag string) {
	socketPath := "\x00my_local_socket"

	// 建立连接
	conn, err := net.DialTimeout("unix", socketPath, 3*time.Second)
	if err != nil {
		logger.LogE(tag, fmt.Sprintf("连接失败：%s", err))
		return
	}
	defer conn.Close()

	logger.LogI(tag, "成功接入系统")

	go func() {
		for {
			logger.LogI(tag, "开始发送消息")
			// 写入数据（根据系统服务的具体协议）
			_, err = conn.Write([]byte("HELLO\n"))
			if err != nil {
				logger.LogE(tag, fmt.Sprintf("发送消息失败：%s", err))
			}
			time.Sleep(10 * time.Second)
		}
	}()

	// 读取返回
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n') // 阻塞直到收到换行符
		if err != nil {
			if err == io.EOF {
				logger.LogE(tag, "服务端断开连接 (EOF)")
			} else {
				logger.LogE(tag, "读取错误: "+err.Error())
			}
			break
		}
		logger.LogD(tag, fmt.Sprintf("收到响应：%s", line))
	}
}
