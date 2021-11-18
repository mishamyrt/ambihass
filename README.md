# AmbiHASS [![Lint](https://github.com/mishamyrt/ambihass/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/mishamyrt/ambihass/actions/workflows/golangci-lint.yml)

Ambilight for Home Assistant.

The app runs on the user's computer and communicates with the Home Assistant server, setting the right colours for the lighting.

## Build

```sh
make dist/ambihass
```

## Configure

Create a `config.json` file

```js
{
    "address": "http://hass.local:8123", // HASS address
    "token": "", // Long lived access token
    "display": 1, // Display ID
    "lights": [
        {
            "id": "light.screen_back_middle",
            "brightnessMin": 60,
            "interval": 300,
            "color": 0,
            "brightness": 1
        },
        {
            "id": "light.desk_backlight",
            "brightnessMin": 20,
            "interval": 2000,
            "color": 1,
            "brightness": 0.3
        }
    ]
}
```

### Lights params

* `id` — Light entitity_id.
* `brightnessMin` — Minimal brightness.
* `interval` — Device update interval. Some devices do not respond well to too frequent updates.
* `color` — Color level. If more than one color is recognised in the image, these additional colors can be applied to the lights. Counts from zero.
* `brightness` — Brightness level. Allows to better adapt lighting to different brightness levels.

## Run

```sh
./dist/ambihass -c config.json
```