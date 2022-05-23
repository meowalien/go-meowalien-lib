package code

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/runtime"
	"log"
)

type Option struct {
	UseLanguage             string
	DefaultLanguage         string
	CodeStringMapping       map[string]map[Code]string
	UndefinedCodeText       string
	StatusCodeMagnification int
}

var option = Option{
	UseLanguage:             "en",
	DefaultLanguage:         "en",
	CodeStringMapping:       nil,
	UndefinedCodeText:       "undefined code",
	StatusCodeMagnification: 1,
}

func SetOption(op Option) {
	option = op
}

type Context interface {
	AbortWithStatus(code int)
	AbortWithStatusJSON(code int, body interface{})
}

type StatusCode interface {
	error
	String() string
	ResponseAs(c Context, resp ResponseStructStyle, data ...interface{})
	HTTPCode() int
	MyCode() int
	WithMsg(s string) StatusCode
	responseStruct() ResponseStructStyle
}

type ResponseStructStyle func(d StatusCode, c Context, data ...interface{})

type Code int

func (s Code) Error() string {
	return s.String()
}

func (s Code) String() string {
	lg, ok := option.CodeStringMapping[option.UseLanguage]
	if !ok {
		// 無指定語言，使用預設語言
		lg = option.CodeStringMapping[option.DefaultLanguage]
	}

	text, exist := lg[s]
	if exist {
		return text
	} else {
		lg = option.CodeStringMapping[option.DefaultLanguage]
		text, exist = lg[s]
		if exist {
			return text
		}
	}

	return option.UndefinedCodeText
}

func (s Code) ResponseAs(c Context, resp ResponseStructStyle, data ...interface{}) {
	resp(s, c, data...)
}

func (s Code) HTTPCode() int {
	res := s.MyCode() / option.StatusCodeMagnification
	if res == 0 {
		return s.MyCode()
	}
	return res
}

func (s Code) MyCode() int {
	return int(s)
}

func (s Code) WithMsg(st string) StatusCode {
	return withMsg(s, st)
}

func (s Code) responseStruct() ResponseStructStyle {
	return nil
}

func withMsg(code Code, msg string) StatusCode {
	return CustomizedStatusCode{
		MyStatusCode:   code.MyCode(),
		HttpCode:       code.HTTPCode(),
		Message:        fmt.Sprint(code.String(), " ", msg),
		ResponseStruct: code.responseStruct(),
	}
}

type CustomizedStatusCode struct {
	Code
	MyStatusCode   int
	HttpCode       int
	Message        string
	ResponseStruct ResponseStructStyle
}

func (s CustomizedStatusCode) responseStruct() ResponseStructStyle {
	return s.ResponseStruct
}

func (s CustomizedStatusCode) String() string {
	return s.Message
}

func (s CustomizedStatusCode) HTTPCode() int {
	return s.HttpCode
}

func (s CustomizedStatusCode) MyCode() int {
	return s.MyStatusCode
}

func (s CustomizedStatusCode) Response(c Context, data ...interface{}) {
	if s.ResponseStruct == nil {
		log.Println(runtime.CallerFileAndLine(1), " - the ResponseStruct is  nil")
		return
	}
	s.ResponseAs(c, s.ResponseStruct, data...)
}
