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

package configures_test

import (
	"fmt"
	"github.com/aacfactory/configures"
	"github.com/aacfactory/json"
	"path/filepath"
	"testing"
	"time"
)

func Test_JsonConfig(t *testing.T) {

	path, err := filepath.Abs("./_example/json")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("config dir", path)

	store := configures.NewFileStore(path, "app", '.')

	retriever, retrieverErr := configures.NewRetriever(configures.RetrieverOption{
		Active: "dev",
		Format: "JSON",
		Store:  store,
	})

	if retrieverErr != nil {
		t.Error(retrieverErr)
		return
	}

	config, configErr := retriever.Get()
	if configErr != nil {
		t.Error(configErr)
		return
	}

	raw := json.RawMessage{}
	_ = config.As(&raw)

	t.Log(string(raw))

	node, _ := config.Node("http")
	t.Log(string(node.Raw()))
	http := Http{}
	has, getErr := config.Get("http", &http)
	t.Log(fmt.Sprintf("%+v", http), has, getErr)
}

func Test_YamlConfig(t *testing.T) {

	path, err := filepath.Abs("./_example/yaml")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("config dir", path)

	store := configures.NewFileStore(path, "app", '-')

	retriever, retrieverErr := configures.NewRetriever(configures.RetrieverOption{
		Active: "prod",
		Format: "YAML",
		Store:  store,
	})

	if retrieverErr != nil {
		t.Error(retrieverErr)
		return
	}

	config, configErr := retriever.Get()
	if configErr != nil {
		t.Error(configErr)
		return
	}

	raw := json.RawMessage{}
	_ = config.As(&raw)

	t.Log(string(raw))

	http := Http{}
	has, getErr := config.Get("http", &http)
	t.Log(fmt.Sprintf("%+v", http), has, getErr)
}

type Http struct {
	Timeout time.Duration `json:"timeout"`
	Port    int           `json:"port"`
	Host    string        `json:"host"`
}
