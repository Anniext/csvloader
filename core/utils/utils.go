package utils

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

var GRepeatDefinedType = make(map[string]bool)

func GetTableInfo(fileName string) []string {
	fileExt := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, fileExt)
	items := strings.Split(baseName, "-")
	return items
}

func ValidFlag(flagStr string) bool {
	flagStr = strings.ToLower(flagStr)
	if strings.Contains(flagStr, "a") || strings.Contains(flagStr, "s") {
		return true
	} else {
		return false
	}
}

func ParseTag(tag string) map[string][]string {
	tagDict := make(map[string][]string)
	var indexName string
	if tag == "" {
		return tagDict
	}

	trimItemsStr := strings.TrimSpace(tag)
	items := strings.Split(trimItemsStr, " ")
	for _, item := range items {
		if strings.Contains(item, "::") {
			return tagDict
		}
		keys := strings.Split(item, ":")
		key := keys[0]
		values := keys[1]
		if strings.Contains(key, "hash") && strings.Contains(values, "group") {
			indexPDict := ParseIndexValue(strings.Split(values[1:len(values)-1], ";"))
			group := indexPDict["group"]
			indexName = key + group
		} else {
			indexName = key
		}
		splitValues := strings.Split(values[1:len(values)-1], ";")
		tagDict[indexName] = append(tagDict[indexName], splitValues...)
	}
	return tagDict
}

func ParseIndexValue(items []string) map[string]string {
	retDict := make(map[string]string)
	for _, item := range items {
		if !strings.Contains(item, "=") {
			retDict["group"] = item
			continue
		}
		group := strings.Split(item, "=")
		key := group[0]
		value := group[1]
		retDict[key] = value
	}
	return retDict
}

func GoDsToLegalName(name string) string {
	// 替换所有需要转换的字符
	name = strings.ReplaceAll(name, " ", "W")
	name = strings.ReplaceAll(name, "}", "H")
	name = strings.ReplaceAll(name, "{", "G")
	name = strings.ReplaceAll(name, "[", "B")
	name = strings.ReplaceAll(name, "]", "E")
	name = strings.ReplaceAll(name, "(", "Q")
	name = strings.ReplaceAll(name, ")", "P")
	name = strings.ReplaceAll(name, "*", "X")

	return name
}

func GetClassSetStr(className string, unique int) string {
	if unique == 1 {
		return fmt.Sprintf("*%s", className)
	}
	return fmt.Sprintf("[]*%s", className)
}

func SplitTitle(s string, sep string) string {
	var title string
	words := strings.Split(s, sep)
	if len(words) > 1 {
		for i, word := range words {
			lowerStr := strings.ToLower(word)
			words[i] = strings.Title(lowerStr)
		}
	} else {
		lowerStr := strings.ToLower(words[0])
		words[0] = strings.Title(lowerStr)
	}
	title = strings.Join(words, sep)
	return title
}

func FormulaCsvToGo(s string) string {
	s = strings.ReplaceAll(s, "[", "(")
	s = strings.ReplaceAll(s, "]", ")")
	s = strings.ReplaceAll(s, "format", "string")
	s = strings.ReplaceAll(s, "Rand", "RandFloat64")
	s = strings.ReplaceAll(s, "RAND", "RandFloat64")
	s = strings.ReplaceAll(s, "Floor", "math.Floor")
	s = strings.ReplaceAll(s, "FLOOR", "math.Floor")
	s = strings.ReplaceAll(s, "Ceil", "math.Ceil")
	s = strings.ReplaceAll(s, "CEIL", "math.Ceil")
	s = strings.ReplaceAll(s, "Max", "math.Max")
	s = strings.ReplaceAll(s, "MAX", "math.Max")
	s = strings.ReplaceAll(s, "Min", "math.Min")
	s = strings.ReplaceAll(s, "MIN", "math.Min")
	s = strings.ReplaceAll(s, "Lg", "math.Log10")
	s = strings.ReplaceAll(s, "LG", "math.Log10")
	s = strings.ReplaceAll(s, "Pow", "math.Pow")
	s = strings.ReplaceAll(s, "POW", "math.Pow")
	s = strings.ReplaceAll(s, "pow", "math.Pow")
	s = strings.ReplaceAll(s, "IF", "IF") // 这里实际上没有改变
	s = strings.ReplaceAll(s, ">", ",")
	return s
}

func JudgmentCommand(cmd string) error {
	if _, err := exec.LookPath(cmd); err != nil {
		return err
	}
	return nil
}
