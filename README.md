Demo: Nomad Python Driver Plugin
==========

Created for demo purposes for [HashiTalks 2020](https://events.hashicorp.com/hashitalks2020)
from the [skeleton project](https://www.nomadproject.io/docs/drivers/index.html) 
and the Nomad [Java driver](https://nomadproject.io/docs/drivers/java/).

Requirements
-------------------

- [Nomad](https://www.nomadproject.io/downloads.html) v0.9+
- [Go](https://golang.org/doc/install) v1.11 or later (to build the plugin)

Building the Skeleton Plugin
-------------------

Clone the repository somewhere in your computer. This project uses
[Go modules](https://blog.golang.org/using-go-modules) so you will need to set
the environment variable `GO111MODULE=on` or work outside your `GOPATH` if it
is set to `auto` or not declared.

```sh
$ git clone git@github.com:cgbaker/hashitalk-nomad-python-driver.git
```

Build the skeleton plugin.

```sh
$ make build
```

## Deploying Driver Plugins in Nomad

The initial version of the skeleton is a simple task that outputs a greeting.
You can try it out by starting a Nomad agent and running the job provided in
the `example` folder:

```sh
$ make build
$ nomad agent -dev -config=./example/agent.hcl -plugin-dir=$(pwd)

# in another shell
$ nomad run ./example/example.nomad
$ nomad logs <ALLOCATION ID>
```

Code Organization
-------------------
Follow the comments marked with a `TODO` tag to implement your driver's logic.
For more information check the
[Nomad documentation on plugins](https://www.nomadproject.io/docs/internals/plugins/index.html).
