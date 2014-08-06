# quantastic-api

## Time Entry Rules

* There can only be one active time entry, indicated by `{end: null}`. The
  active entry's `end` time is considered to be the current time when it comes
  to reporting or shadowing.
* Creating a new active time entry sets the `end` value of the previously
  active entry to the `start` value of the new entry.
* If an entry overlaps with one or more entries, the entry with the higher
  `end` value shadows the older entries. This means the older entries' `end`
  values may be adjusted, or they may be removed from a result set entirely.
  However, the original times remain stored internally and an API for showing
  all shadowing will be exposed in the future.

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
    "end": "2014-08-06T16:39:00Z",
    "start": "2014-08-06T15:23:00Z",
    "note": "Buy razor to shave the yak."
  }
}
```
