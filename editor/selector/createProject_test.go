package selector

import (
	"testing"
)

func Test_concatPathWithProjectName(t *testing.T) {
	concatPathWithProjectNameTest(t, "/", "Test", "/Test")
	concatPathWithProjectNameTest(t, "/Temp", "Te st", "/Temp/Test")
	concatPathWithProjectNameTest(t, "/Temp/ABC", "Tes t", "/Temp/ABC/Test")
	concatPathWithProjectNameTest(t, "/Temp/ABC/", "Test", "/Temp/ABC/Test")
	concatPathWithProjectNameTest(t, "C:\\", "Test", "C:/Test")
	concatPathWithProjectNameTest(t, "C:\\Temp", "T est", "C:/Temp/Test")
	concatPathWithProjectNameTest(t, "C:\\Temp/ABC", " Test", "C:/Temp/ABC/Test")
	concatPathWithProjectNameTest(t, "C:\\Temp/ABC/", "Test ", "C:/Temp/ABC/Test")
}

func concatPathWithProjectNameTest(t *testing.T, path, projectName, answer string) {
	result := concatPathWithProjectName(path, projectName)
	if result != answer {
		t.Errorf("%s concat %s is %s not %s", path, projectName, result, answer)
	}
}
