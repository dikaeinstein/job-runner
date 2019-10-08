# job-runner

An asynchronous job runner written in golang with redis and nats.io

[![Build Status](https://travis-ci.com/dikaeinstein/job-runner.svg?branch=master)](https://travis-ci.com/dikaeinstein/job-runner)
[![Coverage Status](https://coveralls.io/repos/github/dikaeinstein/job-runner/badge.svg?branch=master)](https://coveralls.io/github/dikaeinstein/job-runner?branch=master)

## Run Locally (With docker-compose)

Prerequisites:

- make
- Go 1.10+
- docker

To start the REST api and the sample `new.product` worker

Run:

```sh
make docker
```

To test you can use apache bench to send sample requests to the REST api `/jobs` endpoint:

```sh
ab -c 50 -n 10000 -p job.json -r -T application/json http://localhost:8901/jobs
```

You can tweak the `ab` flags as you wish
