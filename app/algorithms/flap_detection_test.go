package algorithms

import (
	"testing"
	"time"
)

func TestFlapDetection_Detect(t *testing.T) {
	f := &FlapDetection{
		HalfLifeSeconds: 60,
		Threshold:       0.9,
		MinimumScore:    0.01,
	}

	oneSecondAgo := time.Now().Add(-time.Second)
	severityTimeChanges := []time.Time{oneSecondAgo}

	isFlapping, flapScore, remainingSeverityTimeChanges := f.Detect(severityTimeChanges)

	if !isFlapping {
		t.Errorf("should be flapping")
	}

	if flapScore < 0.9 {
		t.Error("flap score should be > 0.9\nflapScore: %f", flapScore)
	}

	if len(remainingSeverityTimeChanges) != 1 {
		t.Error("oneSecondAgo severityTimeChange should not have been decayed from the list")
	}
}

func TestFlapDetection_DetectMultiple(t *testing.T) {
	f := &FlapDetection{
		HalfLifeSeconds: 60,
		Threshold:       1,
		MinimumScore:    0.01,
	}

	oneSecondAgo := time.Now().Add(-time.Second)
	twoSecondsAgo := time.Now().Add(-time.Second * 2)
	severityTimeChanges := []time.Time{oneSecondAgo, twoSecondsAgo}

	isFlapping, flapScore, remainingSeverityTimeChanges := f.Detect(severityTimeChanges)

	if !isFlapping {
		t.Errorf("should be flapping")
	}

	if !(flapScore > 1 && flapScore < 2) {
		t.Errorf("flap score should be > 1 and flapScore < 2\nflapScore: %f", flapScore)
	}

	if len(remainingSeverityTimeChanges) != 2 {
		t.Error("no severityTimeChanges should have decayed from the list")
	}
}

func TestFlapDetection_DetectDecay(t *testing.T) {
	f := &FlapDetection{
		HalfLifeSeconds: 60,
		Threshold:       1,
		MinimumScore:    0.51,
	}

	sixtySecondsAgo := time.Now().Add(-time.Second * 60)
	severityTimeChanges := []time.Time{sixtySecondsAgo}

	isFlapping, flapScore, remainingSeverityTimeChanges := f.Detect(severityTimeChanges)

	if isFlapping {
		t.Errorf("should not be flapping")
	}

	if flapScore > 0.5 {
		t.Errorf("flap score should not be > 0.5\nflapScore: %f", flapScore)
	}

	if len(remainingSeverityTimeChanges) != 0 {
		t.Error("severityTimeChange should have decayed from the list")
	}
}
