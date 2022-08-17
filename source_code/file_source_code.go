package source_code

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func NewFile() FileSourceCode {
	return fileSourceCode{
		File: &ast.File{},
	}
}

func ParseCode(code string) (theFile FileSourceCode, err error) {
	f, err := parser.ParseFile(token.NewFileSet(), "", code, parser.ParseComments)
	if err != nil {
		err = errs.New(err)
		return
	}
	theFile = fileSourceCode{File: f}
	return
}

type FileSourceCode interface {
	PartOfCode
	Appendable
	GetImports() ImportSet
	RemoveImports()
	FindTaggedStruct(targetComment string) (genDecl []StructSourceCode, err error)
	BuildAsPlugin(file string) error
	WriteToFile(filePath string) (err error)
	FindFunctions(s2 string) (sc []FunctionSourceCode, err error)
	Import(s string)
	SetPackage(s string)
	TagExist(tag string) bool
	Package() string
	RemoveTag(tag string)
}

type fileSourceCode struct {
	*ast.File
}

func (f fileSourceCode) RemoveTag(tag string) {
	inspect(f.File, func(nt *ast.CommentGroup) bool {
		for i := 0; i < len(nt.List); i++ {
			if nt.List[i].Text == tag {
				nt.List = append(nt.List[:i], nt.List[i+1:]...)
				i--
			}
		}
		return true
	})
}

func (f fileSourceCode) TagExist(tag string) (exist bool) {
	inspect(f.File, func(nt *ast.Comment) bool {
		if nt.Text == tag {
			exist = true
			return false
		}
		return true
	})
	return
}

func (f fileSourceCode) Package() string {
	return f.Name()
}

func (f fileSourceCode) SetPackage(s string) {
	if f.File.Name != nil {
		f.File.Name.Name = s
	} else {
		f.File.Name = &ast.Ident{Name: s}
	}
}

func (f fileSourceCode) Import(pkgName string) {
	ip := &ast.ImportSpec{
		Path: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("\"%s\"", pkgName)},
	}
	f.Decls = append([]ast.Decl{&ast.GenDecl{
		Tok:   token.IMPORT,
		Specs: []ast.Spec{ip},
	}}, f.Decls...)
}

func (f fileSourceCode) FindFunctions(funcName string) (sc []FunctionSourceCode, err error) {
	inspect(f.File, func(nt *ast.FuncDecl) bool {
		if nt.Name == nil {
			err = errs.New("unexpected nil nt.Name")
			return false
		}
		if nt.Name.Name == funcName {
			sc = append(sc, funcSourceCode{FuncDecl: nt})
		}
		return true
	})
	return
}

func (f fileSourceCode) Name() string {
	return f.File.Name.Name
}

const buildAsPluginTempFileName = "plugin_source_code_*.go"

func (f fileSourceCode) BuildAsPlugin(file string) (err error) {
	tempFile, err := ioutil.TempFile("", buildAsPluginTempFileName)
	if err != nil {
		err = errs.New(err)
		return
	}
	fmt.Println("source file: ", tempFile.Name())
	err = f.WriteToFile(tempFile.Name())
	if err != nil {
		err = errs.New(err)
		return
	}
	file, err = filepath.Abs(file)
	if err != nil {
		err = errs.New(err)
		return
	}
	fmt.Println("BuildAsPlugin new plugin: ", file)

	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", file, tempFile.Name()) //nolint: gosec
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		err = errs.New(err)
		return
	}
	ourPipe, err := cmd.StdoutPipe()
	if err != nil {
		err = errs.New(err)
		return
	}

	go func() {
		_, err2 := io.Copy(os.Stdout, errPipe)
		if err2 != nil {
			err = errs.New(err, err2)
		}
	}()
	go func() {
		_, err2 := io.Copy(os.Stdout, ourPipe)
		if err2 != nil {
			err = errs.New(err, err2)
		}
	}()

	err = cmd.Start()
	if err != nil {
		err = errs.New(err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}

func (f fileSourceCode) Append(set ...PartOfCode) (err error) {
	for _, code := range set {
		fmt.Printf("coee type %T\n", code)
		switch cdeType := code.(type) {
		case fileSourceCode:
			importsInputs := cdeType.GetImports()
			cdeType.RemoveImports()

			currentImports := f.GetImports()
			f.RemoveImports()

			f.appendFront(append(currentImports, importsInputs...).Merge())

			f.appendEnd(cdeType.File.Decls...)
		case ast.Spec:
			switch cdeType1 := cdeType.(type) {
			case structSourceCode:
				f.appendEnd(&ast.GenDecl{
					Tok:   token.TYPE,
					Specs: []ast.Spec{cdeType1.TypeSpec},
				})
				continue
			default:
				return errs.New("unsupported type: %T", cdeType)
			}
		case ast.Decl:
			switch codeType1 := cdeType.(type) {
			case funcSourceCode:
				f.appendEnd(codeType1.FuncDecl)
			default:
				return errs.New("unsupported type: %T", cdeType)
			}
			continue
		default:
			return errs.New("unexpected type %T", cdeType)
		}
	}
	return
}

func (f fileSourceCode) Rename(newName string) {
	f.File.Name.Name = newName
}

func (f fileSourceCode) FindTaggedStruct(targetComment string) (genDecl []StructSourceCode, err error) {
	inspect(f.File, func(nt *ast.GenDecl) bool {
		if nt.Doc != nil {
			for _, comment := range nt.Doc.List {
				if comment.Text == targetComment {
					if nt.Specs != nil {
						if len(nt.Specs) > 1 {
							panic("unexpected len(nt.Specs) > 1")
						}
						str, ok := nt.Specs[0].(*ast.TypeSpec)
						if ok {
							if _, ok = str.Type.(*ast.StructType); ok {
								genDecl = append(genDecl, structSourceCode{TypeSpec: str})
							}
						}
					}
				}
			}
		}
		return true
	})
	return
}

func (f fileSourceCode) WriteToFile(filePath string) (err error) {
	defer func() {
		if rec := recover(); rec != nil {
			err = errs.New(err, rec)
		}
	}()
	theFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		err = errs.New(err)
		return
	}
	defer func() {
		err2 := theFile.Close()
		if err2 != nil {
			err = errs.New(err, err2)
		}
	}()
	//pretty.Println("theFile: ", f.File)
	if err = printer.Fprint(theFile, token.NewFileSet(), f.File); err != nil {
		return errs.New(err)
	}
	return nil
}

func (f fileSourceCode) GetImports() (set ImportSet) {
	for _, decl := range f.Decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		if gd.Tok != token.IMPORT {
			continue
		}
		set = append(set, NewImportSourceCode(gd))
	}
	return
}

func (f fileSourceCode) RemoveImports() {
	for i := 0; i < len(f.File.Decls); i++ {
		gd, ok := f.File.Decls[i].(*ast.GenDecl)
		if !ok {
			continue
		}
		if gd.Tok != token.IMPORT {
			continue
		}
		f.removeIndex(i)
		i--
	}
	return
}

func (f fileSourceCode) removeIndex(i int) {
	f.File.Decls = append(f.File.Decls[:i], f.File.Decls[i+1:]...)
}

func (f fileSourceCode) appendFront(imp ...ast.Decl) {
	f.File.Decls = append(imp, f.File.Decls...)
}

func (f fileSourceCode) appendEnd(decls ...ast.Decl) {
	f.File.Decls = append(f.File.Decls, decls...)
}
