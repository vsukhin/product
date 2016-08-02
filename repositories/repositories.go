package repositories

import (
	"errors"
	"sort"
	"sync"

	"github.com/vsukhin/product/models"
)

// ProductRepositoryInterface is interface of product model managing
type ProductRepositoryInterface interface {
	New(product models.Product) (int, error)
	Update(product models.Product) error
	Delete(id int) error
	Get(id int) (models.Product, error)
	GetAll() ([]models.Product, error)
}

// ProductRepositoryImplementation is implementation of product model managing
type ProductRepositoryImplementation struct {
	id       int
	products map[int]models.Product

	sync.RWMutex
}

var _ ProductRepositoryInterface = &ProductRepositoryImplementation{}

// NewProductRepositoryImplementation is a constructor for ProductRepositoryImplementation
func NewProductRepositoryImplementation() (productrepository ProductRepositoryInterface) {
	return &ProductRepositoryImplementation{id: 1, products: make(map[int]models.Product)}
}

// New creates new product
func (productrepository *ProductRepositoryImplementation) New(product models.Product) (int, error) {
	productrepository.Lock()
	defer productrepository.Unlock()

	product.ID = productrepository.id
	productrepository.products[product.ID] = product
	productrepository.id++

	return product.ID, nil
}

// Update updates product content
func (productrepository *ProductRepositoryImplementation) Update(product models.Product) error {
	if _, ok := productrepository.products[product.ID]; !ok {
		return errors.New("Product is not found")
	}

	productrepository.Lock()
	defer productrepository.Unlock()

	productrepository.products[product.ID] = product

	return nil
}

// Delete deletes a product
func (productrepository *ProductRepositoryImplementation) Delete(id int) error {
	if _, ok := productrepository.products[id]; !ok {
		return errors.New("Product is not found")
	}

	productrepository.Lock()
	defer productrepository.Unlock()

	delete(productrepository.products, id)

	return nil
}

// Get returns a product
func (productrepository *ProductRepositoryImplementation) Get(id int) (models.Product, error) {
	if _, ok := productrepository.products[id]; !ok {
		return models.Product{}, errors.New("Product is not found")
	}

	productrepository.RLock()
	defer productrepository.RUnlock()

	return productrepository.products[id], nil
}

// GetAll returns all products
func (productrepository *ProductRepositoryImplementation) GetAll() ([]models.Product, error) {
	var products ByProduct

	productrepository.RLock()
	for id := range productrepository.products {
		products = append(products, productrepository.products[id])
	}
	productrepository.RUnlock()

	sort.Sort(products)

	return products, nil
}

// ByProduct is a slice for sorting by product
type ByProduct []models.Product

// Len is a length of the slice
func (s ByProduct) Len() int { return len(s) }

// Swap is a swapping for slice elements
func (s ByProduct) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less is a comparision of slice elements
func (s ByProduct) Less(i, j int) bool { return s[i].ID < s[j].ID }
