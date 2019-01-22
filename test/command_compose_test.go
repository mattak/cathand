package cathandtest

import (
	"github.com/mattak/cathand/pkg/cathand"
	"github.com/stretchr/testify/assert"
	"testing"
)

type ComposeTestContext struct {
}

func (ComposeTestContext) setup() {
}

func (ComposeTestContext) tearDown() {
}

func (ComposeTestContext) containsKey(maps map[string][]cathand.Event, key string) bool {
	_, ok := maps[key]
	return ok
}

func TestParseEvent(t *testing.T) {
	context := ComposeTestContext{}
	context.setup()
	defer context.tearDown()

	data := `add device 1: /dev/input/event8
  name:     "uinput-fpc"
add device 2: /dev/input/event7
  name:     "uinput-folio"
add device 3: /dev/input/event6
  name:     "msm8996-tasha-marlin-snd-card Button Jack"
add device 4: /dev/input/event5
  name:     "msm8996-tasha-marlin-snd-card Headset Jack"
add device 5: /dev/input/event1
  name:     "HDMI CEC User or Deck Control"
could not get driver version for /dev/input/mice, Not a typewriter
add device 6: /dev/input/event2
  name:     "synaptics_dsxv26"
add device 7: /dev/input/event3
  name:     "STM VL53L0 proximity sensor"
add device 8: /dev/input/event0
  name:     "qpnp_pon"
add device 9: /dev/input/event4
  name:     "gpio-keys"
[  149840.463459] /dev/input/event2: 0003 0039 0000107b
[  149840.463459] /dev/input/event2: 0003 0035 0000026c
[  149840.463459] /dev/input/event2: 0003 0036 000002bc
[  149840.463459] /dev/input/event2: 0003 003a 00000036
[  149840.463459] /dev/input/event2: 0000 0000 00000000
[  149840.479355] /dev/input/event2: 0003 003a 0000003a
[  149840.479355] /dev/input/event2: 0000 0000 00000000
[  149840.487901] /dev/input/event2: 0003 003a 00000039
[  149840.487901] /dev/input/event2: 0000 0000 00000000
[  149840.496388] /dev/input/event2: 0003 003a 00000038
[  149840.496388] /dev/input/event2: 0000 0000 00000000
[  149840.504466] /dev/input/event2: 0003 003a 00000036
[  149840.504466] /dev/input/event2: 0000 0000 00000000
[  149840.512122] /dev/input/event2: 0003 0039 ffffffff
[  149840.512122] /dev/input/event2: 0000 0000 00000000
`
	eventsMap, err := cathand.ParseEvent(data)
	if err != nil {
		t.Fatal("ParseEvent failed")
	}

	if !context.containsKey(eventsMap, "event2") {
		t.Fatal("event2 not exists")
	}

	events := eventsMap["event2"]

	if len(events) != 15 {
		t.Fatal("parsed events length is not 15")
	}

	assert.Equal(t, int64(149840), events[0].EpochSec)
	assert.Equal(t, int64(463459), events[0].EpochUsec)
	assert.Equal(t, uint16(3), events[0].Type)
	assert.Equal(t, uint16(0x39), events[0].Code)
	assert.Equal(t, uint32(0x107b), events[0].Value)

	assert.Equal(t, int64(149840), events[14].EpochSec)
	assert.Equal(t, int64(512122), events[14].EpochUsec)
	assert.Equal(t, uint16(0), events[14].Type)
	assert.Equal(t, uint16(0x00), events[14].Code)
	assert.Equal(t, uint32(0x00), events[14].Value)
}
