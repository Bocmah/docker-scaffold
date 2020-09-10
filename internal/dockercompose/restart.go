package dockercompose

import "fmt"

// RestartPolicy is one of the restart policies supported by docker
type RestartPolicy string

// All supported restart policies
const (
	RestartPolicyNo            RestartPolicy = "no"
	RestartPolicyAlways        RestartPolicy = "always"
	RestartPolicyOnFailure     RestartPolicy = "on-failure"
	RestartPolicyUnlessStopped RestartPolicy = "unless-stopped"
)

// Render formats RestartPolicy as YAML string
func (r RestartPolicy) Render() string {
	if r == "" {
		return ""
	}

	if r != RestartPolicyNo && r != RestartPolicyAlways && r != RestartPolicyOnFailure && r != RestartPolicyUnlessStopped {
		return ""
	}

	if r == RestartPolicyNo {
		return fmt.Sprintf(`restart: "%s"`, string(r))
	}

	return fmt.Sprintf("restart: %s", string(r))
}
