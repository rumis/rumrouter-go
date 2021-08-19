package parser

import (
	"bufio"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strings"
)

type Parser struct{}

// NewParser create a new annotation paerser instance
func NewParser() *Parser {
	return &Parser{}
}

// parse source code for comments
func (p *Parser) ParseSource(srcPath string, includeRegex string, excludeRegex string) (*Visitor, error) {
	var iReg = regexp.MustCompile(includeRegex)
	var eReg = regexp.MustCompile(excludeRegex)
	packages, err := parseDir(srcPath, iReg, eReg)
	if err != nil {
		return nil, err
	}
	v := &Visitor{}

	basePackageName, err := getBasePackageName(srcPath)
	if err != nil {
		return nil, err
	}
	v.BasePackageName = basePackageName
	v.BaseSrcPath = srcPath

	for _, pack := range packages {
		parsePackage(pack, v)
	}
	return v, nil
}

// parse source code dirctory
func parseDir(dirPath string, iReg, eReg *regexp.Regexp) (map[string]*ast.Package, error) {
	fileSet := token.NewFileSet()
	// parse current dirctory
	packageMap, err := parser.ParseDir(fileSet, dirPath, func(fi os.FileInfo) bool {
		if eReg.MatchString(fi.Name()) {
			return false
		}
		return iReg.MatchString(fi.Name())
	}, parser.ParseComments)
	if err != nil {
		// current dirctory parse failed , create new ast.package map
		packageMap = make(map[string]*ast.Package)
	}
	// parse the child dirctory
	fileInfos, e := ioutil.ReadDir(dirPath)
	if e != nil {
		return packageMap, nil
	}
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			continue
		}
		childPackage, ce := parseDir(dirPath+"/"+fileInfo.Name(), iReg, eReg)
		if ce != nil {
			continue
		}
		// add the child package to packagemap
		for k, v := range childPackage {
			packageMap[k] = v
		}
	}
	return packageMap, nil
}

// parsePackage get the package
func parsePackage(pack *ast.Package, v *Visitor) {
	entries := getFileEntries(pack.Files)
	for _, entry := range entries {
		v.CurrentFileName = entry.Key
		v.CurrentFilePath = path.Dir(entry.Key)

		ast.Walk(v, &entry.File)
	}
}

// getFileEntries get the entry of source file
func getFileEntries(fileMap map[string]*ast.File) FileEntries {
	var entries FileEntries = make([]FileEntry, 0, len(fileMap))
	for k, file := range fileMap {
		if file == nil {
			continue
		}
		entries = append(entries, FileEntry{
			Key:  k,
			File: *file,
		})
	}
	sort.Sort(entries)
	return entries
}

// parse go.mod file ,get the base package name
func getBasePackageName(workspace string) (string, error) {
	file, err := os.Open(path.Join(workspace, "go.mod"))
	if err != nil {
		return "", errors.New("failed to get base package name,must begain with the go.mod dirctory")
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			return "", errors.New("failed to get base package name,invalid go.mod")
		}
		line = strings.TrimSpace(line)
		words := strings.Split(line, " ")
		if words[0] == "module" {
			return words[1], nil
		}
	}
}
