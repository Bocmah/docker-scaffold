package dockercompose

import "strings"

// Image represents 'image' directive in docker-compose file
type Image struct {
	Name string
	Tag  string
}

// Render formats Image as YAML string
func (i *Image) Render() string {
	if i.Name == "" {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("image: ")

	if i.Tag == "" {
		sb.WriteString(i.Name)
		return sb.String()
	}

	sb.WriteString(mapping(i.Name, i.Tag))

	return sb.String()
}
