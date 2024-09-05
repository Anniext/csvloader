package test

import (
	"csv_loader/module"
	"fmt"
	"testing"
)

func TestCsvLoader(t *testing.T) {
	module.Load()
	for _, v := range module.GXianzhuTowerRewardManager.GetAllXianzhuTower_Rewards() {
		fmt.Println(*v)
	}
}
