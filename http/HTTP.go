package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/meowalien/go-meowalien-lib/format/convert"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func ConvertFormRequestToMap(r *http.Request) (m map[string]interface{}) {
	cMap := make(map[string]interface{})
	switch r.Method {
	case "GET":
		query := r.URL.Query()
		for k, v := range query {
			cMap[k] = v[0]
		}

	case "POST":
		r.ParseForm()
		for k, v := range r.Form {
			cMap[k] = v[0]
		}
	}
	return cMap
}

// 發送urlencodedFORM
func DoURLEncodedFormRequest(endpoint string, req map[string]interface{}) ([]byte, error) {
	data := url.Values{}
	for s, i := range req {
		data.Add(s, fmt.Sprintf("%v", i))
	}

	dataEncode := data.Encode()

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(dataEncode)) // URL-encoded payload
	if err != nil {
		return nil, fmt.Errorf("error when NewRequest: %w", err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(dataEncode)))

	res, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error when client.Start: %w", err)
	}

	defer func(Body io.ReadCloser) {
		e := Body.Close()
		if e != nil {
			log.Printf("error when close Body: %s\n", e.Error())
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error StatusCode, res: %v", res.Body)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error when read Body: %w", err)
	}

	return body, nil
}

func JsonRequest(endpoint string, req interface{}) ([]byte, error) {
	jj, err := json.Marshal(req)
	if err != nil {
		return nil, errs.WithLine(err)
	}

	r, err := http.NewRequest("POST", endpoint, bytes.NewReader(jj)) // URL-encoded payload
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(jj)))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, errs.WithLine(err)
	}

	defer func(Body io.ReadCloser) {
		e := Body.Close()
		if e != nil {
			log.Printf("error when close Body: %s\n", e.Error())
		}
	}(res.Body)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error when read Body: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error StatusCodes: %v, body:%s", res.StatusCode, string(body))
	}
	return body, nil
}

func CloseResponse(response *http.Response) error {
	defer response.Body.Close()
	_, err := io.Copy(ioutil.Discard, response.Body)
	return err
}

type PostForm interface {
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

var SHOW_DEBUG_MESSAGE = false

func PostURLFormWithClient(c PostForm, baseURL string, requestmap map[string]interface{}, rep interface{}) (err error) {
	if SHOW_DEBUG_MESSAGE {
		fmt.Println("baseURL: ", baseURL)
		fmt.Println(" ----- ")
		for k, v := range requestmap {
			fmt.Printf("%s:%v\n", k, v)
		}
		fmt.Println(" ----- ")
	}

	form := convert.ConvertMapToURLForm(requestmap)

	res, err := c.PostForm(baseURL, form)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	defer CloseResponse(res)

	if res.StatusCode == http.StatusNoContent {
		log.Println("StatusNoContent ...")
		return
	} else if res.StatusCode != 200 {
		var all []byte
		all, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = errs.WithLine("http response code: %d , rep: %v", res.StatusCode, string(all))
		return
	}

	err = convert.DecodeJsonResponseToStruct(res, rep)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	return
}

type GetForm interface {
	Get(url string) (resp *http.Response, err error)
}

func GetURLFormWithClient(c GetForm, baseURL string, requestmap map[string]interface{}, rep interface{}) (err error) {
	uu, err := url.Parse(baseURL)
	qq := uu.Query()
	for key, value := range requestmap {
		qq.Add(key, fmt.Sprint(value))
	}
	uu.RawQuery = qq.Encode()
	fmt.Println("get url: ", uu.String())
	res, err := c.Get(uu.String())
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	defer CloseResponse(res)

	if res.StatusCode != 200 {
		var all []byte
		all, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = errs.WithLine("http response code: %d , rep: %v", res.StatusCode, string(all))
		return
	}

	err = convert.DecodeJsonResponseToStruct(res, rep)
	if err != nil {
		err = errs.WithLine(err)
		return
	}
	return
}
