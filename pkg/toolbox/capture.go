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
	"gocv.io/x/gocv"
	"log"
)

func (toolbox *Toolbox) capture() {
	for {
		toolbox.webcam.Read(toolbox.img)

		buffer, _ := gocv.IMEncode(".png", *toolbox.img)
		toolbox.imgAsPNG = &buffer
		if toolbox.recorder.isRecording {
			toolbox.recorder.SavePNG(toolbox.imgAsPNG)
		}
	}
}

func (toolbox *Toolbox) startCapture() {
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		log.Printf("toolbox.startCapture: %v", err)
		return
	}
	toolbox.hasCamera = false
	toolbox.webcam = webcam
	if err == nil && toolbox.webcam != nil && toolbox.webcam.IsOpened() {
		toolbox.hasCamera = true
	}
	if !toolbox.hasCamera {
		log.Print("toolbox.startCapture: No camera found")
		return
	}

	if ok := webcam.Read(toolbox.img); !ok {
		log.Printf("Could not read from capture device")
		return
	}
	log.Printf("VideoCaptureFPS = %v", webcam.Get(gocv.VideoCaptureFPS))
	go toolbox.capture()
}
