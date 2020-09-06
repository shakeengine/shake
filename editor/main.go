// #cgo LDFLAGS: -lole32 -loleaut32 -limm32 -lversion
package main

import (
	"os"

	"github.com/therecipe/qt/widgets"
	"github.com/veandco/go-sdl2/sdl"
	vk "github.com/vulkan-go/vulkan"

	"github.com/shakeengine/shake/editor/dock"
	"github.com/shakeengine/shake/editor/dock/scene"
	"github.com/shakeengine/shake/editor/menu"
)

func initUIFrame() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := widgets.NewQMainWindow(nil, 0)
	window.SetMinimumSize2(250, 200)
	window.SetWindowTitle("Shake Engine Editor")

	menu.SetMainWindow(window)
	menu.InitDefaultMenu()

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	window.SetCentralWidget(widget)

	input := widgets.NewQLineEdit(nil)
	input.SetPlaceholderText("Write something ...")
	widget.Layout().AddWidget(input)

	button := widgets.NewQPushButton2("and click me!", nil)
	button.ConnectClicked(func(bool) {
		widgets.QMessageBox_Information(nil, "OK", input.Text(), widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
	widget.Layout().AddWidget(button)

	dock.Init(window)

	canvas := widgets.NewQWidget(nil, 0)
	widget.Layout().AddWidget(canvas)

	scene.Init(canvas)

	window.Show()
	app.Exec()
}

func initVulkan() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	if err := sdl.VulkanLoadLibrary(""); err != nil {
		sdl.Quit()
		panic(err)
	}

	procAddr := sdl.VulkanGetVkGetInstanceProcAddr()
	if procAddr == nil {
		sdl.VulkanUnloadLibrary()
		sdl.Quit()
		panic("GetInstanceProcAddr is nil")
	}
	vk.SetGetInstanceProcAddr(procAddr)
	if err := vk.Init(); err != nil {
		sdl.VulkanUnloadLibrary()
		sdl.Quit()
		panic(err)
	}
}

func releaseVulkan() {
	sdl.VulkanUnloadLibrary()
	sdl.Quit()
}

func main() {
	initVulkan()

	initUIFrame()

	releaseVulkan()
}
