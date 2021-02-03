package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func initFolders() {
	runnerLog("InitFolders")

	path := tmpPath()

	runnerLog("mkdir %s", path)

	if _, errDir := os.Stat(path); os.IsNotExist(errDir) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			runnerLog(err.Error())
		}
	}
}

func isTmpDir(path string) bool {
	absolutePath, _ := filepath.Abs(path)
	absoluteTmpPath, _ := filepath.Abs(tmpPath())

	return absolutePath == absoluteTmpPath
}

func isIgnoredFolder(path string) bool {
	paths := strings.Split(path, "/")
	if len(paths) <= 0 {
		return false
	}

	for _, e := range strings.Split(settings["ignored"], ",") {
		ignoredPaths := strings.Split(strings.TrimSpace(e), "/")
		if len(ignoredPaths) <= len(paths) {
			i := 0
			for ; i < len(ignoredPaths); i++ {
				if paths[i] != ignoredPaths[i] {
					break
				}
			}
			if i == len(ignoredPaths) {
				return true
			}
		}
	}
	return false
}

func isWatchedFile(path string) bool {
	absolutePath, _ := filepath.Abs(path)
	absoluteTmpPath, _ := filepath.Abs(tmpPath())
	fmt.Println(absolutePath, absoluteTmpPath)
	if strings.HasPrefix(absolutePath, absoluteTmpPath+string(filepath.Separator)) {
		return false
	}

	ext := filepath.Ext(path)

	for _, e := range strings.Split(settings["valid_ext"], ",") {
		if strings.TrimSpace(e) == ext {
			return true
		}
	}

	return false
}

func shouldRebuild(eventName string) bool {
	for _, e := range strings.Split(settings["no_rebuild_ext"], ",") {
		e = strings.TrimSpace(e)
		fileName := strings.ReplaceAll(strings.Split(eventName, ":")[0], `"`, "")
		if strings.HasSuffix(fileName, e) {
			return false
		}
	}

	return true
}

func createBuildErrorsLog(message string) bool {
	file, err := os.Create(buildErrorsFilePath())
	if err != nil {
		return false
	}

	_, err = file.WriteString(message)
	if err != nil {
		return false
	}

	return true
}

func removeBuildErrorsLog() error {
	err := os.Remove(buildErrorsFilePath())
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
