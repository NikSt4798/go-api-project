package task

import (
	"fmt"
)

type Task interface {
	Solve(parameters map[string]interface{}) (interface{}, error)
}

type SimpleNumber struct{}

func (s *SimpleNumber) Solve(parameters map[string]interface{}) (interface{}, error) {
	numberString, ok := parameters["Number"]
	if !ok {
		return nil, fmt.Errorf("wrong parameter")
	}

	numberFloat, ok := numberString.(float64)
	if !ok {
		return nil, fmt.Errorf("not int")
	}

	number := int(numberFloat)

	return s.eratosfetn(number), nil
}

func (s *SimpleNumber) eratosfetn(number int) bool {
	boolMap := make([]bool, number+1)
	for i := range boolMap {
		boolMap[i] = true
	}
	for i := range boolMap {
		if i == 0 {
			boolMap[i] = false
			continue
		}
		if i == 1 {
			boolMap[i] = true
			continue
		}
		if boolMap[i] {
			for j := 2; j*(i) < len(boolMap); j++ {
				boolMap[j*(i)] = false
			}
		}
	}

	return boolMap[len(boolMap)-1]
}
