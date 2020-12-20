package middleware

import (
	"fmt"
	"gin-web/tests"
	"github.com/panjf2000/ants"
	"sync"
	"testing"
)

func TestGenIdempotenceToken(t *testing.T) {
	tests.InitTestEnv()
	token := GenIdempotenceToken()
	fmt.Println(token)

	defer ants.Release()
	// set size of goroutine pool
	p, _ := ants.NewPool(10000)

	runTimes := 100000
	// Use the common pool.
	var wg sync.WaitGroup
	syncCalculateSum := func() {
		CheckIdempotenceToken(token)
		wg.Done()
	}
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = p.Submit(syncCalculateSum)
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p.Running())
	fmt.Printf("finish all tasks.\n")
}
