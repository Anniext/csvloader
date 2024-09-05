package global

import (
	"csvloader/core/loader"
	"csvloader/internal/config"
)

var CsvLoaderMap map[string]*loader.Tables

var Config *config.Config
