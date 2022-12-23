package main

import (
	"context"
	"fmt"
	"github.com/Frees0u1/async/src/future"
	"math/rand"
	"time"
)

type UserProfile struct {
	Id     int
	Name   string
	Age    int
	Gender string
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	ctx := context.Background()
	id := 1
	timeoutMs := 80
	nameFuture := future.NewFutureWithTimeout(func() (GetNameResponse, error) {
		return GetName(GetNameRequest{Id: id})
	}, timeoutMs,
	)
	ageFuture := future.NewFutureWithTimeout(func() (GetAgeResponse, error) {
		return GetAge(GetAgeRequest{Id: id})
	}, timeoutMs,
	)
	genderFuture := future.NewFutureWithTimeout(func() (GetGenderResponse, error) {
		return GetGender(GetGenderRequest{Id: id})
	}, timeoutMs)

	var err error
	nameResp, err := nameFuture.Await(ctx)
	if err != nil {
		fmt.Printf("GetName error: %v\n", err)
	}

	ageResp, err := ageFuture.Await(ctx)
	if err != nil {
		fmt.Printf("GetAge error: %v\n", err)
	}
	genderResp, err := genderFuture.Await(ctx)
	if err != nil {
		fmt.Printf("GetGendor error: %v\n", err)
	}

	profile := UserProfile{
		Id:     id,
		Name:   nameResp.Name,
		Age:    ageResp.Age,
		Gender: genderResp.Gender,
	}

	elapsed := time.Since(start)

	fmt.Println(profile)
	fmt.Println(elapsed)
}

type GetNameRequest struct {
	Id int
}

type GetNameResponse struct {
	Id   int
	Name string
}

type GetAgeRequest struct {
	Id int
}

type GetAgeResponse struct {
	Id  int
	Age int
}

type GetGenderRequest struct {
	Id int
}

type GetGenderResponse struct {
	Id     int
	Gender string
}

func GetName(req GetNameRequest) (GetNameResponse, error) {
	sleep := rand.Intn(100) + 100
	fmt.Printf("GetName sleep %d ms\n", sleep)
	time.Sleep(time.Duration(sleep) * time.Millisecond)

	return GetNameResponse{
		Id:   req.Id,
		Name: "Alice",
	}, nil
}

func GetAge(req GetAgeRequest) (GetAgeResponse, error) {
	sleep := rand.Intn(100)
	fmt.Printf("GetAge sleep %d ms\n", sleep)
	time.Sleep(time.Duration(sleep) * time.Millisecond)

	return GetAgeResponse{
		Id:  req.Id,
		Age: 18,
	}, nil
}

func GetGender(req GetGenderRequest) (GetGenderResponse, error) {
	sleep := rand.Intn(100)
	fmt.Printf("GetGendor sleep %d ms\n", sleep)
	time.Sleep(time.Duration(sleep) * time.Millisecond)

	return GetGenderResponse{
		Id:     req.Id,
		Gender: "Female",
	}, nil
}
