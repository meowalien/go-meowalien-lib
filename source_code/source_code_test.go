package source_code

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"testing"
)

func TestJsonMarshal(t *testing.T) {
	sourceCode := "package main\n\nimport (\n\t\"encoding/json\"\n\n\t\"errors\"\n)\n\n//data\ntype Example struct {\n\tData string `json:\"data\"`\n}\nfunc Decoder(key string, data []byte) (b any, err error) {\n\tswitch key {\n\tdefault:\n\t\treturn nil, errors.New(\"no struct found for key \" + key)\n\t}\n}\n\nfunc Marshal()([]byte , error){\n\treturn json.Marshal(Example{})\n}\n\t"
	codeFile, err := ParseCode(sourceCode)
	if err != nil {
		err = errs.New(err)
		panic(err)
	}

	taggedStruct, err := codeFile.FindTaggedStruct("//data")
	if err != nil {
		err = errs.New(err)
		panic(err)
	}
	js, sourceFilePath, err := JsonMarshal(taggedStruct[0])
	if err != nil {
		err = errs.New(err)
		panic(err)
	}
	fmt.Println("sourceFilePath: ", sourceFilePath)
	fmt.Println(string(js))
}
