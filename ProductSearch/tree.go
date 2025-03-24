package productsearch

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
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

var mutex = sync.RWMutex{}

// Добавление разделов каталога
func (c *Catalog) AddToCatalog(wg *sync.WaitGroup) {
	defer wg.Done()
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
				mutex.Lock()
				current.NextCatalog = append(current.NextCatalog, &Chapter{Name: name, Previous: current})
				mutex.Unlock()
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
				mutex.Lock()
				current.NextCatalog = append(current.NextCatalog, current.Previous)
				mutex.Unlock()
			}
		case "q":
			return
		}
	}
}

//Добавление продуктов в каталог

func (c *Catalog) AddProduct(wg *sync.WaitGroup) {
	var addMutex = sync.Mutex{}
	defer wg.Done()
	var name string
	var price float64
	var numberOfProduct int
	visited := c.DFS()
	current := c.HeadCatalogs
	fmt.Println(visited)
MainLoop:
	for {
		fmt.Println("Введите название раздела каталога,чтобы добавить товары.Чтобы выйти из функции нажите q:")
		a := false
		fmt.Scanln(&name)
		switch name {
		case "q":
			break MainLoop
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
			_, err1 := fmt.Scanln(&price)
			if err1 != nil {
				fmt.Println("Неверный формат ввода данных. Товар не будет добавлен в раздел")
				break AddProduct
			}
			fmt.Print("Кол-во товара:")
			_, err2 := fmt.Scanln(&numberOfProduct)
			if err2 != nil {
				fmt.Println("Неверный формат ввода данных. Товар не будет добавлен в раздел")
				break AddProduct
			}
			addMutex.Lock()
			current.Products = append(current.Products, Product{Name: nameproduct, Price: price, NumberOfProduct: numberOfProduct})
			addMutex.Unlock()
		}
	}
}

func OrderingProducts() {

}

// Обзор каталога
func (c Catalog) ViewCatalog(wg *sync.WaitGroup) {
	defer wg.Done()
	visited := c.DFS()
	if len(visited) == 0 {
		fmt.Println("В каталоге Yandex Market (beta) нет разделов")
		fmt.Println()
		return
	}
	mutex.RLock()
	for _, element := range visited {
		fmt.Printf("%s\n", element.Name)
		if element.Products != nil {
			fmt.Println("____________________")
			for _, product := range element.Products {
				fmt.Printf("Товар: %s\nЦена: %2.f\nКол-во товара %d\n", product.Name, product.Price, product.NumberOfProduct)
				fmt.Println("------------------")
			}
			fmt.Println("____________________")
		}
	}
	mutex.RUnlock()

}

// Вывод элементов определенного раздела
func (c Catalog) PrintProducts(wg *sync.WaitGroup) {
	defer wg.Done()
	var name string
	visited := c.DFS()
	if len(visited) == 0 {
		fmt.Println("В каталоге Yandex Market (beta) нет разделов")
		fmt.Println()
		return
	}
	fmt.Println("Введите название каталога")
	fmt.Scanln(&name)
	mutex.RLock()
	for _, chapter := range visited {
		if name == chapter.Name {
			if len(chapter.Products) == 0 {
				fmt.Println("В разделе еще нет товаров")
				fmt.Println()
				return
			}
			fmt.Println("____________________")
			for _, product := range chapter.Products {
				fmt.Printf("Товар:%s\nЦена:%2.f\nКол-во:%d\n", product.Name, product.Price, product.NumberOfProduct)
				fmt.Println()
			}
			fmt.Println("____________________")
			mutex.RUnlock()
			return
		}
	}

}

// Удаление каталога
func (c *Catalog) RemoveChapter(wg *sync.WaitGroup) {
	defer wg.Done()
	visited := c.DFS()
	if len(visited) == 0 {
		fmt.Println("В каталоге Yandex Market (beta) нет разделов, которые можно удалить")
		fmt.Println()
		return
	}
	var name string
	fmt.Println("Введите название раздела для удаления")
	fmt.Scanln(&name)
	for _, element := range visited {
		if element.Name == name {
			current := element.Previous
			for i := 0; i < len(current.NextCatalog); i++ {
				if current.NextCatalog[i].Name == name {
					mutex.Lock()
					current.NextCatalog[i], current.NextCatalog[len(current.NextCatalog)-1] = current.NextCatalog[len(current.NextCatalog)-1], current.NextCatalog[i]
					current.NextCatalog = current.NextCatalog[:len(current.NextCatalog)-1]
					current.Next = nil
					mutex.Unlock()
					return
				}
			}
		}

	}
	fmt.Println("Раздела не существует")
}

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
	if current.NextCatalog == nil {
		return visited
	}
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
