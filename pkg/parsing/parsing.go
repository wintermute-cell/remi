package parsing

import (
	"regexp"
	"remi/pkg/model"
	"remi/pkg/remi_errors"
	"strconv"
	"strings"
	"time"
)

// CollectMessageFromArgs collects the message from the args,
// to catch the case where the message contains spaces and is not quoted.
func CollectMessageFromArgs(args []string, startIdx int) (string, error) {
	ret := ""
	for i := startIdx; i < len(args); i++ {
		ret += args[i] + " "
	}
	ret = strings.TrimSpace(ret)
	if ret == "" {
		return "", remi_errors.MissingMessageError
	}
	return ret, nil
}

// ParseDuration parses a string input and returns a time.Duration struct.
// Possible formats: "111d" (111 day), "111h" (111 hour), "111m" (111 minute)
func ParseDuration(input string) (time.Duration, error) {
	isDuration, _ := regexp.MatchString("^\\d+[dhm]$", input)
	if !isDuration {
		return 0, remi_errors.NotADurationError
	}
	durationNumStr := input[:len(input)-1]
	durationNum, err := strconv.ParseInt(durationNumStr, 10, 64)
	if err != nil {
		return 0, err
	}
	durationUnit := input[len(input)-1]
	switch durationUnit {
	case 'd':
		return time.Duration(durationNum) * 24 * time.Hour, nil
	case 'h':
		return time.Duration(durationNum) * time.Hour, nil
	case 'm':
		return time.Duration(durationNum) * time.Minute, nil
	}

	// Should never reach this point.
	return 0, nil
}

// ParseTimestamp parses a string input and returns a Timestamp struct.
// Possible formats: "dd.mm.yy" or "dd.mm.yy@hh:mm"
func ParseTimestamp(input string) (model.Timestamp, error) {
	layoutOptions := []string{"02.01.06", "02.01.06@15:04"}
	match, _ := regexp.MatchString("^\\d{2}\\.\\d{2}\\.\\d{2}(@\\d{2}:\\d{2})?$", input)
	if !match {
		return model.Timestamp{}, remi_errors.InvalidTimestampError
	}

	t := model.Timestamp{}
	layout := layoutOptions[0]
	if strings.Contains(input, "@") {
		layout = layoutOptions[1]
		t.HasTime = true
	}

	var err error
	t.DateTime, err = time.Parse(layout, input)
	return t, err
}
