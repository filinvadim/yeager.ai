package main

import (
	"errors"
	"math/rand"
	"testing"
)

type testCase struct {
	name      string
	fragments map[sequence]Fragment
	expected  string
}

func intToPtr(i int) *int {
	return &i
}

func TestReconstructOrder_Positive(t *testing.T) {
	testCases := []testCase{
		{
			name: "simple",
			fragments: map[sequence]Fragment{
				intToPtr(0): {dataKey: "Hello", hashKey: simpleHash("Hello")},
				intToPtr(2): {dataKey: "World", hashKey: simpleHash("World")},
				intToPtr(3): {dataKey: "!", hashKey: simpleHash("!")},
			},
			expected: "HelloWorld!",
		},
		{
			name: "mixed",
			fragments: map[sequence]Fragment{
				intToPtr(4): {dataKey: "!", hashKey: simpleHash("!")},
				intToPtr(2): {dataKey: "the", hashKey: simpleHash("the")},
				intToPtr(1): {dataKey: "save", hashKey: simpleHash("save")},
				intToPtr(3): {dataKey: "Queen", hashKey: simpleHash("Queen")},
				intToPtr(0): {dataKey: "God", hashKey: simpleHash("God")},
			},
			expected: "GodsavetheQueen!",
		},
		{
			name: "big",
			fragments: map[sequence]Fragment{
				intToPtr(0):   {dataKey: "Hasta", hashKey: simpleHash("Hasta")},
				intToPtr(22):  {dataKey: "la", hashKey: simpleHash("la")},
				intToPtr(26):  {dataKey: "vista", hashKey: simpleHash("vista")},
				intToPtr(666): {dataKey: ",", hashKey: simpleHash(",")},
				intToPtr(999): {dataKey: "baby", hashKey: simpleHash("baby")},
			},
			expected: "Hastalavista,baby",
		},
	}

	for _, cs := range testCases {
		t.Run(cs.name, func(t *testing.T) {
			reconstructed, err := reconstructData(cs.fragments)
			if err != nil {
				t.Fatal(err)
			}
			if reconstructed != cs.expected {
				t.Errorf("got %s, want %s", reconstructed, cs.expected)
				t.Fatal("reconstructed != expected")
			}
			t.Logf("reconstructed: %s", reconstructed)
		})
	}
}

func TestReconstructOrder_Negative(t *testing.T) {
	testCases := []struct {
		name      string
		fragments map[sequence]Fragment
		expected  error
	}{
		{
			name: "missing fragment",
			fragments: map[sequence]Fragment{
				intToPtr(0): {dataKey: "Hello", hashKey: simpleHash("Hello")},
				intToPtr(1): nil,
				intToPtr(2): {dataKey: "World", hashKey: simpleHash("World")},
				intToPtr(3): {dataKey: "!", hashKey: simpleHash("!")},
			},
			expected: ErrMissingFragment,
		},
		{
			name: "broken order",
			fragments: map[sequence]Fragment{
				intToPtr(0): {dataKey: "God", hashKey: simpleHash("God")},
				nil:         {dataKey: "save", hashKey: simpleHash("save")}, // nil key to simulate broken order
				intToPtr(2): {dataKey: "the", hashKey: simpleHash("the")},
				intToPtr(3): {dataKey: "Queen", hashKey: simpleHash("Queen")},
				intToPtr(4): {dataKey: "!", hashKey: simpleHash("!")},
			},
			expected: ErrBrokenOrder,
		},
		{
			name: "hash mismatch",
			fragments: map[sequence]Fragment{
				intToPtr(0): {dataKey: "Hasta", hashKey: simpleHash("Hasta")},
				intToPtr(1): {dataKey: "la", hashKey: "invalid_hash"}, // Incorrect hash
				intToPtr(2): {dataKey: "vista", hashKey: simpleHash("vista")},
			},
			expected: ErrIntegrityVerification,
		},
	}

	for _, cs := range testCases {
		t.Run(cs.name, func(t *testing.T) {
			_, err := reconstructData(cs.fragments)
			if !errors.Is(err, cs.expected) {
				t.Fatalf("unexpected error: got %v, want %v", err, cs.expected)
			}
			t.Logf("expected error occurred: %v", err)
		})
	}
}

