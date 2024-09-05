package test

import (
	"csv_loader/core/utils"
	"fmt"
	"testing"
)

func TestSplitTitle(t *testing.T) {
	a := "FormulaString"
	fmt.Println(utils.SplitTitle(a, "_"))
}
