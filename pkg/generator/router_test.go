package generator

import (
	"testing"

	"github.com/rumis/rumrouter-go/pkg/config"
)

func TestRouterGin(t *testing.T) {

	err := InitGin(config.Options{
		FrameworkType: "gin",
		SourcePath:    "/data/workspace/jiaoyango/courseware_service",
		PackageName:   "test",
		OutputPath:    "/data/workspace/rum/rumrouter-go/test",
		Namespace:     "",
	})

	if err != nil {
		panic(err)
	}

}
