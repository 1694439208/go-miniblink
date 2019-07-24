// miniblink project main.go
package main

import (
	"fmt"
	_ "time"

	"syscall"
	"unsafe"

	"github.com/jthmath/winapi"
)

type wkeWindowType int32

//任务队列,保证所有的API调用都在痛一个线程
var jobQueue = make(chan func())

const (
	WKE_WINDOW_TYPE_POPUP wkeWindowType = iota
	WKE_WINDOW_TYPE_TRANSPARENT
	WKE_WINDOW_TYPE_CONTROL
)

type miniblink struct {
	mbHandle   syscall.Handle
	wkeWebView uintptr
}

func IntPtr(n int) uintptr {
	return uintptr(n)
}

var itob = func(i int) bool {
	if i == 0 {
		return false
	}
	return true
}

var btoi = func(b bool) int {
	if b {
		return 1
	}
	return 0
}

func StrPtr(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(s)))
}
func NewMb() *miniblink {
	mb, _ := syscall.LoadLibrary("node.dll")
	return &miniblink{mbHandle: mb}
}
func (mb *miniblink) FreeMb() {
	syscall.FreeLibrary(mb.mbHandle)
}
func (mb *miniblink) wkeVersion() int {
	wkeVersion, _ := syscall.GetProcAddress(mb.mbHandle, "wkeVersion")
	ret, _, _ := syscall.Syscall(wkeVersion, 0, 0, 0, 0)
	return int(ret)
}
func (mb *miniblink) wkeInitialize() {
	HWND, _ := syscall.GetProcAddress(mb.mbHandle, "wkeInitialize")
	syscall.Syscall(HWND, 0, 0, 0, 0)
}
func (mb *miniblink) wkeCreateWebWindow(Type uintptr, parent int, x int, y int, width int, height int) uintptr {
	HWND, _ := syscall.GetProcAddress(mb.mbHandle, "wkeCreateWebWindow")
	ret, _, _ := syscall.Syscall6(HWND, 6, Type, IntPtr(parent), IntPtr(x), IntPtr(y), IntPtr(width), IntPtr(height))
	mb.wkeWebView = ret
	return ret
}
func (mb *miniblink) wkeShowWindow(showFlag bool) {
	HWND, _ := syscall.GetProcAddress(mb.mbHandle, "wkeShowWindow")
	syscall.Syscall(HWND, 2, mb.wkeWebView, IntPtr(1), 0)
}
func (mb *miniblink) wkeLoadURL(url string) {
	HWND, _ := syscall.GetProcAddress(mb.mbHandle, "wkeLoadURLW")
	syscall.Syscall(HWND, 2, mb.wkeWebView, StrPtr(url), 0)
}

var (
	//    kernel32, _        = syscall.LoadLibrary("kernel32.dll")
	//    getModuleHandle, _ = syscall.GetProcAddress(kernel32, "GetModuleHandleW")

	user32, _     = syscall.LoadLibrary("user32.dll")
	messageBox, _ = syscall.GetProcAddress(user32, "MessageBoxW")
)

func abort(funcname string, err error) {
	panic(funcname + " failed: " + err.Error())
}
func MessageBox(caption, text string, style uintptr) (result int) {
	ret, _, _ := syscall.Syscall9(messageBox,
		4,
		0,
		StrPtr(text),
		StrPtr(caption),
		style,
		0, 0, 0, 0, 0)

	result = int(ret)
	return
}
func main() {
	//num := MessageBox("Do你我他ne Title", "This test is Done.", 0x00000003)
	//退出信号
	mb := NewMb()
	mb.wkeInitialize()
	HWND := mb.wkeCreateWebWindow(uintptr(WKE_WINDOW_TYPE_POPUP), 0, 0, 0, 1080, 680)
	fmt.Println(int(HWND))
	fmt.Println(int(mb.wkeWebView))
	fmt.Println(int(mb.mbHandle))
	mb.wkeShowWindow(true)
	//var fn = syscall.NewCallbackCDecl(cb_my) // 注意调用约定

	mb.wkeLoadURL("http://music.163.com/song?id=1296893537&userid=290764119")
	//启动一个新的协程来处理blink的API调用
	//<-make(chan bool)
	// 3. 主消息循环
	var msg winapi.MSG
	msg.Message = winapi.WM_QUIT + 1 // 让它不等于 winapi.WM_QUIT

	for winapi.GetMessage(&msg, 0, 0, 0) > 0 {
		winapi.TranslateMessage(&msg)
		winapi.DispatchMessage(&msg)
	}
	/*go func() {
		//将这个协程锁在当前的线程上
		runtime.LockOSThread()
		//消费API调用,同时处理好windows消息
		for {
			select {
			case job := <-jobQueue:
				job()
			default:
				//消息循环
				msg := &win.MSG{}
				if win.GetMessage(msg, 0, 0, 0) != 0 {
					win.TranslateMessage(msg)
					//是否传递下去
					next := true
					//拿到对应的webview

					if next {
						win.DispatchMessage(msg)
					}
				}
			}
		}
	}()*/
	fmt.Println("0") //mb.wkeVersion()
}
