package echo

import (
	"path"

	"github.com/rumis/rumrouter-go/pkg/config"
	"github.com/rumis/rumrouter-go/pkg/generator/builder"
	"github.com/rumis/rumrouter-go/pkg/generator/interim"
)

const echoTmplStr = `package {{.PackageName}}

import (
    {{- range $idx,$ele := .Imports}}
    "{{$ele}}"
    {{- end}}
)

func InitRouter(app *echo.Echo) {
	{{- range $idx,$group := .RouterGroups}}

	g{{$idx}} := app.Group("{{$group.Prefix}}"{{MiddleName $group.Middleware}})
	g{{$group.PackName}}Inst := {{$group.Name}}{}
	{{- range $i,$router := $group.Routers}}
	{{- range $j,$method := $router.Methods}}
	g{{$idx}}.{{$method}}("{{$router.Path}}",g{{$group.PackName}}Inst.{{$router.Name}}{{MiddleName $router.Middleware}})
	{{- end}}	
	{{- end}}
	{{- end}}
}
`

// InitEcho generate echo
func InitEcho(tmplAnno interim.TemplateAnnotation, opts config.Options) error {
	tmplAnno.Imports = append(tmplAnno.Imports, "github.com/labstack/echo")
	return initEcho(tmplAnno, path.Join(opts.OutputPath, "echo.gen.go"))
}

// InitEchov4 generate echo v4
func InitEchov4(tmplAnno interim.TemplateAnnotation, opts config.Options) error {
	tmplAnno.Imports = append(tmplAnno.Imports, "github.com/labstack/echo/v4")
	return initEcho(tmplAnno, path.Join(opts.OutputPath, "echo.v4.gen.go"))
}

// initEcho generate echo source code
func initEcho(tmplAnno interim.TemplateAnnotation, out string) error {
	opts := interim.GenerateOption{
		Name:        "echo",
		Tmpl:        echoTmplStr,
		TmplAnno:    tmplAnno,
		OutFileName: out,
	}
	return builder.GenerateFromTempAnnotation(opts)
}
