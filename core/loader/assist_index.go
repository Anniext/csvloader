package loader

import (
	"bytes"
	"csvloader/core/utils"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"strings"
)

type AssistIndex struct {
	Cols   []string
	Unique int
	Order  int
	Auto   int
}

func (a *AssistIndex) GenBlockDetailByKey(ar *ArgvData, isHash bool) string {
	indexKind := "Hash"
	if !isHash {
		indexKind = "Tree"
	}
	if len(a.Cols) == 1 {
		key := a.Cols[0]
		return fmt.Sprintf("type %sKeyType%s%s %s", ar.Class, indexKind, key, ar.TableArgs.KeyTypeDict[key])
	}
	if len(a.Cols) > 1 {
		keyStr := strings.Join(a.Cols, "")
		keyTypeStr := fmt.Sprintf("%s%s", ar.Class, keyStr)
		keyTypeDefStr := fmt.Sprintf("type %sKeyType%s%s %s", ar.Class, indexKind, keyStr, keyTypeStr)
		if isHash {
			var keyFieldLines []string
			for i, _ := range a.Cols {
				keyFieldLines = append(keyFieldLines, fmt.Sprintf("    %s %s\n", utils.SplitTitle(a.Cols[i], "_"), ar.TableArgs.KeyTypeDict[a.Cols[i]]))
			}
			tupleKeyDef := fmt.Sprintf(`
type %s%s struct {
%s
}
`, ar.Class, keyStr, strings.Join(keyFieldLines, ""))
			return tupleKeyDef + keyTypeDefStr
		}

		// Tree index only allows one key
		panic(fmt.Sprintf("tree index takes exactly 1 key (%d given)", len(a.Cols)))
	}
	panic("len(index) == 0?")
}

func (a *AssistIndex) GenBlockDetail(ar *ArgvData, isHash bool) string {
	indexKind := "Hash"
	if !isHash {
		indexKind = "Tree"
	}
	classSet := utils.GetClassSetStr(ar.Class, a.Unique)
	keyStr := strings.Join(a.Cols, "")

	var err error
	var tmpl *jet.Template

	views := jet.NewSet(jet.NewInMemLoader())
	if isHash {
		tmpl, err = views.Parse("template.jet", `type {{.Class}}{{.indexKind}}{{.keyStr}} map[{{.Class}}KeyType{{.indexKind}}{{.keyStr}}]{{.classSet}}`)
	} else {
		tmpl, err = views.Parse("template.jet", `type {{.Class}}{{.indexKind}}{{.keyStr}} *redblacktree.Tree`)
	}
	if err != nil {
		panic(err)
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, nil, map[string]interface{}{
		"indexKind": indexKind,
		"keyStr":    keyStr,
		"Class":     ar.Class,
		"classSet":  classSet,
	})

	if err != nil {
		panic(err)
	}

	return buf.String()
}

func (a *AssistIndex) GenBlockDetailByInsert(ar *ArgvData, isHash bool) string {
	keyStr := strings.Join(a.Cols, "")

	keyList := make([]string, len(a.Cols))
	for i, str := range a.Cols {
		keyList[i] = utils.SplitTitle(str, "_")
	}
	keyUnderLineStr := strings.Join(keyList, "")

	var dotKeyList string
	for i, v := range a.Cols {
		dotKeyList += fmt.Sprintf("cls.%s", utils.SplitTitle(v, "_"))
		if i != len(keyList)-1 {
			dotKeyList += ", "
		}
	}

	createKey := fmt.Sprintf("{%s}", dotKeyList)
	if len(a.Cols) == 1 {
		createKey = fmt.Sprintf("(cls.%s)", utils.SplitTitle(a.Cols[0], "_"))
	}

	var buf bytes.Buffer
	var err error
	var tmpl *jet.Template

	views := jet.NewSet(jet.NewInMemLoader())

	if isHash {
		if a.Unique == 1 {
			tmpl, err = views.Parse("template.jet", `
		m.hash{{.keyUnderlineStr}}[{{.Class}}KeyTypeHash{{.keyStr}}{{.createKey}}] = cls`)
			err = tmpl.Execute(&buf, nil, map[string]interface{}{
				"createKey":       createKey,
				"keyUnderlineStr": keyUnderLineStr,
				"Class":           ar.Class,
				"keyStr":          keyStr,
			})
			if err != nil {
				panic(err)
			}
		} else {
			tmpl, err = views.Parse("template.jet", `
		tmp{{.Class}}KeyTypeHash{{.keyStr}} := {{.Class}}KeyTypeHash{{.keyStr}}{{.createKey}}
		m.hash{{.keyUnderlineStr}}[tmp{{.Class}}KeyTypeHash{{.keyStr}}] = append(m.hash{{.keyUnderlineStr}}[tmp{{.Class}}KeyTypeHash{{.keyStr}}], cls)`)
			err = tmpl.Execute(&buf, nil, map[string]interface{}{
				"Class":           ar.Class,
				"keyStr":          keyStr,
				"createKey":       createKey,
				"keyUnderlineStr": keyUnderLineStr,
			})
			if err != nil {
				panic(err)
			}
		}
	} else {
		if a.Unique != 1 {
			panic("nonsupport non-uniqueness tree index")
		}
		tmpl, err = views.Parse("template.jet", `(*m.tree{{.keyUnderlineStr}}).Put(cls.{{.key}}, cls)`)
		err = tmpl.Execute(&buf, nil, map[string]interface{}{
			"keyUnderlineStr": keyUnderLineStr,
			"key":             utils.SplitTitle(a.Cols[0], "_"),
		})
		if err != nil {
			panic(err)
		}
	}
	return buf.String()
}

