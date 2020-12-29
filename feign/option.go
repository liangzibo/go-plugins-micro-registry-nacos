package feign

import "github.com/micro/go-micro/v2/registry"

type Options struct {
	Registry registry.Registry
	Service  string
	Headers map[string]string
}

type Option func(*Options)
