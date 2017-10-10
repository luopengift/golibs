package sys

import "testing"

func Test_user(t *testing.T) {
    t.Log(UserExist("root"))
}
