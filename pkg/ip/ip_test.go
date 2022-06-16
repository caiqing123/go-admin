package ip

import (
	"fmt"
	"testing"
)

func TestGetClientIP(t *testing.T) {
	fmt.Println(GetLocation("223.104.67.113"))
}
