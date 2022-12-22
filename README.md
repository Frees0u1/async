# Async
Write golang code in async/await way

## Example
For most production hydration usage, we need fire fetch request asynchronously, and gathering all the response back,
and then hydrated all them to a final response. **Future** will helps such task a lot! 

```go
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
```

