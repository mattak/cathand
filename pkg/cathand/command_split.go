package cathand

import "path"

func GenerateImagesFromVideo(project Project) {
	videoFileNames, err := ListUpFileNamesWithoutExtension(project.VideoDir)
	AssertError(err)

	MakeDirectory(project.ImageDir)

	for _, videoFileNameWithoutExtension := range videoFileNames {
		videoFilePath := path.Join(project.VideoDir, videoFileNameWithoutExtension+".mp4")
		imageFilePathFormat := project.ImageFileFormat(videoFileNameWithoutExtension)

		// ffmpeg -i sample.record/video/00.mp4 -r 1.0 -vf scale=320:-1 sample.record/image/00_%04d.jpg
		RunWait("ffmpeg", "-i", videoFilePath, "-r", "1.0", "-vf", "scale=256:-1", imageFilePathFormat)
	}
}

func CommandSplit(projects ...Project) {
	for _, project := range projects {
		GenerateImagesFromVideo(project)
	}
}