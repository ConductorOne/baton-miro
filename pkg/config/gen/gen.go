package main

import (
	cfg "github.com/conductorone/baton-miro/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("miro", cfg.Config)
}
