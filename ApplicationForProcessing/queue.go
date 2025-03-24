package applicationforprocessing

import (
	"fmt"
	"math/rand"
	productsearch "multi-threadeddataprocessingsysytem/ProductSearch"
	"strconv"
	"sync"
	"time"
)

//В данном пакете обрабатываются разные виды заявок, такие как покупка, отслеживаение, отмена заказа

//Структура заказа

type Order struct {
	NameProduct []string
	Price       float64
	Number      int
}

func intn() {
	rand.Seed(time.Now().Unix())
}

// Новый заказ
func NewOrder(name []string, price float64, number int) Order {
	return Order{NameProduct: name, Price: price, Number: number}
}

func MakeOrder(orderMap *map[int]Order, catalog productsearch.Catalog, wg *sync.WaitGroup, hashChanel chan int) {
	names := []string{}  // Массив под названия продукта
	price := 0.0         // Конечная сумма заказа
	var n, num, pnum int // Номер заказа, номер раздела, номер продукта
	var action string

	defer wg.Done()
	visited := catalog.DFS() //Visited массив
	if len(visited) == 0 {
		fmt.Println("В каталоге Yandex Market (beta) нет разделов")
		fmt.Println()
		return
	}
	for i, chapter := range visited {
		fmt.Printf("%d. %s\n", i, chapter.Name)
	}
	for {
		fmt.Println("Введите номер интересующего раздела.Для выхода нажмите q")
		fmt.Scanln(&action)
		if action == "q" {
			break
		} else {
			num, _ = strconv.Atoi(action)
		}
		for j, product := range visited[num].Products {
			fmt.Printf("%d.Товар: %s\nЦена:%2.f\nКол-во:%d\n", j, product.Name, product.Price, product.NumberOfProduct)
			fmt.Println("----------------------")
		}
		fmt.Println("Введите номера интересующих вас товаров. Для выхода нажмите q")
		for {
			fmt.Scanln(&action)
			if action == "q" {
				break
			} else {
				pnum, _ = strconv.Atoi(action)
			}
			if visited[num].Products[pnum].NumberOfProduct == 0 {
				fmt.Printf("Товар %s закончился\n", visited[num].Products[pnum].Name)
			} else {
				names = append(names, visited[num].Products[pnum].Name)
				fmt.Println(names)
				price += visited[num].Products[pnum].Price
				visited[num].Products[pnum].NumberOfProduct--
			}
		}
	}
	if names != nil {
		n = rand.Intn(100) + 1
		order := NewOrder(names, price, n)
		(*orderMap)[n] = order
		hashChanel <- n
	}
}
