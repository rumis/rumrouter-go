package interim

import "strings"

type TemplateAnnotation struct {
	PackageName  string
	Imports      []string
	RouterGroups []RouterGroup
	Middlewares  map[string]Middleware
}

func (ta TemplateAnnotation) GetMiddleware(name string) string {
	if name == "" {
		return ""
	}
	names := strings.Split(name, ",")
	funcNames := make([]string, 0)
	for _, mname := range names {
		middleware, ok := ta.Middlewares[mname]
		if !ok {
			continue
		}
		funcNames = append(funcNames, middleware.FuncName)
	}
	if len(funcNames) == 0 {
		return ""
	}
	return " ," + strings.Join(funcNames, ",")
}
