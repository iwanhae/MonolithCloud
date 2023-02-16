package main

import (
	"os"
	"time"

	"github.com/iwanhae/monolithcloud/cmd"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4 -config oapi-codegen.yaml api/oas3.yaml
func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Caller().Logger()
}

func main() {
	cmd.Execute()
}
