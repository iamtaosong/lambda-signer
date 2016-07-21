package lambda

import (
	"testing"
)

func TestDetermineAutoScalingEventType(t *testing.T) {
	samples := []struct {
		Value  AutoScalingEventType
		Source string
	}{
		{Value: AutoScalingEventNone, Source: "EC2 None"},
		{Value: AutoScalingEventLaunch, Source: "EC2 Instance Launch Successful"},
		{Value: AutoScalingEventTerminate, Source: "EC2 Instance Terminate Successful"},
	}

	for _, s := range samples {
		if DetermineAutoScalingEventType(s.Source) != s.Value {
			t.Fatal("Error")
		}
	}
}
