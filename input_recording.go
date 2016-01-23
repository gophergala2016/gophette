package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
)

const (
	recordingInput = false
	replayingInput = false
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

func saveRecordedInputs() {
	input := bytes.NewBuffer(nil)
	input.WriteString(`package main

var recordedInputs = []inputRecord{
`)
	for i := range inputs {
		fmt.Fprintf(
			input,
			"\t{%v, InputEvent{%v, %v}},\n",
			inputs[i].frame,
			inputs[i].event.Action,
			inputs[i].event.Pressed,
		)
	}
	input.WriteString(`}
`)
	ioutil.WriteFile("./recorded_inputs.go", input.Bytes(), 0777)
}
