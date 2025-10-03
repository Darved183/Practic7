package main

import (
	"errors"
	"fmt"
	"sync"
)

type Product struct {
	ID       int
	Name     string
	Price    float64
	Quantity int
}

type Inventory struct {
	products map[int]Product
	mux      sync.RWMutex
}

func (inv *Inventory) AddProduct(product Product) {
	inv.mux.Lock()
	defer inv.mux.Unlock()

	if inv.products == nil {
		inv.products = make(map[int]Product)
	}
	inv.products[product.ID] = product
}

func (inv *Inventory) WriteOff(productID int, quantity int) error {
	inv.mux.Lock()
	defer inv.mux.Unlock()

	product, exists := inv.products[productID]
	if !exists {
		return errors.New("продукт не найден")
	}
	if quantity > product.Quantity {
		return errors.New("товара нет")
	}
	product.Quantity -= quantity
	inv.products[productID] = product
	return nil
}

func (inv *Inventory) RemoveProduct(productID int) error {
	inv.mux.Lock()
	defer inv.mux.Unlock()

	if _, exists := inv.products[productID]; !exists {
		return errors.New("продукт не найден")
	}
	delete(inv.products, productID)
	return nil
}

func (inv *Inventory) GetTotalValue() float64 {
	inv.mux.RLock()
	defer inv.mux.RUnlock()

	total := 0.0
	for _, product := range inv.products {
		total += product.Price * float64(product.Quantity)
	}
	return total
}

func main() {
	inventory := Inventory{}

	//ИИ

	inventory.AddProduct(Product{ID: 1, Name: "Ноутбук", Price: 50000, Quantity: 5})
	inventory.AddProduct(Product{ID: 2, Name: "Мышь", Price: 1000, Quantity: 10})
	inventory.AddProduct(Product{ID: 3, Name: "Клавиатура", Price: 2000, Quantity: 8})

	fmt.Printf("Общая стоимость товаров: %.2f\n", inventory.GetTotalValue())

	err := inventory.WriteOff(1, 2)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Println("Списано")
	}

	err = inventory.RemoveProduct(3)
	if err != nil {
		fmt.Println("Ошибка удаления:", err)
	} else {
		fmt.Println("Клавиатура удалена со склада")
	}
}
