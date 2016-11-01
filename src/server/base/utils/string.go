package utils

import (
	"strconv"
	"strings"
)

const (
	StrAbcSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Str123Set = "1234567890"
)

func IsGeneratedBy(s string, set string) bool {
	for i := 0; i < len(s); i++ {
		if !strings.Contains(set, string(s[i])) {
			return false
		}
	}
	return true
}

func ValidEmailFormat(email string) bool {
	if len(email) == 0 {
		return false
	}

	subs := strings.Split(email, "@")
	if subs == nil {
		return false
	}

	if len(subs) != 2 {
		return false
	}

	if !IsGeneratedBy(subs[0], StrAbcSet+Str123Set+"_") {
		return false
	}

	subTails := strings.Split(subs[1], ".")
	if subTails == nil {
		return false
	}

	if len(subTails) == 0 {
		return false
	}

	for _, tail := range subTails {
		if len(tail) == 0 {
			return false
		}

		if !IsGeneratedBy(tail, StrAbcSet+Str123Set+"_") {
			return false
		}
	}

	return true
}

func StrToInt32(s string) (int32, error) {
	i, e := strconv.ParseInt(s, 10, 32)
	if e != nil {
		return 0, e
	}

	return int32(i), nil
}

func StrToInt64(s string) (int64, error) {
	i, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		return 0, e
	}

	return i, nil
}

func StrToFloat32(s string) (float32, error) {
	f, e := strconv.ParseFloat(s, 64)
	if e != nil {
		return 0, e
	}

	return float32(f), nil
}