func (a *AssistIndex) GenBlockDetailByGet(ar *ArgvData, isHash bool) string {
	keyStr := strings.Join(a.Cols, "")

	keyList := make([]string, len(a.Cols))
	for i, str := range a.Cols {
		keyList[i] = utils.SplitTitle(str, "_")
	}
	keyUnderLineStr := strings.Join(keyList, "")

	var dotKeyList string
	for i, v := range a.Cols {
		dotKeyList += fmt.Sprintf("%s", v)
		if i != len(keyList)-1 {
			dotKeyList += ", "
		}
	}

	createKey := fmt.Sprintf("{%s}", dotKeyList)
	if len(a.Cols) == 1 {
		createKey = fmt.Sprintf("(%s)", a.Cols[0])
	}

	classSet := utils.GetClassSetStr(ar.Class, a.Unique)
	many := ""
	if a.Unique != 1 {
		many = "s"
	}

	parameterItems := make([]string, 0)
	formalParameterItems := make([]string, 0)
	for i := 0; i < len(a.Cols); i++ {
		parameterItems = append(parameterItems, a.Cols[i]+" "+ar.TableArgs.KeyTypeDict[a.Cols[i]])
		formalParameterItems = append(formalParameterItems, a.Cols[i]+" "+ar.TableArgs.KeyTypeDict[a.Cols[i]])
	}

	var buf bytes.Buffer
	var err error
	var tmpl *jet.Template

	views := jet.NewSet(jet.NewInMemLoader())

	if isHash {
		formalParameterList := strings.Join(parameterItems, ", ")
		actualParameterList := fmt.Sprintf("%sKeyTypeHash%s%s", ar.Class, keyStr, createKey)
		tmpl, err = views.Parse("template.jet", `
func (m *{{.ClassManage}}) Get{{.Class}}{{.many}}By{{.keyStr}}({{.formalParameterList}}) ({{.classSet}}) {
    if data, ok := m.hash{{.keyUnderlineStr}}[{{.actualParameterList}}]; ok {
        return data
    }
    return nil
}
`)
		err = tmpl.Execute(&buf, nil, map[string]interface{}{
			"ClassManage":         ar.ClassManage,
			"Class":               ar.Class,
			"many":                many,
			"keyStr":              keyStr,
			"formalParameterList": formalParameterList,
			"classSet":            classSet,
			"keyUnderlineStr":     keyUnderLineStr,
			"actualParameterList": actualParameterList,
		})
		if err != nil {
			panic(err)
		}
	} else {
		if len(a.Cols) > 1 {
			panic("nonsupport non-uniqueness tree index")
		}

		formalParameterList := strings.Join(parameterItems, ", ")

		tmpl, err = views.Parse("template.jet", `
func (m *{{.ClassManage}}) Get{{.Class}}{{.many}}By{{.keyStr}}({{.formalParameterList}}) ({{.classSet}}) {
    if data, found := (*m.tree{{.keyUnderlineStr}}).Get({{.key}}); found {
        return data.(*{{.Class}})
    }
    return nil
}

func (m *{{.ClassManage}}) Get{{.Class}}LeftBy{{.keyStr}}() ({{.classSet}}) {
    if node := (*m.tree{{.keyUnderlineStr}}).Left(); node != nil {
        return node.Value.(*{{.Class}})
    }
    return nil
}

func (m *{{.ClassManage}}) Get{{.Class}}RightBy{{.keyStr}}() ({{.classSet}}) {
    if node := (*m.tree{{.keyUnderlineStr}}).Right(); node != nil {
        return node.Value.(*{{.Class}})
    }
    return nil
}


func (m *{{.ClassManage}}) Get{{.Class}}FloorBy{{.keyStr}}({{.formalParameterList}}) ({{.classSet}}) {
    if node, found := (*m.tree{{.keyUnderlineStr}}).Floor({{.key}}); found {
        return node.Value.(*{{.Class}})
    }
    return nil
}

func (m *{{.ClassManage}}) Get{{.Class}}CeilingBy{{.keyStr}}({{.formalParameterList}}) ({{.classSet}}) {
    if node, found := (*m.tree{{.keyUnderlineStr}}).Ceiling({{.key}}); found {
        return node.Value.(*{{.Class}})
    }
    return nil
}
`)
		err = tmpl.Execute(&buf, nil, map[string]interface{}{
			"ClassManage":         ar.ClassManage,
			"Class":               ar.Class,
			"many":                many,
			"keyStr":              keyStr,
			"formalParameterList": formalParameterList,
			"classSet":            classSet,
			"keyUnderlineStr":     keyUnderLineStr,
			"key":                 a.Cols[0],
		})
		if err != nil {
			panic(err)
		}
	}

	return buf.String()
}

