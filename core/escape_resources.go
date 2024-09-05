package core

import (
	"csvloader/global"
	"log"
	"os"
	"os/exec"
)

func EscapeResources() {
	if !global.Config.Embed || global.Config.CsvPath == "" {
		return
	}

	if exists, err := dirExists(global.Config.CsvPath); err != nil {
		log.Fatal(err)
	} else if exists {
		if err := removeDir(global.Config.CsvPath); err != nil {
			log.Fatal(err)
		}
		if err := copyDir(global.Config.WorkPath, global.Config.CsvPath); err != nil {
			log.Fatalf("An error occurred while moving the folder: %v", err)
		}
	} else {
		if err := copyDir(global.Config.WorkPath, global.Config.CsvPath); err != nil {
			log.Fatalf("An error occurred while moving the folder: %v", err)
		}
	}
}

// dirExists 检查目录是否存在
func dirExists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// removeDir 删除目录
func removeDir(name string) error {
	return os.RemoveAll(name)
}

// copyDir 复制目录
func copyDir(src, dst string) error {
	cmd := exec.Command("cp", "-r", src, dst)
	return cmd.Run()
}
