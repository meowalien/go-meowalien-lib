package source_code

import (
	"github.com/meowalien/go-meowalien-lib/errs"
	"go/ast"
	"go/token"
)

type ImportSet []ImportSourceCode

func (s ImportSet) Merge() (newDecl *ast.GenDecl) {
	if len(s) == 0 {
		return nil
	}
	if isc, ok := s[0].(importSourceCode); ok {
		newDecl = &ast.GenDecl{
			Tok:    token.IMPORT,
			TokPos: isc.GenDecl.TokPos,
			Lparen: isc.GenDecl.Lparen,
			Rparen: isc.GenDecl.Rparen,
		}
	} else {
		panic(errs.New("unexpected type: %T", s[0]))
	}

	for _, code := range s {
		if isc, ok := code.(importSourceCode); ok {
			newDecl.Specs = append(newDecl.Specs, isc.Specs...)
		} else {
			panic(errs.New("unexpected type: %T", s[0]))
		}
	}
	return
}

type ImportSourceCode interface {
	PartOfCode
}

func NewImportSourceCode(gd *ast.GenDecl) ImportSourceCode {
	return importSourceCode{GenDecl: gd}
}

type importSourceCode struct {
	*ast.GenDecl
}

func (i importSourceCode) Append(scB ...PartOfCode) (err error) {
	for _, code := range scB {
		switch ct := code.(type) {
		case importSourceCode:
			i.GenDecl.Specs = append(i.GenDecl.Specs, ct.GenDecl.Specs...)
		default:
			err = errs.New("unexpected type %T", ct)
		}
	}
	return
}

func (i importSourceCode) Name() string {
	return "import"
}
