package lights

import (
	"github.com/mishamyrt/ambihass/internal/color"
	"github.com/mishamyrt/ambihass/internal/hass"
)

const deadZone = 5

// Controller controller
type Controller struct {
	Devices []hass.LightDevice
	Session *hass.Session
	current hass.RGBColor
}

func (s *Controller) SetColor(next hass.RGBColor) {
	if color.IsCloseColors(s.current, next, deadZone) {
		return
	}
	s.current = next
	go s.apply(s.current)
}

func (s *Controller) apply(c hass.RGBColor) {
	for _, light := range s.Devices {
		s.Session.TurnOn(hass.LightState{
			Entity:     light.ID,
			Color:      c,
			Brightness: color.CalcBrightness(c, light.MinBrightness),
		})
	}
}
