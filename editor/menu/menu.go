package menu

import (
	"errors"

	"github.com/therecipe/qt/widgets"
)

var mainWindow *widgets.QMainWindow = nil

var mapCategory map[string]*widgets.QMenu = map[string]*widgets.QMenu{}

// SetMainWindow : Set main window for menubar.
func SetMainWindow(window *widgets.QMainWindow) {
	mainWindow = window
}

// AddCategory : Add a category to menubar of main window.
func AddCategory(category string) error {
	if mainWindow == nil {
		return errors.New("mainWindow is null")
	}

	if _, ok := mapCategory[category]; !ok {
		mapCategory[category] = mainWindow.MenuBar().AddMenu2(category)
	}

	return nil
}

// AddMenu : Add a menu and link a invoke function.
func AddMenu(category, menu string, f func(b bool)) error {
	if qmenu, ok := mapCategory[category]; ok {
		qmenu.AddAction(menu).ConnectTriggered(f)
		return nil
	}
	return errors.New("Cannot find the category")
}

// Clear : Clear all menus and categories.
func Clear() {
	mainWindow.MenuBar().Clear()
}
