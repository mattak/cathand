package cathand

import (
	"fmt"
	"os"
	"path"
	"time"
)

func CommandVerify(project1, project2, reportProject Project, verifyOption VerifyOption) {
	project1ImageFiles, err := ListUpFilePathes(project1.ImageDir)
	AssertError(err)
	project2ImageFiles, err := ListUpFilePathes(project2.ImageDir)
	AssertError(err)

	images1, err := ReadImageFiles(project1ImageFiles)
	AssertError(err)
	images2, err := ReadImageFiles(project2ImageFiles)
	AssertError(err)

	startIndex := 0
	sequentialFalseFrameCount := 0
	failure := false

	report := Report{}
	report.Matches = []ReportMatchImage{}
	//report.Matches = make([]ReportMatchImage, len(project1ImageFiles))

	for index, targetImage := range images1 {
		matchIndex, similarity := SearchFirstMatch(
			targetImage, images2, startIndex, verifyOption.ColorThreshold, verifyOption.SimilarityThreshold)

		match := ReportMatchImage{}
		match.Number = index + 1
		match.TotalCount = len(project1ImageFiles)
		match.Similarity = similarity
		match.Path1 = reportProject.RelativeFile("image1", path.Base(project1ImageFiles[index]))
		match.SetSuccess(matchIndex >= 0)

		if matchIndex >= 0 {
			fmt.Printf("%.3f\t%s\t->\t%s\n", similarity, project1ImageFiles[index], project2ImageFiles[matchIndex])
			match.Path2 = reportProject.RelativeFile("image2", path.Base(project2ImageFiles[matchIndex]))
			sequentialFalseFrameCount = 0
		} else {
			sequentialFalseFrameCount++
			fmt.Printf("%.3f\t%s\t->\t%s\n", similarity, project1ImageFiles[index], "null")
		}

		report.Matches = append(report.Matches, match)

		if matchIndex > startIndex {
			startIndex = matchIndex
		}

		if sequentialFalseFrameCount > verifyOption.SequentialFalseFrameThreshold {
			failure = true
			break
		}
	}

	RemoveFile(reportProject.RootDir)
	MakeDirectory(reportProject.RootDir)
	RunWait("rsync", "-av", project1.ImageDir+"/", path.Join(reportProject.RootDir, "image1/"))
	RunWait("rsync", "-av", project2.ImageDir+"/", path.Join(reportProject.RootDir, "image2/"))

	report.CreatedAt = time.Now().Format(time.RFC3339)
	GenerateHtml(
		os.Getenv("GOPATH")+"/src/github.com/mattak/cathand/template/report.html",
		reportProject.ReportHtmlFile(),
		report)

	fmt.Println("report: ", reportProject.ReportHtmlFile())

	if failure {
		panic(fmt.Sprintf("VerificationError: sequencial false count is more than adaptive count (%d)", verifyOption.SequentialFalseFrameThreshold))
	}
}
