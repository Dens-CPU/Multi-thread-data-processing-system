package main

import (
	"fmt"
	applicationforprocessing "multi-threadeddataprocessingsysytem/ApplicationForProcessing"
	productsearch "multi-threadeddataprocessingsysytem/ProductSearch"
	"sync"
)

/*
В контексте программы есть админ, который может удалять или добавлять разделы на маркете,
Есть продавцы, которые могут добавлять товары на маркет в соответвующие разделы.
Есть покупатели, которые могут либо просто посмотреть товар либо заказать его
*/

// Структура покупателя

type Buyer struct {
	Mail          string
	Passwd        string
	NumberOfOrder int
}

func main() {
	var wg sync.WaitGroup
	var user string
	catalog := productsearch.MakeCatalog()
	//Карта email
	emailMap := make(map[string]bool)

	//Карта пользователей платформы
	UserMap := make(map[string]Buyer)

	//Карта заказов
	OrderMap := make(map[int]applicationforprocessing.Order)

	//Канал для хеширования номера заказа
	HashChanel := make(chan int, 5)

	for {
		// СОЗДАНИЕ ФУНКЦИОНЛА АДМИНА
		fmt.Print("Выберите пользователя:\n1.Admin\n2.Salesman\n3.Buyer\nДля выхода из программы нажмите q:")
		fmt.Scanln(&user)
		switch user {
		case "1":
			var password string
			var action int
			fmt.Print("Введите пароль:")
			fmt.Scanln(&password)
			if password == "12345" {
			AdminActions:
				for {
					fmt.Println("1.Создание раздела.\n2.Удаление раздела.\n3.Обзор каталога\n4.Обзор раздела.\n5.Выход")
					fmt.Scanln(&action)
					switch action {
					case 1:
						wg.Add(1)
						go func() {
							catalog.AddToCatalog(&wg)
						}()
						wg.Wait()

					case 2:
						wg.Add(1)
						go func() {
							catalog.RemoveChapter(&wg)
						}()
						wg.Wait()
					case 3:
						wg.Add(1)
						go func() {
							catalog.ViewCatalog(&wg)
						}()
						wg.Wait()
					case 4:
						wg.Add(1)
						go func() {
							catalog.PrintProducts(&wg)
						}()
						wg.Wait()
					case 5:
						break AdminActions
					}
				}
			} else {
				fmt.Println("Неверный пароль")
				continue
			}
		case "2":

			var action int
		SalesmanAction:
			for {
				fmt.Println("1.Обзор каталога\n2.Обзор раздела.\n3.Добавление товара.\n4.Выход")
				fmt.Scanln(&action)
				switch action {
				case 1:
					wg.Add(1)
					go func() {
						catalog.ViewCatalog(&wg)
					}()
					wg.Wait()
				case 2:
					wg.Add(1)
					go func() {
						catalog.PrintProducts(&wg)
					}()
					wg.Wait()
				case 3:
					wg.Add(1)
					go func() {
						catalog.AddProduct(&wg)
					}()
					wg.Wait()
				case 4:
					break SalesmanAction
				}
			}
		case "q":
			return
		case "3":
			var action int
			var mail, password string
			fmt.Println("Добро полжаловать")
		BuyerAction:
			for {
				fmt.Println("1.Регистрация\n2.Вход")
				fmt.Scanln(&action)
				switch action {
				case 1:
					for {
						fmt.Print("Введите e-mail:")
						fmt.Scanln(&mail)
						if emailMap[mail] {
							fmt.Println("Пользователь с данным email уже существет")
						} else {
							fmt.Print("Введите пароль:")
							fmt.Scanln(&password)
							emailMap[mail] = true
							UserMap[mail] = Buyer{Mail: mail, Passwd: password}

							fmt.Println("Вы успешно зарегистрированы")
							// fmt.Println(emailMap)
							// fmt.Println(UserMap)
							break
						}
					}
				case 2:

					for {
						fmt.Print("Введите e-mail:")
						fmt.Scanln(&mail)
						fmt.Print("Введите пароль:")
						fmt.Scanln(&password)
						if mail != UserMap[mail].Mail || password != UserMap[mail].Passwd {
							fmt.Println("Введен неверно email или пароль")
							break
						} else {
							fmt.Println("Вход выполнен успешно")
							for {
								fmt.Println("1.Сделать заказ\n2.Отследить заказ\n3.Отменить заказ\n4.Выход")
								fmt.Scanln(&action)
								switch action {
								case 1:
									wg.Add(1)
									go func() {
										applicationforprocessing.MakeOrder(&OrderMap, catalog, &wg, HashChanel)
									}()
									UserMap[mail] = Buyer{Mail: mail, Passwd: password, NumberOfOrder: Hash(HashChanel, &wg)}
									fmt.Println(UserMap)
									wg.Wait()
									fmt.Println(OrderMap)
								case 2:
									go func() {

									}()
								case 3:
									go func() {

									}()
								case 4:
									break BuyerAction
								}
							}
						}
					}
				}

			}
		default:
			continue

		}
		// СОЗДАНИЕ ФУНКЦИОНЛА ПОКУПАТЕЛЯ
	}
}

// Функция хеширования номера заказа
func Hash(hashChanel chan int, wg *sync.WaitGroup) int {
	var n int
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case n = <-hashChanel:
				n = n*8 + 7
				return
			default:
				continue
			}
		}
	}()
	wg.Wait()
	return n
}
