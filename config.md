# Environment Variables

## Config

 - `LISTEN_ADDRESS` (**required**, non-empty, default: `:8080`) - HTTP server bind address.
 - `REAL_IP_HEADER` (default: `false`) - Get client IP address from the "Real-IP" header.
 - `ACCESS_TOKEN` - Token required to access the API. If set, the value must be provided as a `token` query parameter.
 - `PUBLIC_URL` (**required**, non-empty) - This app's public URL.
 - `NIGHTSCOUT_URL` (**required**, non-empty) - Nightscout base URL
 - `NIGHTSCOUT_TOKEN` - Nightscout token. Using an access token is recommended instead of the API secret.
 - `NIGHTSCOUT_UNITS` - Blood sugar unit. (one of: mg/dL, mmol/L)
 - `FETCH_DELAY` (default: `30s`) - Time to wait before the next reading should be ready. In testing, this seems to be about 20s behind, so the default is 30s to be safe. Your results may vary.
 - `FALLBACK_INTERVAL` (default: `30s`) - Normally, readings will be fetched when ready (after ~5m). This interval will be used if the next reading time cannot be estimated due to sensor warm-up, missed readings, errors, etc.
 - `TIME_FORMAT` (default: `3:04 PM`) - Customize the time format. Use `3:04 PM` for 12-hour time or `15:04` for 24-hour. See [time](https://pkg.go.dev/time) for more details.
 - `GRAPH_DURATION` (default: `6h`) - How far back in time the graph should go.
 - `HIGH_THRESHOLD` (default: `200`) - Where to draw the upper line.
 - `HIGH_BACKGROUND_SHADE` (default: `250`) - Background color above the high threshold line. Value must be between 0-255.
 - `LOW_THRESHOLD` (default: `70`) - Where to draw the lower line.
 - `LOW_BACKGROUND_SHADE` (default: `247`) - Background color below the low threshold line. Value must be between 0-255.
 - `INVERT` - Render with a black background and a white foreground.
 - `INVERT_BELOW` (default: `55`) - Invert colors when below this value. (Stacks with the `INVERT` option)
 - `INVERT_ABOVE` (default: `300`) - Invert colors when above this value. (Stacks with the `INVERT` option)

