```json
{
  "name": "Reddit Post Summarizer",
  "steps": [
    {
      "type": "http_request",
      "method": "GET",
      "URL": "https://www.reddit.com/r/{{subreddit}}/top.json?limit=1",
      "headers": {
        "User-Agent": "Fline/1.0"
      },
      "timeout": 5000,
      "retries": 3
    },
    {
      "type": "transform",
      "ops": [
        {
          "op": "template",
          "to": "postTitle",
          "template": "{{result.data.children.0.data.title}}"
        },
        {
          "op": "template",
          "to": "prompt",
          "template": "Summarize why this Reddit post might be interesting in one sentence: {{postTitle}}"
        },
        {
          "op": "pick",
          "paths": ["prompt"]
        }
      ]
    },
    {
      "type": "http_request",
      "method": "POST",
      "URL": "https://api.openai.com/v1/responses",
      "headers": {
        "Content-Type": "application/json",
        "Authorization": "Bearer YOUR_OPENAI_KEY"
      },
      "body": {
        "model": "gpt-4o-mini",
        "input": "{{prompt}}"
      },
      "timeout": 15000,
      "retries": 3
    }
  ]
}
```

```json
{
    "prompt": "Summarize why this Reddit post might be interesting in one sentence: Go's implicit interface system is there a real solution to the discoverability problem or is it just accepted as a tradeoff",
    "result": {
        "background": false,
        "billing": {
            "payer": "developer"
        },
        "completed_at": 1776968619,
        "created_at": 1776968617,
        "error": null,
        "frequency_penalty": 0,
        "id": "resp_00eae70a1f07c5450069ea63a9aa948199a3310c594c9648c1",
        "incomplete_details": null,
        "instructions": null,
        "max_output_tokens": null,
        "max_tool_calls": null,
        "metadata": {},
        "model": "gpt-4o-mini-2024-07-18",
        "moderation": null,
        "object": "response",
        "output": [
            {
                "content": [
                    {
                        "annotations": [],
                        "logprobs": [],
                        "text": "This Reddit post explores whether Go's implicit interface system effectively addresses the challenges of discoverability in coding, or if programmers simply accept it as an unavoidable drawback of the language.",
                        "type": "output_text"
                    }
                ],
                "id": "msg_00eae70a1f07c5450069ea63ab716881998e3b0b310984c637",
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
            "input_tokens": 45,
            "input_tokens_details": {
                "cached_tokens": 0
            },
            "output_tokens": 35,
            "output_tokens_details": {
                "reasoning_tokens": 0
            },
            "total_tokens": 80
        },
        "user": null
    }
}
```