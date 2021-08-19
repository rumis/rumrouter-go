package interim

type RouterGroup struct {
	Prefix     string
	Name       string
	Routers    []Router
	PackName   string
	Middleware string
}
