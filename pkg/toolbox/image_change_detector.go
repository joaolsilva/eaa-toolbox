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
	"math"
)

var previousImg *gocv.Mat // Gray scale version of the previous image
var currentImg *gocv.Mat  // Gray scale version of the current image
var deltaImg *gocv.Mat
var normHistory []float64 = []float64{}

const maxNormHistorySize = 154

func normChanged() bool {
	if len(normHistory) <= 2 {
		return true
	}

	n := normHistory[len(normHistory)-1]
	n1 := normHistory[len(normHistory)-2]
	n2 := normHistory[len(normHistory)-3]

	r1 := math.Abs(n-n1) / n
	r2 := math.Abs(n-n2) / n
	r := math.Abs(r1-r2) / r1

	//log.Printf("n %v n1 %v n2 %v r1 %v r2 %v r %v", n, n1, n2, r1, r2, r)

	return r < 0.04
}

func ImageChanged(img *gocv.Mat) bool {
	if deltaImg == nil {
		m := img.Clone()
		deltaImg = &m
	}

	if currentImg == nil {
		m := img.Clone()
		currentImg = &m
	} else {
		img.CopyTo(*currentImg)
	}
	gocv.CvtColor(*currentImg, currentImg, gocv.ColorRGBToGray)

	if previousImg == nil {
		m := currentImg.Clone()
		previousImg = &m
	}

	gocv.Subtract(*previousImg, *currentImg, deltaImg)
	currentImg.CopyTo(*previousImg)

	norm := gocv.Norm(*deltaImg, gocv.NormL2)

	if norm == 0 {
		return false
	}

	if len(normHistory) >= maxNormHistorySize {
		normHistory = normHistory[len(normHistory)-maxNormHistorySize+1:]
	}
	normHistory = append(normHistory, norm)
	changed := normChanged()
	//log.Printf("ImageChanged, norm: %v Changed? %v", norm, changed)

	return changed
}
