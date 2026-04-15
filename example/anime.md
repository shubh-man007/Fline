# AnimeChan API Doc

```bash
https://animechan.io/docs/quote/random-via-anime
BASE_URL: https://api.animechan.io/v1   {animechan}
curl -v https://api.animechan.io/v1/quotes/random
```

## Response
```json
{
  "status": "success",
  "data": {
    "content": "I told you before, Komamura. The only paths that I see with these eyes are the ones not dyed with blood. Those paths are the paths to justice. So whichever path I choose...Is justice.",
    "anime": {
      "id": 222,
      "name": "Bleach",
      "altName": "Bleach"
    },
    "character": {
      "id": 2143,
      "name": "Tousen Kaname"
    }
  }
}
```

## Create Workflow
```json
{
    "name":"Random Anime Quote",
    "steps":[
        {
            "type":"http_request",
            "method":"GET",
            "URL":"https://api.animechan.io/v1/quotes/random",
            "headers":{"Accept":"application/json"},
            "timeout":5000,
            "retries":3
        }
    ]
}
```

## Trigger Workflow
```json
{
    "id": "22e0ce3c-f316-4766-aacf-76d1922841a3",
    "name": "Random Anime Quote",
    "enabled": true,
    "steps": [
        {
            "type": "http_request",
            "method": "GET",
            "URL": "https://api.animechan.io/v1/quotes/random",
            "headers": {
                "Accept": "application/json"
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
        "data": {
            "anime": {
                "altName": "Vampire Knight",
                "id": 813,
                "name": "Vampire Knight"
            },
            "character": {
                "id": 2422,
                "name": "Yuki Kuran"
            },
            "content": "If it's something I can only cry about in my heart, it's almost like a sin."
        },
        "status": "success"
    }
}
```