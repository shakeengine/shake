package selector

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"

	"github.com/shakeengine/shake/misc"
)

var tableName string = "project"

func dbGetDB() *sql.DB {
	db, err := sql.Open("sqlite3", "file:selector.data")
	misc.ErrorCheck(err)

	/// --------------------------------------
	/// If there is not a table named [tableName], create the table.
	res, err := db.Query(fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", tableName))
	misc.ErrorCheck(err)
	defer res.Close()
	if !res.Next() {
		sql := fmt.Sprintf("create table %s (id integer not null primary key, name text, path text);", tableName)
		_, err = db.Exec(sql)
		misc.ErrorCheck(err)
	}
	/// ---------------------------------------

	return db
}

// dbAddProject : Add a new project whose name is [projectName] and path is [projectPath].
// It returns 0 and an error if project which is duplicated by [projectPath] exists.
// Otherwise, it returns the projectID and nil.
func dbAddProject(projectName, projectPath string) (int, error) {
	db := dbGetDB()
	defer db.Close()

	/// ----------------------------------
	/// Check whether the same project exists.
	sameRes, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE path='%s';", tableName, projectPath))
	misc.ErrorCheck(err)
	if sameRes.Next() {
		sameRes.Close()
		return 0, errors.New("The given project path was added")
	}
	sameRes.Close()
	/// ------------------------------------

	/// -------------------------------------
	/// Creating ProjectID which is not duplicated.
	projectID := rand.Int()
	dupRes, err := db.Query(fmt.Sprintf("SELECT * FROM %s WHERE id=%d;", tableName, projectID))
	misc.ErrorCheck(err)
	for dupRes.Next() {
		dupRes.Close()
		projectID = rand.Int()
		dupRes, err = db.Query(fmt.Sprintf("SELECT * FROM %s WHERE id=%d;", tableName, projectID))
		misc.ErrorCheck(err)
	}
	dupRes.Close()
	/// -------------------------------------

	// Insert the project.
	res, err := db.Exec(fmt.Sprintf("INSERT INTO %s VALUES(%d, '%s', '%s');", tableName, projectID, projectName, projectPath))
	misc.ErrorCheck(err)
	if count, err := res.RowsAffected(); count <= 0 || err != nil {
		misc.ErrorCheck(err)
		return 0, errors.New("Nothing was added")
	}

	return projectID, nil
}

// dbRemoveProject : Delete the project with [projectID].
// It returns the number of rows which was affected.
func dbRemoveProject(projectID int) int64 {
	db := dbGetDB()
	defer db.Close()

	res, err := db.Exec(fmt.Sprintf("DELETE FROM %s WHERE id=%d;", tableName, projectID))
	misc.ErrorCheck(err)
	count, err := res.RowsAffected()
	misc.ErrorCheck(err)

	return count
}

type projectInfo struct {
	id   int
	name string
	path string
}

// dbGetProjectList : Get a list of projects.
func dbGetProjectList() []*projectInfo {
	db := dbGetDB()
	defer db.Close()

	result := []*projectInfo{}
	res, err := db.Query(fmt.Sprintf("SELECT * FROM %s;", tableName))
	misc.ErrorCheck(err)
	defer res.Close()
	for res.Next() {
		info := &projectInfo{}
		res.Scan(&info.id, &info.name, &info.path)
		result = append(result, info)
	}

	return result
}
