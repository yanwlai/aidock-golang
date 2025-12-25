package control

import (
	"andtest/logger"
	"encoding/binary"
	"fmt"
	"os"
	"time"
)

const (
	EV_SYN = 0x00
	EV_KEY = 0x01
	EV_ABS = 0x03

	BTN_TOUCH          = 0x14a
	ABS_MT_TOUCH_MAJOR = 0x30
	ABS_MT_TOUCH_MINOR = 0x31
	ABS_MT_ORIENTATION = 0x34
	ABS_MT_POSITION_X  = 0x35
	ABS_MT_POSITION_Y  = 0x36
	ABS_MT_TRACKING_ID = 0x39
	ABS_MT_PRESSURE    = 0x3a
	SYN_REPORT         = 0x00
)

func Control(tag string) {
	fd, err := os.OpenFile("/dev/input/event3", os.O_WRONLY, 0666)
	if err != nil {
		logger.LogE(tag, fmt.Sprintf("打开文件错误：%s", err))
		return
	}
	defer fd.Close()

	for {
		logger.LogI(tag, "开始点击")
		tap(fd, 487, 1257)
		time.Sleep(5 * time.Second)
		logger.LogI(tag, "滑动返回")
		swipe(fd, 0, 1200, 800, 1200, 500*time.Millisecond)
		time.Sleep(20 * time.Second)
	}

	// 坐标 0x1e7 = 487, 0x4e9 = 1257
	logger.LogD(tag, "开始点击")
	//tap(fd, 0x1e7, 0x4e9)
	swipe(fd, 800, 1200, 300, 1200, 500*time.Millisecond)
}

func tap(f *os.File, x, y int32) {
	// --- 按下序列 (严格参考日志内容) ---
	writeEvent(f, EV_KEY, BTN_TOUCH, 1)
	writeEvent(f, EV_ABS, ABS_MT_TRACKING_ID, 0x50) // 使用日志中的 ID 0x50
	writeEvent(f, EV_ABS, ABS_MT_POSITION_X, x)
	writeEvent(f, EV_ABS, ABS_MT_POSITION_Y, y)
	writeEvent(f, EV_ABS, ABS_MT_TOUCH_MAJOR, 0x9d)
	writeEvent(f, EV_ABS, ABS_MT_TOUCH_MINOR, 0x9c)
	writeEvent(f, EV_ABS, ABS_MT_PRESSURE, 0x4c)
	writeEvent(f, EV_ABS, ABS_MT_ORIENTATION, -783) // 对应 fffffcf1
	writeEvent(f, EV_SYN, SYN_REPORT, 0)

	time.Sleep(50 * time.Millisecond)

	// --- 抬起序列 ---
	writeEvent(f, EV_ABS, ABS_MT_PRESSURE, 0)
	writeEvent(f, EV_ABS, ABS_MT_TRACKING_ID, -1)
	writeEvent(f, EV_ABS, ABS_MT_ORIENTATION, 0)
	writeEvent(f, EV_KEY, BTN_TOUCH, 0)
	writeEvent(f, EV_SYN, SYN_REPORT, 0)
}

func swipe(f *os.File, x1, y1, x2, y2 int32, duration time.Duration) {
	// 1. 按下 (Start)
	writeEvent(f, EV_ABS, ABS_MT_TRACKING_ID, 0x60) // 换个 ID 避免冲突
	writeEvent(f, EV_KEY, BTN_TOUCH, 1)
	writeEvent(f, EV_ABS, ABS_MT_PRESSURE, 0x4c)
	writeEvent(f, EV_ABS, ABS_MT_POSITION_X, x1)
	writeEvent(f, EV_ABS, ABS_MT_POSITION_Y, y1)
	writeEvent(f, EV_SYN, SYN_REPORT, 0)

	// 2. 移动 (Move)
	steps := 20 // 采样步数，越多越平滑
	sleepStep := duration / time.Duration(steps)

	for i := 1; i <= steps; i++ {
		// 计算当前步的坐标
		currX := x1 + (x2-x1)*int32(i)/int32(steps)
		currY := y1 + (y2-y1)*int32(i)/int32(steps)

		writeEvent(f, EV_ABS, ABS_MT_POSITION_X, currX)
		writeEvent(f, EV_ABS, ABS_MT_POSITION_Y, currY)
		writeEvent(f, EV_SYN, SYN_REPORT, 0) // 每一帧移动都要同步

		time.Sleep(sleepStep)
	}

	// 3. 抬起 (End)
	writeEvent(f, EV_ABS, ABS_MT_PRESSURE, 0)
	writeEvent(f, EV_ABS, ABS_MT_TRACKING_ID, -1)
	writeEvent(f, EV_KEY, BTN_TOUCH, 0)
	writeEvent(f, EV_SYN, SYN_REPORT, 0)
}

func writeEvent(f *os.File, typ, code uint16, value int32) {
	// 64位架构 24字节 buffer
	buf := make([]byte, 24)
	// 前16字节是秒和微秒，留空或填入当前时间。内核注入通常不强制校验。
	// binary.LittleEndian.PutUint64(buf[0:8], uint64(time.Now().Unix()))

	binary.LittleEndian.PutUint16(buf[16:18], typ)
	binary.LittleEndian.PutUint16(buf[18:20], code)
	binary.LittleEndian.PutUint32(buf[20:24], uint32(value))

	f.Write(buf)
}
