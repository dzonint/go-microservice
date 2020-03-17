package data

import "fmt"

var ErrFailedToOpenDB = fmt.Errorf("Unable to open database")
var ErrFailedToAddUser = fmt.Errorf("Unable to add user to database")
var ErrUserNotFound = fmt.Errorf("User not found")
var ErrFailedToUpdateDB = fmt.Errorf("Unable to update database")
var ErrFailedToGetUsers = fmt.Errorf("Unable to fetch users")
var ErrFailedToAddProduct = fmt.Errorf("Unable to add product to database")
var ErrProductNotFound = fmt.Errorf("Product not found")
var ErrFailedToGetProducts = fmt.Errorf("Unable to fetch products")
