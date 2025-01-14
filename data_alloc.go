package main

import (
	"math"
	"sort"
)

func distributeFragments(dataCenters []int, fragments int) int64 {
	if len(dataCenters) == 0 || fragments <= 0 {
		return 0
	}
	// dataCenters: [20, 10, 30]
	sort.Ints(dataCenters) // dataCenters sorted: [10, 20, 30]

	mostRiskyDC := dataCenters[len(dataCenters)-1] // mostRiskyDC = 30
	// binary search bounds
	minRisk := int64(1)
	maxRisk := calculateMaxRisk(mostRiskyDC, fragments)

	// binary search to find the minimal max risk.
	for minRisk < maxRisk {
		mediumRisk := (minRisk + maxRisk) / 2
		if isRiskAchievable(mediumRisk, dataCenters, fragments) {
			maxRisk = mediumRisk
		} else {
			minRisk = mediumRisk + 1
		}
	}

	return minRisk
}

// if a given maximum risk is achievable with the current configuration.
func isRiskAchievable(maxRisk int64, dataCenters []int, fragments int) bool {
	remainingFragments := fragments

	for _, risk := range dataCenters {
		count := 0

		for (math.Pow(float64(risk), float64(count+1))) <= float64(maxRisk) {
			count++ // putting fragment into datacenter if they fit
		}

		remainingFragments -= count

		if remainingFragments <= 0 {
			return true
		}
	}

	return false
}

// otherwise overflow leads to integer being "wrap around": x int8 = 127; x++; x = -128
func calculateMaxRisk(base, exp int) int64 {
	result := int64(1)
	for i := 0; i < exp; i++ {
		if result > (math.MaxInt64 / int64(base)) {
			panic("maxRisk overflow")
		}
		result *= int64(base)
	}
	return result
}
