package excelize

import (
	"fmt"
	"testing"

	"api/app/models/user"
)

func TestRead(t *testing.T) {
	userModel := user.User{}
	rows, err := ImportToPath("./demo.xlsx", userModel)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rows)
	for _, row := range rows {
		for _, colCell := range row {
			fmt.Print(colCell, "\t")
		}
		userModel.Name = row[0]
		fmt.Println()
	}
}
