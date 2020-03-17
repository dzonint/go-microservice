package data

import (
	"encoding/json"
	"io"
	"regexp"
	"time"

	"github.com/asdine/storm"
	"github.com/go-playground/validator/v10"
)

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

type Products []*Product

func (p *Product) Validate() error {
	validate := validator.New()
	validate.RegisterValidation("sku", validateSKU)
	return validate.Struct(p)
}

func validateSKU(fl validator.FieldLevel) bool {
	// SKU format: abc-abcd-abcde
	re := regexp.MustCompile(`[a-z]+-[a-z]+-[a-z]+`)
	matches := re.FindAllString(fl.Field().String(), -1)

	if len(matches) != 1 {
		return false
	}

	return true
}

func (p *Products) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Product) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func PopulateDB() error {
	var err error
	for _, val := range productList {
		err = AddProduct(val)
		if err != nil {
			return err
		}
	}
	return nil
}

func AddProduct(p *Product) error {
	db, err := storm.Open("products.db")
	if err != nil {
		return ErrFailedToOpenDB
	}
	defer db.Close()

	p.CreatedOn = time.Now().String()

	err = db.Save(p)
	if err != nil {
		return ErrFailedToAddProduct
	}

	return nil
}

func UpdateProduct(id int, p *Product) error {
	db, err := storm.Open("products.db")
	if err != nil {
		return ErrFailedToOpenDB
	}
	defer db.Close()

	var prod Product
	err = db.One("ID", id, &prod)
	if err != nil {
		return ErrProductNotFound
	}

	p.ID = id
	p.UpdatedOn = time.Now().String()
	err = db.Update(p)
	if err != nil {
		return ErrFailedToUpdateDB
	}

	return nil
}

func RemoveProduct(id int) error {
	db, err := storm.Open("products.db")
	if err != nil {
		return ErrFailedToOpenDB
	}
	defer db.Close()

	var prod Product
	err = db.One("ID", id, &prod)
	if err != nil {
		return ErrProductNotFound
	}

	prod.ID = id
	err = db.DeleteStruct(&prod)
	if err != nil {
		return ErrFailedToUpdateDB
	}

	return nil
}

func GetProduct(id int) (Product, error) {
	var Prod Product
	db, err := storm.Open("products.db")
	if err != nil {
		return Product{}, ErrFailedToOpenDB
	}
	defer db.Close()

	err = db.One("ID", id, &Prod)
	if err != nil {
		return Product{}, ErrFailedToGetProducts
	}

	return Prod, nil
}

func GetProducts() (Products, error) {
	db, err := storm.Open("products.db")
	if err != nil {
		return nil, ErrFailedToOpenDB
	}
	defer db.Close()

	var Products []*Product
	err = db.All(&Products)
	if err != nil {
		return nil, ErrFailedToGetProducts
	}

	return Products, nil
}

var productList = []*Product{
	&Product{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc-sdfv-opfdm",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
	&Product{
		ID:          2,
		Name:        "Espresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "has-hafg-emrbo",
		CreatedOn:   time.Now().UTC().String(),
		UpdatedOn:   time.Now().UTC().String(),
	},
}
