package ticker

import (
	"context"
	"github.com/julioisaac/daxxer-api/internal/logs"
	"reflect"
	"runtime"
	"time"
)

type ticker struct {}

func NewDaxxerTicker() *ticker {
	return &ticker{}
}

func (*ticker) Run(ctx context.Context, interval time.Duration, f func(context.Context) error) {
	go func() {
		logs.Instance.Log.Debug(context.Background(), "starting daxxer ticker for"+GetFunctionName(f))
		ticker := time.NewTicker(interval * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				f(ctx)
			}
		}
	}()
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}