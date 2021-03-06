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
	"encoding/json"
	"fmt"
	"github.com/soniakeys/meeus/v3/coord"
	"github.com/soniakeys/meeus/v3/globe"
	"github.com/soniakeys/meeus/v3/julian"
	"github.com/soniakeys/meeus/v3/sidereal"
	"github.com/soniakeys/unit"
	"github.com/zserge/webview"
	"gocv.io/x/gocv"
	"log"
	"os"
	"r2discover.com/go/eaa-toolbox/pkg/astrometry"
	"time"
)

type State struct {
	CurrentPosition string
	GoToPosition    string
	Message         string
}

type Toolbox struct {
	webcam        *gocv.VideoCapture
	hasCamera     bool
	ShouldQuit    bool
	hub           *Hub
	serverStarted bool
	recorder      *Recorder
	state         State
	fps           float64
	img           *gocv.Mat
	imgAsPNG      *[]byte
	solver        *astrometry.Solver
	goTo          struct {
		RA  float64
		Dec float64
	}
	appConfig AppConfig
}

const maxResults = 8

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func RADecToHoriz(RA, Dec, latitude, longitude float64) (h *coord.Horizontal) {

	eq := coord.Equatorial{RA: unit.RAFromDeg(RA), Dec: unit.AngleFromDeg(Dec)}
	h = &coord.Horizontal{}
	g := globe.Coord{Lat: unit.AngleFromDeg(latitude), Lon: unit.AngleFromDeg(-longitude)}
	h.EqToHz(&eq, &g, sidereal.Apparent(julian.TimeToJD(time.Now())))
	return h
}

func NewToolbox() *Toolbox {
	toolbox := Toolbox{}

	toolbox.solver = astrometry.NewSolver()
	toolbox.appConfig = LoadAppConfig()
	img := gocv.NewMat()
	toolbox.img = &img
	recorder := Recorder{}
	toolbox.recorder = &recorder

	return &toolbox
}

func (toolbox *Toolbox) Start() {
	log.Print("toolbox.Start()")
	go toolbox.startServer()

	w := webview.New(webview.Settings{
		Title:     "EAA Toolbox",
		URL:       "http://127.0.0.1:32243/index.html",
		Width:     800,
		Height:    600,
		Resizable: true,
		Debug:     toolbox.appConfig.Settings.Debug,
	})
	defer w.Exit()

	toolbox.startCapture()
	w.Run()
}

func (toolbox *Toolbox) processCommand(command string) {
	switch command {
	case "PLATE_SOLVING":
		gocv.IMWrite("/tmp/eaa.png", *toolbox.img)
		toolbox.solver.StartFieldSolver("/tmp/eaa.png")
		toolbox.hub.Broadcast(toolbox.stateJSON())
	case "START_RECORDING":
		if toolbox.recorder.isRecording {
			break
		}

		toolbox.recorder.Initialize(toolbox.appConfig.Paths.Recordings)
		toolbox.recorder.isRecording = true
		toolbox.state.Message = "Recording..."
		toolbox.hub.Broadcast(toolbox.stateJSON())
	case "STOP_RECORDING":
		if !toolbox.recorder.isRecording {
			break
		}

		toolbox.recorder.isRecording = false
		toolbox.state.Message = ""
		toolbox.hub.Broadcast(toolbox.stateJSON())
	case "REFRESH":

		if toolbox.solver.IsSolving() {
			toolbox.state.CurrentPosition = "Solving..."
		} else if toolbox.solver.RA != 0.0 || toolbox.solver.Dec != 0.0 {
			currentPosition := RADecToHoriz(toolbox.solver.RA, toolbox.solver.Dec, toolbox.appConfig.Location.Latitude, toolbox.appConfig.Location.Longitude)
			toolbox.state.CurrentPosition = fmt.Sprintf("Pos: RA %.3f Dec %.3f Alt %.1f Az %.1f", toolbox.solver.RA, toolbox.solver.Dec, currentPosition.Alt.Deg(), currentPosition.Az.Deg())
		} else {
			toolbox.state.CurrentPosition = "Solver failed"
		}

		toolbox.hub.Broadcast(toolbox.stateJSON())
	default:
		log.Printf("toolbox.processCommand: Unknown command: %v", command)
	}
}

func (toolbox *Toolbox) stateJSON() []byte {
	state, err := json.Marshal(toolbox.state)
	if err != nil {
		return []byte{}
	}
	return state
}
