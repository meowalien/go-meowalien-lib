package uuid

import (
	"fmt"
	"github.com/meowalien/go-meowalien-lib/random"
)

func NewUUID(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, random.Snowflake().Base58())
}
