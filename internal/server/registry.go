package server

import (
	"github.com/datpp/go-kratos-based-template/internal/conf"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/hashicorp/consul/api"
)

// NewRegistrar creates a new consul registrar.
func NewRegistrar(c *conf.Registry, logger log.Logger) registry.Registrar {
	if c.GetEnabled() {
		logger = log.With(logger, "module", "server.registry")
		log.NewHelper(logger).Info("creating consul registrar")

		// Create consul client
		config := api.DefaultConfig()
		if c.Consul.Address != "" {
			config.Address = c.Consul.Address
		}
		if c.Consul.Scheme != "" {
			config.Scheme = c.Consul.Scheme
		}

		client, err := api.NewClient(config)
		if err != nil {
			log.NewHelper(logger).Fatalf("failed to create consul client: %v", err)
		}

		// Create consul registrar
		r := consul.New(client, consul.WithHealthCheck(c.Consul.HealthCheck))
		return r
	}

	return nil
}
