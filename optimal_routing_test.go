package main

import (
	"math"
	"testing"
)

func TestFindMinimumLatencyPath(t *testing.T) {
	tests := []struct {
		name             string
		graph            map[string][]Router
		compressionNodes []string
		source           string
		destination      string
		expectedLatency  float64
	}{
		{
			name: "Simple graph without compression",
			graph: map[string][]Router{
				"A": {{"B", 10}, {"C", 20}},
				"B": {{"D", 15}},
				"C": {{"D", 30}},
				"D": {},
			},
			compressionNodes: []string{},
			source:           "A",
			destination:      "D",
			expectedLatency:  25,
		},
		{
			name: "Simple graph with compression",
			graph: map[string][]Router{
				"A": {{"B", 10}, {"C", 20}},
				"B": {{"D", 15}},
				"C": {{"D", 30}},
				"D": {},
			},
			compressionNodes: []string{"B"},
			source:           "A",
			destination:      "D",
			expectedLatency:  17.5,
		},
		{
			name: "Larger graph with multiple compression nodes",
			graph: map[string][]Router{
				"A": {{"B", 10}, {"C", 15}},
				"B": {{"C", 5}, {"D", 20}},
				"C": {{"D", 10}},
				"D": {{"E", 5}},
				"E": {},
			},
			compressionNodes: []string{"B", "C"},
			source:           "A",
			destination:      "E",
			expectedLatency:  22.5,
		},
		{
			name: "Unreachable destination",
			graph: map[string][]Router{
				"A": {{"B", 10}},
				"B": {{"C", 20}},
				"C": {},
			},
			compressionNodes: []string{"A", "B"},
			source:           "A",
			destination:      "D",
			expectedLatency:  math.Inf(1),
		},
		{
			name: "Graph with no edges",
			graph: map[string][]Router{
				"A": {},
				"B": {},
			},
			compressionNodes: []string{"A"},
			source:           "A",
			destination:      "B",
			expectedLatency:  math.Inf(1),
		},
		{
			name: "Source equals destination",
			graph: map[string][]Router{
				"A": {{"B", 10}},
				"B": {},
			},
			compressionNodes: []string{"A"},
			source:           "A",
			destination:      "A",
			expectedLatency:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			latency := findMinimumLatencyPath(tt.graph, tt.compressionNodes, tt.source, tt.destination)
			if latency != tt.expectedLatency {
				t.Errorf("expected latency %.2f, got %.2f", tt.expectedLatency, latency)
			}
		})
	}
}
