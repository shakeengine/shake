package dock

import (
	"github.com/shakeengine/shake/editor/dock/project"
	"github.com/therecipe/qt/widgets"
)

var mainWindow *widgets.QMainWindow = nil
var projectPath string = ""
var projectView *project.ProjectView = nil

// Init : Init the DockManager. DockManager is static.
func Init(window *widgets.QMainWindow) {
	mainWindow = window
}

// SetProjectPath : Set the project path.
func SetProjectPath(path string) {
	projectPath = path
}

// OpenProjectView : Open a project window.
func OpenProjectView() {
	if projectView != nil {
		projectView.Destroy()
		projectView = nil
	}
	projectView = project.NewProjectView(projectPath, mainWindow)
}
