package generator

import (
	"fmt"

	"github.com/rumis/rumrouter-go/pkg/config"
	"github.com/rumis/rumrouter-go/pkg/generator/builder"
	"github.com/rumis/rumrouter-go/pkg/generator/echo"
	"github.com/rumis/rumrouter-go/pkg/generator/gin"
)

// InitRouter 初始化路由
func InitRouter(opts config.Options) error {
	switch opts.FrameworkType {
	case "echo":
		return InitEcho(opts)
	case "echov4":
		return InitEchov4(opts)
	case "gin":
		return InitGin(opts)
	default:
		return fmt.Errorf("framework %s is not available now", opts.FrameworkType)
	}
}

// InitEchov4 echov4
func InitEchov4(opts config.Options) error {
	// ta := interim.TemplateAnnotation{
	// 	Imports: []string{
	// 		"jiaoyan.com/pkg/test",
	// 		"jiaoyan.com/pkg/demo"},
	// 	RouterGroups: []interim.RouterGroup{{
	// 		Prefix: "/test1",
	// 		Routers: []interim.Router{
	// 			{
	// 				Path:       "/t1",
	// 				Methods:    []string{"GET"},
	// 				Name:       "t1",
	// 				Middleware: "auth,mtest",
	// 			},
	// 			{
	// 				Path:    "/t2",
	// 				Methods: []string{"GET", "POST"},
	// 				Name:    "t2",
	// 			},
	// 		},
	// 		Name:       "test.Test",
	// 		PackName:   "Test",
	// 		Middleware: "auth",
	// 	},
	// 		{
	// 			Prefix: "/demo",
	// 			Routers: []interim.Router{
	// 				{
	// 					Path:    "/d1",
	// 					Methods: []string{"GET"},
	// 					Name:    "d1",
	// 				},
	// 				{
	// 					Path:       "/d2",
	// 					Methods:    []string{"PUT", "OPTIONS"},
	// 					Name:       "d2",
	// 					Middleware: "auth",
	// 				},
	// 			},
	// 			Name:       "demo.Demo",
	// 			PackName:   "Demo",
	// 			Middleware: "auth,log",
	// 		},
	// 	},
	// 	Middlewares: map[string]interim.Middleware{
	// 		"auth": {
	// 			Name:     "auth",
	// 			FuncName: "middleware.Auth",
	// 		},
	// 		"log": {
	// 			Name:     "log",
	// 			FuncName: "middleware.Log",
	// 		},
	// 		"mtest": {
	// 			Name:     "mtest",
	// 			FuncName: "middleware.Test",
	// 		},
	// 	},
	// }

	ta, err := builder.ExtractTmplAnnotation(opts)
	if err != nil {
		return err
	}
	return echo.InitEchov4(ta, opts)
}

// InitEcho echo
func InitEcho(opts config.Options) error {
	ta, err := builder.ExtractTmplAnnotation(opts)
	if err != nil {
		return err
	}
	return echo.InitEcho(ta, opts)
}

// InitGin gin
func InitGin(opts config.Options) error {
	ta, err := builder.ExtractTmplAnnotation(opts)
	if err != nil {
		return err
	}
	return gin.InitGin(ta, opts)
}
