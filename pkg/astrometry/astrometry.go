package astrometry /* import "r2discover.com/go/eaa-toolbox/pkg/astrometry" */

// Package to allow the use of the command line utilities from Astrometry.net

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

var mutex = &sync.Mutex{}

type Solver struct {
	cmd            *exec.Cmd
	buff           *bytes.Buffer
	cancelFilename string
	RA, Dec        float64
	FieldSize      string
}

func NewSolver() *Solver {
	solver := Solver{}
	solver.buff = &bytes.Buffer{}
	return &solver
}

func (solver *Solver) StopSolver() {
	err := ioutil.WriteFile(solver.cancelFilename, []byte{0}, 0600)
	if err != nil {
		log.Printf("StopSolver: %v", err)
	}
	for i := 0; i < 200; i++ {
		if solver.cmd == nil {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	mutex.Lock()
	solver.clearCmd()
	mutex.Unlock()
}

func (solver *Solver) StartFieldSolver(filename string) (err error) {
	if solver.IsSolving() {
		solver.StopSolver()
	}

	solver.cancelFilename = filename + ".cancel"
	os.Remove(solver.cancelFilename)
	solver.cmd = exec.Command("/usr/local/astrometry/bin/solve-field", "--no-plots", "-L", "0.1", "-H", "5", "-z", "2", "--objs", "50", "--cancel", solver.cancelFilename, "--overwrite", filename)
	solver.cmd.Stdout = solver.buff
	solver.cmd.Stderr = solver.buff
	if err = solver.cmd.Start(); err != nil {
		log.Printf("astrometry.StartFieldSolver: %v", err)
		return err
	}

	go solver.wait()

	return nil
}

func (solver *Solver) IsSolving() bool {
	return solver.cmd != nil
}

func (solver *Solver) clearCmd() {
	solver.cmd = nil
	solver.buff.Reset()
}

func (solver *Solver) wait() {
	if solver.cmd == nil {
		return
	}
	err := solver.cmd.Wait()
	if err != nil {
		log.Printf("astrometry.wait: %v", err)
		mutex.Lock()
		solver.clearCmd()
		mutex.Unlock()
		return
	}

	solver.processOutput(solver.buff.String())
	mutex.Lock()
	solver.clearCmd()
	mutex.Unlock()
}

func (solver *Solver) processOutput(output string) {
	s1 := "Field center: (RA,Dec) = ("
	for _, line := range strings.Split(output, "\n") {
		//log.Printf("line: >%v<", line)
		if strings.HasPrefix(line, s1) {
			line = line[len(s1):]
			i := strings.Index(line, ", ")
			if i != -1 {
				solver.RA, _ = strconv.ParseFloat(strings.TrimSpace(line[:i]), 64)
				line = line[i+1:]
				i := strings.Index(line, ")")
				if i != -1 {
					solver.Dec, _ = strconv.ParseFloat(strings.TrimSpace(line[:i]), 64)
				}
			}
		} else if strings.HasPrefix(line, "Field size:") {
			solver.FieldSize = strings.TrimSpace(line)
		}
	}
	log.Printf("RA %v Dec %v", solver.RA, solver.Dec)
}
