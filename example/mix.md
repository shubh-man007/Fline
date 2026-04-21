# Multi Service Workflow

## Create Workflow
```json
{
  "name": "Weather Summary",
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
    },
    {
      "type": "transform",
      "ops": [
        {
          "op": "template",
          "to": "openAIPrompt",
          "template": "Summarize this 5-day weather forecast for Ahmedabad in 3 concise sentences, mention the headline and temperature range: Headline: {{result.Headline.Text}}. Day 1: {{result.DailyForecasts.0.Day.IconPhrase}}, High {{result.DailyForecasts.0.Temperature.Maximum.Value}}F Low {{result.DailyForecasts.0.Temperature.Minimum.Value}}F. Day 2: {{result.DailyForecasts.1.Day.IconPhrase}}, High {{result.DailyForecasts.1.Temperature.Maximum.Value}}F Low {{result.DailyForecasts.1.Temperature.Minimum.Value}}F. Day 3: {{result.DailyForecasts.2.Day.IconPhrase}}, High {{result.DailyForecasts.2.Temperature.Maximum.Value}}F Low {{result.DailyForecasts.2.Temperature.Minimum.Value}}F. Day 4: {{result.DailyForecasts.3.Day.IconPhrase}}, High {{result.DailyForecasts.3.Temperature.Maximum.Value}}F Low {{result.DailyForecasts.3.Temperature.Minimum.Value}}F. Day 5: {{result.DailyForecasts.4.Day.IconPhrase}}, High {{result.DailyForecasts.4.Temperature.Maximum.Value}}F Low {{result.DailyForecasts.4.Temperature.Minimum.Value}}F."
        }
      ]
    },
    {
      "type": "http_request",
      "method": "POST",
      "URL": "https://api.openai.com/v1/responses",
      "headers": {
        "Content-Type": "application/json",
        "Authorization": "Bearer $OPENAI_API_KEY"
      },
      "body": {
        "model": "gpt-4o-mini",
        "input": "{{openAIPrompt}}"
      },
      "timeout": 15000,
      "retries": 3
    }
  ]
}
```

## Trigger WorkFlow
```json
{
    "id": "af139ef2-ab01-491e-a636-ad7f44e2b87d",
    "name": "Weather Summary",
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
        },
        {
            "type": "transform",
            "ops": [
                {
                    "op": "template",
                    "to": "openAIPrompt",
                    "template": "Summarize this 5-day weather forecast for Ahmedabad in 3 concise sentences, mention the headline and temperature range: Headline: {{result.Headline.Text}}. Day 1: {{result.DailyForecasts.0.Day.IconPhrase}}, High {{result.DailyForecasts.0.Temperature.Maximum.Value}}F Low {{result.DailyForecasts.0.Temperature.Minimum.Value}}F. Day 2: {{result.DailyForecasts.1.Day.IconPhrase}}, High {{result.DailyForecasts.1.Temperature.Maximum.Value}}F Low {{result.DailyForecasts.1.Temperature.Minimum.Value}}F. Day 3: {{result.DailyForecasts.2.Day.IconPhrase}}, High {{result.DailyForecasts.2.Temperature.Maximum.Value}}F Low {{result.DailyForecasts.2.Temperature.Minimum.Value}}F. Day 4: {{result.DailyForecasts.3.Day.IconPhrase}}, High {{result.DailyForecasts.3.Temperature.Maximum.Value}}F Low {{result.DailyForecasts.3.Temperature.Minimum.Value}}F. Day 5: {{result.DailyForecasts.4.Day.IconPhrase}}, High {{result.DailyForecasts.4.Temperature.Maximum.Value}}F Low {{result.DailyForecasts.4.Temperature.Minimum.Value}}F."
                }
            ]
        },
        {
            "type": "http_request",
            "method": "POST",
            "URL": "https://api.openai.com/v1/responses",
            "headers": {
                "Content-Type": "application/json",
                "Authorization": "Bearer $OPENAI_API_KEY"
            },
            "body": {
                "model": "gpt-4o-mini",
                "input": "{{openAIPrompt}}"
            },
            "timeout": 15000,
            "retries": 3
        }
    ]
}
```

## Trigger Response
```json
{
    "openAIPrompt": "Summarize this 5-day weather forecast for Ahmedabad in 3 concise sentences, mention the headline and temperature range: Headline: Very warm from Monday to Tuesday. Day 1: Intermittent clouds, High 107F Low 77F. Day 2: Hazy sunshine, High 106F Low 79F. Day 3: Hazy sunshine, High 107F Low 79F. Day 4: Sunny, High 105F Low 77F. Day 5: Sunny, High 105F Low 79F.",
    "result": {
        "background": false,
        "billing": {
            "payer": "developer"
        },
        "completed_at": 1776626480,
        "created_at": 1776626478,
        "error": null,
        "frequency_penalty": 0,
        "id": "resp_097f5a9485c3102e0069e52b2e1eac8193a8c600680f054276",
        "incomplete_details": null,
        "instructions": null,
        "max_output_tokens": null,
        "max_tool_calls": null,
        "metadata": {},
        "model": "gpt-4o-mini-2024-07-18",
        "object": "response",
        "output": [
            {
                "content": [
                    {
                        "annotations": [],
                        "logprobs": [],
                        "text": "**Headline: Very warm from Monday to Tuesday.** The weather forecast for Ahmedabad shows high temperatures ranging from 105°F to 107°F throughout the week, with lows between 77°F and 79°F. Conditions will include intermittent clouds and hazy sunshine on most days, culminating in sunny weather by the end of the week.",
                        "type": "output_text"
                    }
                ],
                "id": "msg_097f5a9485c3102e0069e52b2f9b5c81939cc389cc356eba78",
                "role": "assistant",
                "status": "completed",
                "type": "message"
            }
        ],
        "parallel_tool_calls": true,
        "presence_penalty": 0,
        "previous_response_id": null,
        "prompt_cache_key": null,
        "prompt_cache_retention": "in_memory",
        "reasoning": {
            "effort": null,
            "summary": null
        },
        "safety_identifier": null,
        "service_tier": "default",
        "status": "completed",
        "store": true,
        "temperature": 1,
        "text": {
            "format": {
                "type": "text"
            },
            "verbosity": "medium"
        },
        "tool_choice": "auto",
        "tools": [],
        "top_logprobs": 0,
        "top_p": 1,
        "truncation": "disabled",
        "usage": {
            "input_tokens": 123,
            "input_tokens_details": {
                "cached_tokens": 0
            },
            "output_tokens": 68,
            "output_tokens_details": {
                "reasoning_tokens": 0
            },
            "total_tokens": 191
        },
        "user": null
    },
    "weatherURL": "https://dataservice.accuweather.com/forecasts/v1/daily/5day/202438"
}
```