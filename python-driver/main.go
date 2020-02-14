package main

import (
	log "github.com/hashicorp/go-hclog"

	"github.com/cgbaker/hashitalk-nomad-python-driver/python"
	"github.com/hashicorp/nomad/plugins"
)

func main() {
	// Serve the plugin
	plugins.Serve(factory)
}

// factory returns a new instance of a nomad driver plugin
func factory(log log.Logger) interface{} {
	return python.NewPlugin(log)
}
