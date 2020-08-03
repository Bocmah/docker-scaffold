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

type DatabaseConfig struct {
	System      SupportedSystem
	Version     string
	Name        string
	Port        int
	Credentials `yaml:",inline"`
}

func (d *DatabaseConfig) FillDefaultsIfNotSet() {
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

func (d *DatabaseConfig) Validate() error {
	errors := &ValidationErrors{}

	if d.System != MySQL && d.System != PostgreSQL {
		errors.Add("Unsupported database system")
	}

	if d.Port == 0 {
		errors.Add("DatabaseConfig port is required")
	}

	if d.System == MySQL && d.RootPassword == "" {
		errors.Add("DatabaseConfig root password is required for MySQL")
	}

	if d.System == PostgreSQL && d.Password == "" {
		errors.Add("DatabaseConfig password is required for PostgreSQL")
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}

func (d *DatabaseConfig) String() string {
	return fmt.Sprintf(
		"DatabaseConfig{System: %v, Version: %s, Name: %s, HTTPPort: %d, Username: %s, Password: %s, RootPassword: %s}",
		d.System,
		d.Version,
		d.Name,
		d.Port,
		d.Username,
		d.Password,
		d.RootPassword,
	)
}

func (d *DatabaseConfig) Environment() map[string]string {
	switch d.System {
	case MySQL:
		return d.mySQLEnvironment()
	case PostgreSQL:
		return d.postgreSQLEnvironment()
	default:
		return map[string]string{}
	}
}

func (d *DatabaseConfig) mySQLEnvironment() map[string]string {
	env := map[string]string{}

	if d.RootPassword != "" {
		env["MYSQL_ROOT_PASSWORD"] = d.RootPassword
	}

	if d.Name != "" {
		env["MYSQL_DATABASE"] = d.Name
	}

	if d.Username != "" {
		env["MYSQL_USER"] = d.Username
	}

	if d.Password != "" {
		env["MYSQL_PASSWORD"] = d.Password
	}

	return env
}

func (d *DatabaseConfig) postgreSQLEnvironment() map[string]string {
	env := map[string]string{}

	if d.Name != "" {
		env["POSTGRES_DB"] = d.Name
	}

	if d.Username != "" {
		env["POSTGRES_USER"] = d.Username
	}

	if d.Password != "" {
		env["POSTGRES_PASSWORD"] = d.Password
	}

	return env
}
