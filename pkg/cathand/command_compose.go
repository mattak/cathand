package cathand

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type EventInfo struct {
	Epoch    float64
	Touch    bool
	Source   string
	Position int
}

type EventData struct {
	Source string
	Time   float64
	Data   []byte
}

func ParseEventText(eventLogFile string) ([]EventInfo, error) {
	data, err := ioutil.ReadFile(eventLogFile)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")
	regex := regexp.MustCompile("^\\[\\s*([\\d\\.]+)\\s*\\]\\s*/dev/input/(\\w+):\\s+[\\w\\_]+\\s+([\\w\\_]+)\\s+(\\w+)")

	events := []EventInfo{}
	position := 0

	for i := 0; i < len(lines); i++ {
		if !regex.MatchString(lines[i]) {
			continue
		}

		matches := regex.FindStringSubmatch(lines[i])
		epoch, err := strconv.ParseFloat(matches[1], 64)
		if err != nil {
			return nil, err
		}

		eventDriverName := matches[2]
		eventTag := matches[3]
		eventValue := matches[4]
		touched := eventTag == "ABS_MT_TRACKING_ID" && eventValue != "ffffffff"

		events = append(events, EventInfo{Epoch: epoch, Source: eventDriverName, Touch: touched, Position: position})
		position++
	}

	return events, nil
}

func SplitEvents(filename string, info []EventInfo) ([]EventData, error) {
	data := []EventData{}

	allbytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	previousIndex := 0
	chunkSize := 24

	for i := 1; i < len(info); i++ {
		if info[i].Touch || info[i].Position-info[previousIndex].Position >= 32 {
			var timeDiff float64

			if info[i].Touch {
				timeDiff = info[i].Epoch - info[i-1].Epoch
			} else {
				timeDiff = info[i].Epoch - info[previousIndex].Epoch
			}

			bytes := allbytes[(previousIndex * chunkSize):(i * chunkSize)]
			data = append(data, EventData{Data: bytes, Time: timeDiff, Source: info[i].Source})
			previousIndex = i
		}
	}

	lastIndex := len(info) - 1
	bytes := allbytes[(previousIndex * chunkSize):((lastIndex + 1) * chunkSize)]
	data = append(data, EventData{Data: bytes, Time: info[lastIndex].Epoch - info[previousIndex].Epoch, Source: info[lastIndex].Source})

	return data, nil
}

func WriteInputData(project *Project, data []EventData) {
	for i := 0; i < len(data); i++ {
		filename := project.InputFile(i)
		ioutil.WriteFile(filename, data[i].Data, 0644)
	}
}

func WriteShell(project *Project, data []EventData) {
	shellHeader := `#!/bin/sh
cd $(dirname $0)

`
	fileHandle, err := os.Create(project.RunShellFile)
	AssertError(err)
	defer fileHandle.Close()

	fileHandle.Write([]byte(shellHeader))

	for i := 0; i < len(data); i++ {
		command := strings.Join([]string{
			fmt.Sprintf("cat %s > /dev/input/%s", project.InputFileWithoutRootDir(i), data[i].Source),
			fmt.Sprintf("echo sleep %f", data[i].Time),
			fmt.Sprintf("sleep %f", data[i].Time),
			"\n",
		}, "\n")
		fileHandle.Write([]byte(command))
	}
}

func CommandCompose(projectName string) {
	recordProject := NewProject(projectName, ".record", "")

	if !FileExists(recordProject.RootDir) {
		panic("Cannot find record directory: " + recordProject.RootDir)
	}
	if !FileExists(recordProject.EventFile) {
		panic("Cannot find event file : " + recordProject.EventFile)
	}

	infos, err := ParseEventText(recordProject.EventFile)
	AssertError(err)

	recordEventDataFile := recordProject.DeviceFile(infos[0].Source)
	data, err := SplitEvents(recordEventDataFile, infos)
	AssertError(err)

	playProject := NewProject(projectName, ".play", "")

	RemoveFile(playProject.RootDir)
	MakeDirectory(playProject.RootDir)
	MakeDirectory(playProject.InputDir)
	WriteInputData(&playProject, data)
	WriteShell(&playProject, data)
}
