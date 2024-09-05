package jet

import (
	"csvloader/global"
	"github.com/CloudyKit/jet/v6"
)

func NewJetLoader() {
	var jetLoader jet.Loader
	if global.Config.TemplatePath != "" {
		jetLoader = jet.NewOSFileSystemLoader(global.Config.TemplatePath)
		if jetLoader == nil {
			panic("failed to jet loader:")
		}
	} else {
		jetLoader = jet.NewInMemLoader()
	}
	global.Config.DefinedJetLoader = jet.NewSet(jetLoader, jet.WithSafeWriter(nil))
}
