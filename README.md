# quantastic-api

## API

### List Time Entries

Request:

```
GET /times
```

Response:

```json
{
  "times": [
    {
      "category": ["Work", "Acme", "Project X"],
      "end": null,
      "start": "2014-08-06T15:23:00Z",
      "note": "Buy razor to shave the yak."
    },
    {
      "category": ["Sleep"],
      "end": "2014-08-06T15:23:00Z",
      "start": "2014-08-06T07:39:00Z",
      "note": "zzZZzz"
    }
  ]
}
```

### Create Time Entry

Request:

```
POST /times

{
  "category": ["Work", "Acme", "Project X"],
  "end": null,
  "start": "2014-08-06T15:23:00Z",
  "note": "Buy razor to shave the yak."
}
```

Response:

```json
{
  "time": {
    "category": ["Work", "Acme", "Project X"],
    "end": null,
    "start": "2014-08-06T15:23:00Z",
    "note": "Buy razor to shave the yak."
  }
}
```
