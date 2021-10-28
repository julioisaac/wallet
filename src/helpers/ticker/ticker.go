package ticker

import (
	"fmt"
	"reflect"
	"runtime"
	"time"
)

type ticker struct {}

func NewDaxxerTicker() *ticker {
	return &ticker{}
}

func (*ticker) Run(interval time.Duration, f func()) {
	go func() {
		fmt.Printf("starting daxxer ticker for %s\n", GetFunctionName(f))
		ticker := time.NewTicker(interval * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				f()
			}
		}
	}()
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}