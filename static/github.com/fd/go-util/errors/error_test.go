package errors

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os/exec"
	"strings"
	"testing"
)

func TestAnnotateNested(t *testing.T) {
	err1 := New("%s err", "foo")
	err1.AddContext("hello=%s", "world")

	err2 := Annotate(err1, "%s err", "bar")
	err2.AddContext("c=%d", 7)
	err2.AddContext("a=%d", 42)

	diff(t, "new", err1.Error())
	diff(t, "annotate", err2.Error())
}

var write_generated = flag.Bool("write-generated", false, "Write generated error messages")

func diff(t *testing.T, golden, generated string) {
	if write_generated != nil && *write_generated {
		err := ioutil.WriteFile("testdata/"+golden+".txt", []byte(generated), 0644)
		if err != nil {
			panic(err)
		}
	}

	data, err := ioutil.ReadFile("testdata/" + golden + ".txt")
	if err != nil {
		panic(err)
	}

	golden_str := string(data)
	golden_str = strings.TrimRight(golden_str, "\n") // trim off trailing \n

	if golden_str == generated {
		return
	}

	cmd := exec.Command("diff", "--ignore-space-change", "-U", "5", "--to-file", "-", "testdata/"+golden+".txt")
	cmd.Stdin = bytes.NewReader([]byte(generated))
	output, err := cmd.CombinedOutput()
	if _, ok := err.(*exec.ExitError); ok && err.Error() == "exit status 1" {
		t.Errorf("%s failed:\n%s\n", golden, output)
		return
	}
	if err != nil {
		t.Error(string(output))
		panic(err)
	}
}
