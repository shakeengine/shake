package main

import (
	"fmt"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	vk "github.com/vulkan-go/vulkan"

	"github.com/shakeengine/shake/editor/dock"
	"github.com/shakeengine/shake/editor/dock/scene"
	"github.com/shakeengine/shake/editor/menu"
	"github.com/shakeengine/shake/editor/selector"

	_ "github.com/mattn/go-sqlite3"
)

const (
	canvasObjectName  = "Canvas"
	windowTitle       = "Shake Engine Editor"
	windowTitleFormat = "Shake Engine Editor - %s"
)

// ShakeEditorWindow : the main frame editor window.
type ShakeEditorWindow struct {
	main        *widgets.QMainWindow
	projectPath string
	sceneView   *scene.SceneView
}

// NewShakeEditorWindow : Main Shake editor window.
func NewShakeEditorWindow() *ShakeEditorWindow {
	window := &ShakeEditorWindow{}
	window.init()

	selector.NewProjectListWindow(window.main, func(path string) {
		// Callback when user select a project to open.
		window.openProject(path)
	}, func() {
		// If window.projectPath is not empty, open the project because the project has been selected in ProjectListWindow.
		// Otherwise, just close this main ShakeEditorWindow.
		if window.projectPath == "" {
			window.main.Close()
		}
	})

	return window
}

func (window *ShakeEditorWindow) init() {
	///------------------------------------------
	/// MainWindow - mainWindow
	/// -> Widget - widget
	///    -> VBoxLayout
	///       -> Widget - canvas
	window.main = widgets.NewQMainWindow(nil, 0)
	mainWindow := window.main
	mainWindow.SetMinimumSize2(800, 480)
	mainWindow.SetWindowTitle(windowTitle)

	menu.SetMainWindow(mainWindow)
	menu.InitDefaultMenu()

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	mainWindow.SetCentralWidget(widget)

	dock.Init(mainWindow)

	canvas := widgets.NewQWidget(nil, 0)
	canvas.SetObjectName(canvasObjectName)
	widget.Layout().AddWidget(canvas)
	///------------------------------------------

	mainWindow.Show()
}

func (window *ShakeEditorWindow) openProject(projectPath string) {
	window.projectPath = projectPath

	// Change the window title to represent project path.
	window.main.SetWindowTitle(fmt.Sprintf(windowTitleFormat, window.projectPath))

	// Initialize SceneView.
	window.initSceneView()

	// Give project path to DockManager.
	dock.SetProjectPath(window.projectPath)
	dock.OpenProjectView()
}

// initSceneView : Open a SceneView.
// ShakeEditorWindow manage SceneView directly because SceneView is not a dock.
func (window *ShakeEditorWindow) initSceneView() {
	// If SceneView exists, destroy it.
	if window.sceneView != nil {
		window.sceneView.Destroy()
		window.sceneView = nil
	}
	// Find the canvas object.
	canvasObject := window.main.FindChild(canvasObjectName, core.Qt__FindChildrenRecursively)
	canvas := widgets.NewQWidgetFromPointer(canvasObject.Pointer())
	// Create a SceneView instance.
	window.sceneView = scene.NewSceneView(canvas)
}

func initVulkan() {
	// Initialize somethings for vulkan.
	// Do not use defer in this function because instances for vulkan must remain until the shake editor being closed.

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	// In order to use font to show some texts for debugging.
	if err := ttf.Init(); err != nil {
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
	ttf.Quit()
	sdl.Quit()
}

func main() {
	app := widgets.NewQApplication(len(os.Args), os.Args)

	initVulkan()
	defer releaseVulkan()

	NewShakeEditorWindow()

	app.Exec()
}
