package convert

import (
	"encoding/json"
	"fmt"
	go_meowalien_lib "github.com/meowalien/go-meowalien-lib"
	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"net/url"
)

func ConvertMapToURLForm(postMap map[string]interface{}) url.Values {
	form := url.Values{}
	for k, v := range postMap {
		form.Set(k, fmt.Sprintf("%v", v))
	}
	return form
}
func MapstructureOnJsonTag(input interface{}, i interface{}) (err error) {
	return MapstructureOnTag(input, "json", i)
}

func MapstructureOnFormTag(input interface{}, i interface{}) (err error) {
	return MapstructureOnTag(input, "form", i)
}


func MapstructureOnTag(input interface{}, tag string, i interface{}) (err error) {
	jsonDecoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  i,
		TagName: tag,
	})
	if err != nil {
		err = errs.WithLine(err)
		return
	}

	err = jsonDecoder.Decode(input)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	return
}



func DecodeJsonResponseToStruct(res *http.Response, i interface{}) (err error) {
	if go_meowalien_lib.SHOW_DEBUG_MESSAGE{
		var all []byte
		all, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return errs.WithLine(err)
		}

		fmt.Println("response: ",string(all))

		err = json.Unmarshal(all, i)
	}else {
		decoder := json.NewDecoder(res.Body)
		err = decoder.Decode(i)
		if err != nil {
			err = errs.WithLine(err)
			return
		}
	}
	return
}
