package annotation

type RouterGroup struct {
	Prefix     string `json:"prefix"`
	Middleware string `json:"middleware"`
}
