package kuwo

import (
	"fmt"
	"testing"

	"api/pkg/music/kugou"
)

func TestMusic(*testing.T) {
	d, err := kugou.NewKugou("qq", "1")
	fmt.Println(err)
	fmt.Println(d)
}
