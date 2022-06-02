package builder

import (
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/rumis/rumrouter-go/pkg/annotation"
	"github.com/rumis/rumrouter-go/pkg/config"
	"github.com/rumis/rumrouter-go/pkg/generator/interim"
	"github.com/rumis/rumrouter-go/pkg/parser"
)

// get annotation from template
func ExtractTmplAnnotation(opts config.Options) (interim.TemplateAnnotation, error) {
	annoParser := parser.NewParser()
	annoVisitor, err := annoParser.ParseSource(opts.SourcePath, "^.*.go$", "^.*_test.go$")

	ta := interim.TemplateAnnotation{
		PackageName: opts.PackageName,
		Middlewares: make(map[string]interim.Middleware),
	}

	if err != nil {
		return ta, err
	}
	packages := make(map[string]string)
	routerGroupMap := make(map[string]interim.RouterGroup)

	// extra router groups
	for _, group := range annoVisitor.Structs {
		for _, docLine := range group.DocLines {
			routerGroup, err := annotation.ParseRouterGroup(docLine)
			if err != nil {
				continue
			}
			// check is it contains the namespace
			if !opts.IsContains(routerGroup.Namespace) {
				continue
			}

			// keep the package path
			packages[group.PackagePath] = group.PackagePath

			inteRouteGroup := interim.RouterGroup{
				Prefix:     routerGroup.Prefix,
				Middleware: routerGroup.Middleware,
				Name:       group.PackageName + "." + group.Name,
				PackName:   group.Name,
			}
			routerGroupMap[path.Join(group.PackagePath, group.Name)] = inteRouteGroup
		}
	}
	// extra router and middleware
	for _, opera := range annoVisitor.Operations {
		for _, docLine := range opera.DocLines {
			if opera.RelatedStruct != nil {
				// router is a mehtod with receiver
				router, err := annotation.ParseRouter(docLine)
				if err != nil {
					continue
				}
				// check is it contains the namespace
				if !opts.IsContains(router.Namespace) {
					continue
				}

				inteRouter := interim.Router{
					Path:       router.Path,
					Methods:    strings.Split(strings.ToUpper(router.Method), ","),
					Name:       opera.Name,
					Middleware: router.Middleware,
				}
				relRouterGroupName := path.Join(opera.RelatedStruct.PackagePath, opera.RelatedStruct.TypeName)
				routerGroup, ok := routerGroupMap[relRouterGroupName]
				if ok {
					routerGroup.Routers = append(routerGroup.Routers, inteRouter)
					routerGroupMap[relRouterGroupName] = routerGroup
					packages[opera.PackagePath] = opera.PackagePath
				}
			} else {
				// middleware is only a function whithout receiver
				middleware, err := annotation.ParseMiddleware(docLine)
				if err != nil {
					continue
				}
				inteMiddleware := interim.Middleware{
					Name:     middleware.Name,
					FuncName: opera.PackageName + "." + opera.Name,
				}
				ta.Middlewares[middleware.Name] = inteMiddleware
				packages[opera.PackagePath] = opera.PackagePath
			}
		}
	}
	for _, rgroup := range routerGroupMap {
		ta.RouterGroups = append(ta.RouterGroups, rgroup)
	}
	for _, pack := range packages {
		ta.Imports = append(ta.Imports, pack)
	}

	return ta, nil
}

// generate source code from template and annotation
func GenerateFromTempAnnotation(opts interim.GenerateOption) error {
	funcMap := template.FuncMap{
		"MiddleName": func(mname string) string {
			return opts.TmplAnno.GetMiddleware(mname)
		},
	}
	if opts.Funcs != nil {
		for name, fnc := range opts.Funcs {
			funcMap[name] = fnc
		}
	}
	tmpl := template.New(opts.Name)
	tmpl = tmpl.Funcs(funcMap)
	tmpl, err := tmpl.Parse(opts.Tmpl)
	if err != nil {
		return err
	}

	// 检查原文件是否存在，如存在则删除
	_, err = os.Stat(opts.OutFileName)
	if err == nil {
		os.Remove(opts.OutFileName)
	}

	// 构建目录并生成文件
	err = os.MkdirAll(path.Dir(opts.OutFileName), 0755)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(opts.OutFileName, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	err = tmpl.Execute(file, opts.TmplAnno)
	file.Close()
	return err
}
