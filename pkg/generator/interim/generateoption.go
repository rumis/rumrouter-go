package interim

type GenerateOption struct {
	Name        string
	Funcs       map[string]interface{}
	Tmpl        string
	OutFileName string
	TmplAnno    TemplateAnnotation
}
