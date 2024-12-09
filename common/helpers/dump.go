package helpers

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"

	"github.com/gookit/goutil/dump"
)

func Println(vs ...any) {

	pc, file, line, _ := runtime.Caller(1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println()

		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")
		fmt.Println(file + ":" + strconv.Itoa(line))
		fmt.Println(runtime.FuncForPC(pc).Name())
		fmt.Println("++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++")

		dump.Println(vs)

		fmt.Println("____________________________________________________________")

		fmt.Println()
	}()

	wg.Wait()
}
