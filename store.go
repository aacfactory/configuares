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

package configures

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func NewFileStore(configPath string, prefix string, splitter byte) Store {
	return &FileStore{
		configPath: configPath,
		prefix:     prefix,
		splitter:   splitter,
	}
}

type FileStore struct {
	configPath string
	prefix     string
	splitter   byte
}

func (store *FileStore) Read() (root []byte, subs map[string][]byte, err error) {
	file, openErr := os.Open(store.configPath)
	if openErr != nil {
		err = fmt.Errorf("config file store open %s failed, %v", store.configPath, openErr)
		return
	}
	fileStat, statErr := file.Stat()
	if statErr != nil {
		err = fmt.Errorf("config file store get %s file info failed, %v", store.configPath, statErr)
		return
	}
	if !fileStat.IsDir() {
		_ = file.Close()
		fileContent, readErr := os.ReadFile(store.configPath)
		if readErr != nil {
			err = fmt.Errorf("config file store read %s failed, %v", store.configPath, readErr)
			return
		}
		root = fileContent
		return
	}
	subs = make(map[string][]byte)
	dirErr := filepath.Walk(store.configPath, func(path string, info fs.FileInfo, cause error) (err error) {
		if info.IsDir() {
			return
		}
		filename := filepath.Base(path)
		if strings.Index(filename, store.prefix) != 0 {
			return
		}
		fileContent, readErr := os.ReadFile(path)
		if readErr != nil {
			err = fmt.Errorf("read %s failed, %v", path, readErr)
			return
		}
		key := filename[:strings.LastIndexByte(filename, '.')]
		idx := strings.IndexByte(key, store.splitter)
		if idx < 1 {
			root = fileContent
			return
		}
		key = key[idx+1:]
		subs[strings.ToUpper(strings.TrimSpace(key))] = fileContent
		return
	})
	if dirErr != nil {
		err = fmt.Errorf("config file store read %s dir failed, %v", store.configPath, dirErr)
		return
	}
	return
}
