package algorithms

import (
	"math"
	"time"
)

type FlapDetection struct {
	Enabled bool `toml:"enabled"`

	//How many seconds does it take for a timestamp to reach 0.5 flapscore
	HalfLifeSeconds float64 `toml:"half_life_seconds"`

	//Flap score at which the alert is considered to be flapping
	Threshold float64 `toml:"threshold"`

	//individual severity time change score threshold at which a severity time change is dropped from history
	MinimumScore float64 `toml:"minimum_score"`
}

func (f *FlapDetection) Init() {
	if f.HalfLifeSeconds == 0 {

		//one minute is an easy to grasp half life
		f.HalfLifeSeconds = 60
	}
	if f.Threshold == 0 {

		//This implies a minimum of 4 severity change events before flapping is detected
		f.Threshold = 4
	}
	if f.MinimumScore == 0 {
		f.MinimumScore = 0.02
	}
}

//http://nagios.manubulon.com/traduction/docs25en/flapping.html
//http://linuxczar.net/blog/2016/01/31/flap-detection/
//sum the decaying values of the severity change timestamps - values are between 0 and 1
func (f *FlapDetection) Detect(severityChangeTimes []time.Time) (isFlapping bool, currentFlapScore float64, remainingSeverityTimeChanges []time.Time) {
	now := time.Now()

	for _, severityChangeTime := range severityChangeTimes {
		severityChangeTimeScore := f.exponentialDecay(severityChangeTime, now)
		if severityChangeTimeScore > f.MinimumScore {
			remainingSeverityTimeChanges = append(remainingSeverityTimeChanges, severityChangeTime)
			currentFlapScore += severityChangeTimeScore
		}
	}

	isFlapping = currentFlapScore > f.Threshold

	return
}

//Generates values between 0 and 1 using the time difference between two timestamps and the given half life
func (f *FlapDetection) exponentialDecay(t1, t2 time.Time) float64 {
	elapsedSeconds := t2.Sub(t1).Seconds()
	return math.Exp((-elapsedSeconds * math.Ln2) / f.HalfLifeSeconds)
}
