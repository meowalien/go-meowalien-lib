package files

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"io"
	"io/ioutil"
)

func PrintFile(filePath string, w io.Writer) (err error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		err = errs.New(err)
		return
	}
	_, err = fmt.Fprintf(w, "%s:\n%s\n", filePath, string(data))
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}
