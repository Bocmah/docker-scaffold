package dockercompose

import "strings"

type Image struct {
	Name string
	Tag  string
}

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

	sb.WriteString(Mapping(i.Name, i.Tag))

	return sb.String()
}
