# Query Parameters

The following GET parameters can be passed to the JSON or image endpoint. These values will override the env config defined in [`envs.md`](envs.md).

- `token` - Access token. Required if the `ACCESS_TOKEN` env is set.
- `unit` - Blood sugar unit. (one of: mg/dL, mmol/L)
- `time_format` (default `3:04 PM`) - Customize the time format. Use `3:04 PM` for 12-hour time or `15:04` for 24-hour. See [time](https://pkg.go.dev/time) for more details.
- `graph_duration` (default `6h`) - How far back in time the graph should go.
- `graph_min` (default `40`) - Minimum X-axis value.
- `graph_max` (default `300`) - Maximum X-axis value.
- `point_stroke_radius` (default `4`) - Control the plot point stroke radius. Set to 0 to disable.
- `high_threshold` (default `200`) - Where to draw the upper line.
- `low_threshold` (default `70`) - Where to draw the lower line.
- `invert` - Render with a black background and a white foreground.
- `invert_below` (default `55`) - Invert colors when below this value. (Stacks with the `INVERT` option)
- `invert_above` (default `300`) - Invert colors when above this value. (Stacks with the `INVERT` option)
- `color_mode` (default `2bit`) - Output color mode. 2-bit will be antialiased and dithering will be higher quality, but requires TRMNL firmware v1.6.0+. (one of 1bit, 2bit)
