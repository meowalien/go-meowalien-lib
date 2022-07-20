package client

import (
	"net/http"
)

type ClientModifier func(c *http.Client)
