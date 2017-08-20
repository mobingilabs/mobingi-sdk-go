package pretty

import (
	"fmt"
	"testing"
)

func TestIndent(t *testing.T) {
	if i := Indent(4); i != "    " {
		t.Errorf("Expected four(4) whitespaces, got %v", i)
	}
}

type t1 struct {
	S string
}

type t2 struct {
	M   map[string]string
	I   int
	T1  t1
	Pt1 *t1
	St1 []t1
}

func TestJSON(t *testing.T) {
	mck := t2{
		M: map[string]string{"one": "1", "two": "2"},
		I: 100,
		T1: t1{
			S: "struct",
		},
		Pt1: &t1{
			S: "struct pointer",
		},
		St1: make([]t1, 0),
	}

	mck.St1 = append(mck.St1, t1{S: "hello"})
	fmt.Println(JSON(mck, 1, 2))
}
