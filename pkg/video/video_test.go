package video

import (
	"fmt"
	"testing"
)

func TestVideo(*testing.T) {
	d, err := QueryResources("", "", "1", "")
	fmt.Println(err)
	fmt.Println(d)
}
