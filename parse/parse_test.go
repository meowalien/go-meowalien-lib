package parse

import (
	"encoding/json"
	"fmt"
	"github.com/meowalien/go-meowalien-lib/errs"
	"strings"
	"testing"
)

func TestParseNestedMap(t *testing.T) {
	//type Root struct {
	//	StatusCode int `json:"ErrorCode,omitempty"`
	//	Data       *struct {
	//		RoomData *struct {
	//			GameID int `json:"GameID"`
	//		} `json:"RoomData"`
	//		RoutineData *struct {
	//			PartyID int `json:"PartyID"`
	//		} `json:"RoutineData"`
	//	} `json:"Data"`
	//	TraceId string `json:"RequestID"`
	//}
	jsonStr := "{\n  \"ErrorCode\": 4,\n  \"Data\": {\n    \"RoomData\": {\n      \"GameID\": 1\n    },\n    \"RoutineData\":{\n        \"PartyID\": 2\n      }\n  },\n  \"RequestID\": \"123\"\n}"

	m := make(map[string]interface{})
	dc := json.NewDecoder(strings.NewReader(jsonStr))
	dc.UseNumber()
	err := dc.Decode(&m)
	if err != nil {
		err = errs.New(err)
		panic(err)
		return
	}
	//v, err := ParseNestedMap(jsonStr, "Data", "RoutineData", "PartyID")
	//v, err := ParseNestedMap(jsonStr, "Data", "RoomData", "GameID")
	//v, err := ParseNestedMap(jsonStr, "RequestID")
	v, _, err := ParseNestedMap(m, "ErrorCode")
	if err != nil {
		err = errs.New(err)
		panic(err)
	}
	fmt.Println("v: ", v.(int64))
}
