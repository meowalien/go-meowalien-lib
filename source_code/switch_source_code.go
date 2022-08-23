package source_code

import (
	"go/ast"
)

type SwitchSourceCode interface {
	PartOfCode
	InsertCase(callCase *ast.CaseClause)
}

type switchSourceCode struct {
	*ast.SwitchStmt
}

func (s switchSourceCode) Name() string {
	return "switch"
}

func (s switchSourceCode) InsertCase(callCase *ast.CaseClause) {
	s.SwitchStmt.Body.List = append(s.SwitchStmt.Body.List, callCase)
}
