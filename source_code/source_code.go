package source_code

import (
	"github.com/meowalien/go-meowalien-lib/errs"
	"go/ast"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

type PartOfCode interface {
	ast.Node
	Name() string
}

type CanRename interface {
	Rename(newName string)
}

type Appendable interface {
	Append(scB ...PartOfCode) (err error)
}

type CanFindSwitch interface {
	FindSwitch() ([]SwitchSourceCode, error)
}

type FileAddress string

var rep = strings.NewReplacer(".", "_", "/", "__")

func (a FileAddress) FormatAsValidGoIdentifier() (s string, err error) {
	here, err := filepath.Abs("./")
	if err != nil {
		err = errs.New(err)
		return
	}
	path, err := filepath.Abs(string(a))
	if err != nil {
		err = errs.New(err)
		return
	}
	s, err = filepath.Rel(here, path)
	if err != nil {
		err = errs.New(err)
		return
	}
	s = rep.Replace(s)
	return
}

type SourceFileSet map[FileAddress]FileSourceCode

func (s *SourceFileSet) Foreach(each func(filePath FileAddress, sourceCodeFile FileSourceCode) bool) {
	for filePath, sourceCodeFile := range (map[FileAddress]FileSourceCode)(*s) {
		if !each(filePath, sourceCodeFile) {
			return
		}
	}
}

func ParseDir(dir string) (sc SourceFileSet, err error) {
	sc = map[FileAddress]FileSourceCode{}
	err = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		ff, err := os.Open(path)
		if err != nil {
			err = errs.New(err)
			return err
		}
		defer func() {
			err = errs.New(err, ff.Close())
		}()
		b, err := io.ReadAll(ff)
		if err != nil {
			err = errs.New(err)
			return err
		}

		sc[FileAddress(path)], err = ParseCode(string(b))
		if err != nil {
			err = errs.New(err)
			return err
		}
		return nil
	})
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}

const jsonMarshalTempFileName = "json_marshal_plugin_source_code_*.so"

func JsonMarshal(code StructSourceCode) (js []byte, sourceFilePath string, err error) {
	newFile := NewFile()
	newFile.SetPackage("main")
	err = newFile.Append(code)
	if err != nil {
		err = errs.New(err)
		return
	}
	newFile.Import("encoding/json")
	err = newFile.Append(newMarshalFunction(code.Name()))
	if err != nil {
		err = errs.New(err)
		return
	}
	tempFile, err := ioutil.TempFile("", jsonMarshalTempFileName)
	if err != nil {
		err = errs.New(err)
		return
	}
	defer func() {
		err = errs.New(err, tempFile.Close())
	}()
	sourceFilePath, err = newFile.BuildAsPlugin(tempFile.Name())
	if err != nil {
		err = errs.New(err)
		return
	}
	pg, err := plugin.Open(tempFile.Name())
	if err != nil {
		err = errs.New(err)
		return
	}
	sb, err := pg.Lookup("Marshal")
	if err != nil {
		err = errs.New(err)
		return
	}
	marshal, ok := sb.(func() ([]byte, error))
	if !ok {
		err = errs.New("unexpected type")
		return
	}
	js, err = marshal()
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}

func newMarshalFunction(name string) FunctionSourceCode {
	return funcSourceCode{&ast.FuncDecl{
		Name: &ast.Ident{
			Name: "Marshal",
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{
							Elt: &ast.Ident{
								NamePos: 165,
								Name:    "byte",
							},
						},
					},
					{
						Type: &ast.Ident{
							Name: "error",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X: &ast.Ident{
									Name: "json",
								},
								Sel: &ast.Ident{
									Name: "Marshal",
								},
							},
							Args: []ast.Expr{&ast.CompositeLit{
								Type: &ast.Ident{
									Name: name,
								},
							}},
						},
					},
				},
			},
		},
	}}
}

func inspect[T any](code ast.Node, f func(nt T) bool) {
	ast.Inspect(code, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		t, ok := n.(T)
		if !ok {
			return true
		}
		return f(t)

	})
}
