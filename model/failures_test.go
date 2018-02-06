package model

import "testing"

func TestAddMessage(t *testing.T) {
	f := new(Failures)
	f.Messages = append(f.Messages, "one")
	f.Messages = append(f.Messages, "two")

	expected := "one\ntwo\n"
	actual := f.Error()

	if actual != expected {
		t.Errorf("expected '%s' but got '%s'", expected, actual)
	}
}
