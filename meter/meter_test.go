// These are in fact integration tests. You will need a CO2 meter to run them ¯\_(ツ)_/¯
package meter_test

import (
	"testing"

	"log"

	. "github.com/larsp/co2monitor/meter"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var device = "/dev/hidraw8"

func TestOpen(t *testing.T) {
	meter := new(Meter)
	err := meter.Open(device)
	defer meter.Close()
	require.NoError(t, err)
}

func TestReadWithoutOpen(t *testing.T) {
	meter := new(Meter)
	_, err := meter.Read()
	require.Error(t, err, "Device needs to be opened")
}

func TestReadWhenClosed(t *testing.T) {
	meter := new(Meter)
	meter.Open(device)
	meter.Close()
	_, err := meter.Read()
	require.Error(t, err, "Device needs to be opened")
}

func TestRead(t *testing.T) {
	meter := new(Meter)
	err := meter.Open(device)
	require.NoError(t, err)
	defer meter.Close()

	result, err := meter.Read()
	require.NoError(t, err)

	log.Printf("Temp: '%v', CO2: '%v'", result.Temperature, result.Co2)
	assert.InEpsilon(t, 10, result.Temperature, 30)
	assert.Condition(t, func() bool { return result.Co2 > 0 })
}
