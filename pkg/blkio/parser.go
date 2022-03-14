package blkio

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type ByteSize uint64

// Used to convert user input to ByteSize
var unitMap = map[string]ByteSize{
	"KB": 1000,
	"MB": 1000 * 1000,
	"GB": 1000 * 1000 * 1000,
}

// Used to convert user input to ByteSize
var durationMap = map[string]time.Duration{
	"s": 1 * time.Second,
	"m": 1 * time.Minute,
	"h": 1 * time.Hour,
}

func ParseRateParam(args string) (ByteSize, time.Duration) {
	if args == "" {
		return 0, 0
	}

	ps := strings.Split(args, "/")
	if len(ps) != 2 {
		log.Fatal("invalid read rate param")
	}

	rate, err := ParseByteSize(ps[0])
	if err != nil {
		log.Fatal(err.Error())
	}

	dur, err := ParseDuration(ps[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	return rate, dur
}

func ParseDuration(s string) (time.Duration, error) {
	_, ok := durationMap[s]
	if ok {
		return time.ParseDuration("1" + s)
	}
	return time.ParseDuration(s)
}

func ParseByteSize(s string) (ByteSize, error) {
	// Remove leading and trailing whitespace
	s = strings.TrimSpace(s)

	split := make([]string, 0)
	for i, r := range s {
		if !unicode.IsDigit(r) {
			// Split the string by digit and size designator, remove whitespace
			split = append(split, strings.TrimSpace(string(s[:i])))
			split = append(split, strings.TrimSpace(string(s[i:])))
			break
		}
	}

	// Check to see if we split successfully
	if len(split) != 2 {
		return 0, errors.New("Unrecognized size suffix")
	}

	// Check for MB, MEGABYTE, and MEGABYTES
	unit, ok := unitMap[strings.ToUpper(split[1])]
	if !ok {
		return 0, errors.New("Unrecognized size suffix " + split[1])

	}

	value, err := strconv.ParseFloat(split[0], 64)
	if err != nil {
		return 0, err
	}

	bytesize := ByteSize(value * float64(unit))
	return bytesize, nil
}
