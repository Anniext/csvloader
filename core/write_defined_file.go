package core

import (
	"csvloader/core/types"
	"csvloader/global"
	"csvloader/internal/embed"
	"github.com/CloudyKit/jet/v6"
	"os"
	"path/filepath"
)

func WriteDefinedFile() {
	// 检测genPath是否存在
	genPath := global.Config.GenPath
	if _, err := os.Stat(genPath); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(genPath, 0755) // 0755 是默认权限
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	file, err := os.OpenFile(filepath.Join(global.Config.GenPath, global.Config.DefinedFilePath), os.O_WRONLY|os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Load the template file
	view := global.Config.DefinedJetLoader
	var tmpl *jet.Template

	if global.Config.TemplatePath == "" {
		templateString := embed.OpenTemplateFileToString("csv_defined.jet")
		tmpl, err = view.Parse("template.jet", templateString)
	} else {
		tmpl, err = view.GetTemplate("csv_defined.jet")
	}

	if err != nil {
		panic(err)
	}

	tmplData := &types.ModulePackageName{PackageName: global.Config.PackageName}

	// Execute the template with data and print the result
	if err = tmpl.Execute(file, nil, tmplData); err != nil {
		panic(err)
	}
}
