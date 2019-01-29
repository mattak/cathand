package cathand

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func ReadImageFile(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, err := png.Decode(file)

	if err != nil {
		return nil, err
	}

	return image, nil
}

func ReadImageFiles(imageFiles []string) ([]image.Image, error) {
	images := []image.Image{}

	for _, filename := range imageFiles {
		image, err := ReadImageFile(filename)
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}

func uint32max(v1, v2, v3 uint32) (int, uint32) {
	if (v1 >= v2) {
		if (v1 >= v3) {
			return 1, v1
		} else {
			return 3, v3
		}
	} else {
		if (v2 >= v3) {
			return 2, v2
		} else {
			return 3, v3
		}
	}
}

func uint32min(v1, v2, v3 uint32) (int, uint32) {
	if (v1 <= v2) {
		if (v1 <= v3) {
			return 1, v1
		} else {
			return 3, v3
		}
	} else {
		if (v2 <= v3) {
			return 2, v2
		} else {
			return 3, v3
		}
	}
}

func hsl(color color.Color) (float64, float64, float64) {
	r, g, b, _ := color.RGBA()
	r = r / 256
	g = g / 256
	b = b / 256
	rf := float64(r)
	gf := float64(g)
	bf := float64(b)

	_, cmax := uint32max(r, g, b)
	indexmin, cmin := uint32min(r, g, b)
	sum := float64(cmax) + float64(cmin)
	delta := float64(cmax) - float64(cmin)

	h := 0.0
	if delta <= 0.0 {
		h = 0.0
	} else {
		switch indexmin {
		case 1: // red
			h = 60*float64(bf-gf)/delta + 180
			break;
		case 2: // green
			h = 60*float64(rf-bf)/delta + 300
			break;
		case 3: // blue
			h = 60*float64(gf-rf)/delta + 60
			break;
		}
	}

	l := sum / 2
	s := delta

	return h, s, l
}

func hslDistance(h1, s1, l1, h2, s2, l2 float64) float64 {
	hdiff := math.Abs(h1 - h2)
	if hdiff > 180.0 {
		hdiff = hdiff - 180.0
	}
	hdiff = hdiff / 180.0

	sdiff := math.Abs(s1-s2) / 255
	ldiff := math.Abs(l1-l2) / 255

	return hdiff*0.5 + sdiff*0.25 + ldiff*0.25
}

func CalcSimilarityOfImages(image1, image2 image.Image, colorL1Threshold float64) (float64, error) {
	x1_max := image1.Bounds().Max.X
	y1_max := image1.Bounds().Max.Y
	x2_max := image2.Bounds().Max.X
	y2_max := image2.Bounds().Max.Y

	if x1_max != x2_max || y1_max != y2_max {
		return 0, errors.New(fmt.Sprintf("Image size is not match: %dx%d <=> %dx%d", x1_max, y1_max, x2_max, y2_max))
	}

	difference_count := 0

	for y := 0; y < y1_max; y++ {
		for x := 0; x < x1_max; x++ {
			c1 := image1.At(x, y)
			c2 := image2.At(x, y)
			h1, s1, l1 := hsl(c1)
			h2, s2, l2 := hsl(c2)

			distance := hslDistance(h1, s1, l1, h2, s2, l2)
			if distance > colorL1Threshold {
				difference_count++
			}
		}
	}

	return 1.0 - float64(difference_count)/float64(x1_max*y1_max), nil
}

func SearchFirstMatch(targetImage image.Image, images []image.Image, startIndex int, colorThreshold, similarityThreshold float64) (int, float64) {
	matchIndex := -1
	matchSimilarity := 0.0

	for index := startIndex; index < len(images); index++ {
		image := images[index]
		similarity, err := CalcSimilarityOfImages(targetImage, image, colorThreshold)
		AssertError(err)

		if similarity >= similarityThreshold && similarity >= matchSimilarity {
			matchSimilarity = similarity
			matchIndex = index

			if matchSimilarity >= similarityThreshold {
				return matchIndex, matchSimilarity
			}
		}
	}

	return matchIndex, matchSimilarity
}

func SearchBestMatch(targetImage image.Image, images []image.Image, colorThreshold float64) (int, float64) {
	matchIndex := -1
	matchSimilarity := 0.0

	for index, image := range images {
		similarity, err := CalcSimilarityOfImages(targetImage, image, colorThreshold)
		AssertError(err)

		if similarity >= matchSimilarity {
			matchSimilarity = similarity
			matchIndex = index
		}
	}

	return matchIndex, matchSimilarity
}

func CommandVerify(project1 Project, project2 Project) {
	project1ImageFiles, err := ListUpFilePathes(project1.ImageDir)
	AssertError(err)
	project2ImageFiles, err := ListUpFilePathes(project2.ImageDir)
	AssertError(err)

	images1, err := ReadImageFiles(project1ImageFiles)
	AssertError(err)
	images2, err := ReadImageFiles(project2ImageFiles)
	AssertError(err)

	startIndex := 0
	sequencialFalseCount := 0
	sequencialFalseAdaptiveCount := 5

	for index, targetImage := range images1 {
		matchIndex, similarity := SearchFirstMatch(targetImage, images2, startIndex, 0.1, 0.7)

		if matchIndex >= 0 {
			fmt.Printf("%.3f\t%s\t->\t%s\n", similarity, project1ImageFiles[index], project2ImageFiles[matchIndex])
			sequencialFalseCount = 0
		} else {
			sequencialFalseCount++
			fmt.Printf("%.3f\t%s\t->\t%s\n", similarity, project1ImageFiles[index], "null")
		}

		if matchIndex > startIndex {
			startIndex = matchIndex
		}

		if sequencialFalseCount > sequencialFalseAdaptiveCount {
			panic(fmt.Sprintf("VerificationError: sequencial false count is more than adaptive count"))
		}
	}
}
