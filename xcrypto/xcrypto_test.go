package xcrypto

import (
	"fmt"
	"testing"
)

func TestGenerateRSAPEM(t *testing.T) {
	cert, priv, err := GenerateRSAPEM(1024)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%v\n\n%v\n\n", cert, priv)
	GenerateRSA(1024)
}
