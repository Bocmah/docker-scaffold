package service

import "fmt"

type SupportedSystem string

const (
	MySQL      SupportedSystem = "mysql"
	PostgreSQL SupportedSystem = "posgresql"
)

type systemDefaults struct {
	version string
	port    int
}

var defaults = map[SupportedSystem]systemDefaults{
	MySQL: {
		version: "8.0",
		port:    3306,
	},
	PostgreSQL: {
		version: "12.3",
		port:    5432,
	},
}

type Credentials struct {
	Username     string
	Password     string
	RootPassword string `yaml:"rootPassword"`
}

type Database struct {
	System      SupportedSystem
	Version     string
	Name        string
	Port        int
	Credentials `yaml:",inline"`
}

func (d *Database) FillDefaultsIfNotSet() {
	if d.System == "" {
		d.System = MySQL
	}

	if d.Version == "" {
		d.Version = defaults[d.System].version
	}

	if d.Port == 0 {
		d.Port = defaults[d.System].port
	}
}

func (d *Database) Validate() error {
	errors := &ValidationErrors{}

	if d.System != MySQL && d.System != PostgreSQL {
		errors.Add("Unsupported database system")
	}

	if d.Port == 0 {
		errors.Add("Database port is required")
	}

	if d.RootPassword == "" {
		errors.Add("Database root password is required")
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

func (d *Database) String() string {
	return fmt.Sprintf(
		"Database{System: %v, Version: %s, Name: %s, HttpPort: %d, Username: %s, Password: %s, RootPassword: %s}",
		d.System,
		d.Version,
		d.Name,
		d.Port,
		d.Username,
		d.Password,
		d.RootPassword,
	)
}
