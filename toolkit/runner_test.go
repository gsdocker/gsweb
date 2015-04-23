package toolkit

import "testing"

func TestRunner(t *testing.T) {
	app, err := LoadApp("github.com/gsdocker/gsweb/sample")

	if err != nil {
		t.Fatal(err)
	}

	runner, err := NewAppRunner(app)

	if err != nil {
		t.Fatal(err)
	}

	runner.Run()
}
