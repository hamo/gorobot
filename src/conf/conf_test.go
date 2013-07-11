package conf

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	filename, _ := filepath.Abs("freenode.yaml")
	cs, err := ParseFile(filename)
	if err != nil {
		t.Errorf("ParseFile failed!\n")
	}

	fmt.Printf("%#+v\n", cs)

	fmt.Printf("%#+v\n", cs.Conn)
	fmt.Printf("%#+v\n", cs.Channels)
}
