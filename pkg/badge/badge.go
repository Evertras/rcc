package badge

import (
	"bytes"
	"fmt"
	"text/template"
)

type coverageTemplateData struct {
	CoveragePercent string
}

var coverageTemplate = template.Must(
	template.New("coverage").
		Parse(`<svg xmlns="http://www.w3.org/2000/svg" width="99" height="20"><linearGradient id="a" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><rect rx="3" width="99" height="20" fill="#555"/><rect rx="3" x="63" width="36" height="20" fill="#97CA00"/><path fill="#97CA00" d="M63 0h4v20h-4z"/><rect rx="3" width="99" height="20" fill="url(#a)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="32.5" y="15" fill="#010101" fill-opacity=".3">coverage</text><text x="32.5" y="14">coverage</text><text x="80" y="15" fill="#010101" fill-opacity=".3">{{ .CoveragePercent }}</text><text x="80" y="14">{{ .CoveragePercent }}</text></g></svg>`),
)

func GenerateCoverageSVG(value1000 int) (string, error) {
	var buf bytes.Buffer

	err := coverageTemplate.Execute(&buf, coverageTemplateData{
		CoveragePercent: fmt.Sprintf("%d.%d%%", value1000/10, value1000%10),
	})

	if err != nil {
		return "", fmt.Errorf("template failed: %w", err)
	}

	return buf.String(), nil
}
