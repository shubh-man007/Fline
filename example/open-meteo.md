```json
{
  "name": "Open-Meteo Forecast",
  "steps": [
    {
      "type": "http_request",
      "method": "GET",
      "URL": "https://api.open-meteo.com/v1/forecast?latitude={{lat}}&longitude={{lon}}&daily=temperature_2m_max&timezone=auto",
      "headers": {},
      "timeout": 5000,
      "retries": 3
    },
    {
      "type": "transform",
      "ops": [
        {
          "op": "pick",
          "paths": ["result"]
        }
      ]
    }
  ]
}
```

```json
{
    "result": {
        "daily": {
            "temperature_2m_max": [
                41.2,
                42.3,
                43.3,
                44.3,
                43.6,
                42.5,
                42.6
            ],
            "time": [
                "2026-04-23",
                "2026-04-24",
                "2026-04-25",
                "2026-04-26",
                "2026-04-27",
                "2026-04-28",
                "2026-04-29"
            ]
        },
        "daily_units": {
            "temperature_2m_max": "°C",
            "time": "iso8601"
        },
        "elevation": 53,
        "generationtime_ms": 0.04029273986816406,
        "latitude": 23,
        "longitude": 72.625,
        "timezone": "Asia/Kolkata",
        "timezone_abbreviation": "GMT+5:30",
        "utc_offset_seconds": 19800
    }
}
```