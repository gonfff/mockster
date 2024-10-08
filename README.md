
<p align="center">
    <img src="mockster.svg" alt="Mockster" width="200">
</p>
<p align="center">
  <a href="https://github.com/gonfff/mockster/actions/workflows/ci.yaml">
    <img src="https://github.com/gonfff/mockster/actions/workflows/ci.yaml/badge.svg" alt="Build Status">
  </a>
  <a href="https://goreportcard.com/report/github.com/gonfff/mockster">
    <img src="https://goreportcard.com/badge/github.com/gonfff/mockster" alt="Go Report Card">
  </a>
  <a href="https://coveralls.io/github/gonfff/mockster?branch=master">
    <img src="https://coveralls.io/repos/github/gonfff/mockster/badge.svg?branch=master" alt="Coveralls">
  </a>
  <a href="https://codecov.io/gh/gonfff/mockster">
    <img src="https://codecov.io/gh/gonfff/mockster/branch/master/graph/badge.svg" alt="codecov">
  </a>
</p>




Mockster is a minimalistic, lightweight and simple mock server for testing your http clients. It can be embedded in your CI pipeline for tests or used as a standalone server while developing your application.

---
## Run app

```bash
docker run -p 8080:8080 --rm ghcr.io/gonfff/mockster:<version>
```

## Create mocks
* by YAML-file with environment variable MOCK_FILE_PATH can be provided ([example](examples/mocks.yaml))
* by sending a POST request to the `/management/mocks` endpoint with the following body ([example](examples/curl.txt))
* by importing YAML-file with UI
* by creating mock with UI

## Request mocks

By sending a request to the `/mock/<your_path>` endpoint and providing all nessessary headers, body, cookies, query params.

## UI
Mockster has a simple UI for creating and managing mocks. It is available on the `/` endpoint.

![UI](examples/ui1.png)
![UI2](examples/ui2.png)

## Environment variables
- `MOCK_FILE_PATH` - path to YAML-file with mocks
- `LOG_FORMATTER` - log formatter (text or json)
- `PORT` - port on which the server will be running
- `STATIC_PATH` - path to static files (UI)

---
#### Additional info
Thats my first Go, Bootstrap, JS project, so I will be glad to hear any feedback and suggestions.
Feel free to open issues and PRs.
