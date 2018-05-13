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

package toolbox

// Read ~/.eaatoolboxrc

import (
	"github.com/naoina/toml"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const defaultLimitingMagnitude = 9.0

type AppConfig struct {
	Location struct {
		Latitude  float64
		Longitude float64
	}

	GoTo struct {
		LimitingMagnitude float64
	}
	Paths struct {
		VSOP87 string
		Web    string
	}
}

func LoadAppConfig() (appConfig AppConfig) {
	appConfig.GoTo.LimitingMagnitude = defaultLimitingMagnitude
	u, err := user.Current()
	if err != nil {
		return appConfig
	}
	f, err := os.Open(filepath.Join(u.HomeDir, ".eaatoolboxrc"))
	if err != nil {
		log.Printf("ReadAppConfig: Open failed: %v", err)
		return appConfig
	}
	defer f.Close()
	if err := toml.NewDecoder(f).Decode(&appConfig); err != nil {
		log.Printf("ReadAppConfig: %v", err)
		return appConfig
	}
	return appConfig
}
