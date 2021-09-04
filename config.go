package configuares

import (
	"errors"
	"fmt"
	"github.com/aacfactory/json"
	"github.com/tidwall/gjson"
	"reflect"
)

type config struct {
	raw []byte
}

func (config *config) As(v interface{}) (err error) {
	switch v.(type) {
	case *Raw:
		p := v.(*Raw)
		*p = append((*p)[0:0], config.raw...)
	case *json.RawMessage:
		p := v.(*json.RawMessage)
		*p = append((*p)[0:0], config.raw...)
	case *[]byte:
		p := v.(*[]byte)
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
	case *Raw:
		p := v.(*Raw)
		*p = append((*p)[0:0], config.raw...)
	case *json.RawMessage:
		p := v.(*json.RawMessage)
		*p = append((*p)[0:0], config.raw...)
	case *[]byte:
		p := v.(*[]byte)
		*p = append((*p)[0:0], config.raw...)
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

func (r Raw) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}
	if !json.Validate(r) {
		return nil, errors.New("configuares.Raw: MarshalJSON on invalid message")
	}
	return r, nil
}

func (r *Raw) UnmarshalJSON(data []byte) error {
	if r == nil {
		return errors.New("configuares.Raw: UnmarshalJSON on nil pointer")
	}
	if !json.Validate(data) {
		return errors.New("configuares.Raw: UnmarshalJSON on invalid message")
	}
	*r = append((*r)[0:0], data...)
	return nil
}
