package config

import (
	"time"
)

type DatabaseBackend string

const (
	MongoBackend  DatabaseBackend = "mongo"
	MemoryBackend DatabaseBackend = "memory"
)

type Database struct {
	Uri      string
	Username string
	Password string
	Name     string
	Timeout  time.Duration
}

type Config struct {
	Debug           bool
	SecretKey       string
	Database        Database
	DatabaseBackend DatabaseBackend
}
