package main

import (
	"fmt"
	"sync"
)

type Customer struct {
	ID   int
	Name string
}

type OrderItem struct {
	ProductID int
	Name      string
	Price     float64
	Quantity  int
}

type Order struct {
	ID       int
	Customer Customer
	Items    []OrderItem
	Status   string
	mux      sync.RWMutex
}

func (o *Order) AddItem(item OrderItem) {
	o.mux.Lock()
	defer o.mux.Unlock()
	o.Items = append(o.Items, item)
}

func (o *Order) RemoveItem(productID int) {
	o.mux.Lock()
	defer o.mux.Unlock()
	for i, item := range o.Items {
		if item.ProductID == productID {
			o.Items = append(o.Items[:i], o.Items[i+1:]...)
			return
		}
	}
}

func (o *Order) GetTotal() float64 {
	o.mux.RLock()
	defer o.mux.RUnlock()
	total := 0.0
	for _, item := range o.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func (o *Order) SetStatus(status string) {
	o.mux.Lock()
	defer o.mux.Unlock()
	o.Status = status
}

func (o *Order) GetStatus() string {
	o.mux.RLock()
	defer o.mux.RUnlock()
	return o.Status
}

func main() {
	customer := Customer{ID: 1, Name: "Петр"}
	order := Order{ID: 1, Customer: customer, Status: "Создан"}

	//ИИ

	order.AddItem(OrderItem{ProductID: 1, Name: "Ноутбук", Price: 50000, Quantity: 1})
	order.AddItem(OrderItem{ProductID: 2, Name: "Мышь", Price: 1000, Quantity: 2})

	fmt.Printf("Заказ №%d, Клиент: %s\n", order.ID, order.Customer.Name)
	fmt.Printf("Статус: %s\n", order.GetStatus())
	fmt.Printf("Общая сумма: %.2f\n", order.GetTotal())

	order.RemoveItem(2)
	fmt.Printf("После удаления мыши: %.2f\n", order.GetTotal())

	order.SetStatus("Оплачен")
	fmt.Printf("Новый статус: %s\n", order.GetStatus())
}
