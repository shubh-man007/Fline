# OpenAI API Doc

```bash
curl https://api.openai.com/v1/responses \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $OPENAI_API_KEY" \
  -d '{
    "model": "gpt-5.4",
    "input": "Tell me a three sentence bedtime story about a unicorn."
  }'
```

## Create Workflow
```json
{
    "name":"SLOPenAI",
    "steps":[
        {
            "type":"http_request",
            "method":"post",
            "URL":"https://api.openai.com/v1/responses",
            "headers":{"Accept":"application/json", "Authorization": "Bearer $OPENAI_API_KEY"},
            "body":{"model": "gpt-5.4", "input": "Tell me a three sentence bedtime story about a unicorn."},
            "timeout":5000,
            "retries":3
        }
    ]
}
```

## Trigger Workflow
```json
{
    "id": "22831e01-247a-4902-9c6e-9f9553a8659f",
    "name": "SLOPenAI",
    "enabled": true,
    "steps": [
        {
            "type": "http_request",
            "method": "post",
            "URL": "https://api.openai.com/v1/responses",
            "headers": {
                "Accept": "application/json",
                "Authorization": "Bearer $OPENAI_API_KEY"
            },
            "body": {
                "model": "gpt-5.4",
                "input": "Tell me a three sentence bedtime story about a unicorn."
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
        "background": false,
        "billing": {
            "payer": "developer"
        },
        "completed_at": 1776277862,
        "created_at": 1776277859,
        "error": null,
        "frequency_penalty": 0,
        "id": "resp_09cc0f38a8fd52410069dfd9636b848193ac0a756b191ac9f4",
        "incomplete_details": null,
        "instructions": null,
        "max_output_tokens": null,
        "max_tool_calls": null,
        "metadata": {},
        "model": "gpt-5.4-2026-03-05",
        "object": "response",
        "output": [
            {
                "content": [
                    {
                        "annotations": [],
                        "logprobs": [],
                        "text": "A gentle unicorn with a silver mane tiptoed through a moonlit meadow, leaving tiny sparkles in the grass wherever her hooves touched. She curled up beside a sleepy brook and used her glowing horn to paint soft stars in the sky for all the forest animals to dream under. Soon the whole meadow was quiet and cozy, and the unicorn closed her eyes, smiling as the night hummed a lullaby.",
                        "type": "output_text"
                    }
                ],
                "id": "msg_09cc0f38a8fd52410069dfd963ed94819399d0383e54c068b0",
                "phase": "final_answer",
                "role": "assistant",
                "status": "completed",
                "type": "message"
            }
        ],
        "parallel_tool_calls": true,
        "presence_penalty": 0,
        "previous_response_id": null,
        "prompt_cache_key": null,
        "prompt_cache_retention": null,
        "reasoning": {
            "effort": "none",
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
        "top_p": 0.98,
        "truncation": "disabled",
        "usage": {
            "input_tokens": 17,
            "input_tokens_details": {
                "cached_tokens": 0
            },
            "output_tokens": 87,
            "output_tokens_details": {
                "reasoning_tokens": 0
            },
            "total_tokens": 104
        },
        "user": null
    }
}
```