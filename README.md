# `taskey`
[![Go](https://github.com/LassiHeikkila/taskey/actions/workflows/go.yml/badge.svg)](https://github.com/LassiHeikkila/taskey/actions/workflows/go.yml)
[![CodeQL](https://github.com/LassiHeikkila/taskey/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/LassiHeikkila/taskey/actions/workflows/codeql-analysis.yml)
[![DeepSource](https://deepsource.io/gh/LassiHeikkila/taskey.svg/?label=active+issues&show_trend=true&token=HV16nyHJUUw1Gw8R_CF3Ezq-)](https://deepsource.io/gh/LassiHeikkila/taskey/?ref=repository-badge)
[![DeepSource](https://deepsource.io/gh/LassiHeikkila/taskey.svg/?label=resolved+issues&show_trend=true&token=HV16nyHJUUw1Gw8R_CF3Ezq-)](https://deepsource.io/gh/LassiHeikkila/taskey/?ref=repository-badge)

# PWP SPRING 2022
# Project name: `taskey`
# Group information
* Student 1. Lassi Heikkila (Lassi.Heikkila@student.oulu.fi)

__Remember to include all required documentation and HOWTOs, including how to create and populate the database, how to run and test the API, the url to the entrypoint and instructions on how to setup and run the client__

# Running the server
If you have Go installed locally, you can simply run `go run ./cmd/taskey/` at the root of the repository to start the server.

You must have a running PostgreSQL server on your localhost, e.g. in a Docker container.

You may also build a Docker image using `docker image build -t taskey .` and run everything (requires `docker-compose`) by doing `docker-compose up -d`.

# Running tests
If you have Go installed locally, you can simply run `go test -v ./...` at the root of the repository to execute all tests.

# Deployment
Service is deployed with [Heroku](https://taskey-service.herokuapp.com).

You can view the API documentation at [/api/v1/](https://taskey-service.herokuapp.com/api/v1/).

# taskeyd
`taskeyd` is the daemon that will run on a machine intended for executing tasks. It is a small program that downloads the defined schedule and tasks from the server, executes them based on the schedule and uploads results back to the server.

Below is a small demo of how it looks in action:

[![asciicast](https://asciinema.org/a/MrTAIV70UIcXkbyHj9qJhI193.svg)](https://asciinema.org/a/MrTAIV70UIcXkbyHj9qJhI193)
