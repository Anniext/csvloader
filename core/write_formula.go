package core

import (
	"csvloader/core/types"
	"csvloader/core/utils"
	"csvloader/global"
	"csvloader/internal/embed"
	"encoding/csv"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type FormulaBlock struct {
	BlockInitFormula    string
	BlockFormulaDefine  string
	BlockFormulaDefine1 string
}

func (f *FormulaBlock) SetBlockInitFormula(fId []int) {
	blockInitFormula := make([]string, 0)
	for _, id := range fId {
		template := fmt.Sprintf("    formulaManager.GetFormula_FormulaById(%d).Formulastring = FormulaMap_%d\n", id, id)
		blockInitFormula = append(blockInitFormula, template)
	}
	f.BlockInitFormula = strings.Join(blockInitFormula, "")
}

func (f *FormulaBlock) SetBlockFormulaDefine(fId []int, formula []string) {
	blockFormulaDefine := make([]string, 0)
	for index, value := range fId {
		template := fmt.Sprintf(`
func Formula_%d(PA, PB, PC, PD, PE, PF, PG, PH, PI, PJ,PK,PL,PM,PN,PO,PP,PQ,PR,PS,PT float64) float64 { 
    return %s
}
`, value, utils.FormulaCsvToGo(formula[index]))
		blockFormulaDefine = append(blockFormulaDefine, template)
	}
	f.BlockFormulaDefine = strings.Join(blockFormulaDefine, "")
}

func (f *FormulaBlock) SetBlockFormulaDefine1(fId []int, formula []string) {
	blockFormulaDefine1 := make([]string, 0)
	for index, value := range fId {
		template := fmt.Sprintf(`
func FormulaMap_%d(pMp map[int32]float64) float64 {
    PA := pMp[0]
    PB := pMp[1]
    PC := pMp[2]
    PD := pMp[3]
    PE := pMp[4]
    PF := pMp[5]
    PG := pMp[6]
    PH := pMp[7]
    PI := pMp[8]
    PJ := pMp[9]
    PK := pMp[10]
    PL := pMp[11]
    PM := pMp[12]
    PN := pMp[13]
    PO  := pMp[14]
    PP  := pMp[15]
    PQ  := pMp[16]
    PR  := pMp[17]
    PS  := pMp[18]
    PT  := pMp[19]
    Formula_%d(PA, PB, PC, PD, PE, PF, PG, PH, PI, PJ,PK,PL,PM,PN,PO,PP,PQ,PR,PS,PT)
    return %s
}
`, value, value, utils.FormulaCsvToGo(formula[index]))
		blockFormulaDefine1 = append(blockFormulaDefine1, template)
	}
	f.BlockFormulaDefine1 = strings.Join(blockFormulaDefine1, "")
}

func WriteFormulaFile() {
	if file, err := os.OpenFile(filepath.Join(global.Config.GenPath, global.Config.FormulaFilePath), os.O_WRONLY|os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_APPEND, 0644); err != nil {
		panic(err)
	} else {
		defer file.Close()

		// Load the template file
		view := global.Config.DefinedJetLoader

		csvFile, err := os.Open(global.Config.WorkPath + "/Formula_Formula-公式表-cs.csv")
		if err != nil {
			panic(err)
		}

		defer csvFile.Close()

		rows, err := csv.NewReader(csvFile).ReadAll()

		fIdList := make([]int, 0)
		formulaList := make([]string, 0)

		for col, row := range rows {
			if col <= 4 {
				continue
			}
			idInt, err := strconv.Atoi(row[0])
			if err != nil {
				panic(err)
			}
			fIdList = append(fIdList, idInt)
			formulaList = append(formulaList, row[1])
		}

		formulaBlock := &FormulaBlock{}
		formulaBlock.SetBlockInitFormula(fIdList)
		formulaBlock.SetBlockFormulaDefine(fIdList, formulaList)
		formulaBlock.SetBlockFormulaDefine1(fIdList, formulaList)

		tmplData := &types.FormulaPackageData{
			BlockInitFormula:    formulaBlock.BlockInitFormula,
			BlockFormulaDefine:  formulaBlock.BlockFormulaDefine,
			BlockFormulaDefine1: formulaBlock.BlockFormulaDefine1,
		}
		tmplName := &types.ModulePackageName{
			PackageName: global.Config.PackageName,
		}

		var tmplN *jet.Template
		var tmplD *jet.Template

		if global.Config.TemplatePath == "" {
			templateNStr := embed.OpenTemplateFileToString("formula_head.jet")
			templateDStr := embed.OpenTemplateFileToString("formula_tail.jet")
			tmplN, err = view.Parse("template.jet", templateNStr)
			tmplD, err = view.Parse("template.jet", templateDStr)
		} else {
			tmplN, err = view.GetTemplate("formula_head.jet")
			tmplD, err = view.GetTemplate("formula_tail.jet")
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
