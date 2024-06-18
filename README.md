# Awasm

Build, deploy, and run your application on the edge. Written in Go.

## Demo

![demo image](./assets/demo.gif)

## Usage

```
$ awasm -h
Awasm is the application that you can build, deploy, and run your application on the edge.
To run awasm, use:
  - 'awasm run' to launch application.
  - 'awasm endpoints' to manage endpoints.

Usage:
  awasm [command]

Available Commands:
  deployments Used to manage deployments
  endpoints   Used to manage endpoints
  help        Help about any command
  keys        Used to manage api-keys
  login       Login into your Awasm account
  reset       Used to delete all Awasm related data on your machine
  run         Run the application
  signup      Signup into your Awasm account

Flags:
      --debug           Indicate whether the debug mode is turned on
      --domain string   Point the CLI to your own backend [can also set via environment var
iable name: AWASM_API_URL] (default "http://127.0.0.1:8080/api")
  -h, --help            help for awasm
  -v, --version         version for awasm
```
