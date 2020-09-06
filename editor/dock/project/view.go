package project

import (
	"os"

	"github.com/shakeengine/shake/misc"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

var opennedDir string = ""

// OpenProjectDirectory : Open the project directory in ShakeEditor
func OpenProjectDirectory(dir string) {
	opennedDir = dir
}

// OpenProjectView : Open the project view
func OpenProjectView() *widgets.QDockWidget {
	info, err := os.Stat(opennedDir)
	misc.ErrorCheck(err)
	if !info.IsDir() {
		panic("The given path is not a directory")
	}

	return createNewDockWidget()
}

func createNewDockWidget() *widgets.QDockWidget {
	// Project View is a dock widget
	dock := widgets.NewQDockWidget("Project", nil, 0)
	dock.SetAllowedAreas(core.Qt__AllDockWidgetAreas)
	dock.SetLayout(widgets.NewQVBoxLayout())

	// We need a widget that can contain a layout
	layout := widgets.NewQWidget(nil, 0)
	layout.SetLayout(widgets.NewQVBoxLayout())
	dock.SetWidget(layout)

	// FileSystemModel with extension filters
	model := widgets.NewQFileSystemModel(nil)
	model.SetRootPath(opennedDir)
	model.SetNameFilters([]string{"*.go"})
	model.SetNameFilterDisables(false)

	// Use TreeView to show the folder structure in easy way
	view := widgets.NewQTreeView(nil)
	view.SetModel(model)
	view.SetRootIndex(model.Index2(opennedDir, 0))
	layout.Layout().AddWidget(view)

	return dock
}
