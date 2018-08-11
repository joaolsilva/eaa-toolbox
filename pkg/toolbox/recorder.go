/*
Copyright 2018 EAA Toolbox Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package toolbox /* import "r2discover.com/go/eaa-toolbox/pkg/toolbox" */

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type Recorder struct {
	isRecording bool
	frameNumber int64
	directory   string
}

func (recorder *Recorder) Initialize(baseDir string) {
	recorder.frameNumber = 0
	recorder.directory = ""
	var i int64
	for i = 0; i < 100000; i++ {
		dir := path.Join(expandHomeDir(baseDir), fmt.Sprintf("eaa-%05d", i))
		if fileExists(dir) {
			continue
		}
		os.Mkdir(dir, 0755)
		recorder.directory = dir
		break
	}
}

func (recorder *Recorder) SavePNG(pngData *[]byte) {
	filename := path.Join(recorder.directory, fmt.Sprintf("eaa-%08d.png", recorder.frameNumber))
	err := ioutil.WriteFile(filename, *pngData, 0644)
	if err != nil {
		log.Printf("recorder.SavePNG: %v", err)
	}
	recorder.frameNumber += 1
}
