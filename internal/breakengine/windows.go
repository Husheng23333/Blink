package breakengine

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// windowCmd 是从状态机 goroutine 发送给窗口控制 goroutine 的消息（非阻塞）。
// 状态机从不直接操作 AppKit：窗口的 Show/Hide/Close 和 NewWithOptions
// 内部都会调度到主线程，因此在持有引擎互斥锁时执行这些操作可能导致
// 与主线程上的绑定调用死锁。通过将所有窗口操作集中到这个单一 goroutine，
// 我们确保互斥锁不会因等待主线程而阻塞。
type windowCmd int

const (
	cmdShowBreakIcon windowCmd = iota
	cmdHideBreakIcon
	cmdShowNotice
	cmdHideNotice
)

// sendCmd 非阻塞地入队一个窗口命令。在引擎互斥锁下调用；
// 丢弃命令是无害的，因为后续命令会覆盖前面的（例如 hide 跟在 show 后面会胜出）。
func (e *Engine) sendCmd(c windowCmd) {
	select {
	case e.cmdCh <- c:
	default:
	}
}

func (e *Engine) showBreakIcon() { e.sendCmd(cmdShowBreakIcon) }
func (e *Engine) hideBreakIcon() { e.sendCmd(cmdHideBreakIcon) }
func (e *Engine) showNotice()    { e.sendCmd(cmdShowNotice) }
func (e *Engine) hideNotice()    { e.sendCmd(cmdHideNotice) }

// windowState 是 windowLoop 独占的窗口句柄集合。
//
// 它必须是每代循环各自持有的一份局部状态，而不是 Engine 的字段：
// Stop() 之后旧循环还要留下来排空退出命令（hide break icon，见 windowLoop），
// 而 Start() 可能已经启动了新一代循环，两代会在一段时间内并存。
// 若句柄挂在 Engine 上，两代就会同时读写同一组字段——这是数据竞争。
// 各自持有一份后，"windowLoop 独占窗口句柄"的不变式在跨代时依然成立，
// 因此这些字段不需要任何锁。
type windowState struct {
	breakIcon application.Window
	notice    application.Window
}

// windowLoop 拥有休息图标和通知窗口的控制权。它是唯一读写 ws.breakIcon
// 和 ws.notice 的 goroutine，因此这些字段不需要互斥锁。
//
// stopCh / cmdCh 由 Start 传入并只对应本代循环：本循环必须只消费自己
// 这一代的窗口命令。若读的是 e.cmdCh 字段，Stop→Start 换代后本循环会
// 开始窃取新一代引擎的命令，新一代自己的 windowLoop 反而收不到——
// 表现为休息图标该显示时不显示、该关闭时关不掉。
func (e *Engine) windowLoop(stopCh chan struct{}, cmdCh chan windowCmd) {
	var ws windowState
	for {
		select {
		case <-stopCh:
			// 退出前排空 Stop 入队的命令（如 hide break icon），
			// 否则上面的 select 可能先选中 stopCh，导致休息图标/通知窗口
			// 残留在屏幕上。
			for {
				select {
				case c := <-cmdCh:
					e.execWindowCmd(&ws, c)
				default:
					return
				}
			}
		case c := <-cmdCh:
			e.execWindowCmd(&ws, c)
		}
	}
}

// execWindowCmd 执行单个窗口命令。
func (e *Engine) execWindowCmd(ws *windowState, c windowCmd) {
	switch c {
	case cmdShowBreakIcon:
		// 图标窗口创建即可见（见 breakIconOptions），眨眼动画随页面挂载开始。
		// 不激活应用、不抢焦点：图标只是一个角落里的提示，不打断用户当前的窗口。
		e.ensureBreakIcon(ws)
		if ws.breakIcon != nil {
			ws.breakIcon.Show()
		}
	case cmdHideBreakIcon:
		// 销毁图标窗口而非隐藏：隐藏会保持 WKWebView 存活，
		// 其 WebKit 渲染进程仍驻留内存。下次休息时在 ensureBreakIcon
		// 中重建的开销很小。
		if ws.breakIcon != nil {
			ws.breakIcon.Close()
			ws.breakIcon = nil
		}
	case cmdShowNotice:
		e.ensureNotice(ws)
		if ws.notice != nil {
			ws.notice.Show()
		}
	case cmdHideNotice:
		if ws.notice != nil {
			ws.notice.Close()
			ws.notice = nil
		}
	}
}

// ensureBreakIcon 懒创建唯一的休息图标窗口，放置在主显示器右下角。
func (e *Engine) ensureBreakIcon(ws *windowState) {
	if ws.breakIcon != nil {
		return
	}
	a := e.app_()
	if a == nil {
		return
	}
	ws.breakIcon = a.Window.NewWithOptions(breakIconOptions(a))
}

// ensureNotice 懒创建唯一的休息前提醒窗口，放置在主显示器顶部居中。
func (e *Engine) ensureNotice(ws *windowState) {
	if ws.notice != nil {
		return
	}
	a := e.app_()
	if a == nil {
		return
	}
	ws.notice = a.Window.NewWithOptions(noticeOptions(a))
}
