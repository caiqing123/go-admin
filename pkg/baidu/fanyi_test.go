package baidu

import (
	"fmt"
	"testing"
	"time"
)

func TestFanyi(t *testing.T) {
	tk, err := Gtk()
	if err != nil {
		fmt.Println(err)
		return
	}

	options := NewOptions(EN, ZH, tk, "")

	r, err := Do("Save the field to set the <b>Related services</b> settings", options)
	if err != nil {
		fmt.Println("err")
		fmt.Println(err)
		return
	}

	time.Sleep(time.Second * 1)

	fmt.Printf("%s -> %s \n", r.TransResult.Data[0].Src, r.TransResult.Data[0].Dst)

	options.To = KOR
	r, err = Do("Save the field to set the <b>Related services</b> settings", options)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s -> %s \n", r.TransResult.Data[0].Src, r.TransResult.Data[0].Dst)

}
