package user

import "testing"

func TestAuth(t *testing.T) {
	_, token := Login("1234", "abcd")
	t.Log(token)
}
