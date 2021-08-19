package model

// Controller struct define
type Controller struct {
	PackagePath string
	PackageName string
	Filename    string
	DocLines    []string
	Name        string
}

// Receiver function receiver
type Receiver struct {
	PackagePath string
	PackageName string
	Name        string
	TypeName    string
	Star        bool
}

// Operation function define
type Operation struct {
	PackagePath   string
	PackageName   string
	Filename      string
	DocLines      []string
	RelatedStruct *Receiver
	Name          string
}
