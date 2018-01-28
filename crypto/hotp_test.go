package crypto

import (
	"fmt"
	"testing"
)

func Test_totp(t *testing.T) {
	token, remain, err := TotpToken("EDN5UCREKI4E76IH", 30)
	fmt.Println(token, remain, err)
}
