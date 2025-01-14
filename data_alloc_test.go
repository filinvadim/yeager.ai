package main

import "testing"

func TestDistributeFragments(t *testing.T) {
	t.Run("Simple", TestDistributeFragments_Simple)
	t.Run("EqualRisks", TestDistributeFragments_EqualRisks)
	t.Run("OneCenter", TestDistributeFragments_OneCenter)
	t.Run("ManyFragments", TestDistributeFragments_ManyFragments)
	t.Run("EmptyCenters", TestDistributeFragments_EmptyCenters)
	t.Run("NoFragments", TestDistributeFragments_NoFragments)
	t.Run("NegativeFragments", TestDistributeFragments_NegativeFragments)
	t.Run("HighRisk", TestDistributeFragments_HighRisk)
	t.Run("LargeInput", TestDistributeFragments_LargeInput)
}

func TestDistributeFragments_Simple(t *testing.T) {
	dataCenters := []int{10, 30, 20}
	fragments := 5
	expected := int64(400)
	result := distributeFragments(dataCenters, fragments)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

func TestDistributeFragments_EqualRisks(t *testing.T) {
	dataCenters := []int{10, 10, 10}
	fragments := 4
	expected := int64(100)
	result := distributeFragments(dataCenters, fragments)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

func TestDistributeFragments_OneCenter(t *testing.T) {
	dataCenters := []int{10}
	fragments := 3
	expected := int64(1000) // 10^3
	result := distributeFragments(dataCenters, fragments)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

func TestDistributeFragments_ManyFragments(t *testing.T) {
	dataCenters := []int{5, 10, 7}
	fragments := 10
	expected := int64(1000)

	result := distributeFragments(dataCenters, fragments)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

func TestDistributeFragments_EmptyCenters(t *testing.T) {
	dataCenters := []int{}
	fragments := 5
	expected := int64(0) // Нет центров, поэтому результат - 0
	result := distributeFragments(dataCenters, fragments)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

func TestDistributeFragments_NoFragments(t *testing.T) {
	dataCenters := []int{20, 30, 10}
	fragments := 0
	expected := int64(0)
	result := distributeFragments(dataCenters, fragments)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

func TestDistributeFragments_NegativeFragments(t *testing.T) {
	dataCenters := []int{10, 20, 30}
	fragments := -3
	expected := int64(0)
	result := distributeFragments(dataCenters, fragments)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

func TestDistributeFragments_HighRisk(t *testing.T) {
	dataCenters := []int{1000000}
	fragments := 2
	expected := int64(1000000000000) // 1000000^2
	result := distributeFragments(dataCenters, fragments)
	if result != expected {
		t.Errorf("Expected %d, but got %d", expected, result)
	}
}

func TestDistributeFragments_LargeInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("this code must panic")
		}
	}()
	dataCenters := []int{5, 10, 15, 20, 25, 30, 35, 40}
	fragments := 50
	_ = distributeFragments(dataCenters, fragments)
}
