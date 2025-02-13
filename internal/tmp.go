/*
 *    Copyright 2024 Han Li and contributors
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/version-fox/vfox/internal/util"
)

type Temp struct {
	dirPath        string
	CurProcessPath string
}

func (t *Temp) Remove() {
	dir, err := os.ReadDir(t.dirPath)
	if err == nil {
		_ = os.RemoveAll(t.CurProcessPath)
		for _, file := range dir {
			if !file.IsDir() {
				continue
			}
			names := strings.SplitN(file.Name(), "-", 2)
			if len(names) != 2 {
				continue
			}
			timestamp := names[0]
			i, err := strconv.ParseInt(timestamp, 10, 64)
			if err != nil {
				continue
			}
			if util.IsBeforeToday(i) {
				_ = os.Remove(filepath.Join(t.dirPath, file.Name()))
			}
		}
	}
}

func NewTemp(dirPath string, pid int) (*Temp, error) {
	timestamp := util.GetBeginOfToday()
	name := fmt.Sprintf("%d-%d", timestamp, pid)
	path := filepath.Join(dirPath, name)
	if !util.FileExists(path) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return nil, fmt.Errorf("create temp dir failed: %w", err)
		}
	}
	t := &Temp{
		dirPath:        dirPath,
		CurProcessPath: path,
	}
	return t, nil
}
