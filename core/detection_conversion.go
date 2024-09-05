package core

import (
	"bytes"
	"csvloader/core/utils"
	"csvloader/global"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func DetectionConversion() {
	if global.Config.XlsxFilePath == "" {
		return
	}

	if _, err := exec.LookPath("excel2csv"); err != nil {
		panic(err)
	}
	if err := utils.JudgmentCommand("excel2csv"); err != nil {
		panic(err)
	}
	csvFilePath := filepath.Join(filepath.Dir(global.Config.XlsxFilePath), "csv")

	// 不存在文件夹则创建
	if _, err := os.Stat(csvFilePath); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(csvFilePath, 0755) // 0755 是默认权限
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	fmt.Println("-------------------------------------")
	fmt.Println("开始解析,请稍候........")

	// 覆盖更新
	cmd := exec.Command("excel2csv", fmt.Sprintf("-i=%s", global.Config.XlsxFilePath), fmt.Sprintf("-o=%s", csvFilePath))

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(out.String())

	fmt.Print("转换完毕 ----------> 更新csv地址:")
	global.Config.WorkPath = csvFilePath
	fmt.Println(csvFilePath)
}
