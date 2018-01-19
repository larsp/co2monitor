package meter

import (
	"os"
	"syscall"
	"unsafe"

	"log"

	"crypto/rand"

	"sync/atomic"

	"github.com/pkg/errors"
)

const (
	meterTemp       byte    = 0x42
	meterCO2        byte    = 0x50
	hidiocsfeature9 uintptr = 0xc0094806
)

var key = [8]byte{}

// Meter gives access to the CO2 Meter. Make sure to call Open before Read.
type Meter struct {
	file   *os.File
	opened int32
}

// Measurement is the result of a Read operation.
type Measurement struct {
	Temperature float64
	Co2         int
}

// Open will open the device file specified in the path which is usually something like /dev/hidraw2.
func (m *Meter) Open(path string) (err error) {
	atomic.StoreInt32(&m.opened, 1)
	m.initKey()

	m.file, err = os.OpenFile(path, os.O_RDWR, 0644)

	if err != nil || m.file == nil {
		return errors.Wrapf(err, "Failed to open '%v'", path)
	}

	log.Printf("Device '%v' opened", m.file.Name())
	return m.ioctl()
}

// initKey writes 8 bytes entropy to the global key variable. A static key would be sufficient, but lets stick with
// real randomness
func (m *Meter) initKey() {
	_, err := rand.Read(key[:])
	if err != nil {
		panic(err)
	}
}

// ioctl writes into the device file. We need to write 9 bytes where the first byte specifies the report number.
// In this case 0x00.
func (m *Meter) ioctl() error {
	data := [9]byte{}
	copy(data[1:], key[0:]) // remember, first byte needs to be 0
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, m.file.Fd(), hidiocsfeature9, uintptr(unsafe.Pointer(&data)))

	if err != 0 {
		return errors.Wrap(syscall.Errno(err), "ioctl failed")
	}
	return nil
}

// Read will read from the device file until it finds a temperature and co2 measurement. Before it can be used the
// device file needs to be opened via Open.
func (m *Meter) Read() (*Measurement, error) {
	if atomic.LoadInt32(&m.opened) != 1 {
		return nil, errors.New("Device needs to be opened")
	}

	result := make([]byte, 8)
	measurement := &Measurement{Co2: 0, Temperature: -273.15}

	for {
		_, err := m.file.Read(result)
		if err != nil {
			return nil, errors.Wrapf(err, "Could not read from: '%v'", m.file.Name())
		}

		decrypted := m.decrypt(result)

		operation := decrypted[0]
		value := decrypted[1]<<8 | decrypted[2]

		switch byte(operation) {
		case meterCO2:
			measurement.Co2 = int(value)
		case meterTemp:
			measurement.Temperature = float64(value)/16.0 - 273.15
		}

		if measurement.Co2 != 0 && measurement.Temperature != -273.15 {
			return measurement, nil
		}
	}
}

// decrypt is a clone of the python decrypt function of the original article: https://hackaday.io/project/5301-reverse-engineering-a-low-cost-usb-co-monitor/log/17909-all-your-base-are-belong-to-us
func (m *Meter) decrypt(data []byte) []uint {
	state := []uint{0x48, 0x74, 0x65, 0x6D, 0x70, 0x39, 0x39, 0x65}
	shuffle := []int{2, 4, 0, 7, 1, 6, 5, 3}

	phase1 := make([]uint, 8)
	for i := range shuffle {
		phase1[shuffle[i]] = uint(data[i])
	}

	phase2 := make([]uint, 8)
	for i := 0; i < 8; i++ {
		phase2[i] = phase1[i] ^ uint(key[i])
	}

	phase3 := make([]uint, 8)
	for i := 0; i < 8; i++ {
		phase3[i] = ((phase2[i] >> 3) | (phase2[(i-1+8)%8] << 5)) & 0xff
	}

	tmp := make([]uint, 8)
	for i := 0; i < 8; i++ {
		tmp[i] = ((state[i] >> 4) | (state[i] << 4)) & 0xff
	}

	result := make([]uint, 8)
	for i := 0; i < 8; i++ {
		result[i] = (0x100 + phase3[i] - tmp[i]) & 0xff
	}
	return result
}

// Close will close the device file.
func (m *Meter) Close() error {
	log.Printf("Closing '%v'", m.file.Name())
	atomic.StoreInt32(&m.opened, 0)
	return m.file.Close()
}
