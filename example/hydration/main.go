package main

import (
	"async/src/future"
	"context"
	"fmt"
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
	nameFuture := future.NewFuture(func() (GetNameResponse, error) {
		return GetName(GetNameRequest{Id: id})
	})
	ageFuture := future.NewFuture(func() (GetAgeResponse, error) {
		return GetAge(GetAgeRequest{Id: id})
	})
	genderFuture := future.NewFuture(func() (GetGenderResponse, error) {
		return GetGender(GetGenderRequest{Id: id})
	})

	timeout := 80 * time.Millisecond
	nameResp, _ := nameFuture.Await(ctx, &timeout)
	ageResp, _ := ageFuture.Await(ctx, &timeout)
	genderResp, _ := genderFuture.Await(ctx, &timeout)

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
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return GetNameResponse{
		Id:   req.Id,
		Name: "Alice",
	}, nil
}

func GetAge(req GetAgeRequest) (GetAgeResponse, error) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return GetAgeResponse{
		Id:  req.Id,
		Age: 18,
	}, nil
}

func GetGender(req GetGenderRequest) (GetGenderResponse, error) {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return GetGenderResponse{
		Id:     req.Id,
		Gender: "Female",
	}, nil
}
