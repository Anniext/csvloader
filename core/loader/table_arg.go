package loader

import (
	"csvloader/core/types"
	"csvloader/core/utils"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type TableArgs struct {
	TablePath           string
	TableName           string
	TableFlag           string
	TableDesc           string
	CsvName             string
	SpamReader          *csv.Reader
	Cols                int
	Rows                int
	DescList            []string
	FlagList            []string
	TypeList            []string
	KeyList             []string
	TreeIndexList       map[string]*AssistIndex
	HashIndexList       map[string]*AssistIndex
	ExternalLoadDict    map[int]bool
	KeyTypeDict         map[string]string
	UserDefinedTypeList [][]string
	HashIndexPk         []string
	ArgvData            *ArgvData
}

func (t *TableArgs) init() {
	t.SpamReader = nil
	t.DescList = make([]string, 0)
	t.FlagList = make([]string, 0)
	t.TypeList = make([]string, 0)
	t.KeyList = make([]string, 0)
	t.TreeIndexList = make(map[string]*AssistIndex)
	t.HashIndexList = make(map[string]*AssistIndex)
	t.ExternalLoadDict = make(map[int]bool)
	t.KeyTypeDict = make(map[string]string)
	t.UserDefinedTypeList = [][]string{
		{"array_int", "[]int32"},
		{"array_float", "[]float32"},
		{"array_int64", "[]int64"},
		{"array_string", "[]string"},
	}
}

func (t *TableArgs) ParseName(fileName string) {
	t.isCsv()
	items := utils.GetTableInfo(fileName)
	if len(items) < 3 {
		t.TableFlag = ""
		return
	}
	t.CsvName = fileName
	t.TableName = items[0]
	t.TableDesc = items[1]
	t.TableFlag = items[2]
}

func (t *TableArgs) OpenCsvFil() {
	csvFile, err := os.Open(t.TablePath)
	if err != nil {
		fmt.Println("Failed to open file:", err)
		return
	}
	t.SpamReader = csv.NewReader(csvFile)
}
func (t *TableArgs) CloseCsvFile() {
}

func (t *TableArgs) InitArgvDict() {
	args := ArgvData{}
	args.SetTableArgs(t)
	args.SetClass(t.TableName)
	args.SetClassManage(t.TableName)
	args.SetClassManageLower(t.TableName)
	args.SetCsvFileName(t.CsvName)
	args.SetPkVar(t.HashIndexPk)
	args.SetKeyStrList(t.KeyList)
	args.SetClsDotKeyList(t.KeyList)
	args.InitBlockClassFieldLines(t)
	args.InitBlockIndexKeyTypeDef(t)
	args.InitBlockIndexTypeDef(t)
	args.InitBlockIndexInsert(t)
	args.InitBlockGetMethod(t)
	args.InitBlockIndexVar(t)
	args.InitBlockInitMethod(t)
	args.InitBlockUserDefinedType()

	t.ArgvData = &args
}

// Parse 解析csv结构
func (t *TableArgs) Parse() {
	t.init()       // 初始化数据
	t.OpenCsvFil() // 打开csvLoader
	// 解析csv头部
	Rows, err := t.SpamReader.ReadAll()
	if err != nil {
		fmt.Println("Failed to read line from csv: ", err)
		return
	}
	fmt.Println("开始解析csv:", t.CsvName)
	t.ParseCsvHead(Rows)
	t.ParseCsvRow(Rows)
}

func (t *TableArgs) ParseCsvRow(Rows [][]string) {
	for i := 5; i < len(Rows); i++ {
		//row := Rows[i]
		t.Rows += 1
		//t.Dat
	}
}

func (t *TableArgs) ParseCsvHead(Rows [][]string) {
	tagInfo := &TagInfo{}
	tagInfo.init(Rows)
	copy(tagInfo.DescRow, Rows[0])
	copy(tagInfo.FlagRow, Rows[1])
	copy(tagInfo.TypeRow, Rows[2])
	copy(tagInfo.KeyRow, Rows[3])
	copy(tagInfo.TagRow, Rows[4])

	for index, value := range tagInfo.FlagRow {
		if utils.ValidFlag(value) {
			t.Cols += 1
			t.DescList = append(t.DescList, tagInfo.DescRow[index])
			t.FlagList = append(t.FlagList, tagInfo.FlagRow[index])
			t.TypeList = append(t.TypeList, tagInfo.TypeRow[index])
			t.KeyList = append(t.KeyList, tagInfo.KeyRow[index])
			tagInfo.TagDict = utils.ParseTag(tagInfo.TagRow[index])
			t.parseCsvTag(tagInfo, index)
		}
	}

	tmpDict := make(map[string]string)
	for i, key := range t.KeyList {
		if i < len(t.TypeList) {
			tmpDict[key] = types.TypeDict[t.TypeList[i]]
		}
	}
	for key, value := range t.KeyTypeDict {
		tmpDict[key] = value
	}
	t.KeyTypeDict = tmpDict
	t.HashIndexList = tagInfo.AssistHashIndexDict
	t.TreeIndexList = tagInfo.AssistTreeIndexDict

	leftInx := 99999
	keyDict := make([]string, 0)
	for key := range t.HashIndexList {
		keyDict = append(keyDict, key)
	}

	for idx := 0; idx < len(keyDict); idx++ {
		propDict := t.HashIndexList[keyDict[idx]]
		if propDict.Unique == 1 && idx < leftInx {
			leftInx = idx
		}
	}
	if leftInx == 99999 {
		panic("invalid pk:" + t.TableName)
	}
	t.HashIndexPk = make([]string, len(t.HashIndexList[keyDict[leftInx]].Cols))
	t.HashIndexPk = t.HashIndexList[keyDict[leftInx]].Cols
}

// Valid 验证是否有效表格
func (t *TableArgs) Valid() bool {
	return strings.Contains(t.TableFlag, "s")
}

// Correctness 判断是否可以转换
func (t *TableArgs) Correctness() bool {
	tableName := t.CsvName
	if filepath.Ext(tableName) != ".csv" {
		fmt.Println(tableName, ":不是csv文件, 排除!!!!!")
		return false
	}
	return true
}

func (t *TableArgs) isCsv() {
	if filepath.Ext(t.TableName) != ".csv" {
		t.TableFlag = ""
	}
	return
}

func NewTableArgs() *TableArgs {
	return &TableArgs{}
}

func (t *TableArgs) parseCsvTag(tagInfo *TagInfo, index int) {
	if tagInfo.TagDict == nil {
		return
	}

	for key, value := range tagInfo.TagDict {
		var assistIndexDict map[string]*AssistIndex

		if strings.Contains(key, "hash") || strings.Contains(key, "tree") {
			if strings.Contains(key, "hash") {
				assistIndexDict = t.HashIndexList
			} else {
				assistIndexDict = t.TreeIndexList
			}
			indexPDict := utils.ParseIndexValue(value)
			newIndexName := key + indexPDict["group"]
			if _, ok := assistIndexDict[newIndexName]; !ok {
				assistIndexDict[newIndexName] = &AssistIndex{
					Cols:   make([]string, 0),
					Unique: 1,
					Order:  0,
					Auto:   0,
				}
			}
			for pKey, pValue := range indexPDict {
				if pKey == "group" {
					assistIndexDict[newIndexName].Cols = append(assistIndexDict[newIndexName].Cols, tagInfo.KeyRow[index])
				} else {
					pValueInt, err := strconv.Atoi(pValue)
					if err != nil {
						fmt.Println("strconv.Atoi error : ", err)
					}
					switch pKey {
					case "auto":
						assistIndexDict[newIndexName].Auto = pValueInt
					case "order":
						assistIndexDict[newIndexName].Order = pValueInt
					case "unique":
						assistIndexDict[newIndexName].Unique = pValueInt
					}
				}
			}
		} else if strings.Contains(key, "type") {
			userType := value[0]
			if _, ok := types.TypeDict[userType]; ok {
				t.KeyTypeDict[tagInfo.KeyRow[index]] = types.TypeDict[userType]
			} else {
				t.UserDefinedTypeList = append(t.UserDefinedTypeList, []string{t.TableName, userType})
				t.KeyTypeDict[tagInfo.KeyRow[index]] = utils.GoDsToLegalName(userType)
			}
		} else if strings.Contains(key, "external") {
			t.ExternalLoadDict[index] = true
		}

		if strings.Contains(key, "hash") {
			for akey, aValue := range assistIndexDict {
				tagInfo.AssistHashIndexDict[akey] = aValue
			}
		} else {
			for akey, aValue := range assistIndexDict {
				tagInfo.AssistTreeIndexDict[akey] = aValue
			}
		}
	}
}
