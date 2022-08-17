package source_code

import (
	"go/ast"
)

type StructSourceCode interface {
	PartOfCode
	CanRename
	RemoveTaggedField(targetComment string)
}

type structSourceCode struct {
	*ast.TypeSpec
}

func (s structSourceCode) Rename(newName string) {
	s.TypeSpec.Name.Name = newName
}

func (s structSourceCode) Name() string {
	return s.TypeSpec.Name.Name
}

func (s structSourceCode) RemoveTaggedField(targetComment string) {
	for i, field := range s.TypeSpec.Type.(*ast.StructType).Fields.List {
		if field.Tag != nil {
			if field.Tag.Value == targetComment {
				s.TypeSpec.Type.(*ast.StructType).Fields.List = append(s.TypeSpec.Type.(*ast.StructType).Fields.List[:i], s.TypeSpec.Type.(*ast.StructType).Fields.List[i+1:]...)
				continue
			}
		}
	}
}
