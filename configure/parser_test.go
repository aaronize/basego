package configure

import (
	"log"
	"testing"
)

func TestParser(t *testing.T) {
	config := &struct {}{}
	cp := NewConfigParser("", "json")

	if err := cp.Parse(config); err != nil {
		log.Fatal(err.Error())
	}
}
