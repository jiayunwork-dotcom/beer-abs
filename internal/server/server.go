package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"beer-abs/internal/band"
	"beer-abs/internal/mixture"
)

type Config struct {
	Addr string
}

func ListenAndServe(cfg Config) error {
	mux := New(cfg)
	return http.ListenAndServe(cfg.Addr, mux)
}

func New(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/absorbance", handleAbsorbance)
	return mux
}

func FormatAddr(addr string) string {
	port := ParsePort(addr)
	if port == 0 {
		return addr
	}
	return fmt.Sprintf("http://0.0.0.0:%d", port)
}

func ParsePort(addr string) int {
	parts := strings.Split(addr, ":")
	if len(parts) < 2 {
		return 0
	}
	p, _ := strconv.Atoi(parts[len(parts)-1])
	return p
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type componentReq struct {
	Label          string  `json:"label"`
	Extinction     float64 `json:"extinction"`
	ExtinctionLow  float64 `json:"extinction_low"`
	ExtinctionHigh float64 `json:"extinction_high"`
	Concentration  float64 `json:"concentration"`
}

type bandReq struct {
	Center    float64 `json:"center"`
	HalfWidth float64 `json:"half_width"`
}

type absorbanceRequest struct {
	Components    []componentReq `json:"components"`
	PathLength    float64        `json:"path_length"`
	StrayFraction float64        `json:"stray_fraction"`
	Band          *bandReq       `json:"band"`
}

func handleAbsorbance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		httpError(w, http.StatusMethodNotAllowed, "POST required")
		return
	}
	var req absorbanceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Components) == 0 {
		httpError(w, http.StatusBadRequest, "components array is empty")
		return
	}

	var rect *band.RectBand
	if req.Band != nil {
		b, err := band.NewRectBand(req.Band.Center, req.Band.HalfWidth)
		if err != nil {
			httpError(w, http.StatusBadRequest, err.Error())
			return
		}
		rect = &b
	}

	comps := make([]mixture.Component, len(req.Components))
	for i, cs := range req.Components {
		ext := cs.Extinction
		if rect != nil {
			ext = band.BandAverageExtinction(cs.ExtinctionLow, cs.ExtinctionHigh)
		}
		comps[i] = mixture.Component{
			Label:         cs.Label,
			Extinction:    ext,
			Concentration: cs.Concentration,
		}
	}

	m, err := mixture.New(comps, req.PathLength, req.StrayFraction)
	if err != nil {
		httpError(w, http.StatusBadRequest, err.Error())
		return
	}
	analysis, err := m.Analyze()
	if err != nil {
		httpError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, analysis)
}

func httpError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}
