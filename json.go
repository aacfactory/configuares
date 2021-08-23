/*
 * Copyright 2021 Wang Min Xiang
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package configuares

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
)

type JsonConfig struct {
	raw []byte
}

func (config *JsonConfig) As(v interface{}) (err error) {
	decodeErr := json.Unmarshal(config.raw, v)
	if decodeErr != nil {
		err = fmt.Errorf("decode config as %v failed, %s, %v", v, string(config.raw), decodeErr)
	}
	return
}

func (config *JsonConfig) Get(path string, v interface{}) (has bool, err error) {
	result := gjson.GetBytes(config.raw, path)
	if !result.Exists() {
		return
	}
	decodeErr := json.Unmarshal([]byte(result.Raw), v)
	if decodeErr != nil {
		err = fmt.Errorf("config get %s failed, %v", path, decodeErr)
		return
	}
	has = true
	return
}

func (config *JsonConfig) Raw() (raw []byte) {
	raw = config.raw
	return
}
