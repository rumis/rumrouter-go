package parser

import (
	"go/ast"
	"path/filepath"

	"github.com/rumis/rumrouter-go/pkg/model"
)

type Visitor struct {
	CurrentFileName string
	CurrentFilePath string
	PackagePath     string
	PackageName     string
	BasePackageName string
	BaseSrcPath     string
	Structs         []model.Controller
	Operations      []model.Operation
}

func (v *Visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return v
	}

	// get current package name
	v.extractPackageName(node)
	// get all structs
	v.parseAsStruct(node)
	// get all functions
	v.parseAsOperation(node)

	return v
}

// parse package name
func (v *Visitor) extractPackageName(node ast.Node) {
	file, ok := node.(*ast.File)
	if !ok {
		return
	}
	if file.Name == nil {
		return
	}
	rel, err := filepath.Rel(v.BaseSrcPath, v.CurrentFilePath)
	if err != nil {
		return
	}
	v.PackagePath = filepath.Join(v.BasePackageName, rel)
	v.PackageName = file.Name.Name
}

// get all structs
func (v *Visitor) parseAsStruct(node ast.Node) {
	file, ok := node.(*ast.File)
	if !ok {
		return
	}
	for _, decl := range file.Decls {
		mStruct := extractGenDeclForStruct(decl)
		if mStruct == nil {
			continue
		}
		mStruct.PackagePath = v.PackagePath
		mStruct.PackageName = v.PackageName
		mStruct.Filename = v.CurrentFileName
		v.Structs = append(v.Structs, *mStruct)
	}
}

// gen decl for struct
func extractGenDeclForStruct(node ast.Node) *model.Controller {
	genDecl, ok := node.(*ast.GenDecl)
	if !ok {
		return nil
	}
	mStruct := extractSpecsForStruct(genDecl.Specs)
	if mStruct == nil {
		return nil
	}
	mStruct.DocLines = extractComments(genDecl.Doc)
	return mStruct
}

// struct type
func extractSpecsForStruct(specs []ast.Spec) *model.Controller {
	if len(specs) == 0 {
		return nil
	}
	typeSpec, ok := specs[0].(*ast.TypeSpec)
	if !ok {
		return nil
	}
	_, ok = typeSpec.Type.(*ast.StructType)
	if !ok {
		return nil
	}
	return &model.Controller{
		Name: typeSpec.Name.Name,
	}
}

// get the coments
func extractComments(commentGroup *ast.CommentGroup) []string {
	lines := make([]string, 0)
	if commentGroup == nil {
		return lines
	}
	for _, comment := range commentGroup.List {
		lines = append(lines, comment.Text)
	}
	return lines
}

// get all functions define
func (v *Visitor) parseAsOperation(node ast.Node) {
	file, ok := node.(*ast.File)
	if !ok {
		return
	}
	for _, decl := range file.Decls {
		mOperation := extractOperation(decl)
		if mOperation == nil {
			continue
		}
		mOperation.PackagePath = v.PackagePath
		mOperation.PackageName = v.PackageName
		mOperation.Filename = v.CurrentFileName

		if mOperation.RelatedStruct != nil {
			mOperation.RelatedStruct.PackagePath = v.PackagePath
			mOperation.RelatedStruct.PackageName = v.PackageName
		}

		v.Operations = append(v.Operations, *mOperation)
	}
}

// parse functions
func extractOperation(decl ast.Decl) *model.Operation {
	funcDecl, ok := decl.(*ast.FuncDecl)
	if !ok {
		return nil
	}
	mOperation := model.Operation{
		DocLines: extractComments(funcDecl.Doc),
	}
	if funcDecl.Recv != nil {
		recv := extractReceiver(funcDecl.Recv)
		if recv != nil {
			mOperation.RelatedStruct = recv
		}
	}
	if funcDecl.Name != nil {
		mOperation.Name = funcDecl.Name.Name
	}
	return &mOperation
}

// parse receiver
func extractReceiver(fieldList *ast.FieldList) *model.Receiver {
	if fieldList == nil || fieldList.List[0] == nil {
		return nil
	}

	field := fieldList.List[0]
	recv := extractField(field)
	if recv == nil {
		return nil
	}

	if field.Names[0] != nil {
		recv.Name = field.Names[0].Name
	}

	return recv
}

// parse field
func extractField(field *ast.Field) *model.Receiver {
	fieldType, ok := field.Type.(*ast.StarExpr)
	if ok {
		return extractStarExpr(fieldType)
	}
	return extractIdentExpr(field.Type)
}

// parse receiver pointer
func extractStarExpr(expr *ast.StarExpr) *model.Receiver {
	if expr == nil {
		return nil
	}
	recv := extractIdentExpr(expr.X)
	if recv == nil {
		return nil
	}

	recv.Star = true

	return recv

}

// parse IdentExpr
func extractIdentExpr(expr ast.Expr) *model.Receiver {
	if expr == nil {
		return nil
	}
	fieldType, ok := expr.(*ast.Ident)
	if !ok {
		return nil
	}
	return &model.Receiver{
		TypeName: fieldType.Name,
	}
}
