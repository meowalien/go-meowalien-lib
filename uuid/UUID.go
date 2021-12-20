package uuid

import (
	"core1/src/pkg/meowalien_lib/random"
	"fmt"
)

func NewUUID(prefix string) string {
	return fmt.Sprintf("%s%s",prefix, random.Snowflake().Base58())
}