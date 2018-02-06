package model

import "bytes"

type Failures struct {
	Messages []string
}

func (f *Failures) Error() string {
	var buffer bytes.Buffer
	for _, m := range f.Messages {
		buffer.WriteString(m)
		buffer.WriteString("\n")
	}
	return buffer.String()
}

func (f *Failures) IsEmpty() bool {
	return len(f.Messages) == 0
}
