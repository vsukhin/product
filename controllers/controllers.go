package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/astaxie/beego/validation"
	"github.com/gorilla/mux"

	"github.com/vsukhin/product/logger"
	"github.com/vsukhin/product/models"
	"github.com/vsukhin/product/repositories"
)

// ProductControllerInterface is interface of product managing controller
type ProductControllerInterface interface {
	List(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Get(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	SetPrices(w http.ResponseWriter, r *http.Request)
}

// ProductControllerImplementation is implementation of product managing controller
type ProductControllerImplementation struct {
	ProductRepository repositories.ProductRepositoryInterface
}

var _ ProductControllerInterface = &ProductControllerImplementation{}

// NewProductControllerImplementation is a constructor for ProductControllerImplementation
func NewProductControllerImplementation(productRepository repositories.ProductRepositoryInterface) (productController ProductControllerInterface) {
	return &ProductControllerImplementation{ProductRepository: productRepository}
}

// List returns a list of products
func (productcontroller *ProductControllerImplementation) List(w http.ResponseWriter, r *http.Request) {
	products, err := productcontroller.ProductRepository.GetAll()
	if err != nil {
		logger.Log.Println(err)
		http.Error(w, "can't get products", http.StatusInternalServerError)
		return
	}

	err = productcontroller.renderJSON(w, http.StatusOK, products)
	if err != nil {
		http.Error(w, "can't return JSON", http.StatusInternalServerError)
		return
	}
}

// Create creates a new product
func (productcontroller *ProductControllerImplementation) Create(w http.ResponseWriter, r *http.Request) {
	product, err := productcontroller.readProduct(r)
	if err != nil {
		http.Error(w, "can't retrieve product", http.StatusBadRequest)
		return
	}

	errs, err := productcontroller.validateProduct(product)
	if err != nil {
		http.Error(w, "can't validate product", http.StatusBadRequest)
		return
	}
	if len(errs) != 0 {
		err = productcontroller.renderJSON(w, http.StatusBadRequest, errs)
		if err != nil {
			http.Error(w, "can't return JSON", http.StatusInternalServerError)
		}
		return
	}

	id, err := productcontroller.ProductRepository.New(product)
	if err != nil {
		http.Error(w, "can't create product", http.StatusInternalServerError)
		return
	}

	err = productcontroller.renderJSON(w, http.StatusOK, models.ID{ID: id})
	if err != nil {
		http.Error(w, "can't return JSON", http.StatusInternalServerError)
		return
	}
}

// Get returns a product ref. by id
func (productcontroller *ProductControllerImplementation) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	product, err := productcontroller.findProduct(vars["id"])
	if err != nil {
		http.Error(w, "can't find product", http.StatusNotFound)
		return
	}

	err = productcontroller.renderJSON(w, http.StatusOK, product)
	if err != nil {
		http.Error(w, "can't return JSON", http.StatusInternalServerError)
		return
	}
}

// Update updates product data
func (productcontroller *ProductControllerImplementation) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oldproduct, err := productcontroller.findProduct(vars["id"])
	if err != nil {
		http.Error(w, "can't find product", http.StatusNotFound)
		return
	}

	newproduct, err := productcontroller.readProduct(r)
	newproduct.ID = oldproduct.ID
	newproduct.Prices = oldproduct.Prices
	if err != nil {
		http.Error(w, "can't retrieve product", http.StatusBadRequest)
		return
	}

	errs, err := productcontroller.validateProduct(newproduct)
	if err != nil {
		http.Error(w, "can't validate product", http.StatusBadRequest)
		return
	}
	if len(errs) != 0 {
		err = productcontroller.renderJSON(w, http.StatusBadRequest, errs)
		if err != nil {
			http.Error(w, "can't return JSON", http.StatusInternalServerError)
		}
		return
	}

	err = productcontroller.ProductRepository.Update(newproduct)
	if err != nil {
		http.Error(w, "can't update product", http.StatusInternalServerError)
		return
	}
}

// Delete deletes product
func (productcontroller *ProductControllerImplementation) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	product, err := productcontroller.findProduct(vars["id"])
	if err != nil {
		http.Error(w, "can't find product", http.StatusNotFound)
		return
	}

	err = productcontroller.ProductRepository.Delete(product.ID)
	if err != nil {
		logger.Log.Println(err)
		http.Error(w, "can't delete product", http.StatusInternalServerError)
	}
}

// SetPrices sets additional prices for a product
func (productcontroller *ProductControllerImplementation) SetPrices(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	product, err := productcontroller.findProduct(vars["id"])
	if err != nil {
		http.Error(w, "can't find product", http.StatusNotFound)
		return
	}

	var prices map[string]float64
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Log.Println(err)
		http.Error(w, "can't retrieve prices", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &prices)
	if err != nil {
		logger.Log.Println(err)
		http.Error(w, "can't retrieve prices", http.StatusBadRequest)
		return
	}

	if _, ok := prices[models.DefaultCurrency]; ok {
		logger.Log.Println("only for additional currencies")
		http.Error(w, "only for additional currencies", http.StatusBadRequest)
		return
	}

	for _, price := range prices {
		if price < 0 {
			logger.Log.Println("price can't be negative")
			http.Error(w, "price can't be negative", http.StatusBadRequest)
			return
		}
	}

	product.Prices = prices
	err = productcontroller.ProductRepository.Update(product)
	if err != nil {
		http.Error(w, "can't update prices", http.StatusInternalServerError)
		return
	}
}

// findProduct returns product ref. by id
func (productcontroller *ProductControllerImplementation) findProduct(id string) (product models.Product, err error) {
	productid, err := strconv.Atoi(id)
	if err != nil {
		logger.Log.Println(err)
		return models.Product{}, err
	}

	product, err = productcontroller.ProductRepository.Get(productid)
	if err != nil {
		logger.Log.Println(err)
		return models.Product{}, err
	}

	return product, nil
}

// readProduct returns product from http request
func (productcontroller *ProductControllerImplementation) readProduct(r *http.Request) (product models.Product, err error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logger.Log.Println(err)
		return models.Product{}, err
	}

	err = json.Unmarshal(body, &product)
	if err != nil {
		logger.Log.Println(err)
		return models.Product{}, err
	}

	return product, nil
}

// validateProduct validates product
func (productcontroller *ProductControllerImplementation) validateProduct(product models.Product) ([]models.ValidationError, error) {
	validator := validation.Validation{}
	valid, err := validator.Valid(product)
	if err != nil {
		logger.Log.Println(err)
		return []models.ValidationError{}, err
	}
	if !valid {
		var errs []models.ValidationError
		for _, vErr := range validator.Errors {
			logger.Log.Println(vErr.Field + ":" + vErr.Message)
			errs = append(errs, models.ValidationError{Field: vErr.Field, Message: vErr.Message})
		}
		return errs, nil
	}

	if product.Price < 0 {
		logger.Log.Println("price can't be negative")
		return []models.ValidationError{models.ValidationError{Field: "Price", Message: "price can't be negative"}}, nil
	}

	return []models.ValidationError{}, nil
}

// renderJSON renders object in JSON
func (productcontroller *ProductControllerImplementation) renderJSON(w http.ResponseWriter, code int, object interface{}) error {
	js, err := json.Marshal(object)
	if err != nil {
		logger.Log.Println(err)
		return err
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(js)

	return nil
}
