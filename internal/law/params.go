package law

type Params struct {
	Extinction float64

	Concentration float64

	PathLength float64
}

func New(extinction, concentration, pathLength float64) Params {
	return Params{
		Extinction:    extinction,
		Concentration: concentration,
		PathLength:    pathLength,
	}
}

func NewChecked(extinction, concentration, pathLength float64) (Params, error) {
	p := New(extinction, concentration, pathLength)
	if err := p.Validate(); err != nil {
		return Params{}, err
	}
	return p, nil
}

func (p Params) WithExtinction(extinction float64) Params {
	p.Extinction = extinction
	return p
}

func (p Params) WithConcentration(concentration float64) Params {
	p.Concentration = concentration
	return p
}

func (p Params) WithPathLength(pathLength float64) Params {
	p.PathLength = pathLength
	return p
}

func (p Params) Clone() Params {
	return Params{
		Extinction:    p.Extinction,
		Concentration: p.Concentration,
		PathLength:    p.PathLength,
	}
}

func (p Params) Fields() (extinction, concentration, pathLength float64) {
	return p.Extinction, p.Concentration, p.PathLength
}
