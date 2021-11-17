# Ambihass

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
        "light.screen_back_middle" // Lights entities
    ]
}
```

## Run

```sh
./dist/ambihass -c config.json
```