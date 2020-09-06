package menu

import (
	"github.com/shakeengine/shake/editor/dock"
	"github.com/shakeengine/shake/misc"
	"github.com/therecipe/qt/widgets"
)

// InitDefaultMenu : Add default menus
func InitDefaultMenu() {
	misc.ErrorCheck(AddCategory("File"))
	misc.ErrorCheck(AddMenu("File", "Open", func(b bool) {
		dir := widgets.QFileDialog_GetExistingDirectory(nil, "OK", "", 0)
		dock.OpenProjectDirectory(dir)
		dock.OpenProjectView()
	}))

	misc.ErrorCheck(AddMenu("File", "Save Scene", func(b bool) {

	}))

	misc.ErrorCheck(AddCategory("View"))
	misc.ErrorCheck(AddMenu("View", "Hierarchy", func(checked bool) {

	}))
	misc.ErrorCheck(AddMenu("View", "Project", func(checked bool) {

	}))
	misc.ErrorCheck(AddMenu("View", "Inspector", func(checked bool) {

	}))
}
