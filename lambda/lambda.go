package lambda

var (
	// DefaultConfigFile stores the default config file name
	DefaultConfigFile = "config.json"

	autoScalingEventTypes = []string{
		"AutoScalingEventNone",
		"AutoScalingEventLaunch",
		"AutoScalingEventTerminate",
	}
)

// AutoScalingEventType represents an autoscaling action
type AutoScalingEventType int

func (asgEvt AutoScalingEventType) String() string {
	return autoScalingEventTypes[int(asgEvt)]
}

const (
	// AutoScalingEventNone for an unknown event
	AutoScalingEventNone AutoScalingEventType = iota
	// AutoScalingEventLaunch for launch action
	AutoScalingEventLaunch
	// AutoScalingEventTerminate for terminate action
	AutoScalingEventTerminate
)

// DetermineAutoScalingEventType detects the autoscaling action and returns a native type
func DetermineAutoScalingEventType(detailType string) AutoScalingEventType {
	switch detailType {
	case "EC2 Instance Launch Successful":
		return AutoScalingEventLaunch
	case "EC2 Instance Terminate Successful":
		return AutoScalingEventTerminate
	default:
		return AutoScalingEventNone
	}
}
