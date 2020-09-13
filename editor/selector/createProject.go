package selector

import (
	"errors"
	"os"
	"strings"

	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// CreateProjectDialog : A dialog for creating a new project.
type CreateProjectDialog struct {
	qdialog *widgets.QDialog
}

// NewCreateProjectDialog : Dialog for creating a new project.
func NewCreateProjectDialog(parent *widgets.QDialog, onClose func()) *CreateProjectDialog {
	dialog := &CreateProjectDialog{}
	dialog.init(parent, onClose)

	return dialog
}

// openNewProejctModalDialog : Open the ProjectSelector window.
func (dialog *CreateProjectDialog) init(parent *widgets.QDialog, onClose func()) {
	dialog.qdialog = widgets.NewQDialog(parent, 0)
	qdialog := dialog.qdialog
	qdialog.SetWindowTitle("New Project")
	qdialog.SetMinimumSize2(500, 100)
	qdialog.SetModal(true)

	qdialog.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		onClose()
	})

	mainLayout := widgets.NewQVBoxLayout()

	/// --------------------------------------
	/// Add widgets into formLayout
	formLayout := widgets.NewQGridLayout2()
	labelName := widgets.NewQLabel2("Project Name", nil, 0)
	labelPath := widgets.NewQLabel2("Project Path", nil, 0)

	textName := widgets.NewQLineEdit(nil)
	textPath := widgets.NewQLineEdit(nil)
	buttonPath := widgets.NewQPushButton2("...", nil)
	buttonPath.ConnectClicked(func(checked bool) {
		dir := widgets.QFileDialog_GetExistingDirectory(nil, "Select a folder which contains a Shake project", "", widgets.QFileDialog__ShowDirsOnly)
		textPath.SetText(dir)
	})

	formLayout.AddWidget2(labelName, 0, 0, 0)
	formLayout.AddWidget2(labelPath, 1, 0, 0)

	formLayout.AddWidget2(textName, 0, 1, 0)
	formLayout.AddWidget2(textPath, 1, 1, 0)

	formLayout.AddWidget2(buttonPath, 1, 2, 0)
	mainLayout.AddLayout(formLayout, 0)
	/// --------------------------------------

	/// --------------------------------------
	/// Add button widgets into buttonLayout
	buttonLayout := widgets.NewQHBoxLayout()
	buttonCancel := widgets.NewQPushButton2("Cancel", nil)
	buttonCancel.ConnectClicked(func(checked bool) {
		qdialog.Close()
	})
	buttonCreate := widgets.NewQPushButton2("Create", nil)
	buttonCreate.ConnectClicked(func(checked bool) {
		path := textPath.Text()
		name := textName.Text()
		if err := pathCheck(path, name); err != nil {
			widgets.QMessageBox_Critical(nil, "Failed to create", err.Error(), widgets.QMessageBox__Close, widgets.QMessageBox__Close)
		}
		if createNewProject(path, name) {
			qdialog.Close()
		}
	})
	buttonLayout.AddStretch(10)
	buttonLayout.AddWidget(buttonCancel, 0, 0)
	buttonLayout.AddWidget(buttonCreate, 0, 0)
	mainLayout.AddLayout(buttonLayout, 0)
	/// --------------------------------------

	qdialog.SetLayout(mainLayout)
	qdialog.Show()
}

func concatPathWithProjectName(path, projectName string) string {
	replaced := strings.ReplaceAll(path, "\\", "/")
	trimPath := strings.TrimSuffix(replaced, "/")

	spaceRemoved := strings.ReplaceAll(projectName, " ", "")
	tabRemoved := strings.ReplaceAll(spaceRemoved, "	", "")

	return trimPath + "/" + tabRemoved
}

func pathCheck(path, projectName string) error {
	/// --------------------------------------------
	/// Check whether a folder with given path exists or not.
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return errors.New("The given path does not exist")
	}
	if !stat.IsDir() {
		return errors.New("Given path is not a folder")
	}
	/// ------------------------------------------

	/// ---------------------------------------------
	/// Check whether the project path is duplicated or not.
	stat, err = os.Stat(concatPathWithProjectName(path, projectName))
	if os.IsExist(err) {
		return errors.New("The other project exists with the same name in the path")
	}
	/// ---------------------------------------------

	return nil
}

func createNewProject(path, projectName string) bool {
	projectPath := concatPathWithProjectName(path, projectName)
	var err error = nil

	// Make project directory.
	err = os.Mkdir(projectPath, os.ModeDir)
	if !errorCheck(err) {
		return false
	}

	// Make resource project.
	err = os.Mkdir(projectPath+"/Resource", os.ModeDir)
	if !errorCheck(err) {
		return false
	}

	// Make project settings.
	err = os.Mkdir(projectPath+"/ProjectSettings", os.ModeDir)
	if !errorCheck(err) {
		return false
	}

	// Add the project to the project list.
	if _, err = dbAddProject(projectName, projectPath); !errorCheck(err) {
		return false
	}

	return true
}

func errorCheck(err error) bool {
	if err != nil {
		widgets.QMessageBox_Critical(nil, "Failed to create", err.Error(), widgets.QMessageBox__Close, widgets.QMessageBox__Close)
		return false
	}
	return true
}
