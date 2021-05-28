package main

import "github.com/danesparza/fxtrigger/cmd"

// @title fxTrigger
// @version 1.0
// @description fxTrigger REST based management for GPIO/Sensor -> endpoint triggers (on Raspberry Pi)

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /v1
func main() {
	cmd.Execute()
}
