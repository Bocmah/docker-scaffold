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
	tests := map[string]struct {
		conf     *service.NginxConfig
		wantErrs []string
	}{
		"empty config": {
			conf: &service.NginxConfig{},
			wantErrs: []string{
				"nginx port is required",
				"nginx FastCGI pass port is required",
				"nginx FastCGI read timeout is required",
			},
		},
		"without FastCGI pass port": {
			conf: &service.NginxConfig{
				HTTPPort:   80,
				HTTPSPort:  443,
				ServerName: "test-server",
				FastCGI:    &service.FastCGI{ReadTimeoutSeconds: 60},
			},
			wantErrs: []string{
				"nginx FastCGI pass port is required",
			},
		},
		"without FastCGI read timeout": {
			conf: &service.NginxConfig{
				HTTPPort:   80,
				HTTPSPort:  443,
				ServerName: "test-server",
				FastCGI:    &service.FastCGI{PassPort: 9000},
			},
			wantErrs: []string{
				"nginx FastCGI read timeout is required",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			errs := tc.conf.Validate()

			if errs != nil {
				res := validationResult{
					wantErrs:     tc.wantErrs,
					actualErrs:   errs,
					validatedVal: tc.conf,
				}

				failTestOnUnspottedError(res, t)
			} else {
				t.Errorf("Did not return any errors for value %v", tc.conf)
			}
		})
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
