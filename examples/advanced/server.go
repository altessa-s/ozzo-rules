package main

import (
	ozzo_rules "github.com/altessa-s/ozzo-rules"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// Network configuration validation
type ServerConfig struct {
	ListenAddr string
	BackendURL string
	PortRange  string
}

func (c ServerConfig) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ListenAddr,
			ozzo_rules.ListenAddress()),
		validation.Field(&c.BackendURL,
			ozzo_rules.URI()),
		validation.Field(&c.PortRange,
			ozzo_rules.PortWithRange(1, 65535)),
	)
}
