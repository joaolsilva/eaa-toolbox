package sacdb /* import "r2discover.com/go/eaa-toolbox/pkg/sacdb" */
import (
	"github.com/soniakeys/meeus/v3/angle"
	"github.com/soniakeys/unit"
	"log"
	"sort"
	"strconv"
	"strings"
)

// Saguaro Astronomy Club Database, http://www.saguaroastro.org/content/downloads.htm

type SACDB struct {
	ObjectID          string  `json:"object_id"`
	OtherID           string  `json:"other_id"`
	Type              string  `json:"type"`
	Constellation     string  `json:"constellation"`
	RA                string  `json:"ra"`
	Dec               string  `json:"dec"`
	Magnitude         float64 `json:"magnitude"`
	SurfaceBrightness float64 `json:"surface_brightness"`
	U2K               string  `json:"u2k"`
	Ti                string  `json:"ti"`
	SizeMax           string  `json:"size_max"`
	SizeMin           string  `json:"size_min"`
	PositionAngle     int     `json:"position_angle"`
	Classification    string  `json:"classification"`
	NStars            int     `json:"n_stars"`
	BrightestStar     float64 `json:"brightest_star"`
	Catalogs          string  `json:"catalogs"`
	NGCDesc           string  `json:"ngc_desc"`
	Notes             string  `json:"notes"`
}

type DeepSkyObject struct {
	Desc      string
	RA        float64
	Dec       float64
	Magnitude float64
	distance  float64
}

func sacdbRAToDegrees(sacdbRA string) (RA float64) {
	l := strings.Split(sacdbRA, " ")
	if len(l) == 0 || len(l) > 2 {
		return 0.0
	}

	RA, err := strconv.ParseFloat(l[0], 64)
	if err != nil {
		log.Printf("sacdbRAToDegrees: %v", err)
	}

	RA = 360.0 * (RA / 24.0)

	if len(l) == 2 {
		minutes, err := strconv.ParseFloat(l[1], 64)
		if err != nil {
			log.Printf("sacdbRAToDegrees: %v", err)
		}
		RA += (minutes / 60.0) * (360.0 / 24.0)
	}
	return RA
}

func sacdbDecToDegrees(sacdbDec string) (Dec float64) {

	l := strings.Split(sacdbDec, " ")
	if len(l) == 0 || len(l) > 2 {
		return 0.0
	}

	Dec, err := strconv.ParseFloat(l[0], 64)
	if err != nil {
		log.Printf("sacdbDecToDegrees: %v", err)
	}

	if len(l) == 2 {
		minutes, err := strconv.ParseFloat(l[1], 64)
		if err != nil {
			log.Printf("sacdbDecToDegrees: %v", err)
		}
		if Dec > 0.0 {
			Dec += (minutes / 60.0)
		} else {
			Dec -= (minutes / 60.0)
		}
	}

	return Dec
}

func ObjectsCloseTo(RA, Dec float64, limitingMagnitude float64, maxResults int) []DeepSkyObject {
	objects := []DeepSkyObject{}
	for _, v := range DeepSky {
		if v.Magnitude > limitingMagnitude {
			continue
		}

		obj := DeepSkyObject{}
		obj.Desc = v.ObjectID
		if v.OtherID != "" {
			obj.Desc = obj.Desc + " " + v.OtherID
		}
		obj.Magnitude = v.Magnitude
		obj.RA = sacdbRAToDegrees(v.RA)
		obj.Dec = sacdbDecToDegrees(v.Dec)
		obj.distance = float64(angle.SepHav(unit.AngleFromDeg(RA), unit.AngleFromDeg(Dec), unit.AngleFromDeg(obj.RA), unit.AngleFromDeg(obj.Dec)))
		objects = append(objects, obj)
	}
	sort.Slice(objects, func(i, j int) bool { return objects[i].distance < objects[j].distance })
	if len(objects) > maxResults {
		return objects[:maxResults]
	}

	return objects
}

func searchMatch(s1, s2 string) bool {
	s1 = strings.Replace(strings.ToLower(s1), " ", "", -1)
	s2 = strings.Replace(strings.ToLower(s2), " ", "", -1)
	return strings.Index(s1, s2) != -1
}

func Search(RA, Dec float64, searchText string, maxResults int) []DeepSkyObject {
	objects := []DeepSkyObject{}
	for _, v := range DeepSky {
		obj := DeepSkyObject{}
		obj.Desc = v.ObjectID
		if v.OtherID != "" {
			obj.Desc = obj.Desc + " " + v.OtherID
		}
		if !searchMatch(obj.Desc, searchText) {
			continue
		}

		obj.Magnitude = v.Magnitude
		obj.RA = sacdbRAToDegrees(v.RA)
		obj.Dec = sacdbDecToDegrees(v.Dec)
		obj.distance = float64(angle.SepHav(unit.AngleFromDeg(RA), unit.AngleFromDeg(Dec), unit.AngleFromDeg(obj.RA), unit.AngleFromDeg(obj.Dec)))
		objects = append(objects, obj)
	}
	sort.Slice(objects, func(i, j int) bool { return objects[i].distance < objects[j].distance })
	if len(objects) > maxResults {
		return objects[:maxResults]
	}

	return objects
}
