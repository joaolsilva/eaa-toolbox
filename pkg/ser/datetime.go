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

import "time"

type DateTime uint64

const (
	secondsPerDay       = 24 * 3600
	unixOffset    int64 = (1969*365 + 1969/4 - 1969/100 + 1969/400) * secondsPerDay
)

func DateTimeNow() DateTime {
	return DateTime((time.Now().UnixNano() / 100) + unixOffset*1E7)
}

func TimeFromDateTime(dateTime DateTime) time.Time {
	return time.Unix(0, (int64(dateTime)-unixOffset*1E7)*100)
}
