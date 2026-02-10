# Aircraft System for Monitoring and Resource Management


This is the API for the aircraft system that is able to manage the resources/parts of an airplane, which sends an alert to the mechanic when a specific or multiple parts have an indication of near limit/or outdated format. 


## Need to Install

Make sure you have the following installed:

- **Go** (1.22+ recommended) (version used is go 1.24.4)
- **Git**
- **make**
- **air** (for live reload during development)
- **docker**

## Installations

Install `air` if you don’t have it yet:

```bash
go install github.com/cosmtrek/air@latest
```

## Available Make 
Make is a popular build automation tool that runs tasks defined in a file called a Makefile. It automates building executables by tracking file changes and only recompiling what is necessary. Makefiles serve as documentation and task runners for various projects, including Docker and software deployment. 

#### Start dev server with live reload
```bash
make dev
```
Runs the app using Air for hot reloading.

#### Set global GitHub username and email
```bash
make set-github name="Your Name" email="you@example.com"
```
Notice: This updates your global Git configuration, not just this repository.

#### Build the application
```bash
make build
```

Compiles the Go app and outputs the binary to: /bin/app

#### Commit and push to main
```bash
make push m="your commit message"
```

This will:
- Add all files 
- Commit with the provided message
- Push to the main branch


#### Create a new branch and push it
```bash
make branch b="branch-name" m="initial commit message"
```

This will:
- Create and switch to a new branch
- Add all files
- Commit with the provided message
- Push the branch and set the upstream


#### Ping the local server
```bash
make ping
```
Useful for quick health checks.

