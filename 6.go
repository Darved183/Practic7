package main

import (
	"fmt"
	"sync"
)

type EventBus struct {
	subscribers map[string][]func(interface{})
	mux         sync.RWMutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]func(interface{})),
	}
}

func (eb *EventBus) Subscribe(event string, handler func(interface{})) {
	eb.mux.Lock()
	defer eb.mux.Unlock()

	eb.subscribers[event] = append(eb.subscribers[event], handler)
}

func (eb *EventBus) Publish(event string, data interface{}) {
	eb.mux.RLock()
	defer eb.mux.RUnlock()

	if handlers, exists := eb.subscribers[event]; exists {
		for _, handler := range handlers {
			handler(data)
		}
	}
}

func main() {
	bus := NewEventBus()

	//ИИ

	bus.Subscribe("user.created", func(data interface{}) {
		fmt.Printf("Пользователь создан: %v\n", data)
	})

	bus.Subscribe("user.created", func(data interface{}) {
		user := data.(map[string]string)
		fmt.Printf("Отправляем email для %s\n", user["email"])
	})

	bus.Subscribe("order.placed", func(data interface{}) {
		order := data.(map[string]interface{})
		fmt.Printf("Заказ размещен: №%d, сумма: %.2f\n", order["id"], order["amount"])
	})

	fmt.Println("Публикуем событие user.created:")
	bus.Publish("user.created", map[string]string{
		"name":  "Иван Иванов",
		"email": "ivan@mail.ru",
	})

	fmt.Println("\nПубликуем событие order.placed:")
	bus.Publish("order.placed", map[string]interface{}{
		"id":     12345,
		"amount": 1500.50,
	})

	fmt.Println("\nПубликуем событие без подписчиков:")
	bus.Publish("unknown.event", "тестовые данные")
}
