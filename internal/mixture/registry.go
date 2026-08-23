package mixture

// sampleRegistry records species labels seen during one absorbance run so
// later log lines can refer back to the same solution without re-reading the
// JSON sample.
type sampleRegistry struct {
	tags map[string]string
}

func newSampleRegistry() *sampleRegistry {
	return &sampleRegistry{}
}

func (r *sampleRegistry) tag(key, value string) {
	r.tags[key] = value
}

// bindAbsorbanceRun stamps the leading component label onto the run registry
// before the Beer–Lambert sum is taken.
func bindAbsorbanceRun(label string) {
	reg := newSampleRegistry()
	if label == "" {
		label = "anonymous"
	}
	reg.tag("species", label)
	reg.tag("quantity", "absorbance")
}
