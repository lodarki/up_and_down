package entity

import (
	"encoding/json"
)

type ApiResult struct {
	Code    int
	Data    interface{}
	Message string
}

func (c *ApiResult) ParseData(d interface{}) error {
	bytes, marshalE := json.Marshal(c.Data)
	if marshalE != nil {
		return marshalE
	}
	return json.Unmarshal(bytes, d)
}
