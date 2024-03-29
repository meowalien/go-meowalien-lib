package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	go_meowalien_lib "github.com/meowalien/go-meowalien-lib"
	"github.com/meowalien/go-meowalien-lib/errs"
	"github.com/meowalien/go-meowalien-lib/format/convert"
)

func ConvertRequestToMap(r *http.Request) (cMap map[string]interface{}, err error) {
	cMap = make(map[string]interface{})
	switch r.Method {
	case "GET":
		query := r.URL.Query()
		for k, v := range query {
			cMap[k] = v[0]
		}
	case "POST":
		err = r.ParseForm()
		for k, v := range r.Form {
			cMap[k] = v[0]
		}
	}
	return
}

// 發送urlencodedFORM
func DoURLEncodedFormRequest(endpoint string, method string, req map[string]interface{}) ([]byte, error) {
	data := url.Values{}
	for s, i := range req {
		data.Add(s, fmt.Sprintf("%v", i))
	}

	dataEncode := data.Encode()

	client := &http.Client{}
	r, err := http.NewRequestWithContext(context.TODO(), method, endpoint, strings.NewReader(dataEncode)) // URL-encoded payload
	if err != nil {
		return nil, fmt.Errorf("error when NewRequest: %w", err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(dataEncode)))

	res, err := client.Do(r) //nolint:bodyclose
	if err != nil {
		return nil, fmt.Errorf("error when client.Start: %w", err)
	}

	defer func(Body io.Closer) {
		e := Body.Close()
		if e != nil {
			log.Printf("error when close Body: %s\n", e.Error())
		}
	}(res.Body)
	if res.StatusCode != http.StatusOK {
		var body []byte
		body, err = ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(errs.New("error when read Body: %w", err).Error())
		}

		return nil, fmt.Errorf("error StatusCode: %d, res: %v", res.StatusCode, string(body))
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
		return nil, errs.New(err)
	}

	r, err := http.NewRequestWithContext(context.TODO(), "POST", endpoint, bytes.NewReader(jj))
	if err != nil {
		return nil, errs.New(err)
	}
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Content-Length", strconv.Itoa(len(jj)))

	client := &http.Client{}
	res, err := client.Do(r) //nolint:bodyclose
	if err != nil {
		return nil, errs.New(err)
	}

	defer func(Body io.Closer) {
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

func closeAndDrainResponseBody(response *http.Response) (err error) {
	defer func(body io.Closer) {
		err1 := body.Close()
		if err1 != nil {
			err = errs.New(err, err1)
		}
	}(response.Body)
	_, err = io.Copy(ioutil.Discard, response.Body)
	if err != nil {
		log.Println(errs.New(err).Error())
	}
	return
}

type PostForm interface {
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}

func PostFormWithClient(c PostForm, baseURL string, requestmap map[string]interface{}, rep interface{}) (err error) {
	if go_meowalien_lib.SHOW_DEBUG_MESSAGE {
		fmt.Println("baseURL: ", baseURL)
		fmt.Println(" ----- ")
		for k, v := range requestmap {
			fmt.Printf("%s:%v\n", k, v)
		}
		fmt.Println(" ----- ")
	}

	form := convert.ConvertMapToURLForm(requestmap)

	res, err := c.PostForm(baseURL, form) //nolint:bodyclose
	if err != nil {
		err = errs.New(err)
		return
	}
	defer func(res *http.Response) {
		err1 := closeAndDrainResponseBody(res)
		if err1 != nil {
			err = errs.New(err, err1)
		}
	}(res)

	if res.StatusCode == http.StatusNoContent {
		log.Println("StatusNoContent ...")
		return
	} else if res.StatusCode != 200 {
		var all []byte
		all, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = errs.New("http response code: %d , rep: %v", res.StatusCode, string(all))
		return
	}

	err = convert.DecodeJsonResponseToStruct(res, rep)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}

type GetForm interface {
	Get(url string) (resp *http.Response, err error)
}

func GetFormWithClient(c GetForm, baseURL string, requestmap map[string]interface{}, rep interface{}) (err error) {
	uu, err := url.Parse(baseURL)
	qq := uu.Query()
	for key, value := range requestmap {
		qq.Add(key, fmt.Sprint(value))
	}
	uu.RawQuery = qq.Encode()
	//fmt.Println("get url: ", uu.String())
	res, err := c.Get(uu.String()) //nolint:bodyclose
	if err != nil {
		err = errs.New(err)
		return
	}
	defer func(res *http.Response) {
		err1 := closeAndDrainResponseBody(res)
		if err1 != nil {
			err = errs.New(err, err1)
		}
	}(res)
	if res.StatusCode != 200 {
		var all []byte
		all, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = errs.New("fail when GET %s , http response code: %d , rep: %v", uu.String(), res.StatusCode, string(all))
		return
	}

	err = convert.DecodeJsonResponseToStruct(res, rep)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}

type PostBody interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}

func PostJsonWithClient(c PostBody, baseURL string, request interface{}, rep interface{}) (err error) {
	if go_meowalien_lib.SHOW_DEBUG_MESSAGE {
		fmt.Println("baseURL: ", baseURL)
		fmt.Println(" ----- ")
		fmt.Printf("%+v\n", request)
		fmt.Println(" ----- ")
	}

	b, err := json.Marshal(request)

	res, err := c.Post(baseURL, "application/json", bytes.NewReader(b)) //nolint:bodyclose
	if err != nil {
		err = errs.New(err)
		return
	}

	defer func(res *http.Response) {
		err1 := closeAndDrainResponseBody(res)
		if err1 != nil {
			err = errs.New(err, err1)
		}
	}(res)

	if res.StatusCode == http.StatusNoContent {
		log.Println("StatusNoContent ...")
		return
	} else if res.StatusCode != 200 {
		var all []byte
		all, err = ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		err = errs.New("http response code: %d , rep: %v", res.StatusCode, string(all))
		return
	}

	err = convert.DecodeJsonResponseToStruct(res, rep)
	if err != nil {
		err = errs.New(err)
		return
	}
	return
}
