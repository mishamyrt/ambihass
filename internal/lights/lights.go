package lights

import (
	"sync"
	"time"

	"github.com/mishamyrt/ambihass/internal/color"
	"github.com/mishamyrt/ambihass/internal/hass"
)

const deadZone = 7

type deviceController struct {
	Device      hass.LightDevice
	lastUpdated time.Time
	wai         sync.Mutex
}

func (s *deviceController) SetColor(next hass.RGBColor) {

}

// Controller controller
type Controller struct {
	Devices     []hass.LightDevice
	Session     *hass.Session
	current     []hass.RGBColor
	nextUpdate  []time.Time
	needsUpdate []bool
}

func (s *Controller) Start(interval int, ch <-chan []hass.RGBColor) {
	s.needsUpdate = make([]bool, len(s.Devices))
	s.nextUpdate = make([]time.Time, len(s.Devices))
	go s.listenColors(ch)
	s.mainLoop(interval)
}

func (s *Controller) mainLoop(interval int) {
	duration := time.Duration(interval) * time.Millisecond
	for {
		time.Sleep(duration)
		for i, device := range s.Devices {
			if s.needsUpdate[i] && s.nextUpdate[i].Before(time.Now()) {
				colorIndex := 0
				if device.Color > len(s.current)-1 {
					colorIndex = len(s.current) - 1
				} else {
					colorIndex = device.Color
				}
				s.apply(device, s.current[colorIndex])
				s.nextUpdate[i] = time.Now().Add(time.Duration(device.Interval) * time.Millisecond)
				s.needsUpdate[i] = false
			}
		}
	}
}

func (s *Controller) listenColors(ch <-chan []hass.RGBColor) {
	for {
		select {
		case colors := <-ch:
			s.setColor(colors)
		}
	}
}

func (s *Controller) setColor(next []hass.RGBColor) {
	if len(s.current) > 0 && color.IsCloseColors(s.current[0], next[0], deadZone) {
		return
	}
	s.current = next
	s.setDirty()
}

func (s *Controller) setDirty() {
	for i := range s.Devices {
		s.needsUpdate[i] = true
	}
}

func (s *Controller) apply(light hass.LightDevice, c hass.RGBColor) {
	s.Session.TurnOn(hass.LightState{
		Entity: light.ID,
		Color:  c,
		Brightness: uint32(
			float32(color.CalcBrightness(c, light.BrightnessMin)) * light.Brightness,
		),
	})
}
