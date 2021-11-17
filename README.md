# AmbiHASS

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
    "address": "http://hass.local:8123", // HASS adress
    "token": "", // Long lived access token
    "display": 1, // Display ID
    "lights": [
        {
            "id": "light.screen_back_middle", // Light entitity_id
            "minBrightness": 60, // minimal brightness
            "interval": 300 // update interval
        },
        {
            "id": "light.desk_backlight",
            "minBrightness": 60,
            "interval": 2000
        }
    ]
}
```

## Run

```sh
./dist/ambihass -c config.json
```