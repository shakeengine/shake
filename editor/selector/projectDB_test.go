package selector

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_ProjectDB(t *testing.T) {
	var projectList []*projectInfo

	projectList = dbGetProjectList()
	beforeCount := len(projectList)

	projectID, err := dbAddProject("DB Test Project", "/Temporary/Project/Path/")
	if err != nil {
		t.Errorf("There is an error in adding. %s", err.Error())
	}

	projectList = dbGetProjectList()
	addedCount := len(projectList)
	if beforeCount+1 != addedCount {
		t.Errorf("Count is mismatched after adding. beforeCount : %d, addedCount : %d", beforeCount, addedCount)
	}

	affected := dbRemoveProject(projectID)
	if affected != 1 {
		t.Errorf("Affected count is not one in removing. affected : %d", affected)
	}

	projectList = dbGetProjectList()
	removedCount := len(projectList)
	if beforeCount != removedCount {
		t.Errorf("Count is mismatched after removing. beforeCount : %d, removedCount : %d", beforeCount, removedCount)
	}
}
