package core

import (
	"csvloader/core/loader"
	"csvloader/global"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ParseLoaderMap() {
	global.CsvLoaderMap = loader.NewCsvLoaderMap()

	err := filepath.Walk(global.Config.WorkPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		if !info.IsDir() {
			table := loader.NewTableArgs()
			table.TablePath = path
			table.ParseName(info.Name())
			if !table.Valid() || !table.Correctness() {
				fmt.Println("不可被转换 -------> ", table.TableName)
				return nil
			}
			table.Parse() // 解析argInfo
			tableClass := strings.Split(table.TableName, "_")[0]
			fileLoaderPath := filepath.Join(global.Config.GenPath, "g_"+tableClass+".go")
			if global.CsvLoaderMap[tableClass] == nil {
				global.CsvLoaderMap[tableClass] = loader.NewCsvLoader()
			}
			global.CsvLoaderMap[tableClass].InitConfig(global.Config)
			global.CsvLoaderMap[tableClass].InitLoader(fileLoaderPath)
			global.CsvLoaderMap[tableClass].InitTables()
			table.InitArgvDict()
			global.CsvLoaderMap[tableClass].Append(table)
		} else {
			// TODO 递归文件夹处理
		}
		return nil
	})

	if err != nil {
		fmt.Println("Walk error:", err)
	}
}
