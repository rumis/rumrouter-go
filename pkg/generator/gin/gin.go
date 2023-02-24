package gin

import (
	"fmt"
	"path"
	"strings"

	"github.com/rumis/rumrouter-go/pkg/config"
	"github.com/rumis/rumrouter-go/pkg/generator/builder"
	"github.com/rumis/rumrouter-go/pkg/generator/interim"
)

const ginTmplStr = `package {{.PackageName}}

import (
    {{- range $idx,$ele := .Imports}}
    "{{$ele}}"
    {{- end}}
)

func InitRouter(app *gin.Engine) {
	{{- range $idx,$group := .RouterGroups}}

	g{{$idx}} := app.Group("{{$group.Prefix}}"{{MiddleName $group.Middleware}})
	g{{$group.PackName}}Inst := {{$group.Name}}{}
	{{- range $i,$router := $group.Routers}}
	{{- range $j,$method := $router.Methods}}
	g{{$idx}}.{{$method}}("{{$router.Path}}"{{Context $router.Context}}{{MiddleName $router.Middleware}},g{{$group.PackName}}Inst.{{$router.Name}})
	{{- end}}	
	{{- end}}
	{{- end}}
}
`

// InitGin generate gin
func InitGin(tmplAnno interim.TemplateAnnotation, opts config.Options) error {
	tmplAnno.Imports = append(tmplAnno.Imports, "github.com/gin-gonic/gin")
	return initGin(tmplAnno, path.Join(opts.OutputPath, "gin.gen.go"))
}

// initGin generate gin source code
func initGin(tmplAnno interim.TemplateAnnotation, out string) error {

	opts := interim.GenerateOption{
		Name:        "gin",
		Tmpl:        ginTmplStr,
		TmplAnno:    tmplAnno,
		OutFileName: out,
		Funcs: map[string]interface{}{
			"Context": func(ctx string) string {
				if ctx == "" {
					return ""
				}
				kvs := strings.Split(ctx, ",")
				opts := []string{}
				for _, v := range kvs {
					idx := strings.Index(v, "=")
					opts = append(opts, fmt.Sprintf("c.Set(\"%s\",%s)", v[:idx], v[idx+1:]))
				}
				return fmt.Sprintf("%s%s%s", ",func(c *gin.Context){", strings.Join(opts, ";"), `}`)
			},
		},
	}
	return builder.GenerateFromTempAnnotation(opts)
}
