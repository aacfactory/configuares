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
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"strings"
)

type RetrieverOption struct {
	Active string
	Format string
	Store  Store
}

func NewRetriever(option RetrieverOption) (retriever *Retriever, err error) {
	format := strings.ToUpper(strings.TrimSpace(option.Format))
	if format == "" || !(format == "JSON" || format == "YAML") {
		err = fmt.Errorf("create config retriever failed, format is not support")
		return
	}
	store := option.Store
	if store == nil {
		err = fmt.Errorf("create config retriever failed, store is nil")
		return
	}
	retriever = &Retriever{
		active: strings.ToUpper(strings.TrimSpace(option.Active)),
		format: format,
		store:  store,
	}
	return
}

type Retriever struct {
	active string
	format string
	store  Store
}

func (retriever *Retriever) Get() (config Config, err error) {
	root, subs, readErr := retriever.store.Read()
	if readErr != nil {
		err = fmt.Errorf("config retriever get failed, %v", readErr)
		return
	}
	if root == nil || len(root) == 0 {
		err = fmt.Errorf("config retriever get failed, not found")
		return
	}
	if retriever.active == "" {
		if retriever.format == "JSON" {
			if !gjson.ValidBytes(root) {
				err = fmt.Errorf(" config retriever get failed, bad json content")
				return
			}
			config = &JsonConfig{
				raw: root,
			}
		} else if retriever.format == "YAML" {
			_, validErr := yaml.YAMLToJSON(root)
			if validErr != nil {
				err = fmt.Errorf("config retriever get failed, bad yaml content")
				return
			}
			config = &YamlConfig{
				raw: root,
			}
		} else {
			err = fmt.Errorf("config retriever get failed, format is unsupported")
			return
		}
	}

	if subs == nil || len(subs) == 0 {
		err = fmt.Errorf("config retriever get failed, ative(%s) is not found", retriever.active)
		return
	}

	sub, hasSub := subs[retriever.active]
	if !hasSub {
		err = fmt.Errorf("config retriever get failed, ative(%s) is not found", retriever.active)
		return
	}

	mergedConfig, mergeErr := retriever.merge(retriever.format, root, sub)
	if mergeErr != nil {
		err = fmt.Errorf("config retriever get failed, merge ative failed %v", mergeErr)
		return
	}

	config = mergedConfig
	return
}

func (retriever *Retriever) merge(format string, root []byte, sub []byte) (config Config, err error) {
	if format == "JSON" {
		config, err = retriever.mergeJson(root, sub)
	} else if format == "YAML" {
		config, err = retriever.mergeYaml(root, sub)
	} else {
		err = fmt.Errorf("format is unsupported")
		return
	}
	return
}

func (retriever *Retriever) mergeJson(root []byte, sub []byte) (config Config, err error) {
	if !gjson.ValidBytes(root) {
		err = fmt.Errorf("merge failed, bad json content")
		return
	}
	if !gjson.ValidBytes(sub) {
		err = fmt.Errorf("merge failed, bad json content")
		return
	}
	subResult := gjson.ParseBytes(sub)
	subResult.ForEach(func(key gjson.Result, value gjson.Result) bool {
		root0, setErr := sjson.SetRawBytes(root, key.String(), []byte(value.Raw))
		if setErr != nil {
			return false
		}
		root = root0
		return true
	})
	config = &JsonConfig{
		raw: root,
	}
	return
}

func (retriever Retriever) mergeYaml(root []byte, sub []byte) (config Config, err error) {
	rootJson, rootToJsonErr := yaml.YAMLToJSON(root)
	if rootToJsonErr != nil {
		err = fmt.Errorf("merge failed, content format is not supported, %v", rootToJsonErr)
		return
	}
	subJson, subToJsonErr := yaml.YAMLToJSON(sub)
	if subToJsonErr != nil {
		err = fmt.Errorf("merge failed, content format is not supported, %v", subToJsonErr)
		return
	}
	jsonConfig, mergeJsonErr := retriever.mergeJson(rootJson, subJson)
	if mergeJsonErr != nil {
		err = mergeJsonErr
		return
	}
	yamlContent, toYamlErr := yaml.JSONToYAML(jsonConfig.Raw())
	if toYamlErr != nil {
		err = fmt.Errorf("merge failed, transfer failed, %v", toYamlErr)
		return
	}
	config = &YamlConfig{
		raw: yamlContent,
	}
	return
}
