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

package configuares_test

import (
	"encoding/json"
	"github.com/aacfactory/configuares"
	"path/filepath"
	"testing"
)

func Test_JsonConfig(t *testing.T) {

	path, err := filepath.Abs("./_example/json")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("config dir", path)

	store := configuares.NewFileStore(path, "app", '.')

	retriever, retrieverErr := configuares.NewRetriever(configuares.RetrieverOption{
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
}

func Test_YamlConfig(t *testing.T) {

	path, err := filepath.Abs("./_example/yaml")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("config dir", path)

	store := configuares.NewFileStore(path, "app", '-')

	retriever, retrieverErr := configuares.NewRetriever(configuares.RetrieverOption{
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

	http := configuares.Raw{}
	has, getErr := config.Get("http", &http)
	t.Log(string(http), has, getErr)
}
