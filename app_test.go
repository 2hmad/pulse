package pulse

import (
	"testing"
)

var app *Pulse

func init() {
	app = New(Config{
		AppName: "Test App",
	})
}

func TestPulse_Run(t *testing.T) {
	err := app.Run("127.0.0.1:8082")
	if err != nil {
		return
	}
}
