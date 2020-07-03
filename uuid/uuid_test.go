package uuid

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	fmt.Println(New())
	fmt.Println(MID())
}
