package sectigo

import (
	"fmt"
	"net/http"
	"testing"

	"go.uber.org/zap"
)

func TestClientService_List(t *testing.T) {

	logger, _ := zap.NewProduction()
	c := NewClient(http.DefaultClient, logger, "jiobxn", "id38huAQ23@#sws!", "secure")

	list, err := c.ClientService.List(nil)
	fmt.Println(err)
	fmt.Println(list)

}
