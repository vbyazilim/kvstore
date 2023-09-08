![Version](https://img.shields.io/badge/version-0.1.3-orange.svg)
[![Golang Tests](https://github.com/vbyazilim/kvstore/actions/workflows/go-test.yml/badge.svg)](https://github.com/vbyazilim/kvstore/actions/workflows/go-test.yml)
[![Golang CI Lint](https://github.com/vbyazilim/kvstore/actions/workflows/go-lint.yml/badge.svg)](https://github.com/vbyazilim/kvstore/actions/workflows/go-lint.yml)
[![codecov](https://codecov.io/gh/vbyazilim/kvstore/graph/badge.svg?token=514LHYMOA4)](https://codecov.io/gh/vbyazilim/kvstore)
![Powered by Rake](https://img.shields.io/badge/powered_by-rake-blue?logo=ruby)

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
GET    /healthz/live/
GET    /healthz/ready/

POST   /api/v1/set/
GET    /api/v1/get/?key={key}
PUT    /api/v1/update/
DELETE /api/v1/delete/?key={key}
GET    /api/v1/list/
```

Also, you can use [postman](postman/KVStore.postman_collection.json) collection.

---

## Development

### Requirements

- `go1.21.0`
- `bumpversion`
- `pre-commit`

You can create `.env` file inside of the project root for environment variables

Environment variables information:

| Variable Name | Description | Default Value |
|:--------------|:------------|:------------|
| `SERVER_ENV` | Server environment information for run-time | `local` |
| `LOG_LEVEL` | Logging level | `INFO` |

### Install `pre-commit`

https://pre-commit.com/

```bash
$ cd /path/to/kvstore
$ pre-commit install       # do only once!
```

Available tasks:

```bash
$ rake -T

rake default                        # default task
rake docker:build                   # Build image (locally)
rake docker:run                     # Run image (locally)
rake lint                           # run golangci lint
rake release[revision]              # release new version major,minor,patch, default: patch
rake run:server                     # run server
rake test:run_all                   # run all tests
rake test:run_all_display_coverage  # run all tests and display coverage
```

Run all tests via;

```bash
rake test:run_all
rake test:run_all_display_coverage  # macos only!
```

---

## Docker

```bash
# build
rake docker:build

# run
rake docker:run
```

---

