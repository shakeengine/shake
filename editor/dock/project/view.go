package project

import (
	"os"

	"github.com/shakeengine/shake/misc"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// ProjectView : View for resource tree.
type ProjectView struct {
	dock         *widgets.QDockWidget
	resourcePath string
}

// NewProjectView : Create a project view.
func NewProjectView(projectPath string, mainWindow *widgets.QMainWindow) *ProjectView {
	view := &ProjectView{}

	view.init(projectPath)
	mainWindow.AddDockWidget(core.Qt__BottomDockWidgetArea, view.dock)

	return view
}

func (view *ProjectView) init(dir string) {
	view.resourcePath = dir + "/Resource"
	checkPath(view.resourcePath)

	/// -----------------------------------
	/// DockWidget - dock
	/// -> Widget - layout
	///    -> VBoxLayout
	///       -> TreeView
	// Project View is a dock widget.
	view.dock = widgets.NewQDockWidget("Project", nil, 0)
	dock := view.dock
	dock.SetAllowedAreas(core.Qt__AllDockWidgetAreas)
	// We need a widget that can contain a layout.
	layout := widgets.NewQWidget(nil, 0)
	layout.SetLayout(widgets.NewQVBoxLayout())
	dock.SetWidget(layout)
	// FileSystemModel with extension filters.
	model := widgets.NewQFileSystemModel(nil)
	model.SetRootPath(view.resourcePath)
	model.SetNameFilters([]string{"*.go"})
	model.SetNameFilterDisables(false)
	// Use TreeView to show the folder structure in easy way.
	tree := widgets.NewQTreeView(nil)
	tree.SetModel(model)
	tree.SetRootIndex(model.Index2(view.resourcePath, 0))
	layout.Layout().AddWidget(tree)
	/// -----------------------------------
}

// Destroy : Destroy all members in this instance.
func (view *ProjectView) Destroy() {
	if view.dock != nil {
		view.dock.Close()
	}
}

// checkPath : Check whether the given path is valid or not.
func checkPath(resourcePath string) {
	info, err := os.Stat(resourcePath)
	misc.ErrorCheck(err)
	if !info.IsDir() {
		panic("The given path is not a directory")
	}
}
