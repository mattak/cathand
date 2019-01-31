package cathand

import (
	"html/template"
	"os"
	"path"
)

type ReportMatchImage struct {
	Path1      string
	Path2      string
	Similarity float64
	Color      string
	Number     int
	TotalCount int
}

type Report struct {
	Matches   []ReportMatchImage
	Option    VerifyOption
	CreatedAt string
}

func GenerateHtml(templatePath string, writePath string, report Report) {
	MakeDirectory(path.Dir(writePath))
	file, err := os.Create(writePath)
	AssertError(err)
	defer file.Close()

	tmpl := template.Must(template.ParseFiles(templatePath))
	err = tmpl.Execute(file, report)
	AssertError(err)
}

func (r *ReportMatchImage) SetSuccess(success bool) {
	if success {
		r.Color = "#B2FF59"
	} else {
		r.Color = "#FF8F00"
	}
}
