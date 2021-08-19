package annotation

const ANNO_ROUTERGROUP = "routergroup"
const ANNO_ROUTER = "router"
const ANNO_MIDDLEWARE = "middleware"

type AnnoType int

const (
	Initiale AnnoType = iota
	AnnoName
	AttrName
	AttrValue
	Done
)

type Annotation struct {
	Name  string            `json:"name"`
	Attrs map[string]string `json:"attrs"`
}
