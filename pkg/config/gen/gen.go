//go:build ignore

package main

import (
	"github.com/conductorone/baton-sdk/pkg/config"

	cfg "github.com/ConductorOne/baton-google-identity-platform/pkg/config"
)

func main() {
	config.Generate("google-identity-platform", cfg.Config)
}
