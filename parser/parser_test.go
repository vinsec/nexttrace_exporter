package parser

import (
	"testing"
)

func TestParseNextTraceOutput(t *testing.T) {
	jsonData := []byte(`{
		"target": "8.8.8.8",
		"hops": [
			{
				"ttl": 1,
				"ip": "192.168.1.1",
				"hostname": "gateway.local",
				"rtt": [1.23, 1.45, 1.34],
				"loss": 0.0,
				"asn": "AS0",
				"location": ""
			},
			{
				"ttl": 2,
				"ip": "10.0.0.1",
				"hostname": "isp.router",
				"rtt": [5.67, 5.89, 6.01],
				"loss": 0.0,
				"asn": "AS12345",
				"location": "City, Country"
			}
		]
	}`)

	result, err := ParseNextTraceOutput(jsonData)
	if err != nil {
		t.Fatalf("ParseNextTraceOutput failed: %v", err)
	}

	if result.Target != "8.8.8.8" {
		t.Errorf("Expected target 8.8.8.8, got %s", result.Target)
	}

	if len(result.Hops) != 2 {
		t.Errorf("Expected 2 hops, got %d", len(result.Hops))
	}

	// Test first hop
	hop1 := result.Hops[0]
	if hop1.TTL != 1 {
		t.Errorf("Expected hop1 TTL 1, got %d", hop1.TTL)
	}
	if hop1.IP != "192.168.1.1" {
		t.Errorf("Expected hop1 IP 192.168.1.1, got %s", hop1.IP)
	}

	// Test average RTT
	avgRTT := hop1.AverageRTT()
	expectedAvg := (1.23 + 1.45 + 1.34) / 3
	tolerance := 0.001
	if avgRTT < expectedAvg-tolerance || avgRTT > expectedAvg+tolerance {
		t.Errorf("Expected average RTT %.4f (±%.4f), got %.4f", expectedAvg, tolerance, avgRTT)
	}
}

func TestHopAverageRTT(t *testing.T) {
	tests := []struct {
		name     string
		hop      Hop
		expected float64
	}{
		{
			name:     "normal RTT values",
			hop:      Hop{RTT: []float64{1.0, 2.0, 3.0}},
			expected: 2.0,
		},
		{
			name:     "single RTT value",
			hop:      Hop{RTT: []float64{5.5}},
			expected: 5.5,
		},
		{
			name:     "empty RTT values",
			hop:      Hop{RTT: []float64{}},
			expected: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hop.AverageRTT()
			if result != tt.expected {
				t.Errorf("Expected %.2f, got %.2f", tt.expected, result)
			}
		})
	}
}

func TestHopHasValidIP(t *testing.T) {
	tests := []struct {
		name     string
		hop      Hop
		expected bool
	}{
		{
			name:     "valid IP",
			hop:      Hop{IP: "192.168.1.1"},
			expected: true,
		},
		{
			name:     "empty IP",
			hop:      Hop{IP: ""},
			expected: false,
		},
		{
			name:     "asterisk IP",
			hop:      Hop{IP: "*"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.hop.HasValidIP()
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
