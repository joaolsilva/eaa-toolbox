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
