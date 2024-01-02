package values

type ProductCategory = string

const (
	ProductCategoryGoods    ProductCategory = "GOODS"
	ProductCategoryServices ProductCategory = "SERVICES"
	ProductCategoryExpenses ProductCategory = "EXPENSES"
)

var validProductCategories = map[string]ProductCategory{
	ProductCategoryGoods:    ProductCategoryGoods,
	ProductCategoryServices: ProductCategoryServices,
	ProductCategoryExpenses: ProductCategoryExpenses,
}
