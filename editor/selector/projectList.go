package selector

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const projectListObjectName string = "ProjectList"

// ProjectListDialog : A dialog which has a list of projects.
type ProjectListDialog struct {
	qdialog *widgets.QDialog
	onOpen  func(path string)
	onClose func()
}

// NewProjectListWindow : Make a ProjectListWindow instance.
// [parent] is the parent window.
// [onOpen] is a callback which is called when use click the open button of a project.
// [onClose] is a callback which is called when the window closed including after calling of [onOpen].
func NewProjectListWindow(parent *widgets.QMainWindow, onOpen func(path string), onClose func()) *ProjectListDialog {
	window := &ProjectListDialog{}
	window.init(parent)
	window.onOpen = onOpen
	window.onClose = onClose

	return window
}

func (window *ProjectListDialog) init(parent *widgets.QMainWindow) {
	/// -----------------------------------
	/// Dialog - qdialog
	/// -> VBoxLayout - mainLayout
	///    -> Widget - menuWidget
	///       -> HBoxLayout - menuLayout
	///    -> ScrollArea - scrollArea
	///       -> GridLayout - layout
	window.qdialog = widgets.NewQDialog(parent, 0)
	qdialog := window.qdialog
	qdialog.SetModal(true)
	qdialog.SetMinimumSize2(800, 480)
	qdialog.SetWindowTitle("Project Selector")
	qdialog.SetBackgroundRole(gui.QPalette__Dark)
	qdialog.ConnectCloseEvent(func(event *gui.QCloseEvent) {
		window.onClose()
	})

	mainLayout := widgets.NewQVBoxLayout()
	qdialog.SetLayout(mainLayout)

	menuLayout := widgets.NewQHBoxLayout()
	menuWidget := widgets.NewQWidget(nil, 0)
	menuWidget.SetLayout(menuLayout)
	mainLayout.AddWidget(menuWidget, 0, 0)

	buttonAddProject := widgets.NewQPushButton2("Add Project", nil)
	buttonAddProject.ConnectClicked(func(checked bool) {
		dir := widgets.QFileDialog_GetExistingDirectory(nil, "Select a folder which contains a Shake project", "", widgets.QFileDialog__ShowDirsOnly)
		print(dir)
	})
	buttonNewProject := widgets.NewQPushButton2("Create Project", nil)
	buttonNewProject.ConnectClicked(func(checked bool) {
		NewCreateProjectDialog(qdialog, func() {
			window.updateProjectList()
		})
	})
	menuLayout.AddWidget(buttonAddProject, 0, 0)
	menuLayout.AddWidget(buttonNewProject, 0, 0)

	layout := widgets.NewQGridLayout2()
	layout.Widget().SetMinimumSize2(200, 400)
	scrollArea := widgets.NewQScrollArea(nil)
	scrollArea.SetObjectName(projectListObjectName)
	scrollArea.SetBackgroundRole(gui.QPalette__Base)
	scrollArea.SetLayout(layout)
	mainLayout.AddWidget(scrollArea, 0, 0)
	/// -----------------------------------

	window.updateProjectList()

	qdialog.Show()
}

// findProjectListGridLayout : Find the GridLayout object which contains list of projects.
func (window *ProjectListDialog) findProjectListGridLayout() *widgets.QGridLayout {
	scrollAreaObject := window.qdialog.FindChild("ProjectList", core.Qt__FindChildrenRecursively)
	scrollArea := widgets.NewQScrollAreaFromPointer(scrollAreaObject.Pointer())
	layout := scrollArea.Layout()
	gridLayout := widgets.NewQGridLayoutFromPointer(layout.Pointer())

	return gridLayout
}

// updateProjectList : Update project list.
func (window *ProjectListDialog) updateProjectList() {
	window.clearProjectList()
	gridLayout := window.findProjectListGridLayout()

	// Column names - First row
	labelName := widgets.NewQLabel2("Project Name", nil, 0)
	labelPath := widgets.NewQLabel2("Path", nil, 0)
	gridLayout.AddWidget2(labelName, 0, 0, 0)
	gridLayout.AddWidget2(labelPath, 0, 1, 0)
	labelName.SetFixedHeight(30)
	labelPath.SetFixedHeight(30)

	// Get project list from database.
	projectInfoList := dbGetProjectList()

	for i, projectInfo := range projectInfoList {
		labelProjectName := widgets.NewQLabel2(projectInfo.name, nil, 0)
		labelProjectName.SetFixedHeight(30)
		labelProjectPath := widgets.NewQLabel2(projectInfo.path, nil, 0)
		labelProjectPath.SetFixedHeight(30)
		buttonOpen := widgets.NewQPushButton2("Open", nil)
		buttonOpen.SetFixedHeight(30)
		projectPath := projectInfo.path
		buttonOpen.ConnectClicked(func(checked bool) {
			window.onOpen(projectPath)
			window.qdialog.Close()
		})

		gridLayout.AddWidget2(labelProjectName, i+1, 0, 0)
		gridLayout.AddWidget2(labelProjectPath, i+1, 1, 0)
		gridLayout.AddWidget2(buttonOpen, i+1, 2, 0)
	}
}

// clearProjectList : Remove all elements in project list.
func (window *ProjectListDialog) clearProjectList() {
	gridLayout := window.findProjectListGridLayout()
	for i := 0; i < gridLayout.RowCount(); i++ {
		for j := 0; j < gridLayout.ColumnCount(); j++ {
			// Get LayoutItem in (i,j).
			t := gridLayout.ItemAtPosition(i, j)

			// We should remove it from GridLayout and also delete the widget.
			// If we do not delete the widget, it will be shown at unexpected place.
			gridLayout.RemoveItem(t)
			t.Widget().DeleteLater()
		}
	}
}
