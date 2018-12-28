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

package ser /* import "r2discover.com/go/eaa-toolbox/pkg/ser" */

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"log"
)

type ColorID int32
type Endianness int32

const (
	ColorID_MONO       ColorID = 0
	ColorID_BAYER_RGGB ColorID = 8
	ColorID_BAYER_GRBG ColorID = 9
	ColorID_BAYER_GBRG ColorID = 10
	ColorID_BAYER_BGGR ColorID = 11
	ColorID_BAYER_CYYM ColorID = 16
	ColorID_BAYER_YCMY ColorID = 17
	ColorID_BAYER_YMCY ColorID = 18
	ColorID_BAYER_MYYC ColorID = 19
	ColorID_RGB        ColorID = 100
	ColorID_BGR        ColorID = 101

	BigEndian    Endianness = 0
	LittleEndian Endianness = 1
)

var FileSignature [14]byte

type Header struct {
	FileID             [14]byte // File Signature
	LuID               int32
	ColorID            ColorID
	Endianness         Endianness
	ImageWidth         int32
	ImageHeight        int32
	PixelDepthPerPlane int32
	FrameCount         int32
	Observer           [40]byte // Name of observer
	Instrument         [40]byte // Name of camera
	Telescope          [40]byte // Name of telescope
	DateTime           DateTime // Local time
	DateTimeUTC        DateTime // Time (UTC)
}

type SER struct {
	Header   Header
	Filename string
}

func New(filename string, imageWidth int32, imageHeight int32) *SER {
	s := SER{Filename: filename}
	s.Header = Header{}
	s.Header.FileID = FileSignature
	s.Header.ColorID = ColorID_RGB
	s.Header.Endianness = LittleEndian
	s.Header.ImageWidth = imageWidth
	s.Header.ImageHeight = imageHeight
	s.Header.PixelDepthPerPlane = 3

	s.createFile()

	return &s
}

func init() {
	copy(FileSignature[:], "LUNACAM-RECORDER")
}

func (ser *SER) generateHeader() []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, ser.Header.FileID)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.LuID)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.ColorID)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.Endianness)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.ImageWidth)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.ImageHeight)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.PixelDepthPerPlane)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.FrameCount)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.Observer)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.Instrument)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.Telescope)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.DateTime)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	err = binary.Write(buf, binary.LittleEndian, ser.Header.DateTimeUTC)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	return buf.Bytes()
}

func (ser *SER) createFile() {
	ioutil.WriteFile(ser.Filename, ser.generateHeader(), 0640)
}
