package main

import (
	"async/src/future"
	"context"
	"fmt"
	"math/rand"
	"time"
)

var (
	Web   = fakeSearch("web1")
	Image = fakeSearch("image1")
	Video = fakeSearch("video2")
)

type Result string
type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

func Google(query string) (results []Result) {
	ctx := context.Background()

	webSearch := future.NewFuture(func() (Result, error) {
		return fakeSearch("web")(query), nil
	})
	imageSearch := future.NewFuture(func() (Result, error) {
		return fakeSearch("image")(query), nil
	})
	videoSearch := future.NewFuture(func() (Result, error) {
		return fakeSearch("video")(query), nil
	})

	if x, e := webSearch.Await(ctx, nil); e == nil {
		results = append(results, x)
	}

	if x, e := imageSearch.Await(ctx, nil); e == nil {
		results = append(results, x)
	}

	if x, e := videoSearch.Await(ctx, nil); e == nil {
		results = append(results, x)
	}

	return results
}

func main() {
	results := Google("Hello World!")

	for _, r := range results {
		fmt.Println(r)
	}
}
