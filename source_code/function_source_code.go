package source_code

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"go/ast"
)

type FunctionSourceCode interface {
	PartOfCode
	CanFindSwitch
}

type funcSourceCode struct {
	*ast.FuncDecl
}

func (f funcSourceCode) Name() string {
	return f.FuncDecl.Name.Name
}

func (f funcSourceCode) FindSwitch() (ans []SwitchSourceCode, err error) {
	inspect(f.Body, func(nt *ast.SwitchStmt) bool {
		ans = append(ans, switchSourceCode{SwitchStmt: nt})
		return true
	})
	return
}

func (f funcSourceCode) Append(scB ...PartOfCode) (err error) {
	for i, code := range scB {
		switch tc := code.(type) {
		case ast.Stmt:
			switch tcc := tc.(type) {
			case *ast.SwitchStmt:
				f.Body.List = append(f.Body.List, tcc)
			case switchSourceCode:
				f.Body.List = append(f.Body.List, tcc.SwitchStmt)
			default:
				return errs.New(fmt.Errorf("%dth code is not a switch statement", i))
			}
		default:
			err = errs.New("unexpected type %T at index %d", tc, i)
		}
	}
	return
}

func (f funcSourceCode) Rename(newName string) {
	f.FuncDecl.Name.Name = newName
}