func TestSimpleHashLen_Positive(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected int // expected length of the hash
	}{
		{
			name:     "Empty string",
			input:    "",
			expected: hashLength,
		},
		{
			name:     "Short string",
			input:    "test",
			expected: hashLength,
		},
		{
			name:     "Long string",
			input:    "this is a much longer string than usual",
			expected: hashLength,
		},
		{
			name:     "Special characters",
			input:    "!@#$%^&*()_+",
			expected: hashLength,
		},
		{
			name:     "Unicode characters",
			input:    "你好，世界", // "Hello, World" in Chinese
			expected: hashLength,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash := simpleHash(tc.input)
			if len(hash) != tc.expected {
				t.Errorf("expected hash length %d, got %d", tc.expected, len(hash))
			}
		})
	}
}

func TestSimpleHash_Positive(t *testing.T) {
	testCases := []struct {
		name        string
		input1      string
		input2      string
		shouldMatch bool // whether the hashes should match
	}{
		{
			name:        "Same strings",
			input1:      "test",
			input2:      "test",
			shouldMatch: true,
		},
		{
			name:        "Special characters",
			input1:      "!@#$%^&*()_+",
			input2:      "!@#$%^&*()_+",
			shouldMatch: true,
		},
		{
			name:        "Upper case",
			input1:      "TEST",
			input2:      "TEST",
			shouldMatch: true,
		},
		{
			name:        "Empty",
			input1:      "",
			input2:      "",
			shouldMatch: true,
		},
		{
			name:        "Unicode strings",
			input1:      "你好，世界",
			input2:      "你好，世界",
			shouldMatch: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash1 := simpleHash(tc.input1)
			hash2 := simpleHash(tc.input2)

			if tc.shouldMatch && hash1 != hash2 {
				t.Errorf("expected hashes to match for inputs %q and %q, but got %q and %q", tc.input1, tc.input2, hash1, hash2)
			}
		})
	}
}

func TestSimpleHash_Negative(t *testing.T) {
	testCases := []struct {
		name        string
		input1      string
		input2      string
		shouldMatch bool // whether the hashes should match
	}{
		{
			name:        "Same strings",
			input1:      "test",
			input2:      "test",
			shouldMatch: true,
		},
		{
			name:        "Different strings",
			input1:      "test1",
			input2:      "test2",
			shouldMatch: false,
		},
		{
			name:        "Case sensitivity",
			input1:      "Test",
			input2:      "test",
			shouldMatch: false,
		},
		{
			name:        "Empty vs non-empty",
			input1:      "",
			input2:      "non-empty",
			shouldMatch: false,
		},
		{
			name:        "Similar strings",
			input1:      "abcdefg",
			input2:      "abcdeFg",
			shouldMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash1 := simpleHash(tc.input1)
			hash2 := simpleHash(tc.input2)

			if tc.shouldMatch && hash1 != hash2 {
				t.Errorf("expected hashes to match for inputs %q and %q, but got %q and %q", tc.input1, tc.input2, hash1, hash2)
			}

			if !tc.shouldMatch && hash1 == hash2 {
				t.Errorf("expected hashes to differ for inputs %q and %q, but got the same hash %q", tc.input1, tc.input2, hash1)
			}
		})
	}
}

func randomString(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:',.<>?/`~"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func TestSimpleHashCollisions(t *testing.T) {
	collisionCheckMap := make(map[string]string)

	var (
		numTests      = 100000
		stringLength  = 5 + rand.Intn(15)
		duplicatesNum = 0
	)

	for i := 0; i < numTests; i++ {
		input := randomString(stringLength)
		hash := simpleHash(input)

		if existingInput, ok := collisionCheckMap[hash]; ok {
			if existingInput != input {
				t.Errorf("collision detected: input %s and %s produced the same hash %s", existingInput, input, hash)
				return
			}
			t.Logf("input %s and %s produced the same hash %s", existingInput, input, hash)
			duplicatesNum++
		}

		collisionCheckMap[hash] = input
	}

	t.Logf("test completed with %d unique hashes generated", len(collisionCheckMap)-duplicatesNum)
	collisionCheckMap = nil
}
