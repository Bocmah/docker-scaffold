package service

// Environment is a collection of per-service environment variables
type Environment map[SupportedService]map[string]string
