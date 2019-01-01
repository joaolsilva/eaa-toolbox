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
	"errors"
	"gocv.io/x/gocv"
)

var runningAverage *gocv.Mat
var n int64

func AddToRunningAverage(img *gocv.Mat) {
	if runningAverage == nil {
		ResetRunningAverage(img)
		return
	}

	n += 1

	// avg_n = (avg_(n-1) * (n-1) + img ) / n

	gocv.AddWeighted(*runningAverage, float64(n-1.0)/float64(n), *img, float64(1.0)/float64(n), 1.0, runningAverage)
}

func ResetRunningAverage(img *gocv.Mat) {
	if runningAverage == nil {
		m := img.Clone()
		runningAverage = &m
	} else {
		img.CopyTo(runningAverage)
	}
	n = 1
}

func GetRunningAverage() (gocv.Mat, error) {
	if runningAverage == nil {
		return gocv.NewMat(), errors.New("No data")
	}
	return runningAverage.Clone(), nil
}
