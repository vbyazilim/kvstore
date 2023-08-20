![Version](https://img.shields.io/badge/version-0.0.0-orange.svg)

# KVStore

Simple in-memory key/value store for training purpose. Project demonstrates
basic DDD approach (storage/service/http layer logic)

---

## Installation

@wip

---

## Usage

- GET /healthz/live
- GET /healthz/live
- POST /api/v1/set
- GET /api/v1/get?key={key}
- PUT /api/v1/update
- DELETE /api/v1/delete?key={key}
- GET /api/v1/list

---

## Development

### Requirements

- `go1.21.0`
- `bumpversion`

You can create `.env` file inside of the project root for environment variables

Environment variables information:

| Variable Name | Description | Default Value |
|:--------------|:------------|:------------|
| `SERVER_ENV` | Server environment information for run-time | `local` |

Available tasks:

```bash
$ rake -T

rake default            # default task
rake docker:build       # Build image (locally)
rake docker:run         # Run image (locally)
rake lint               # run golangci lint
rake release[revision]  # release new version major,minor,patch, default: patch
rake run:server         # run server
```

---
