package types

type ModulePackageName struct {
	PackageName string
	Embed       bool
}

type ModulePackageTail struct {
	Embed             bool
	BlockTablesInit   string
	BlockTablesLoad   string
	BlockTablesUnload string
	BlockTablesReload string
	BlockTablesDefine string
}

type FormulaPackageData struct {
	BlockInitFormula    string
	BlockFormulaDefine  string
	BlockFormulaDefine1 string
}

type ModulePackageData struct {
	Embed                bool
	Class                string
	ClassManage          string
	BlockUserDefinedType string
	BlockClassFieldLines string
	BlockIndexKeyTypeDef string
	BlockIndexTypeDef    string
	BlockIndexVar        string
	BlockInitMethod      string
	PkVar                string
	CsvFileName          string
	BlockIndexInsert     string
	BlockGetMethod       string
}
