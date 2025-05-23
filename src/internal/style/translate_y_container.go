package style

import (
	"strings"
)

func TranslateYContainer(offsetTop int, containerHeight int, content string) string {
	contentLines := lines(content)

	sb := strings.Builder{}
	for i, l := range contentLines {
		if i < offsetTop || i >= offsetTop+containerHeight {
			continue
		}

		if i-offsetTop > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(l)
	}
	return sb.String()
}
