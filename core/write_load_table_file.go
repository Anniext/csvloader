package core

import (
	"csvloader/core/types"
	"csvloader/global"
	"csvloader/internal/embed"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"os"
	"path/filepath"
	"strings"
)

func WriteLoadTableFile() {
	if file, err := os.OpenFile(filepath.Join(global.Config.GenPath, global.Config.FilePath), os.O_WRONLY|os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644); err != nil {
		panic(err)
	} else {
		defer file.Close()

		// Load the template file
		view := global.Config.DefinedJetLoader

		tmplData := &types.ModulePackageTail{}

		BlockTablesInit := make([]string, 0)
		BlockTablesLoad := make([]string, 0)
		BlockTablesUnload := make([]string, 0)
		BlockTablesReload := make([]string, 0)
		BlockTablesDefine := make([]string, 0)

		for _, tables := range global.CsvLoaderMap {
			for _, table := range tables.Table {
				var tablesInitBuilder strings.Builder
				var tablesLoadBuilder strings.Builder
				var tablesUnloadBuilder strings.Builder
				var tablesReloadBuilder strings.Builder
				var tablesDefineBuilder strings.Builder

				tablesInitBuilder.WriteString(fmt.Sprintf(`
	G%s = &%s{}
	G%s.Init()
`, table.ArgvData.ClassManage, table.ArgvData.ClassManage, table.ArgvData.ClassManage))

				tablesLoadBuilder.WriteString(fmt.Sprintf(`
    if err = G%s.Load(); err != nil { panic(err);return err }
`, table.ArgvData.ClassManage))

				tablesUnloadBuilder.WriteString(fmt.Sprintf(`
    G%s.Unload()
`, table.ArgvData.ClassManage))

				tablesReloadBuilder.WriteString(fmt.Sprintf(`
	tmpG%s := &%s{}
    tmpG%s.Init()
    if err = tmpG%s.Load(); err != nil {
        fmt.Println(err)
        return
    }
`, table.ArgvData.ClassManage, table.ArgvData.ClassManage, table.ArgvData.ClassManage, table.ArgvData.ClassManage))

				tablesDefineBuilder.WriteString(fmt.Sprintf(`
var G%s *%s
`, table.ArgvData.ClassManage, table.ArgvData.ClassManage))

				BlockTablesInit = append(BlockTablesInit, tablesInitBuilder.String())
				BlockTablesLoad = append(BlockTablesLoad, tablesLoadBuilder.String())
				BlockTablesUnload = append(BlockTablesUnload, tablesUnloadBuilder.String())
				BlockTablesReload = append(BlockTablesReload, tablesReloadBuilder.String())
				BlockTablesDefine = append(BlockTablesDefine, tablesDefineBuilder.String())
			}

			for _, table := range tables.Table {
				var tablesReloadBuilder strings.Builder
				tablesReloadBuilder.WriteString(fmt.Sprintf(`
	G%s = tmpG%s
`, table.ArgvData.ClassManage, table.ArgvData.ClassManage))

				BlockTablesReload = append(BlockTablesReload, tablesReloadBuilder.String())
			}
		}

		tmplData.Embed = global.Config.Embed
		tmplData.BlockTablesInit = strings.Join(BlockTablesInit, "\n")
		tmplData.BlockTablesLoad = strings.Join(BlockTablesLoad, "\n")
		tmplData.BlockTablesUnload = strings.Join(BlockTablesUnload, "\n")
		tmplData.BlockTablesReload = strings.Join(BlockTablesReload, "\n")
		tmplData.BlockTablesDefine = strings.Join(BlockTablesDefine, "\n")

		tmplName := types.ModulePackageName{
			Embed:       global.Config.Embed,
			PackageName: global.Config.PackageName,
		}

		var tmplN *jet.Template
		var tmplD *jet.Template
		var err error

		if global.Config.TemplatePath == "" {
			templateNStr := embed.OpenTemplateFileToString("table_module.jet")
			templateDStr := embed.OpenTemplateFileToString("table_module_tail.jet")
			tmplN, err = view.Parse("template.jet", templateNStr)
			tmplD, err = view.Parse("template.jet", templateDStr)
		} else {
			tmplN, err = view.GetTemplate("table_module.jet")
			tmplD, err = view.GetTemplate("table_module_tail.jet")
		}
		if err != nil {
			panic(err)
		}

		// Execute the template with data and print the result
		if err := tmplN.Execute(file, nil, tmplName); err != nil {
			panic(err)
		}

		if err := tmplD.Execute(file, nil, tmplData); err != nil {
			panic(err)
		}
	}
}
