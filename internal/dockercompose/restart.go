package dockercompose

import "fmt"

type RestartPolicy string

const (
	RestartPolicyNo            RestartPolicy = "no"
	RestartPolicyAlways        RestartPolicy = "always"
	RestartPolicyOnFailure     RestartPolicy = "on-failure"
	RestartPolicyUnlessStopped RestartPolicy = "unless-stopped"
)

func (r RestartPolicy) String() string {
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
