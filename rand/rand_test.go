package rand

import (
	"log"
	"testing"
)

func TestRandom(t *testing.T) {

	log.Printf("\n%v\n", RandomString(23))
	if len(RandomString(23)) != 23 {
		log.Printf("length: %v\n", len(RandomString(23)))
		t.Fail()
	}

	if RandomString(23) == RandomString(23) {
		t.Fail()
	}

}
