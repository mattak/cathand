package cathand

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Event struct {
	EpochSec  int64
	EpochUsec int64

	Type  uint16
	Code  uint16
	Value uint32
}

func ParseEventFromFile(eventLogFile string) (map[string][]Event, error) {
	data, err := ioutil.ReadFile(eventLogFile)

	if err != nil {
		return nil, err
	}

	return ParseEvent(string(data))
}

func ParseEvent(data string) (map[string][]Event, error) {
	lines := strings.Split(data, "\n")
	regex := regexp.MustCompile("^\\[\\s*(\\d+)\\.(\\d+)\\s*\\]\\s*/dev/input/(\\w+):\\s+([\\w\\_]+)\\s+([\\w\\_]+)\\s+(\\w+)")

	eventsMap := map[string][]Event{}

	for i := 0; i < len(lines); i++ {
		if !regex.MatchString(lines[i]) {
			continue
		}

		matches := regex.FindStringSubmatch(lines[i])

		epochSec, err := strconv.ParseInt(matches[1], 10, 32)
		if err != nil {
			return nil, err
		}

		epochUsec, err := strconv.ParseInt(matches[2], 10, 32)
		if err != nil {
			return nil, err
		}

		eventDriverName := matches[3]

		eventType, err := strconv.ParseUint(matches[4], 16, 16)
		if err != nil {
			return nil, err
		}

		eventCode, err := strconv.ParseUint(matches[5], 16, 16)
		if err != nil {
			return nil, err
		}

		eventValue, err := strconv.ParseUint(matches[6], 16, 32)
		if err != nil {
			return nil, err
		}

		events, exists := eventsMap[eventDriverName]
		if !exists {
			events = []Event{}
		}
		events = append(events, Event{
			EpochSec:  epochSec,
			EpochUsec: epochUsec,
			Type:      uint16(eventType),
			Code:      uint16(eventCode),
			Value:     uint32(eventValue),
		})
		eventsMap[eventDriverName] = events
	}

	return eventsMap, nil
}

func WriteEvent(project *Project, eventsMap map[string][]Event) {
	for eventDriverName, events := range eventsMap {
		var data bytes.Buffer
		for _, event := range events {
			data.WriteString(fmt.Sprintf("%d.%06d %04x %04x %08x\n",
				event.EpochSec,
				event.EpochUsec,
				event.Type,
				event.Code,
				event.Value))
		}

		filename := project.InputFile(eventDriverName)

		err := ioutil.WriteFile(filename, data.Bytes(), 0644)
		AssertError(err)
	}
}

func WriteShell(project *Project, eventsMap map[string][]Event) {
	var buffer bytes.Buffer

	buffer.WriteString(`#!/bin/sh
cd $(dirname $0)
CPU_ABI=$(getprop ro.product.cpu.abi)
`)

	for eventDriverName, _ := range eventsMap {
		buffer.WriteString(fmt.Sprintf(
			"./bin/${CPU_ABI}/cathand /dev/input/%s %s\n",
			eventDriverName,
			project.InputFileWithoutRootDir(eventDriverName)))
	}

	err := ioutil.WriteFile(project.RunShellFile, buffer.Bytes(), 0755)
	AssertError(err)
}

func CopyExecutable(project *Project) {
	cmd := exec.Command("rsync", "-av", "android_bin/obj/", project.BinDir+"/")
	AssertError(cmd.Run())
}

func CommandCompose(projectName string) {
	recordProject := NewProject(projectName, ".record", "")

	if !FileExists(recordProject.RootDir) {
		panic("Cannot find record directory: " + recordProject.RootDir)
	}
	if !FileExists(recordProject.EventFile) {
		panic("Cannot find event file: " + recordProject.EventFile)
	}

	eventsMap, err := ParseEventFromFile(recordProject.EventFile)
	AssertError(err)

	playProject := NewProject(projectName, ".play", "")
	RemoveFile(playProject.RootDir)
	MakeDirectory(playProject.RootDir)
	MakeDirectory(playProject.InputDir)
	WriteEvent(&playProject, eventsMap)
	WriteShell(&playProject, eventsMap)
	CopyExecutable(&playProject)
}
