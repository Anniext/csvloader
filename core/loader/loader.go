package loader

import (
	"csvloader/core/types"
	"csvloader/internal/config"
	"csvloader/internal/embed"
	"github.com/CloudyKit/jet/v6"
	"os"
)

type ILoader interface {
	SetJetLoader(templatePath string)
	SetFileLoader(fileLoaderPath string)
	SetConfig(config *config.Config)
	WriteFileByLoader(filePath string, tmplData *types.ModulePackageName)
	WriteFileDataByLoader(filePath string, tmplData *types.ModulePackageData)
}

type Loader struct {
	cnf  *config.Config
	file *os.File
	jet  *jet.Set
}

func (t *Loader) SetJetLoader(templatePath string) {
	var jetLoader jet.Loader
	if t.cnf.TemplatePath == "" {
		jetLoader = jet.NewInMemLoader()
	} else {
		jetLoader = jet.NewOSFileSystemLoader(templatePath)
	}
	if jetLoader == nil {
		panic("failed to jet loader:")
	}
	t.jet = jet.NewSet(jetLoader, jet.WithSafeWriter(nil))
}
func (t *Loader) SetFileLoader(fileLoaderPath string) {
	file, err := os.OpenFile(fileLoaderPath, os.O_WRONLY|os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	t.file = file
}

func (t *Loader) WriteFileByLoader(filePath string, tmplData *types.ModulePackageName) {
	var tmpl *jet.Template
	var err error

	if t.cnf.TemplatePath == "" {
		templateStr := embed.OpenTemplateFileToString(filePath)
		tmpl, err = t.jet.Parse("template.jet", templateStr)
	} else {
		tmpl, err = t.jet.GetTemplate(filePath)
	}

	if err != nil {
		panic(err)
	}

	if err = tmpl.Execute(t.file, nil, tmplData); err != nil {
		panic(err)
	}
}
func (t *Loader) WriteFileDataByLoader(filePath string, tmplData *types.ModulePackageData) {
	var tmpl *jet.Template
	var err error

	if t.cnf.TemplatePath == "" {
		templateStr := embed.OpenTemplateFileToString(filePath)
		tmpl, err = t.jet.Parse("template.jet", templateStr)
	} else {
		tmpl, err = t.jet.GetTemplate(filePath)
	}

	if err != nil {
		panic(err)
	}

	if err = tmpl.Execute(t.file, nil, tmplData); err != nil {
		panic(err)
	}
}
func (t *Loader) SetConfig(config *config.Config) {
	t.cnf = config
}

func NewLoader() *Loader {
	return &Loader{}
}
