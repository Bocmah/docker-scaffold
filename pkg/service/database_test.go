package service_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

type validationResult struct {
	wantErrs     []string
	actualErrs   error
	validatedVal interface{}
}

func failTestOnUnspottedError(result validationResult, t *testing.T) {
	for _, e := range result.wantErrs {
		if !strings.Contains(result.actualErrs.Error(), e) {
			t.Errorf("Failed to spot error %s in value %v", e, result.validatedVal)

		}
	}
}

func failTestOnErrorsOnCorrectInput(errs error, t *testing.T) {
	if errs != nil {
		t.Errorf("Following errors were returned despite correct inputs %v", errs)
	}
}

func TestDatabaseConfig_FillDefaultsIfNotSet(t *testing.T) {
	db := service.DatabaseConfig{}

	db.FillDefaultsIfNotSet()

	want := service.DatabaseConfig{
		System:  service.MySQL,
		Port:    3306,
		Version: "8.0",
	}

	if db != want {
		t.Errorf("Incorrect defaults, want %v, got %v", want, db)
	}
}

func TestDatabaseConfig_ValidateIncorrectInput(t *testing.T) {
	db := service.DatabaseConfig{System: "Unsupported"}

	errs := db.Validate()

	if errs != nil {
		res := validationResult{
			wantErrs: []string{
				"Unsupported database system",
				"DatabaseConfig port is required",
			},
			actualErrs:   errs,
			validatedVal: db,
		}

		failTestOnUnspottedError(res, t)
	} else {
		t.Errorf("Did not return any errors for value %v", db)
	}
}

func TestDatabaseConfig_ValidateCorrectInput(t *testing.T) {
	db := service.DatabaseConfig{
		System: service.MySQL,
		Port:   3306,
		Name:   "testdb",
		Credentials: service.Credentials{
			RootPassword: "rootpass",
		},
	}

	errs := db.Validate()

	failTestOnErrorsOnCorrectInput(errs, t)
}

func TestDatabaseConfig_Environment(t *testing.T) {
	tests := map[string]struct {
		conf *service.DatabaseConfig
		want map[string]string
	}{
		"mysql": {
			conf: &service.DatabaseConfig{
				System:  service.MySQL,
				Version: "8.0",
				Name:    "test-db",
				Port:    3306,
				Credentials: service.Credentials{
					Username:     "test-user",
					Password:     "test-password",
					RootPassword: "test-root-password",
				},
			},
			want: map[string]string{
				"MYSQL_USER":          "test-user",
				"MYSQL_DATABASE":      "test-db",
				"MYSQL_ROOT_PASSWORD": "test-root-password",
				"MYSQL_PASSWORD":      "test-password",
			},
		},
		"postgresql": {
			conf: &service.DatabaseConfig{
				System:  service.PostgreSQL,
				Version: "12",
				Name:    "test-db",
				Port:    5432,
				Credentials: service.Credentials{
					Username: "test-user",
					Password: "test-password",
				},
			},
			want: map[string]string{
				"POSTGRES_USER":     "test-user",
				"POSTGRES_DB":       "test-db",
				"POSTGRES_PASSWORD": "test-password",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := tc.conf.Environment()

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Fatalf("conf.Environment() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
