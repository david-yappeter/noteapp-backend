package service

import "github.com/vektah/gqlparser/v2/gqlerror"

type updateArgs struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func gqlError(message string, extCode string, extVal interface{}) error {
	return &gqlerror.Error{
		Message: message,
		Extensions: map[string]interface{}{
			extCode: extVal,
		},
	}
}
