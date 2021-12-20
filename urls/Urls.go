package urls

import (
	"fmt"
	"net/url"
)

func AddValues(rawURL string, values ...[2]string )(ansURL string , err error)  {
	uu , err := url.Parse(rawURL)
	if err != nil{
		return "" , fmt.Errorf("error when Parse: ",err.Error())
	}

	qq:=uu.Query()
	for _, value := range values {
		qq.Add(value[0] , value[1])
	}
	uu.RawQuery = qq.Encode()

	return uu.String() , nil
}