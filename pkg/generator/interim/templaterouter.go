package interim

type Router struct {
	Path       string
	Methods    []string
	Name       string
	Middleware string
	Context    string
}
