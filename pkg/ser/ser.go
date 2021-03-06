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
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

func (c ColorID) String() string {
	switch c {
	case ColorID_MONO:
		return "Mono"
	case ColorID_BAYER_RGGB:
		return "Bayer RGGB"
	case ColorID_BAYER_GRBG:
		return "Bayer GRBG"
	case ColorID_BAYER_GBRG:
		return "Bayer GRBG"
	case ColorID_BAYER_BGGR:
		return "Bayer BGGR"
	case ColorID_BAYER_CYYM:
		return "Bayer CYYM"
	case ColorID_BAYER_YCMY:
		return "Bayer YCMY"
	case ColorID_BAYER_YMCY:
		return "Bayer YMCY"
	case ColorID_BAYER_MYYC:
		return "Bayer MYYC"
	case ColorID_RGB:
		return "RGB"
	case ColorID_BGR:
		return "BGR"
	default:
		return "Unknown ColorID"
	}
}

var FileSignature [14]byte

type Header struct {
	FileID             [14]byte // File Signature
	LuID               int32    // Lumenera camera series ID (unused)
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
	filename string
}

func New(filename string, imageWidth int32, imageHeight int32) *SER {
	s := SER{filename: filename}
	s.Header = Header{}
	s.Header.FileID = FileSignature
	s.Header.ColorID = ColorID_RGB
	s.Header.Endianness = LittleEndian
	s.Header.ImageWidth = imageWidth
	s.Header.ImageHeight = imageHeight
	s.Header.PixelDepthPerPlane = 3
	s.Header.DateTimeUTC = DateTimeNow()

	s.createFile()

	return &s
}

func init() {
	copy(FileSignature[:], "LUCAM-RECORDER")
}

func (ser *SER) generateHeader() []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, ser.Header)
	if err != nil {
		log.Printf("ser.createFile: %v", err)
	}
	return buf.Bytes()
}

func (ser *SER) createFile() {
	ioutil.WriteFile(ser.filename, ser.generateHeader(), 0640)
}

func LoadHeader(filename string) (header Header, err error) {
	header = Header{}
	f, err := os.Open(filename)
	if err != nil {
		log.Printf("LoadHeader: %v", err)
		return header, err
	}

	defer f.Close()

	buf := bufio.NewReader(f)

	err = binary.Read(buf, binary.LittleEndian, &header)
	if err != nil {
		log.Printf("ser.LoadHeader: %v", err)
	}
	return header, nil
}

func Open(filename string) (ser *SER, err error) {

	s := SER{filename: filename}

	s.Header, err = LoadHeader(filename)
	if err != nil {
		return nil, err
	}

	if s.Header.FileID != FileSignature {
		return nil, fmt.Errorf("Invalid .SER file")
	}

	return &s, nil
}

func FixedStringToString(fixed [40]byte) string {
	s := ""
	for _, c := range fixed {
		if c == 0 {
			return s
		}
		if c >= 32 && c <= 126 {
			s += string(c)
		}
	}

	return s
}
