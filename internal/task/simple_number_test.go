package task

import "testing"

func TestSimpleNumber_eratosfen(t *testing.T) {
	primeNumbers := []int{1, 2, 3, 5, 7, 11, 10000079}
	notPrimeNumbers := []int{4, 6, 8, 9, 10000081}

	test := func(s []int, expected bool) {
		pn := &SimpleNumber{}
		for _, v := range s {
			m := make(map[string]interface{})
			m["Number"] = float64(v)

			result, err := pn.Solve(m)
			boolResult := result.(bool)
			if err != nil || boolResult != expected {
				t.Errorf("wrong prime number check: %d", v)
			}
		}
	}

	t.Run("check is prime numbers correct", func(t *testing.T) {
		test(primeNumbers, true)
	})
	t.Run("check is not prime numbers correct", func(t *testing.T) {
		test(notPrimeNumbers, false)
	})
}
