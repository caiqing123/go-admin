package kuwo

import (
	"fmt"
	"testing"

	"api/pkg/music/netease"
)

func TestMusic(*testing.T) {
	d, err := netease.Netease("", "1")
	fmt.Println(err)
	fmt.Println(d)

}
