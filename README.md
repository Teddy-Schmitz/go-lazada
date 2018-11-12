# Go Lazada [![GoDoc](https://godoc.org/github.com/Teddy-Schmitz/go-lazada/lazada?status.png)](https://godoc.org/github.com/Teddy-Schmitz/go-lazada/lazada)
Unofficial Go library for the [Lazada Open Platform]("https://open.lazada.com/")

## Install

```
go get github.com/Teddy-Schmitz/go-lazada
```

## Usage

```go
import "github.com/Teddy-Schmitz/go-lazada/lazada"
```

Construct a client with a specific region

```go
client := lazada.NewClient("AppKey", "AppSecret", lazada.Singapore)
```

Call a service
```go
products, err := client.Products.Get(context.Background,  &lazada.SearchOptions{Filter: "live", Limit: 100, SKUSellerList: &out})
```
See the godoc for a list of available services and methods.

Some services require a client token you can add it to a client like this

```go
userClient := client.NewTokenClient("Token")
```
This returns a new client with the token set.  So you can keep the old generic client and use the token client for a specific user.

You can also change the region if necessary.

```go
client.SetRegion(lazada.Malasyia)
```

### Available APIs

- Products
- Auth (System)

## TODO

Lots of services and methods have not been added as I have only done the ones that I needed.  Pull Requests are welcome to add news.

## License

MIT License