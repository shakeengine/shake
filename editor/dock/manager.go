package dock

import (
	"github.com/shakeengine/shake/editor/dock/project"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

var mainWindow *widgets.QMainWindow = nil

// Init : Init the DockManager. DockManager is a singleton
func Init(window *widgets.QMainWindow) {
	mainWindow = window
}

// OpenProjectDirectory : Open the project directory in ShakeEditor
func OpenProjectDirectory(dir string) {
	project.OpenProjectDirectory(dir)
}

// OpenProjectView : Open a project window
func OpenProjectView() {
	dock := project.OpenProjectView()
	mainWindow.AddDockWidget(core.Qt__BottomDockWidgetArea, dock)
}
