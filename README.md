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
      "id": "9a810cced51b1b29dda47acb1b7af319",
      "category": ["Work", "Acme", "Project X"],
      "end": null,
      "start": "2014-08-06T15:23:00Z",
      "note": "Buy razor to shave the yak."
    },
    {
      "id": "d8a83fd8da3fda45523546f71a4c592e",
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
```
```json
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
    "id": "d8a83fd8da3fda45523546f71a4c592e",
    "category": ["Work", "Acme", "Project X"],
    "end": null,
    "start": "2014-08-06T15:23:00Z",
    "note": "Buy razor to shave the yak."
  }
}
```

### Update Time Entry

Request:

```
PUT /times/d8a83fd8da3fda45523546f71a4c592e
```

```json
{
  "id": "d8a83fd8da3fda45523546f71a4c592e",
  "category": ["Work", "Acme", "Project X"],
  "end": "2014-08-06T16:39:00Z",
  "start": "2014-08-06T15:23:00Z",
  "note": "Buy razor to shave the yak."
}
```

Response:

```json
{
  "time": {
    "id": "d8a83fd8da3fda45523546f71a4c592e",
    "category": ["Work", "Acme", "Project X"],
    "end": null,
    "start": "2014-08-06T15:23:00Z",
    "note": "Buy razor to shave the yak."
  }
}
```
