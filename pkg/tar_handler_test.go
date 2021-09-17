package octopus

import "testing"

type FixedRandomizer struct {
	randomString string
}

func (f *FixedRandomizer) GenerateRandomString(l int) string {
	return f.randomString
}

func TestCopiedTar(t *testing.T) {
	randomizer := &FixedRandomizer{randomString: "random-string"}
	NewTarHandlerWithRandomizer("/tmp", randomizer)
}
