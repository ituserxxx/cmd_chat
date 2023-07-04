package comm

import (
	"encoding/json"
	"testing"
)

func TestName(t *testing.T) {
	v, _ := json.Marshal("encoding/json")
	v1 := B64Encode(v)
	println(v1)

	println(string(B64Encry(v1)))
}
