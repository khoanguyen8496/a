package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/as/shiny/event/key"
	"github.com/as/ui/tag"
)

type ErrGoImports struct {
	Path  string
	GoExt bool
	Err   error
}

func (e ErrGoImports) Error() string {
	m := e.Err.Error()
	if m == "" && !e.GoExt {
		m = "not a go file"
	}
	n := strings.Index(m, "<standard input>")
	if n != -1 {
		m = strings.Replace(m, "<standard input>", e.Path, -1)
	} else {
		m = e.Path + ":" + m
	}
	return fmt.Sprintf("goimports: %s", m)
}

func runGoImports(t *tag.Tag, e key.Event) {
	ee := &ErrGoImports{
		Path:  t.FileName(),
		GoExt: !strings.HasSuffix(t.FileName(), ".go"),
	}
	cmd := exec.Command("goimports")
	cmd.Stdin = bytes.NewReader(t.Bytes())
	b := new(bytes.Buffer)
	berr := new(bytes.Buffer)
	cmd.Stdout = b
	cmd.Stderr = berr
	err := cmd.Run()
	if err != nil || b.Len() < len("package") {
		if err == nil {
			ee.Err = fmt.Errorf("file too short")
		} else {
			if berr.Len() != 0 {
				ee.Err = fmt.Errorf(berr.String())
			} else {
				ee.Err = err
			}
		}
		if ee != nil {
			logf("imports: %s", ee)
		}
		return
	}
	q0, q1 := t.Dot()
	t.Delete(0, t.Len())
	io.Copy(t, b)
	t.Select(q0, q1) // TODO(as): BUG. shouldn't modify selection at all
}
