package scrcpy

import (
	"andtest/logger"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/pion/webrtc/v3/pkg/media/h264reader"
)

const (
	jarPath      = "/data/local/tmp/scrcpy-server"
	fps          = 15
	videoBitRate = 2500000
	maxSize      = 720
)

func Scrcpy(tag string) {
	// 开启服务
	go _start(tag)
	time.Sleep(5 * time.Second)
	_video(tag)
}

func _start(tag string) error {
	args := []string{
		"/", "com.genymobile.scrcpy.Server", "3.3.3",
		"tunnel_forward=true",
		"audio=false",
		"control=false",
		"video=true",
		"raw_stream=true",
		"video_codec=h264",
		fmt.Sprintf("max_fps=%d", fps),
		fmt.Sprintf("video_bit_rate=%d", videoBitRate),
		fmt.Sprintf("max_size=%d", maxSize),
	}
	cmd := exec.Command("app_process", args...)
	// 设置环境变量加载 JAR 包
	cmd.Env = append(os.Environ(), fmt.Sprintf("CLASSPATH=%s", jarPath))

	// 设置标准输出以便调试
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		logger.LogE(tag, fmt.Sprintf("启动程序SCRCPY错误：%s", err))
		return err
	}
	return nil
}

func _video(tag string) {
	conn, err := net.Dial("unix", "@scrcpy")
	if err != nil {
		logger.LogE(tag, fmt.Sprintf("连接错误:%s", err))
		return
	}
	defer conn.Close()
	logger.LogI(tag, "开始读取视频流数据")
	reader, err := h264reader.NewReader(conn)
	if err != nil {
		logger.LogE(tag, fmt.Sprintf("读取视频流错误：%s", err))
		return
	}
	for {
		nal, err := reader.NextNAL()
		if err != nil {
			if err == io.EOF {
				logger.LogE(tag, "服务端已关闭链接")
			} else {
				logger.LogE(tag, fmt.Sprintf("解析NAL错误：%s", err))
			}
			break
		}
		logger.LogD(tag, fmt.Sprintf("帧数据:%x", nal.Data[:32]))
	}
}
