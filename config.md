# Environment Variables

## Config

 - `NIGHTSCOUT_URL` (**required**) - Nightscout base URL
 - `IMAGE_URL` (**required**) - A URL that the TRMNL device can use to download the image from this app. It can be a public URL or an internal IP address as long as the TRMNL device is on the same network.
 - `NIGHTSCOUT_TOKEN` - Nightscout token. Using an access token is recommended instead of the API secret.
 - `LISTEN_ADDRESS` (default: `:8080`) - HTTP server bind address.
 - `ACCESS_TOKEN` - Token required to access the API. If set, the value must be provided as a `token` query parameter.
 - `REAL_IP_HEADER` (default: `false`) - Get client IP address from the "Real-IP" header.
 - `UPDATE_INTERVAL` (default: `5m`) - The interval that new readings are sent to Nightscout.
 - `FETCH_DELAY` (default: `30s`) - Time to wait before the next reading should be ready. In testing, this seems to be about 20s behind, so the default is 30s to be safe. Your results may vary.
 - `FALLBACK_INTERVAL` (default: `30s`) - Normally, readings will be fetched when ready (after ~5m). This interval will be used if the next reading time cannot be estimated due to sensor warm-up, missed readings, errors, etc.
 - `UNIT` - Blood sugar unit. (one of: mg/dL, mmol/L)
 - `TIME_FORMAT` (default: `3:04 PM`) - Customize the time format. Use `3:04 PM` for 12-hour time or `15:04` for 24-hour. See [time](https://pkg.go.dev/time) for more details.
 - `GRAPH_DURATION` (default: `6h`) - How far back in time the graph should go.
 - `GRAPH_MIN` (default: `40`) - Minimum X-axis value.
 - `GRAPH_MAX` (default: `300`) - Maximum X-axis value.
 - `HIGH_THRESHOLD` (default: `200`) - Where to draw the upper line.
 - `LOW_THRESHOLD` (default: `70`) - Where to draw the lower line.
 - `INVERT` - Render with a black background and a white foreground.
 - `INVERT_BELOW` (default: `55`) - Invert colors when below this value. (Stacks with the `INVERT` option)
 - `INVERT_ABOVE` (default: `300`) - Invert colors when above this value. (Stacks with the `INVERT` option)
 - `COLOR_MODE` (default: `1bit`) - Output color mode. 2-bit will be antialiased and dithering will be higher quality, but requires TRMNL firmware v1.6.0+. (one of 1bit, 2bit)

