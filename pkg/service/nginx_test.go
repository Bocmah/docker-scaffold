package service_test

import (
	"testing"

	"github.com/Bocmah/phpdocker-scaffold/pkg/service"
)

func TestNginx_FillDefaultsIfNotSet(t *testing.T) {
	nginx := service.Nginx{}

	nginx.FillDefaultsIfNotSet()

	want := service.Nginx{
		HttpPort:  80,
		HttpsPort: 443,
		FastCGI: service.FastCGI{
			PassPort:           9000,
			ReadTimeoutSeconds: 60,
		},
	}

	if nginx != want {
		t.Errorf("Incorrect defaults, want %v, got %v", want, nginx)
	}
}

func TestNginx_ValidateIncorrectInput(t *testing.T) {
	nginx := service.Nginx{}

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
	nginx := service.Nginx{
		HttpPort: 80,
		FastCGI: service.FastCGI{
			PassPort:           9000,
			ReadTimeoutSeconds: 60,
		},
		ServerName: "testserv",
	}

	errs := nginx.Validate()

	failTestOnErrorsOnCorrectInput(errs, t)
}
