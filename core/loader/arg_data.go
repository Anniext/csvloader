package loader

import (
	"bytes"
	"csvloader/core/utils"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"strings"
)

type ArgvData struct {
	TableArgs            *TableArgs
	Class                string
	ClassManage          string
	LowerClassManage     string
	CsvFileName          string
	PkVar                string
	KeyStrList           string
	ClsDotKeyList        string
	BlockClassFieldLines string
	BlockIndexKeyTypeDef string
	BlockIndexTypeDef    string
	BlockIndexInsert     string
	BlockGetMethod       string
	BlockIndexVar        string
	BlockInitMethod      string
	BlockUserDefinedType string
}

func (a *ArgvData) SetTableArgs(t *TableArgs) {
	a.TableArgs = t
}

func (a *ArgvData) SetClass(tableName string) {
	a.Class = utils.SplitTitle(tableName, "_")
}

func (a *ArgvData) SetClassManage(tableName string) {
	tableNameTitle := utils.SplitTitle(tableName, "_")
	classManage := strings.ReplaceAll(tableNameTitle, "_", "") + "Manager"
	a.ClassManage = classManage
}
func (a *ArgvData) SetClassManageLower(tableName string) {
	lowerClassManage := strings.ReplaceAll(strings.ToLower(tableName), "_", "") + "Manager"
	a.LowerClassManage = lowerClassManage
}

func (a *ArgvData) SetCsvFileName(csvFileName string) {
	a.CsvFileName = csvFileName
}

func (a *ArgvData) SetPkVar(hashIndexPk []string) {
	var pk string
	for _, indexPk := range hashIndexPk {
		pk += utils.SplitTitle(indexPk, "_")
	}
	a.PkVar = "hash" + pk
}

func (a *ArgvData) SetKeyStrList(keyList []string) {
	var str string
	for i, v := range keyList {
		str += fmt.Sprintf("%s", v)
		if i != len(keyList)-1 {
			str += ", "
		}
	}
	a.KeyStrList = str
}

func (a *ArgvData) SetClsDotKeyList(keyList []string) {
	var cls string
	for i, v := range keyList {
		cls += fmt.Sprintf("cls.%s", v)
		if i != len(keyList)-1 {
			cls += ", "
		}
	}
	a.ClsDotKeyList = cls
}

func (a *ArgvData) InitBlockClassFieldLines(t *TableArgs) {
	classFieldLines := make([]string, 0)
	for i, key := range t.KeyList {
		line := "    " + utils.SplitTitle(key, "_") + " " + t.KeyTypeDict[key]
		if !t.ExternalLoadDict[i] {
			line += " `csv:\"" + key + "\" property:\"readonly\"`"
		}
		classFieldLines = append(classFieldLines, line)
	}
	a.BlockClassFieldLines = strings.Join(classFieldLines, "\n")
}

func (a *ArgvData) InitBlockIndexKeyTypeDef(t *TableArgs) {
	indexTypeDefList := make([]string, 0)
	for _, indexItems := range t.HashIndexList {
		indexTypeDefList = append(indexTypeDefList, indexItems.GenBlockDetailByKey(a, true))
	}
	for _, indexItems := range t.TreeIndexList {
		indexTypeDefList = append(indexTypeDefList, indexItems.GenBlockDetailByKey(a, false))
	}
	a.BlockIndexKeyTypeDef = strings.Join(indexTypeDefList, "\n")
}

func (a *ArgvData) InitBlockIndexTypeDef(t *TableArgs) {
	indexTypeDefList := make([]string, 0)
	for _, indexItems := range t.HashIndexList {
		indexTypeDefList = append(indexTypeDefList, indexItems.GenBlockDetail(a, true))
	}
	for _, indexItems := range t.TreeIndexList {
		indexTypeDefList = append(indexTypeDefList, indexItems.GenBlockDetail(a, false))
	}
	a.BlockIndexTypeDef = strings.Join(indexTypeDefList, "\n")
}

func (a *ArgvData) InitBlockIndexInsert(t *TableArgs) {
	indexInsertList := make([]string, 0)
	for _, indexItems := range t.HashIndexList {
		indexInsertList = append(indexInsertList, indexItems.GenBlockDetailByInsert(a, true))
	}
	for _, indexItems := range t.TreeIndexList {
		indexInsertList = append(indexInsertList, indexItems.GenBlockDetailByInsert(a, false))
	}
	a.BlockIndexInsert = strings.Join(indexInsertList, "")
}

func (a *ArgvData) InitBlockGetMethod(t *TableArgs) {
	indexMethodList := make([]string, 0)

	for _, indexItems := range t.HashIndexList {
		indexMethodList = append(indexMethodList, indexItems.GenBlockDetailByGet(a, true))
	}
	for _, indexItems := range t.TreeIndexList {
		indexMethodList = append(indexMethodList, indexItems.GenBlockDetailByGet(a, false))
	}

	var buf bytes.Buffer
	tmpl, err := jet.NewSet(jet.NewInMemLoader()).Parse("template.jet", `
func (m *{{.ClassManage}}) GetAll{{.Class}}s() ([]*{{.Class}}) {
    return m.data
}
`)
	err = tmpl.Execute(&buf, nil, map[string]interface{}{
		"ClassManage": a.ClassManage,
		"Class":       a.Class,
	})
	if err != nil {
		panic(err)
	}
	indexMethodList = append(indexMethodList, buf.String())
	a.BlockGetMethod = strings.Join(indexMethodList, "")
}

func (a *ArgvData) InitBlockIndexVar(t *TableArgs) {
	indexVarList := make([]string, 0)

	for _, indexItems := range t.HashIndexList {
		indexVarList = append(indexVarList, indexItems.GenBlockDetailByVar(a, true))
	}
	for _, indexItems := range t.TreeIndexList {
		indexVarList = append(indexVarList, indexItems.GenBlockDetailByVar(a, false))
	}

	var buf bytes.Buffer
	tmpl, err := jet.NewSet(jet.NewInMemLoader()).Parse("template.jet", `
	data []*{{.Class}}
`)
	err = tmpl.Execute(&buf, nil, map[string]interface{}{
		"Class": a.Class,
	})
	if err != nil {
		panic(err)
	}
	indexVarList = append(indexVarList, buf.String())
	a.BlockIndexVar = strings.Join(indexVarList, "")
}

func (a *ArgvData) InitBlockInitMethod(t *TableArgs) {
	indexVarList := make([]string, 0)

	for _, indexItems := range t.HashIndexList {
		indexVarList = append(indexVarList, indexItems.GenBlockDetailByInit(a, true))
	}
	for _, indexItems := range t.TreeIndexList {
		indexVarList = append(indexVarList, indexItems.GenBlockDetailByInit(a, false))
	}

	a.BlockInitMethod = strings.Join(indexVarList, "")
}

func (a *ArgvData) InitBlockUserDefinedType() {
	blockDetailList := make([]string, 0)
	rangeByItem := func(item []string) bool {
		for _, v := range item {
			if strings.Contains(v, "FormulaFunc") {
				return true
			}
		}
		return false
	}

	for _, item := range a.TableArgs.UserDefinedTypeList {
		if rangeByItem(item) {
			continue
		}
		if _, exists := utils.GRepeatDefinedType[item[1]]; exists {
			continue
		}
		utils.GRepeatDefinedType[item[1]] = true
		blockDetailList = append(blockDetailList, GenBlockDetailByUser(item))
	}

	a.BlockUserDefinedType = strings.Join(blockDetailList, "")
}
