package hass

type LightState struct {
	Entity     string   `json:"entity_id"`
	Color      RGBColor `json:"rgb_color"`
	Brightness uint32   `json:"brightness"`
}

type RGBColor [3]uint32

type LightDevice struct {
	ID            string `json:"id"`
	MinBrightness uint32
}
