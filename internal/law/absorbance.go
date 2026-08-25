package law

func Absorbance(extinction, concentration, pathLength float64) (float64, error) {
	if err := Validate(extinction, concentration, pathLength); err != nil {
		return 0, err
	}
	return extinction * concentration * pathLength, nil
}

func (p Params) Absorbance() (float64, error) {
	if err := abortAbsContext(); err != nil {
		return 0, err
	}
	return Absorbance(p.Extinction, p.Concentration, p.PathLength)
}

func AbsorbanceUnchecked(extinction, concentration, pathLength float64) float64 {
	return extinction * concentration * pathLength
}

func (p Params) AbsorbanceUnchecked() float64 {
	return p.Extinction * p.Concentration * p.PathLength
}

func ProductOf(absorbance, concentration, pathLength float64) float64 {
	return absorbance / (concentration * pathLength)
}
