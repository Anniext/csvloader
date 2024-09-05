package core

import (
	"csvloader/core/types"
	"csvloader/global"
	"fmt"
)

func WriteTableFile() {
	embed := global.Config.Embed
	for tableName, tables := range global.CsvLoaderMap {
		fmt.Println("开始解析 csv 文件 --------> ", tableName)
		moduleData := types.ModulePackageName{
			Embed:       embed,
			PackageName: tables.Cnf.PackageName,
		}
		tables.Loader.WriteFileByLoader("table_module_head.jet", &moduleData)
		for _, table := range tables.Table {
			var tableData = types.ModulePackageData{
				Embed:                embed,
				Class:                table.ArgvData.Class,
				ClassManage:          table.ArgvData.ClassManage,
				BlockUserDefinedType: table.ArgvData.BlockUserDefinedType,
				BlockClassFieldLines: table.ArgvData.BlockClassFieldLines,
				BlockIndexKeyTypeDef: table.ArgvData.BlockIndexKeyTypeDef,
				BlockIndexTypeDef:    table.ArgvData.BlockIndexTypeDef,
				BlockIndexVar:        table.ArgvData.BlockIndexVar,
				BlockInitMethod:      table.ArgvData.BlockInitMethod,
				PkVar:                table.ArgvData.PkVar,
				CsvFileName:          table.ArgvData.CsvFileName,
				BlockIndexInsert:     table.ArgvData.BlockIndexInsert,
				BlockGetMethod:       table.ArgvData.BlockGetMethod,
			}
			tables.Loader.WriteFileDataByLoader("table_module_data.jet", &tableData)
		}
	}
}
