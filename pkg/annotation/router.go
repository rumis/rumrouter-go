package annotation

type Router struct {
	Path       string `json:"path"`
	Method     string `json:"method"`
	Middleware string `json:"middleware"`
	Namespace  string `json:"namespace"`
}
