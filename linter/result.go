package linter

import "fmt"

type Result struct {
	Linter string
	Errors []Error
}

func (r Result) HasLintErrors() bool {
	return len(r.Errors) > 0
}

func (r Result) String() string {
	returnStr := r.Linter
	if len(r.Errors) == 0 {
		returnStr += fmt.Sprintf("\n\t%s\n", `No errors \o/`)
	} else {
		for _, e := range r.Errors {
			returnStr += e.String()
		}
	}
	return returnStr
}

type Error struct {
	Message       string
	Documentation string
}

func (e Error) String() string {
	return fmt.Sprintf("\n\t%s\n\t%s\n", e.Message, e.Documentation)
}
