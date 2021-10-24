package common

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	for i := 30; i <= 37; i++ {
		fmt.Println(GetColorStr("hello", Color(i)))
	}
}
