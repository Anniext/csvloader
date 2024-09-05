package loader

import (
	"csvloader/internal/config"
)

type Tables struct {
	Table  []*TableArgs
	Loader ILoader
	Cnf    *config.Config
}

func (t *Tables) InitLoader(fileLoaderPath string) {
	t.Loader = NewLoader()
	t.Loader.SetConfig(t.Cnf)
	t.Loader.SetJetLoader(t.Cnf.TemplatePath)
	t.Loader.SetFileLoader(fileLoaderPath)
}

func (t *Tables) InitTables() {
	if t.Table == nil {
		t.Table = make([]*TableArgs, 0)
	}
}

func (t *Tables) InitConfig(config *config.Config) {
	t.Cnf = config
}

func (t *Tables) Append(args *TableArgs) {
	t.Table = append(t.Table, args)
}

func NewCsvLoaderMap() map[string]*Tables {
	return make(map[string]*Tables)
}
func NewCsvLoader() *Tables {
	return &Tables{}
}
