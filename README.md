![Version](https://img.shields.io/badge/version-0.0.0-orange.svg)

# KVStore

Simple in-memory key/value store for training purpose. Project demonstrates
basic DDD approach (storage/service/http layer logic)

---

## Usage

To run server locally;

```bash
rake
```

or `cd` to project root;

```bash
go run cmd/server/main.go
```

Endpoints:

```http
GET    /healthz/live
GET    /healthz/ready

POST   /api/v1/set
GET    /api/v1/get?key={key}
PUT    /api/v1/update
DELETE /api/v1/delete?key={key}
GET    /api/v1/list
```

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
| `LOG_LEVEL` | Logging level | `INFO` |

Available tasks:

```bash
$ rake -T

rake default            # default task
rake docker:build       # Build image (locally)
rake docker:run         # Run image (locally)
rake lint               # run golangci lint
rake release[revision]  # release new version major,minor,patch, default: patch
rake run:server         # run server
rake test:run_all       # run all tests
```

Run tests via;

```bash
rake test:run_all
```

---

## Docker

```bash
# build
rake docker:build

# run
rake docker:run
```


