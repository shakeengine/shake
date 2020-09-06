package scene

import (
	"unsafe"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/veandco/go-sdl2/sdl"
)

func Init(widget *widgets.QWidget) {
	widgetID := widget.WinId()
	sdlWindow, err := sdl.CreateWindowFrom(unsafe.Pointer(widgetID))
	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(sdlWindow, 0, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}

	timer := core.NewQTimer(widget)
	timer.ConnectTimeout(func() {
		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.FillRect(nil)
		rect := sdl.Rect{X: 0, Y: 0, W: 200, H: 200}
		renderer.SetDrawColor(0xff, 0xff, 0x00, 0x00)
		renderer.FillRect(&rect)
		renderer.Present()
	})
	timer.Start(1000 / 60)
}
