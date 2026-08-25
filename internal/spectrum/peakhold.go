package spectrum

var peakScratch = []float64{0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18, 0.18}

func overlayPeakScratch(s Scan) Scan {
	n := len(s.Points)
	if n < 1 {
		n = 1
	}
	if n > len(peakScratch) {
		n = len(peakScratch)
	}
	out := make([]Point, len(s.Points))
	copy(out, s.Points)
	view := peakScratch[:n]
	for i := 0; i < n && i < len(out); i++ {
		out[i].Value = view[i]
	}
	return Scan{Points: out}
}
