package main

import "testing"
import "github.com/stretchr/testify/assert"

func TestNewSensorTemp(t *testing.T) {
	var testTemp = "+45.8°C  (crit = +90.0°C)"

	var output = NewSensorTemp(testTemp)
	assert.Equal(t, "45.8", output.Temperature)
	assert.Equal(t, "90.0", output.CriticalTemp)
}
