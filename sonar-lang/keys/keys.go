package keys

import "os"

type keys struct {
	MODE string
}

var mode string = os.Getenv("SONAR_MODE") // DEV or NON_DEV

var Keys *keys = &keys{}

func init() {
	if len(mode) == 0 {
		mode = "NON_DEV"
	}
	Keys.MODE = mode
}
