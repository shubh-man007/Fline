# AccuWeather API Doc

## Create Workflow
```json
{
  "name": "Accuweather Forecast Ahmedabad",
  "steps": [
    {
      "type": "http_request",
      "method": "GET",
      "URL": "https://dataservice.accuweather.com/locations/v1/cities/search?q=Ahmedabad",
      "headers": {
        "Authorization": "Bearer $ACCUWEATHER_API_KEY"
      },
      "timeout": 5000,
      "retries": 3
    },
    {
      "type": "transform",
      "ops": [
        {
          "op": "template",
          "to": "weatherURL",
          "template": "https://dataservice.accuweather.com/forecasts/v1/daily/5day/{{result.0.Key}}"
        }
      ]
    },
    {
      "type": "http_request",
      "method": "GET",
      "URL": "{{weatherURL}}",
      "headers": {
        "Authorization": "Bearer $ACCUWEATHER_API_KEY"
      },
      "timeout": 5000,
      "retries": 3
    }
  ]
}
```

## Trigger Workflow
```json
{
    "id": "dd6d8c72-efe3-44d3-bf4d-12e370b18deb",
    "name": "Accuweather Forecast Ahmedabad",
    "enabled": true,
    "steps": [
        {
            "type": "http_request",
            "method": "GET",
            "URL": "https://dataservice.accuweather.com/locations/v1/cities/search?q=Ahmedabad",
            "headers": {
                "Authorization": "Bearer $ACCUWEATHER_API_KEY"
            },
            "timeout": 5000,
            "retries": 3
        },
        {
            "type": "transform",
            "ops": [
                {
                    "op": "template",
                    "to": "weatherURL",
                    "template": "https://dataservice.accuweather.com/forecasts/v1/daily/5day/{{result.0.Key}}"
                }
            ]
        },
        {
            "type": "http_request",
            "method": "GET",
            "URL": "{{weatherURL}}",
            "headers": {
                "Authorization": "Bearer $ACCUWEATHER_API_KEY"
            },
            "timeout": 5000,
            "retries": 3
        }
    ]
}
```

## Trigger Response
```json
{
    "result": {
        "DailyForecasts": [
            {
                "Date": "2026-04-19T07:00:00+05:30",
                "Day": {
                    "HasPrecipitation": false,
                    "Icon": 4,
                    "IconPhrase": "Intermittent clouds"
                },
                "EpochDate": 1776562200,
                "Link": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?day=1&lang=en-us",
                "MobileLink": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?day=1&lang=en-us",
                "Night": {
                    "HasPrecipitation": false,
                    "Icon": 33,
                    "IconPhrase": "Clear"
                },
                "Sources": [
                    "AccuWeather"
                ],
                "Temperature": {
                    "Maximum": {
                        "Unit": "F",
                        "UnitType": 18,
                        "Value": 107
                    },
                    "Minimum": {
                        "Unit": "F",
                        "UnitType": 18,
                        "Value": 77
                    }
                }
            },
            {
                "Date": "2026-04-20T07:00:00+05:30",
                "Day": {
                    "HasPrecipitation": false,
                    "Icon": 5,
                    "IconPhrase": "Hazy sunshine"
                },
                "EpochDate": 1776648600,
                "Link": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?day=2&lang=en-us",
                "MobileLink": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?day=2&lang=en-us",
                "Night": {
                    "HasPrecipitation": false,
                    "Icon": 33,
                    "IconPhrase": "Clear"
                },
                "Sources": [
                    "AccuWeather"
                ],
                "Temperature": {
                    "Maximum": {
                        "Unit": "F",
                        "UnitType": 18,
                        "Value": 107
                    },
                    "Minimum": {
                        "Unit": "F",
                        "UnitType": 18,
                        "Value": 79
                    }
                }
            },
            {
                "Date": "2026-04-21T07:00:00+05:30",
                "Day": {
                    "HasPrecipitation": false,
                    "Icon": 5,
                    "IconPhrase": "Hazy sunshine"
                },
                "EpochDate": 1776735000,
                "Link": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?day=3&lang=en-us",
                "MobileLink": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?day=3&lang=en-us",
                "Night": {
                    "HasPrecipitation": false,
                    "Icon": 33,
                    "IconPhrase": "Clear"
                },
                "Sources": [
                    "AccuWeather"
                ],
                "Temperature": {
                    "Maximum": {
                        "Unit": "F",
                        "UnitType": 18,
                        "Value": 106
                    },
                    "Minimum": {
                        "Unit": "F",
                        "UnitType": 18,
                        "Value": 78
                    }
                }
            },
            {
                "Date": "2026-04-22T07:00:00+05:30",
                "Day": {
                    "HasPrecipitation": false,
                    "Icon": 1,
                    "IconPhrase": "Sunny"
                },
                "EpochDate": 1776821400,
                "Link": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?day=4&lang=en-us",
                "MobileLink": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?day=4&lang=en-us",
                "Night": {
                    "HasPrecipitation": false,
                    "Icon": 33,
                    "IconPhrase": "Clear"
                },
                "Sources": [
                    "AccuWeather"
                ],
                "Temperature": {
                    "Maximum": {
                        "Unit": "F",
                        "UnitType": 18,
                        "Value": 106
                    },
                    "Minimum": {
                        "Unit": "F",
                        "UnitType": 18,
                        "Value": 77
                    }
                }
            },
            {
                "Date": "2026-04-23T07:00:00+05:30",
                "Day": {
                    "HasPrecipitation": false,
                    "Icon": 1,
                    "IconPhrase": "Sunny"
                },
                "EpochDate": 1776907800,
                "Link": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?day=5&lang=en-us",
                "MobileLink": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?day=5&lang=en-us",
                "Night": {
                    "HasPrecipitation": false,
                    "Icon": 33,
                    "IconPhrase": "Clear"
                },
                "Sources": [
                    "AccuWeather"
                ],
                "Temperature": {
                    "Maximum": {
                        "Unit": "F",
                        "UnitType": 18,
                        "Value": 105
                    },
                    "Minimum": {
                        "Unit": "F",
                        "UnitType": 18,
                        "Value": 78
                    }
                }
            }
        ],
        "Headline": {
            "Category": "heat",
            "EffectiveDate": "2026-04-20T07:00:00+05:30",
            "EffectiveEpochDate": 1776648600,
            "EndDate": "2026-04-22T19:00:00+05:30",
            "EndEpochDate": 1776864600,
            "Link": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?lang=en-us",
            "MobileLink": "http://www.accuweather.com/en/in/ahmedabad/202438/daily-weather-forecast/202438?lang=en-us",
            "Severity": 7,
            "Text": "Very warm from Monday to Wednesday"
        }
    },
    "weatherURL": "https://dataservice.accuweather.com/forecasts/v1/daily/5day/202438"
}
```