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

package main

import (
	"fmt"
	"log"
	"os"
	"r2discover.com/go/eaa-toolbox/pkg/ser"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v file.ser\n", os.Args[0])
		return
	}

	filename := os.Args[1]

	serFile, err := ser.Open(filename)
	if err != nil {
		log.Printf("%v", err)
		return
	}

	fmt.Printf("\nFile %v:\n\n", filename)
	fmt.Printf("LuID: %v\n", serFile.Header.LuID)
	fmt.Printf("ColorID: %v\n", serFile.Header.ColorID)
	fmt.Printf("Endianness: %v\n", serFile.Header.Endianness)
	fmt.Printf("ImageWidth: %v\n", serFile.Header.ImageWidth)
	fmt.Printf("ImageHeight: %v\n", serFile.Header.ImageHeight)
	fmt.Printf("PixelDepthPerPlane: %v\n", serFile.Header.PixelDepthPerPlane)
	fmt.Printf("FrameCount: %v\n", serFile.Header.FrameCount)
	fmt.Printf("Observer: %v\n", ser.FixedStringToString(serFile.Header.Observer))
	fmt.Printf("Instrument: %v\n", ser.FixedStringToString(serFile.Header.Instrument))
	fmt.Printf("Telescope: %v\n", ser.FixedStringToString(serFile.Header.Telescope))
	fmt.Printf("DateTime: %v\n", ser.TimeFromDateTime(serFile.Header.DateTime))
	fmt.Printf("DateTimeUTC: %v\n", ser.TimeFromDateTime(serFile.Header.DateTimeUTC))
}
