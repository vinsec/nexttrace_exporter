package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
)

// NextTraceRawResult represents the raw JSON output from nexttrace -j
type NextTraceRawResult struct {
	Hops        [][]HopDetail `json:"Hops"`
	TraceMapUrl string        `json:"TraceMapUrl"`
}

// HopDetail represents a single probe result at a specific TTL
type HopDetail struct {
	Success  bool     `json:"Success"`
	Address  *IPAddr  `json:"Address"`
	Hostname string   `json:"Hostname"`
	TTL      int      `json:"TTL"`
	RTT      int64    `json:"RTT"` // RTT in nanoseconds
	Error    any      `json:"Error"`
	Geo      *GeoInfo `json:"Geo"`
	Lang     string   `json:"Lang"`
	MPLS     any      `json:"MPLS"`
}

// IPAddr represents an IP address
type IPAddr struct {
	IP   string `json:"IP"`
	Zone string `json:"Zone"`
}

// GeoInfo represents geographic information
type GeoInfo struct {
	IP        string  `json:"ip"`
	ASNumber  string  `json:"asnumber"`
	Country   string  `json:"country"`
	CountryEn string  `json:"country_en"`
	Prov      string  `json:"prov"`
	ProvEn    string  `json:"prov_en"`
	City      string  `json:"city"`
	CityEn    string  `json:"city_en"`
	District  string  `json:"district"`
	Owner     string  `json:"owner"`
	ISP       string  `json:"isp"`
	Domain    string  `json:"domain"`
	Whois     string  `json:"whois"`
	Lat       float64 `json:"lat"`
	Lng       float64 `json:"lng"`
	Prefix    string  `json:"prefix"`
}

// NextTraceResult represents the processed result
type NextTraceResult struct {
	Target string `json:"target"`
	Hops   []Hop  `json:"hops"`
}

// Hop represents aggregated data for a single hop (TTL level)
type Hop struct {
	TTL      int       `json:"ttl"`
	IP       string    `json:"ip"`
	Hostname string    `json:"hostname"`
	RTT      []float64 `json:"rtt"` // RTT in milliseconds
	Loss     float64   `json:"loss"`
	ASN      string    `json:"asn"`
	Location string    `json:"location"`
}

// ParseNextTraceOutput parses the JSON output from nexttrace -j command
func ParseNextTraceOutput(data []byte) (*NextTraceResult, error) {
	// Clean the output: remove ANSI escape sequences and extract only the JSON part
	cleanedData := cleanNextTraceOutput(data)

	var raw NextTraceRawResult
	if err := json.Unmarshal(cleanedData, &raw); err != nil {
		return nil, fmt.Errorf("failed to parse nexttrace JSON: %w", err)
	}

	result := &NextTraceResult{
		Hops: make([]Hop, 0, len(raw.Hops)),
	}

	// Process each TTL level
	for _, probes := range raw.Hops {
		if len(probes) == 0 {
			continue
		}

		hop := Hop{
			RTT: make([]float64, 0, len(probes)),
		}

		successCount := 0
		var firstValidIP string
		var firstValidHostname string
		var firstValidASN string
		var firstValidLocation string

		// Aggregate data from all probes at this TTL
		for _, probe := range probes {
			hop.TTL = probe.TTL

			if probe.Success && probe.Address != nil {
				successCount++

				// Convert RTT from nanoseconds to milliseconds
				if probe.RTT > 0 {
					rttMs := float64(probe.RTT) / 1_000_000.0
					hop.RTT = append(hop.RTT, rttMs)
				}

				// Use the first valid IP/hostname/ASN we see
				if firstValidIP == "" && probe.Address.IP != "" {
					firstValidIP = probe.Address.IP
				}
				if firstValidHostname == "" && probe.Hostname != "" {
					firstValidHostname = probe.Hostname
				}

				// Extract geo information
				if probe.Geo != nil {
					if firstValidASN == "" && probe.Geo.ASNumber != "" {
						firstValidASN = probe.Geo.ASNumber
					}
					if firstValidLocation == "" {
						if probe.Geo.CityEn != "" && probe.Geo.CountryEn != "" {
							firstValidLocation = probe.Geo.CityEn + ", " + probe.Geo.CountryEn
						} else if probe.Geo.CountryEn != "" {
							firstValidLocation = probe.Geo.CountryEn
						}
					}
				}
			}
		}

		hop.IP = firstValidIP
		hop.Hostname = firstValidHostname
		hop.ASN = firstValidASN
		hop.Location = firstValidLocation

		// Calculate packet loss ratio
		totalProbes := len(probes)
		if totalProbes > 0 {
			hop.Loss = float64(totalProbes-successCount) / float64(totalProbes)
		}

		// Only add hops with valid data
		if hop.TTL > 0 {
			result.Hops = append(result.Hops, hop)
		}
	}

	return result, nil
}

// AverageRTT calculates the average RTT from a slice of RTT values
func (h *Hop) AverageRTT() float64 {
	if len(h.RTT) == 0 {
		return 0.0
	}

	var sum float64
	for _, rtt := range h.RTT {
		sum += rtt
	}
	return sum / float64(len(h.RTT))
}

// HasValidIP checks if the hop has a valid IP address
func (h *Hop) HasValidIP() bool {
	return h.IP != "" && h.IP != "*"
}

// cleanNextTraceOutput removes ANSI escape sequences and extracts the JSON part
func cleanNextTraceOutput(data []byte) []byte {
	// Remove ANSI escape sequences (color codes)
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	cleaned := ansiRegex.ReplaceAll(data, []byte(""))

	// Find the first { which marks the start of JSON
	jsonStart := bytes.IndexByte(cleaned, '{')
	if jsonStart == -1 {
		return data // Return original if no JSON found
	}

	// Extract from { to the end
	return cleaned[jsonStart:]
}
