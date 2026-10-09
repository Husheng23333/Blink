package breakengine

import (
	"github.com/wailsapp/wails/v3/pkg/application"
)

// breakIconName 是休息图标窗口的名称。
const breakIconName = "pm-break-icon"

// breakIconOptions 构建休息图标窗口的配置项：主显示器右下角一个透明的小笑脸。
// 窗口不抢焦点、点击穿透（IgnoreMouseEvents），不遮挡用户的操作。
// 窗口层级设为 Floating，高于普通窗口但不压过 Dock 与菜单栏。
func breakIconOptions(a *application.App) application.WebviewWindowOptions {
	const (
		iconSize   = 96
		iconMargin = 24
		// 图标在水平方向整体左移一点，避免贴得太近右侧边缘。
		iconMarginX = iconMargin + 20
	)
	x, y := 0, 0
	if primary := a.Screen.GetPrimary(); primary != nil {
		wa := primary.WorkArea
		x = wa.X + wa.Width - iconSize - iconMarginX
		y = wa.Y + wa.Height - iconSize - iconMargin
	}
	return application.WebviewWindowOptions{
		Name:              breakIconName,
		Title:             "",
		Frameless:         true,
		AlwaysOnTop:       true,
		DisableResize:     true,
		IgnoreMouseEvents: true,
		URL:               "/#break",
		Width:             iconSize,
		Height:            iconSize,
		X:                 x,
		Y:                 y,
		InitialPosition:   application.WindowXY,
		BackgroundType:    application.BackgroundTypeTransparent,
		BackgroundColour:  application.NewRGBA(0, 0, 0, 0),
		Mac: application.MacWindow{
			WindowLevel:        application.MacWindowLevelFloating,
			CollectionBehavior: application.MacWindowCollectionBehaviorCanJoinAllSpaces | application.MacWindowCollectionBehaviorStationary,
			TitleBar:           application.MacTitleBar{Hide: true, AppearsTransparent: true, FullSizeContent: true},
			DisableShadow:      true,
		},
		// Windows：与休息前提醒一致，不占任务栏，去掉 OS 阴影/圆角（阴影交给 CSS）。
		Windows: application.WindowsWindow{
			HiddenOnTaskbar:                   true,
			DisableFramelessWindowDecorations: true,
		},
		CloseButtonState:      application.ButtonHidden,
		MinimiseButtonState:   application.ButtonHidden,
		MaximiseButtonState:   application.ButtonHidden,
		FullscreenButtonState: application.ButtonHidden,
	}
}

// noticeOptions 构建休息前提醒窗口的配置项。
// 窗口透明，使用 macOS 原生半透明效果，放置在主显示器顶部居中。
func noticeOptions(a *application.App) application.WebviewWindowOptions {
	const (
		noticeWidth  = 380
		noticeHeight = 122
	)
	x, y := 0, 80
	if primary := a.Screen.GetPrimary(); primary != nil {
		x = primary.WorkArea.X + (primary.WorkArea.Width-noticeWidth)/2
		y = primary.WorkArea.Y + 80
	}
	return application.WebviewWindowOptions{
		Name:             "pm-notice",
		Title:            "",
		Frameless:        true,
		AlwaysOnTop:      true,
		DisableResize:    true,
		Hidden:           true,
		URL:              "/#notice",
		Width:            noticeWidth,
		Height:           noticeHeight,
		X:                x,
		Y:                y,
		InitialPosition:  application.WindowXY,
		BackgroundType:   application.BackgroundTypeTransparent,
		BackgroundColour: application.NewRGBA(0, 0, 0, 0),
		Mac: application.MacWindow{
			Backdrop:    application.MacBackdropTranslucent,
			WindowLevel: application.MacWindowLevelFloating,
			TitleBar:    application.MacTitleBar{Hide: true, AppearsTransparent: true, FullSizeContent: true},
		},
		// Windows：透明窗口走 WS_EX_NOREDIRECTIONBITMAP（DirectComposition），
		// 仍正常接收鼠标事件，所以"推迟"按钮可用；不设置 IgnoreMouseEvents，
		// 否则会加 WS_EX_TRANSPARENT 导致点击穿透。HiddenOnTaskbar 让提醒不出现在
		// 任务栏；DisableFramelessWindowDecorations 去掉 OS 阴影/圆角，圆角交给 CSS。
		Windows: application.WindowsWindow{
			HiddenOnTaskbar:                   true,
			DisableFramelessWindowDecorations: true,
		},
		CloseButtonState:      application.ButtonHidden,
		MinimiseButtonState:   application.ButtonHidden,
		MaximiseButtonState:   application.ButtonHidden,
		FullscreenButtonState: application.ButtonHidden,
	}
}
