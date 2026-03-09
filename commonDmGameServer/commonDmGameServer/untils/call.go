package untils

import (
	"runtime/debug"
)

// Call 安全地调用函数
func Call(fn func()) {
	if fn == nil {
		return
	}
	defer func() {
		if err := recover(); err != nil {
			stack := debug.Stack()
			PanicPoss(err, stack)
		}
	}()

	fn()
}

// 执行单个协程
func Go2(fn func()) {
	go Call(fn)
}
