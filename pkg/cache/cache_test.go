package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestFileGet(t *testing.T) {
	var rds interface{ Store }
	rds = NewRedisStore(
		fmt.Sprintf("%v:%v", "127.0.0.1", "6379"),
		"",
		"",
		2,
	)
	InitWithCacheStore(rds)

	data := map[string]interface{}{
		"name": "test",
		"age":  18,
	}
	d := Caches("test", data, time.Second*60)
	fmt.Println(d)
	//b, _ := json.Marshal(d)
	//_ = json.Unmarshal(b, &data)
	//fmt.Println(data["name"])

	data1 := `{"name":"test","age":18}`
	d1 := Caches("test1", data1, time.Second*60)
	fmt.Println(d1)

	data2 := func() interface{} {
		return map[string]interface{}{
			"name": "test22",
			"age":  18,
		}
	}
	d2 := Caches("test2", data2, time.Second*60)
	fmt.Println(d2)

	data3 := func() interface{} {
		return "test22"
	}
	d3 := Caches("test3", data3, time.Second*60)
	fmt.Println(d3)
}
