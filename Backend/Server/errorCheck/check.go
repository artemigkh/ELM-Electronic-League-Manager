package errorCheck

import (
	"github.com/kataras/iris/context"
	"log"
	"github.com/kataras/iris"
)

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckCtx(ctx context.Context, err error) {
	if err != nil {
		log.Fatal(err)
		ctx.StatusCode(iris.StatusBadRequest)
	}
}
