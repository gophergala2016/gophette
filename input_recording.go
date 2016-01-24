package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

var (
	recordingInput         = false
	replayingInput         = false
	recordedCharacterIndex = 1
)

type inputRecord struct {
	frame int
	event InputEvent
}

var (
	inputs []inputRecord
	frame  int
)

func init() {
	if replayingInput {
		inputs = recordedInputs
	}
}

func recordInput(event InputEvent) {
	if recordingInput && event.CharacterIndex == recordedCharacterIndex {
		inputs = append(inputs, inputRecord{frame: frame, event: event})
	}
}

func saveRecordedInputs() {
	input := bytes.NewBuffer(nil)
	input.WriteString(`package main

var recordedInputs = []inputRecord{
`)
	for i := range inputs {
		fmt.Fprintf(
			input,
			"\t{%v, InputEvent{%v, %v, %v}},\n",
			inputs[i].frame,
			inputs[i].event.Action,
			inputs[i].event.Pressed,
			inputs[i].event.CharacterIndex,
		)
	}
	input.WriteString(`}
`)
	ioutil.WriteFile("./recorded_inputs.go", input.Bytes(), 0777)
}
