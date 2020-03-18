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
- [AMQP](https://github.com/streadway/amqp) for RabbitMQ integration
- [Logrus](https://github.com/sirupsen/logrus) for logging

## Other resources
- [Mockaroo](https://https://www.mockaroo.com) for mock API
- [CloudAMQP](https://www.cloudamqp.com) for RabbitMQ hosting


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

### RabbitMQ
You can run the service with the following command
```sh
go run main.go -rabbitmq
```

This will enable RabbitMQ service and will run both the consumer and producer.
- Service uses the following data model:
```go
type User struct {
	ID        int       `json:"id" storm:"id,increment=1"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Gender    string    `json:"gender"`
	IPAddress string    `json:"ip_address"`
	CreatedAt time.Time `json:"created_at"`
}
```
- Producer will publish the message to queue every 10 seconds.
- The message is fetched from mock API.
- Consumer will consume the message and insert the user into BoltDB database.

#### Fetch all users
```sh
curl localhost:9090/users
```

