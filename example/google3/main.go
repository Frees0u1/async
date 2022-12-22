package main

import (
	"async/src/future"
	"context"
	"fmt"
	"math/rand"
	"time"
)

var (
	Web    = fakeSearch("web")
	Web2   = fakeSearch("web")
	Image  = fakeSearch("image")
	Image2 = fakeSearch("image")
	Video  = fakeSearch("video")
	Video2 = fakeSearch("video")
)

type Result string
type Search func(query string) Result

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	results := Google("golang") // collate results
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

func webSearch(ctx context.Context, query string) (Result, error) {
	fmt.Println("coming to web search")
	return future.First(ctx, nil, []future.Future[Result]{
		future.NewFuture(func() (Result, error) {
			return Web(query), nil
		}),
		future.NewFuture(func() (Result, error) {
			return Web2(query), nil
		}),
	})
}

func imageSearch(ctx context.Context, query string) (Result, error) {
	return future.First(ctx, nil, []future.Future[Result]{
		future.NewFuture(func() (Result, error) {
			return Image(query), nil
		}),
		future.NewFuture(func() (Result, error) {
			return Image2(query), nil
		}),
	})
}

func videoSearch(ctx context.Context, query string) (Result, error) {
	return future.First(ctx, nil, []future.Future[Result]{
		future.NewFuture(func() (Result, error) {
			return Video(query), nil
		}),
		future.NewFuture(func() (Result, error) {
			return Video2(query), nil
		}),
	})
}

func Google(query string) (results []Result) {
	ctx := context.Background()
	timeout := 80 * time.Millisecond

	searches := []future.Future[Result]{
		future.NewFuture(func() (Result, error) {
			return webSearch(ctx, query)
		}),
		future.NewFuture(func() (Result, error) {
			return imageSearch(ctx, query)
		}),
		future.NewFuture(func() (Result, error) {
			return videoSearch(ctx, query)
		}),
	}

	results, err := future.AwaitAll(ctx, &timeout, searches)
	if err == nil {
		return results
	} else {
		fmt.Printf("err happens! %v\n", err)
	}

	return []Result{}
}
