# DevOps Technical Check-In #1

This is a technical check-in for DevOps engineers. 

# Pre-requisites

You will need the following installed:

- Go to run the application (check with `go version`)
- Docker for image building/publishing (check with `docker version`)
- Docker Compose for environment provisioning (check with `docker-compose version`)
- Git for source control (check with `git -v`)
- Make for simple convenience scripts (check with `make -v`)

You will also need the following accounts:

- GitLab.com ([click here to register/login](https://gitlab.com/users/sign_in))

- - -

# Orientation & Instruction

This repository contains a simple program named `pinger` with a basic `Makefile` that allows you to do pull in dependencies, builds, test runs, and run tests using the following respectively:

- `make dep`
- `make build`
- `make run`
- `make test`

Directory structure:

| Directory | Description |
| --- | --- |
| `/bin` | Contains binaries |
| `/build` | Contains packaging/bundling related files |
| `/cmd` | Contains source code for CLI interfaces |
| `/deployments` | Contains manifests for deployments |
| `/vendor` | Contains dependencies |

- - -

# Get Started

1. Clone this repository from GitHub
2. Create a repository on GitLab
3. Set your local repository's remote to point to your GitLab
4. Push to your GitLab

For the submission, send us the link to your repository in GitLab.

- - -

# Containerisation

## Context

Not everyone has Go installed locally! Let's make it easier for developers to run this without installing anything.

## Task

Create a `Dockerfile` in the `./build` directory according to your best practices.

## Deliverable

Running `docker build -f ./build/Dockerfile -t devops/pinger:latest .` should result in a successful image named `devops/pinger:latest` which is reflected in the output of `docker image ls`.

Running `docker run -it -p 8000:8000 devops/pinger:latest` should result in the same behaviour as running `go run ./cmd/pinger`.

You can test if this works by running:

- `make docker_image`
- `make docker_testrun`

- - -

# Environment

## Context

Developers have been running this manually forever in an isolated setting, let's put a use case to it and demonstrate how it maybe used downstream the value chain!

## Task

Create a `docker-compose.yml` in the `./deployments` to demonstrate two `pinger` services that ping each other

## Deliverable

Running `docker-compose up -f ./deployments/docker-compose.yml` should result in a network of Docker containers that are pinging each other. Exposing the logs should reveal them pinging each other at their different ports.

You can test if this works by running:

- `make testenv`

- - -

# Documentation

## Context

Now that you've added some DevOps tooling to this project, it's time to document it!

## Task

Write a README.md in the `./docs` directory that contains instructions on how to operate this repository. The README should be as concise as possible while enabling anyone new to this project to get started as quickly as possible.

## Deliverable

README.md in the `./docs` directory.

- - -

# Pipeline (optional)

> This section is optional (but a bonus) to complete. If it is left undone, we will cover this during the face-to-face technical assessment so do read up a little on GitLab pipelines and how they work regardless!

## Context

Automation is key in DevOps to deliver value continuously and the first step we can take for this poor un-automated repository is to create a sensible pipeline that automates the build/test/release process.

## Task

Create a pipeline that results in:

1. The binary being built
2. Docker image being built
3. Docker image being published to DockerHub

The following should also be exposed as GitLab job artifacts:

1. The binary itself
2. Docker image in `.tar` format (see `Makefile` recipe `docker_tar`)

## Deliverbale

`.gitlab-ci.yml` in the root of this directory that results in a successful build on your own repository.

- - -

# License

Code is licensed under the [MIT license](./LICENSE-CODE).

Content is licensed under the [Creative Commons 4.0 (Attribution) license](./LICENSE-CONTENT).
