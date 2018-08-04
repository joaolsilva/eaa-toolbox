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

type Screen struct {
	Position string
	GoTo     string
}

type Toolbox struct {
	webcam        *gocv.VideoCapture
	hasCamera     bool
	ShouldQuit    bool
	serverStarted bool
	screen        Screen
	writer        *gocv.VideoWriter
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
