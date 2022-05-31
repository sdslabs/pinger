package central

import (
	"fmt"

	"github.com/sdslabs/pinger/pkg/util/appcontext"
)

func Run(ctx *appcontext.Context) error {
	fmt.Println("Central server is running!")
	return nil
}
