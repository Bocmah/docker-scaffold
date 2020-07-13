package service

type PHP struct {
	Version    string
	Extensions []string
}

func (p *PHP) FillDefaultsIfNotSet() {
	if p.Version == "" {
		p.Version = "7.4"
	}

	if p.Extensions == nil {
		p.Extensions = []string{"mbstring", "zip", "exif", "pcntl", "gd"}
	}
}

func (p *PHP) AddDatabaseExtension(db SupportedSystem) {
	switch db {
	case MySQL:
		p.Extensions = append(p.Extensions, "pdo_mysql")
	case PostgreSQL:
		p.Extensions = append(p.Extensions, "pdo_pgsql")
	}
}

func (p *PHP) Validate() *ValidationErrors {
	errors := &ValidationErrors{}

	if p.Version == "" {
		errors.Add("PHP version is required")
	}

	if errors.IsEmpty() {
		return nil
	}

	return errors
}
