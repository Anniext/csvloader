//go:generate go install csvloader/cmd/csvloader
package main

import (
	"csvloader/core"
	"csvloader/core/jet"
	"csvloader/global"
	"csvloader/internal/config"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "csv_loader",
		Short: "A simple, clear csv import software.",
		Run: func(_ *cobra.Command, _ []string) {
			global.Config = instanceConfig
			// 检测转换
			core.DetectionConversion()
			// 创建模板引擎
			jet.NewJetLoader()
			// 写入defined文件
			core.WriteDefinedFile()
			// 处理loader
			core.ParseLoaderMap()
			// 写入表格到文件
			core.WriteTableFile()
			// 写入导入文件
			core.WriteLoadTableFile()
			// 写入公式表
			core.WriteFormulaFile()
			// 转移嵌入资源
			core.EscapeResources()
		},
	}
	instanceConfig  *config.Config
	genPath         string
	filePath        string
	definedFilePath string
	formulaFilePath string
	workPath        string
	csvPath         string
	embed           bool
	packageName     string
	templatePath    string
	xlsxFilePath    string
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&genPath, "genPath", "g", "", "generate file path?")
	rootCmd.PersistentFlags().StringVarP(&filePath, "filePath", "f", "01_csv_table.go", "generate csv_table file name")
	rootCmd.PersistentFlags().StringVarP(&definedFilePath, "definedFilePath", "d", "02_csv_defined.go", "generate csv definedFile file name")
	rootCmd.PersistentFlags().StringVarP(&formulaFilePath, "formulaFilePath", "F", "03_formula.go", "generate formulaFile file name")
	rootCmd.PersistentFlags().StringVarP(&workPath, "workPath", "w", "", "csv file path ")
	rootCmd.PersistentFlags().StringVarP(&csvPath, "csvPath", "c", "", "if embed is turned on, it indicates the resource address.")
	rootCmd.PersistentFlags().StringVarP(&packageName, "packageName", "p", "csv", "package name")
	rootCmd.PersistentFlags().BoolVarP(&embed, "embed", "e", false, "embed resources into binary files if embed is turned on (default false)")
	rootCmd.PersistentFlags().StringVarP(&templatePath, "templatePath", "t", "", "template file path")
	rootCmd.PersistentFlags().StringVarP(&xlsxFilePath, "xlsxFilePath", "x", "", "convert with xlsx")

	err := viper.BindPFlag("genPath", rootCmd.PersistentFlags().Lookup("genPath"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("filePath", rootCmd.PersistentFlags().Lookup("filePath"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("definedFilePath", rootCmd.PersistentFlags().Lookup("definedFilePath"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("formulaFilePath", rootCmd.PersistentFlags().Lookup("formulaFilePath"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("workPath", rootCmd.PersistentFlags().Lookup("workPath"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("csvPath", rootCmd.PersistentFlags().Lookup("csvPath"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("packageName", rootCmd.PersistentFlags().Lookup("packageName"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("templatePath", rootCmd.PersistentFlags().Lookup("templatePath"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("xlsxFilePath", rootCmd.PersistentFlags().Lookup("xlsxFilePath"))
	if err != nil {
		panic(err)
	}
	err = viper.BindPFlag("embed", rootCmd.PersistentFlags().Lookup("embed"))
	if err != nil {
		panic(err)
	}

	viper.SetDefault("genPath", "")
	viper.SetDefault("filePath", "01_csv_table.go")
	viper.SetDefault("definedFilePath", "02_csv_defined.go")
	viper.SetDefault("formulaFilePath", "03_formula.go")
	viper.SetDefault("workPath", "")
	viper.SetDefault("csvPath", "")
	viper.SetDefault("xlsxFilePath", "")
	viper.SetDefault("templatePath", "")
	viper.SetDefault("embed", false)
	viper.SetDefault("packageName", "csv")
	viper.SetEnvPrefix("csv_loader")

}

func initConfig() {
	viper.AutomaticEnv()
	var err error
	instanceConfig, err = config.GetConfig()
	if err != nil {
		fmt.Println("failed to config", err)
	}

	fmt.Printf(`-----------------------------------------
Csv Loader Config
version: %s
gen_path: %s
file_path: %s
defined_file_path: %s
formula_file_path: %s
work_path: %s
embed: %v
package_name: %s
csv_path: %s
templatge_path: %s
xlsx_file_path: %s
`, instanceConfig.Version, instanceConfig.GenPath, instanceConfig.FilePath, instanceConfig.DefinedFilePath, instanceConfig.FormulaFilePath, instanceConfig.WorkPath, instanceConfig.Embed, instanceConfig.PackageName, instanceConfig.CsvPath, instanceConfig.TemplatePath, instanceConfig.XlsxFilePath)

}

func Execute() error {
	return rootCmd.Execute()
}

func main() {
	err := Execute()
	if err != nil {
		panic(err)
	}
}
