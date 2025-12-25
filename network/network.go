package network

import (
	"andtest/logger"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

func Network(tag string) {
	targetUrl := "https://www.baidu.com"

	logger.LogI(tag, "准备发起请求1")

	// 跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: tr,
	}

	resp, err := client.Get(targetUrl)
	if err != nil {
		logger.LogE(tag, fmt.Sprintf("请求错误:%v", err))
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		successMsg := fmt.Sprintf("请求成功! 状态: %s, 协议: %s", resp.Status, resp.Proto)
		logger.LogI(tag, successMsg)
	} else {
		warnMsg := fmt.Sprintf("请求完成但状态异常: %d", resp.StatusCode)
		logger.LogW(tag, warnMsg)
	}
}
