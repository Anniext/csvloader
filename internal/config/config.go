package config

import (
	"csvloader/internal/version"
	"github.com/CloudyKit/jet/v6"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"sync"
)

var (
	config *Config
	once   sync.Once
)

type Config struct {
	GenPath          string   `json:"gen_path"`
	FilePath         string   `json:"file_path"`
	DefinedFilePath  string   `json:"defined_file_path"`
	FormulaFilePath  string   `json:"formula_file_path"`
	WorkPath         string   `json:"work_path"`
	CsvPath          string   `json:"csv_path"`
	Embed            bool     `json:"embed"`
	Version          string   `json:"version"`
	PackageName      string   `json:"package_name"`
	TemplatePath     string   `json:"template_path"`
	XlsxFilePath     string   `json:"xlsx_file_path"`
	DefinedJetLoader *jet.Set `json:""`
	//CsvLoaderMap     map[string]*loader.Tables `json:""`
}

// GetConfig 单例懒汉加载配置文件
func GetConfig() (*Config, error) {
	err := viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	loadConfig := func() {
		goPath := os.Getenv("GOPATH")
		var genPath string
		if config.GenPath == "" {
			genPath = filepath.Join(goPath, "src/server/csv")
			config.GenPath = genPath
		}
		if config.WorkPath == "" {
			config.WorkPath = filepath.Join(goPath, "src/server/bin/csv")
		}
		if config.Embed {
			config.CsvPath = filepath.Join(config.GenPath, "module")
		}
		config.Version = version.GetCurrentVersion()
	}
	once.Do(loadConfig)
	return config, err
}
