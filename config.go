package configuares

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"reflect"
)

type config struct {
	raw []byte
}

func (config *config) As(v interface{}) (err error) {
	switch v.(type) {
	case *Raw, *json.RawMessage, *[]byte:
		p := v.(*Raw)
		*p = append((*p)[0:0], config.raw...)
	default:
		decodeErr := json.Unmarshal(config.raw, v)
		if decodeErr != nil {
			err = fmt.Errorf("decode config as %v failed", reflect.TypeOf(v))
			return
		}
	}
	return
}

func (config *config) Get(path string, v interface{}) (has bool, err error) {
	result := gjson.GetBytes(config.raw, path)
	if !result.Exists() {
		return
	}
	switch v.(type) {
	case *Raw, *json.RawMessage, *[]byte:
		p := v.(*Raw)
		*p = append((*p)[0:0], result.Raw...)
	default:
		decodeErr := json.Unmarshal([]byte(result.Raw), v)
		if decodeErr != nil {
			err = fmt.Errorf("config get %s failed for decoding failed", path)
			return
		}
	}
	has = true
	return
}

type Raw json.RawMessage

func (r Raw) As(v interface{}) (err error) {
	decodeErr := json.Unmarshal(r, v)
	if decodeErr != nil {
		err = fmt.Errorf("raw config as %s failed", reflect.TypeOf(v).String())
	}
	return
}
