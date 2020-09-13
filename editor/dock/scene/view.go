package scene

import (
	"time"
	"unsafe"

	"github.com/shakeengine/shake/misc"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// Time interval between frames of rendering.
const renderInterval = 1000 / 60

// SceneView : View for rendering the hierarchy.
type SceneView struct {
	qwidget   *widgets.QWidget
	sdlWindow *sdl.Window
	renderer  *sdl.Renderer
	// For debugging.
	font *ttf.Font
	// For dealing with resize event of window.
	resizeTimer *core.QTimer
	// For rendering a frame.
	renderTimer *core.QTimer
}

// NewSceneView : Create a scene view
func NewSceneView(widget *widgets.QWidget) *SceneView {
	view := &SceneView{}
	view.qwidget = widget
	view.resizeTimer = core.NewQTimer(view.qwidget)
	view.renderTimer = core.NewQTimer(view.qwidget)
	view.qwidget.ConnectResizeEvent(func(event *gui.QResizeEvent) {
		size := event.Size()
		view.OnResize(size.Width(), size.Height())
	})

	// We should use goroutine or we get a freeze at 'sdl.CreateWindowFrom'.
	go view.init()

	return view
}

func (view *SceneView) init() {
	view.initSDLWindow()
	view.initRenderer()
	view.initDefaultFont()
	view.initRenderTimer()
}

func (view *SceneView) initRenderTimer() {
	view.renderTimer.ConnectTimeout(func() {
		view.render()
	})
	view.renderTimer.Start(renderInterval)
}

func (view *SceneView) stopRender() {
	if view.renderTimer != nil {
		view.renderTimer.Stop()
	}
}

func (view *SceneView) resumeRender() {
	if view.renderTimer != nil {
		view.renderTimer.Start(renderInterval)
	}
}

func (view *SceneView) destorySDLWindow() {
	// We should destory the renderer first, because it came from SDLWindow.
	view.destroyRenderer()
	if view.sdlWindow != nil {
		view.sdlWindow.Destroy()
		view.sdlWindow = nil
	}
}

func (view *SceneView) initSDLWindow() {
	var err error
	view.destorySDLWindow()
	widgetID := view.qwidget.WinId()
	view.sdlWindow, err = sdl.CreateWindowFrom(unsafe.Pointer(widgetID))
	if err != nil {
		panic(err)
	}
}

func (view *SceneView) destroyRenderer() {
	// We should stop the [renderTimer] to prevent error.
	view.stopRender()
	if view.renderer != nil {
		view.renderer.Destroy()
		view.renderer = nil
	}
}

func (view *SceneView) initRenderer() {
	var err error
	view.destroyRenderer()
	view.renderer, err = sdl.CreateRenderer(view.sdlWindow, 0, sdl.RENDERER_ACCELERATED)
	misc.ErrorCheck(err)
}

func (view *SceneView) initDefaultFont() {
	var err error
	if view.font != nil {
		view.font.Close()
		view = nil
	}
	view.font, err = ttf.OpenFont("../font/OpenSans-Regular.ttf", 15)
	misc.ErrorCheck(err)
}

func (view *SceneView) render() {
	renderer := view.renderer
	renderer.SetDrawColor(0, 0, 0, 0)
	renderer.FillRect(nil)
	rect := sdl.Rect{X: 0, Y: 0, W: 200, H: 200}
	renderer.SetDrawColor(0xff, 0xff, 0x00, 0x00)
	renderer.FillRect(&rect)
	renderer.SetDrawColor(0xff, 0, 0, 0xff)

	/// ------------------------------------
	/// For debugging, especially to know that frames are rendered properly at proper time.
	surface, err := view.font.RenderUTF8Solid(time.Now().String(), sdl.Color{B: 255})
	if err != nil {
		panic(err)
	}
	texture, err := view.renderer.CreateTextureFromSurface(surface)
	defer texture.Destroy()
	_, _, w, h, err := texture.Query()
	if err != nil {
		panic(err)
	}
	renderer.Copy(texture, &sdl.Rect{X: 0, Y: 0, W: w, H: h}, &sdl.Rect{X: 0, Y: 0, W: w, H: h})
	/// ------------------------------------

	renderer.Present()
}

// Destroy : Destroy all members in this instance
func (view *SceneView) Destroy() {
	view.font.Close()
	view.destroyRenderer()
	view.destorySDLWindow()
}

// OnResize : Callback when the resize event is catched
func (view *SceneView) OnResize(w, h int) {
	view.sdlWindow.SetSize(int32(w), int32(h))
	view.renderer.SetLogicalSize(int32(w), int32(h))
	view.resizeTimer.ConnectTimeout(func() {
		view.resizeTimer.Stop()
		view.initRenderer()
		view.resumeRender()
	})
	// 1.0s later, destory ther renderer and recreate the renderer.
	// This method is not super sexy, but reasonable using both Qt and SDL.
	view.resizeTimer.Start(1000)
}
