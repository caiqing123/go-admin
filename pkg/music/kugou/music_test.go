package kugou

import (
	"fmt"
	"testing"
)

func TestMusic(*testing.T) {
	d, _ := Kugou("qq")
	fmt.Println(d)
}
