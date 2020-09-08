package service_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/Bocmah/phpdocker-gen/pkg/service"
)

func TestNginx_FillDefaultsIfNotSet(t *testing.T) {
	nginx := service.NginxConfig{}

	nginx.FillDefaultsIfNotSet()

	want := service.NginxConfig{
		HTTPPort:  80,
		HTTPSPort: 443,
		FastCGI: &service.FastCGI{
			PassPort:           9000,
			ReadTimeoutSeconds: 60,
		},
	}

	if diff := cmp.Diff(want, nginx); diff != "" {
		t.Fatalf("Incorrect defaults (-want +got):\n%s", diff)
	}
}

func TestNginx_ValidateIncorrectInput(t *testing.T) {
	nginx := service.NginxConfig{}

	errs := nginx.Validate()

	if errs != nil {
		res := validationResult{
			wantErrs: []string{
				"nginx port is required",
				"nginx FastCGI pass port is required",
				"nginx FastCGI read timeout is required",
			},
			actualErrs:   errs,
			validatedVal: nginx,
		}

		failTestOnUnspottedError(res, t)
	} else {
		t.Errorf("Did not return any errors for value %v", nginx)
	}
}

func TestNginx_ValidateCorrectInput(t *testing.T) {
	nginx := service.NginxConfig{
		HTTPPort: 80,
		FastCGI: &service.FastCGI{
			PassPort:           9000,
			ReadTimeoutSeconds: 60,
		},
		ServerName: "testserv",
	}

	errs := nginx.Validate()

	failTestOnErrorsOnCorrectInput(errs, t)
}
