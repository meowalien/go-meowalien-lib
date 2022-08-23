package urls

import (
	"github.com/meowalien/go-meowalien-lib/errs"
	"net/url"
)

func AddKeyValues(rawURL string, values ...[2]string) (ansURL string, err error) {
	uu, err := url.Parse(rawURL)
	if err != nil {
		err = errs.New(err)
		return "", err
	}

	qq := uu.Query()
	for _, value := range values {
		qq.Add(value[0], value[1])
	}
	uu.RawQuery = qq.Encode()

	return uu.String(), nil
}

func Join(uu ...string) (u *url.URL, err error) {
	u = &url.URL{}
	for _, s := range uu {
		u, err = u.Parse(s)
		if err != nil {
			err = errs.New(err)
			return
		}
	}
	return
}
