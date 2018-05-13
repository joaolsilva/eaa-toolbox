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
	ShouldQuit    bool
	serverStarted bool
	screen        Screen
	window        *gocv.Window
	writer        *gocv.VideoWriter
	fps           float64
	img           *gocv.Mat
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

	return &toolbox
}

func (toolbox *Toolbox) Start() {
	go toolbox.startServer()
	err := webview.Open("EAA Toolbox", "http://127.0.0.1:32243/index.html", 800, 600, true)
	if err != nil {
		log.Printf("%v", err)
		return
	}
}
