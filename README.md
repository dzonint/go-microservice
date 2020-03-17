# go-microservice
REST Microservice developed in Go. Works with following struct:
```go
type Product struct {
	ID          int     `json:"id" storm:"id,increment=3"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float32 `json:"price" validate:"required,gt=0"`
	SKU         string  `json:"sku" validate:"required,sku"`
	CreatedOn   string  `json:"-"`
	UpdatedOn   string  `json:"-"`
	DeletedOn   string  `json:"-"`
}
```

## Packages used
- [Gorilla/mux](https://github.com/gorilla/mux) for routing
- [Storm](https://github.com/asdine/storm) for DB storage
- [Validator](https://github.com/go-playground/validator) for data validation

## Usage
#### Fetch all products
```sh
curl localhost:9090
```

#### Fetch a specific product
```sh
curl localhost:9090/{id}
```
#### Add product
```sh
curl localhost:9090 -X POST -d `{"name":"New product", "price":3.22, "sku":"abc-abcd-abcde"...}`
```

#### Update product
```sh
curl localhost:9090/{id} -X PUT -d `{"name":"Updated", "price":1.99, "sku":"abc-abcd-abcde"...}`
```

#### Remove product
```sh
curl localhost:9090/{id} -X DELETE
```
