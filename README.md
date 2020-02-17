HashiTalk 2020 Demo
==========

Nomad Python Driver and USB Device Driver
-------------------

Created for demo purposes for [HashiTalks 2020](https://events.hashicorp.com/hashitalks2020).

Python driver plugin was created using:
* [skeleton driver project](https://github.com/hashicorp/nomad-skeleton-driver-plugin) 
* Nomad [Java driver](https://nomadproject.io/docs/drivers/java/)

USB device plugin was created using:
* [skeleton device project](https://github.com/hashicorp/nomad-skeleton-device-plugin) 

Requirements
-------------------

- [Nomad](https://www.nomadproject.io/downloads.html) v0.9+
- [Go](https://golang.org/doc/install) v1.11 or later (to build the plugin)

Build the Plugins
-------------------

Clone the repository somewhere in your computer. This project uses
[Go modules](https://blog.golang.org/using-go-modules) so you will need to set
the environment variable `GO111MODULE=on` or work outside your `GOPATH` if it
is set to `auto` or not declared.

```sh
$ git clone https://github.com/cgbaker/hashitalk-2020-nomad-plugins
```

Build the plugins into the `plugins` directory:
```sh
$ make install
```

## Deploying Driver Plugins in Nomad

Once the plugins are installed to the `plugins` directory using the above steps, 
you can try them out by starting a Nomad agent using the included config and job spec.

```sh
$ nomad agent -dev -plugin-dir=$(PWD)/plugins -config=example/nomad.hcl

# in another shell
$ nomad run ./example/demo.hcl
```

