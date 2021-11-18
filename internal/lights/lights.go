package lights

import (
	"math"
	"time"

	"github.com/mishamyrt/ambihass/internal/color"
	"github.com/mishamyrt/ambihass/internal/hass"
	"github.com/mishamyrt/ambihass/internal/log"
)

const deadZone = 7

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
	var err error
	for {
		time.Sleep(duration)
		for i, device := range s.Devices {
			if s.needsUpdate[i] && s.nextUpdate[i].Before(time.Now()) {
				err = s.apply(device, s.current[int(
					math.Min(
						float64(len(s.current)-1),
						float64(device.Color),
					),
				)])
				if err != nil {
					log.Debug("Update fail ", device.ID)
				}
				s.nextUpdate[i] = time.Now().Add(time.Duration(device.Interval) * time.Millisecond)
				s.needsUpdate[i] = false
			}
		}
	}
}

func (s *Controller) listenColors(ch <-chan []hass.RGBColor) {
	for {
		s.setColor(<-ch)
	}
}

func (s *Controller) setColor(next []hass.RGBColor) {
	if len(s.current) > 0 && color.IsCloseColors(s.current[0], next[0], deadZone) {
		return
	}
	log.Debug("New colors: ", next)
	s.current = next
	s.setDirty()
}

func (s *Controller) setDirty() {
	for i := range s.Devices {
		s.needsUpdate[i] = true
	}
}

func (s *Controller) apply(light hass.LightDevice, c hass.RGBColor) error {
	return s.Session.TurnOn(hass.LightState{
		Entity: light.ID,
		Color:  c,
		Brightness: uint32(
			math.Min(
				float64(color.CalcBrightness(c, light.BrightnessMin))*light.Brightness,
				255,
			),
		),
	})
}
