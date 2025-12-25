package logger

/*
#cgo LDFLAGS: -llog
#include <android/log.h>
#include <stdlib.h>
void android_log(int prio, const char* tag, const char* text) {
	__android_log_write(prio, tag, text);
}
*/
import "C"
import "unsafe"

const (
	ANDROID_LOG_DEBUG = 3
	ANDROID_LOG_INFO  = 4
	ANDROID_LOG_WARN  = 5
	ANDROID_LOG_ERROR = 6
)

func LogD(tag, text string) {
	cTag := C.CString(tag)
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cTag))
	defer C.free(unsafe.Pointer(cText))
	C.android_log(ANDROID_LOG_DEBUG, cTag, cText)
}

func LogI(tag, text string) {
	cTag := C.CString(tag)
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cTag))
	defer C.free(unsafe.Pointer(cText))
	C.android_log(ANDROID_LOG_INFO, cTag, cText)
}

func LogW(tag, text string) {
	cTag := C.CString(tag)
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cTag))
	defer C.free(unsafe.Pointer(cText))
	C.android_log(ANDROID_LOG_WARN, cTag, cText)
}

func LogE(tag, text string) {
	cTag := C.CString(tag)
	cText := C.CString(text)
	defer C.free(unsafe.Pointer(cTag))
	defer C.free(unsafe.Pointer(cText))
	C.android_log(ANDROID_LOG_ERROR, cTag, cText)
}
