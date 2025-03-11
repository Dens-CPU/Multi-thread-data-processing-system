package main

import (
	productsearch "multi-threadeddataprocessingsysytem/ProductSearch"
)

/*
В контексте программы есть админ, который может удалять или добавлять разделы на маркете,
Есть продавцы, которые могут добавлять товары на маркет в соответвующие разделы.
Есть покупатели, которые могут либо просто посмотреть товар либо заказать его
*/
func main() {
	catalog := productsearch.MakeCatalog()

	catalog.AddToCatalog()
	catalog.AddProduct()
	catalog.ViewCatalog()

}
