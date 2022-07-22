package kuwo

import (
	"fmt"
	"testing"

	"api/pkg/music/qqmusic"
)

func TestMusic(*testing.T) {
	d, _ := qqmusic.QQMusic("qq", "1")
	fmt.Println(d)

}
