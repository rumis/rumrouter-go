package parser

import "go/ast"

// FileEntry
type FileEntry struct {
	Key  string
	File ast.File
}

// FileEntries
type FileEntries []FileEntry

// Len get length of the files
func (e FileEntries) Len() int {
	return len(e)
}

// Less compare two file entry
func (e FileEntries) Less(i, j int) bool {
	return e[i].Key < e[j].Key
}

//Swap swap the elements
func (e FileEntries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}
