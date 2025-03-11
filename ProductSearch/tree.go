package productsearch

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

//Структура товара

type Product struct {
	Name            string
	Price           float64
	NumberOfProduct int
}

//Структура раздела каталога

type Chapter struct {
	Name        string
	Products    []Product
	NextCatalog []*Chapter
	Next        *Chapter
	Previous    *Chapter
}

//Структура каталога

type Catalog struct {
	HeadCatalogs *Chapter
}

// Создание каталога
func MakeCatalog() Catalog {
	return Catalog{HeadCatalogs: &Chapter{Name: "Yandex Market (beta)"}}
}

// Добавление разделов каталога
func (c *Catalog) AddToCatalog() {
	var name, action string
	var next int
	current := c.HeadCatalogs
	for {
		fmt.Printf("Введите название каталогов. Для прекращения ввода нажмите q, для выхода из функции нажмите e :")
	Add:
		for {
			fmt.Scanln(&name)
			switch name {
			case "q":
				break Add
			case "e":
				return
			default:
				current.NextCatalog = append(current.NextCatalog, &Chapter{Name: name, Previous: current})
			}
		}
		fmt.Println("Для перехода другой каталог нажмите n, для выхода введите q")
		fmt.Scanln(&action)
		switch action {
		case "n":
			fmt.Println("Выберите следующий каталог, указав его номер")
			for i, nextCatalog := range current.NextCatalog {
				fmt.Printf("%d. %s\n", i, nextCatalog.Name)
			}
			fmt.Scanln(&next)
			current.Next = current.NextCatalog[next]
			current.Previous = current
			current = current.Next
			fmt.Printf("%s\n", current.Name)

			if !Checking(current.NextCatalog, *current.Previous) && current.Name != c.HeadCatalogs.Name {
				current.NextCatalog = append(current.NextCatalog, current.Previous)
			}
		case "q":
			return
		}
	}
}

// УДАЛЕНИЕ РАЗДЕЛОВ КАТАЛОГА

//Добавление продуктов в каталог

func (c *Catalog) AddProduct() {
	var name string
	var price float64
	var numberOfProduct int
	visited := c.DFS()
	current := c.HeadCatalogs
	fmt.Println(visited)
	for {
		fmt.Println("Введите название раздела каталога,чтобы добавить товары.Чтобы выйти из функции нажите q:")
		a := false
		fmt.Scanln(&name)
		switch name {
		case "q":
			return
		default:
		Checking:
			for _, element := range visited {
				if element.Name == name {
					current = element
					a = true
					break Checking
				}
			}
			if !a {
				fmt.Println("Каталога не сущесвтует")
				return
			}
		}

		fmt.Println("Введите название товару и укажите цену к нему. Для выхода нажмите q в поле товар")
	AddProduct:
		for {
			fmt.Print("Товар:")
			nameproduct, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			nameproduct = strings.TrimSpace(nameproduct)
			if (nameproduct) == "q" {
				break AddProduct
			}
			fmt.Print("Цена:")
			fmt.Scanln(&price)
			fmt.Print("Кол-во товара:")
			fmt.Scanln(&numberOfProduct)
			current.Products = append(current.Products, Product{Name: nameproduct, Price: price})
		}
	}
}

func OrderingProducts() {

}

// Обзор каталога
func (c Catalog) ViewCatalog() {
	visited := c.DFS()
	for _, element := range visited {
		fmt.Printf("%s\n", element.Name)
		if element.Products != nil {
			for _, product := range element.Products {
				fmt.Printf("Товар: %s\nЦена: %f\nКол-во товара %d\n", product.Name, product.Price, product.NumberOfProduct)
			}
		}
	}

}

//СДЕЛАТЬ СДЕЛАТЬ ВЫВОД ТОВАРОВ ОПРЕДЕЛЕННОГО КАТАЛОГА

// УДАЛЕНИЕ ПРОДУКТА ИЗ КАТАЛОГА

// Функция проверки на наличие элемента в массиве
func Checking(catalogs []*Chapter, chapter Chapter) bool {
	for _, arrayChapter := range catalogs {
		if arrayChapter.Name == chapter.Name {
			return true
		}
	}
	return false
}

// Функция обхода дерева в глубину
func (c Catalog) DFS() []*Chapter {
	visited := []*Chapter{} //Посещенные разделы
	stack := []*Chapter{}   //Разделы, которые требуется посетить
	current := c.HeadCatalogs
	stack = append(stack, current.NextCatalog[0])
	for len(stack) != 0 {
		for _, element := range current.NextCatalog {
			if !Checking(visited, *element) {
				stack = append(stack, element)
			}
		}
		visited = append(visited, current)
		current = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
	}
	return visited
}