func (a *AssistIndex) GenBlockDetailByVar(ar *ArgvData, isHash bool) string {
	keyStr := strings.Join(a.Cols, "")

	keyList := make([]string, len(a.Cols))
	for i, str := range a.Cols {
		keyList[i] = utils.SplitTitle(str, "_")
	}
	keyUnderLineStr := strings.Join(keyList, "")

	var buf bytes.Buffer
	var err error
	var tmpl *jet.Set
	var views *jet.Template

	tmpl = jet.NewSet(jet.NewInMemLoader())

	if isHash {
		views, err = tmpl.Parse("template.jet", `
    hash{{.keyUnderlineStr}} {{.Class}}Hash{{.keyStr}}
`)
	} else {
		views, err = tmpl.Parse("template.jet", `
    tree{{.keyUnderlineStr}} {{.Class}}Tree{{keyStr}}
`)
	}
	err = views.Execute(&buf, nil, map[string]interface{}{
		"keyUnderlineStr": keyUnderLineStr,
		"Class":           ar.Class,
		"keyStr":          keyStr,
	})
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func (a *AssistIndex) GenBlockDetailByInit(ar *ArgvData, isHash bool) string {
	keyStr := strings.Join(a.Cols, "")

	keyList := make([]string, len(a.Cols))
	for i, str := range a.Cols {
		keyList[i] = utils.SplitTitle(str, "_")
	}
	keyUnderLineStr := strings.Join(keyList, "")

	var buf bytes.Buffer

	if isHash {
		tmpl, err := jet.NewSet(jet.NewInMemLoader()).Parse("template.jet", `
    m.hash{{.keyUnderlineStr}} = make({{.Class}}Hash{{.keyStr}})
`)
		err = tmpl.Execute(&buf, nil, map[string]interface{}{
			"keyUnderlineStr": keyUnderLineStr,
			"Class":           ar.Class,
			"keyStr":          keyStr,
		})
		if err != nil {
			panic(err)
		}
	} else {
		if len(a.Cols) > 1 {
			panic("nonsupport composite keys tree index")
		}
		keyTypeTitle := utils.SplitTitle(ar.TableArgs.KeyTypeDict[a.Cols[0]], "_")
		tmpl, err := jet.NewSet(jet.NewInMemLoader()).Parse("template.jet", `
    m.tree{{.keyUnderlineStr}} = redblacktree.NewWith(utils.{{.keyTypeTitle}}Comparator)
`)
		err = tmpl.Execute(&buf, nil, map[string]interface{}{
			"keyUnderlineStr": keyUnderLineStr,
			"keyTypeTitle":    keyTypeTitle,
		})
		if err != nil {
			panic(err)
		}
	}
	return buf.String()
}
func GenBlockDetailByUser(items []string) string {
	//fmt.Println("------- items = ", items)

	var buf bytes.Buffer
	var err error
	var tmpl *jet.Template

	if strings.Contains(items[1], "interface") {
		tmpl, err = jet.NewSet(jet.NewInMemLoader()).Parse("template.jet", `
			type {{.legalName}} struct {
				Json {{.userType}}
			}
			
			func (s *{{.legalName}}) UnmarshalCSV(val string) error {
				return json.Unmarshal([]byte(val), &s.Json)
			}
`)
	} else {
		tmpl, err = jet.NewSet(jet.NewInMemLoader()).Parse("template.jet", `
/*
    TODO 必须实现自定义反序列化接口
    type {{.legalName}} {{.userType}}

    func (s *{{.legalName}}) UnmarshalCSV(val string) error {
		return json.Unmarshal([]byte(val), &s)
    }
*/
`)
	}
	err = tmpl.Execute(&buf, nil, map[string]interface{}{
		"legalName": utils.GoDsToLegalName(items[1]),
		"userType":  items[1],
	})
	if err != nil {
		panic(err)
	}
	return buf.String()
}
