package parser

import (
	"testing"
)

func TestParseNextTraceOutput(t *testing.T) {
	// Real nexttrace -j output format
	jsonData := []byte(`{
		"Hops": [
			[
				{
					"Success": true,
					"Address": {"IP": "192.168.1.1", "Zone": ""},
					"Hostname": "gateway.local",
					"TTL": 1,
					"RTT": 1230000,
					"Error": null,
					"Geo": {
						"ip": "",
						"asnumber": "64512",
						"country": "",
						"country_en": "",
						"prov": "",
						"prov_en": "",
						"city": "",
						"city_en": "",
						"district": "",
						"owner": "",
						"isp": "",
						"domain": "",
						"whois": "RFC1918",
						"lat": 0,
						"lng": 0,
						"prefix": "",
						"router": null,
						"source": ""
					},
					"Lang": "cn",
					"MPLS": null
				},
				{
					"Success": true,
					"Address": {"IP": "192.168.1.1", "Zone": ""},
					"Hostname": "gateway.local",
					"TTL": 1,
					"RTT": 1450000,
					"Error": null,
					"Geo": {
						"ip": "",
						"asnumber": "64512",
						"country": "",
						"country_en": "",
						"prov": "",
						"prov_en": "",
						"city": "",
						"city_en": "",
						"district": "",
						"owner": "",
						"isp": "",
						"domain": "",
						"whois": "RFC1918",
						"lat": 0,
						"lng": 0,
						"prefix": "",
						"router": null,
						"source": ""
					},
					"Lang": "cn",
					"MPLS": null
				},
				{
					"Success": true,
					"Address": {"IP": "192.168.1.1", "Zone": ""},
					"Hostname": "gateway.local",
					"TTL": 1,
					"RTT": 1340000,
					"Error": null,
					"Geo": {
						"ip": "",
						"asnumber": "64512",
						"country": "",
						"country_en": "",
						"prov": "",
						"prov_en": "",
						"city": "",
						"city_en": "",
						"district": "",
						"owner": "",
						"isp": "",
						"domain": "",
						"whois": "RFC1918",
						"lat": 0,
						"lng": 0,
						"prefix": "",
						"router": null,
						"source": ""
					},
					"Lang": "cn",
					"MPLS": null
				}
			],
			[
				{
					"Success": true,
					"Address": {"IP": "10.0.0.1", "Zone": ""},
					"Hostname": "isp.router",
					"TTL": 2,
					"RTT": 5670000,
					"Error": null,
					"Geo": {
						"ip": "",
						"asnumber": "12345",
						"country": "Country",
						"country_en": "Country",
						"prov": "",
						"prov_en": "",
						"city": "City",
						"city_en": "City",
						"district": "",
						"owner": "",
						"isp": "Test ISP",
						"domain": "",
						"whois": "",
						"lat": 0,
						"lng": 0,
						"prefix": "",
						"router": {},
						"source": ""
					},
					"Lang": "cn",
					"MPLS": null
				},
				{
					"Success": true,
					"Address": {"IP": "10.0.0.1", "Zone": ""},
					"Hostname": "isp.router",
					"TTL": 2,
					"RTT": 5890000,
					"Error": null,
					"Geo": {
						"ip": "",
						"asnumber": "12345",
						"country": "Country",
						"country_en": "Country",
						"prov": "",
						"prov_en": "",
						"city": "City",
						"city_en": "City",
						"district": "",
						"owner": "",
						"isp": "Test ISP",
						"domain": "",
						"whois": "",
						"lat": 0,
						"lng": 0,
						"prefix": "",
						"router": {},
						"source": ""
					},
					"Lang": "cn",
					"MPLS": null
				},
				{
					"Success": true,
					"Address": {"IP": "10.0.0.1", "Zone": ""},
					"Hostname": "isp.router",
					"TTL": 2,
					"RTT": 6010000,
					"Error": null,
					"Geo": {
						"ip": "",
						"asnumber": "12345",
						"country": "Country",
						"country_en": "Country",
						"prov": "",
						"prov_en": "",
						"city": "City",
						"city_en": "City",
						"district": "",
						"owner": "",
						"isp": "Test ISP",
						"domain": "",
						"whois": "",
						"lat": 0,
						"lng": 0,
						"prefix": "",
						"router": {},
						"source": ""
					},
					"Lang": "cn",
					"MPLS": null
				}
			]
		],
		"TraceMapUrl": "https://example.com/trace"
	}`)

	result, err := ParseNextTraceOutput(jsonData)
	if err != nil {
		t.Fatalf("ParseNextTraceOutput failed: %v", err)
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
	if hop1.Hostname != "gateway.local" {
		t.Errorf("Expected hop1 hostname gateway.local, got %s", hop1.Hostname)
	}
	if hop1.ASN != "64512" {
		t.Errorf("Expected hop1 ASN 64512, got %s", hop1.ASN)
	}

	// Test average RTT (converted from nanoseconds to milliseconds)
	avgRTT := hop1.AverageRTT()
	expectedAvg := (1.23 + 1.45 + 1.34) / 3
	tolerance := 0.001
	if avgRTT < expectedAvg-tolerance || avgRTT > expectedAvg+tolerance {
		t.Errorf("Expected average RTT %.4f (±%.4f), got %.4f", expectedAvg, tolerance, avgRTT)
	}

	// Test second hop
	hop2 := result.Hops[1]
	if hop2.TTL != 2 {
		t.Errorf("Expected hop2 TTL 2, got %d", hop2.TTL)
	}
	if hop2.IP != "10.0.0.1" {
		t.Errorf("Expected hop2 IP 10.0.0.1, got %s", hop2.IP)
	}
	if hop2.ASN != "12345" {
		t.Errorf("Expected hop2 ASN 12345, got %s", hop2.ASN)
	}
	if hop2.Location != "City, Country" {
		t.Errorf("Expected hop2 location 'City, Country', got %s", hop2.Location)
	}

	// Test packet loss (all successful, should be 0)
	if hop1.Loss != 0.0 {
		t.Errorf("Expected hop1 loss 0.0, got %.2f", hop1.Loss)
	}
	if hop2.Loss != 0.0 {
		t.Errorf("Expected hop2 loss 0.0, got %.2f", hop2.Loss)
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

func TestParseNextTraceOutputWithPacketLoss(t *testing.T) {
	// Test with packet loss (some probes failed)
	jsonData := []byte(`{
		"Hops": [
			[
				{
					"Success": true,
					"Address": {"IP": "192.168.1.1", "Zone": ""},
					"Hostname": "",
					"TTL": 1,
					"RTT": 1000000,
					"Error": null,
					"Geo": null,
					"Lang": "cn",
					"MPLS": null
				},
				{
					"Success": false,
					"Address": null,
					"Hostname": "",
					"TTL": 1,
					"RTT": 0,
					"Error": {},
					"Geo": null,
					"Lang": "",
					"MPLS": null
				},
				{
					"Success": false,
					"Address": null,
					"Hostname": "",
					"TTL": 1,
					"RTT": 0,
					"Error": {},
					"Geo": null,
					"Lang": "",
					"MPLS": null
				}
			]
		],
		"TraceMapUrl": ""
	}`)

	result, err := ParseNextTraceOutput(jsonData)
	if err != nil {
		t.Fatalf("ParseNextTraceOutput failed: %v", err)
	}

	if len(result.Hops) != 1 {
		t.Errorf("Expected 1 hop, got %d", len(result.Hops))
	}

	hop := result.Hops[0]

	// Should have 33% packet loss (1 success out of 3)
	expectedLoss := 2.0 / 3.0
	tolerance := 0.01
	if hop.Loss < expectedLoss-tolerance || hop.Loss > expectedLoss+tolerance {
		t.Errorf("Expected loss %.2f, got %.2f", expectedLoss, hop.Loss)
	}

	// Should only have 1 RTT value (from the successful probe)
	if len(hop.RTT) != 1 {
		t.Errorf("Expected 1 RTT value, got %d", len(hop.RTT))
	}

	// RTT should be 1.0ms (1000000 nanoseconds)
	if len(hop.RTT) > 0 && hop.RTT[0] != 1.0 {
		t.Errorf("Expected RTT 1.0ms, got %.2f", hop.RTT[0])
	}
}

func TestParseNextTraceOutputAllTimeout(t *testing.T) {
	// Test with all probes timing out
	jsonData := []byte(`{
		"Hops": [
			[
				{
					"Success": false,
					"Address": null,
					"Hostname": "",
					"TTL": 1,
					"RTT": 0,
					"Error": {},
					"Geo": null,
					"Lang": "",
					"MPLS": null
				},
				{
					"Success": false,
					"Address": null,
					"Hostname": "",
					"TTL": 1,
					"RTT": 0,
					"Error": {},
					"Geo": null,
					"Lang": "",
					"MPLS": null
				},
				{
					"Success": false,
					"Address": null,
					"Hostname": "",
					"TTL": 1,
					"RTT": 0,
					"Error": {},
					"Geo": null,
					"Lang": "",
					"MPLS": null
				}
			]
		],
		"TraceMapUrl": ""
	}`)

	result, err := ParseNextTraceOutput(jsonData)
	if err != nil {
		t.Fatalf("ParseNextTraceOutput failed: %v", err)
	}

	if len(result.Hops) != 1 {
		t.Errorf("Expected 1 hop, got %d", len(result.Hops))
	}

	hop := result.Hops[0]

	// Should have 100% packet loss
	if hop.Loss != 1.0 {
		t.Errorf("Expected loss 1.0 (100%%), got %.2f", hop.Loss)
	}

	// Should have no RTT values
	if len(hop.RTT) != 0 {
		t.Errorf("Expected 0 RTT values, got %d", len(hop.RTT))
	}

	// IP should be empty
	if hop.IP != "" {
		t.Errorf("Expected empty IP, got %s", hop.IP)
	}
}
