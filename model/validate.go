package model

import "regexp"

func Validate(exp Experiment) bool {

	if !validName(exp.ExperimentId) {
		return false
	}

	for _, arm := range exp.Arms {
		if !validName(arm.Name) {
			return false
		}
	}

	return true

}

func validName(name string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name)
}
