package servicechecker

import "github.com/huanghj78/jepsenFuzz/pkg/core"

type Checker struct{}

func (Checker) Check(_ core.Model, _ []core.Operation) (bool, error) {

	return true, nil
}

func (Checker) Name() string {
	return "service_checker"
}
