# Transform Step Test
Testing the transform step implementation

## Create Workflow
```json
{
  "name": "Transform Test",
  "steps": [
    {
      "type": "transform",
      "ops": [
        { "op": "default", "path": "actor_name", "value": "Unknown" },
        { "op": "template", "to": "title", "template": "Event {{type}} by {{actor_name}}" },
        { "op": "pick", "paths": ["title"] }
      ]
    }
  ]
}
```

## Trigger Workflow
```json
{
    "id": "4bead0eb-a0d1-4335-a75f-597e69104584",
    "name": "Transform Test",
    "enabled": true,
    "steps": [
        {
            "type": "transform",
            "ops": [
                {
                    "op": "default",
                    "path": "actor_name",
                    "value": "Unknown"
                },
                {
                    "op": "template",
                    "to": "title",
                    "template": "Event {{type}} by {{actor_name}}"
                },
                {
                    "op": "pick",
                    "paths": [
                        "title"
                    ]
                }
            ]
        }
    ]
}
```

## Trigger Response
```bash
curl --location 'http://localhost:1324/t/4bead0eb-a0d1-4335-a75f-597e69104584' \
-H 'Content-Type: application/json' \
-d '{"type": "lock.unlock"}'
```

```json
{
    "title": "Event lock.unlock by Unknown"
}
```