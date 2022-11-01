package kuwo

import (
	"fmt"
	"testing"

	"api/pkg/music/qqmusic"
)

func TestMusic(*testing.T) {
	d, err := qqmusic.QQMusic("qq", "1")
	fmt.Println(err)
	fmt.Println(d)
}
