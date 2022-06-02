package config

import "strings"

// Options command line options
type Options struct {
	FrameworkType string
	SourcePath    string
	PackageName   string
	OutputPath    string
	Namespace     string
}

// IsContain 检查命名空间是否匹配
func (o Options) IsContains(ns string) bool {
	if o.Namespace == "" || ns == "" {
		return true
	}
	allNs := strings.Split(o.Namespace, ",")
	for _, n := range allNs {
		if n == ns {
			return true
		}
	}
	return false
}
