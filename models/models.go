package models

const ASC = "ASC"

const DESC = "DESC"

// This structure wraps all data related to database connection.
// See AWS Secrets Manager.
type Secret struct {
	Username            string `json:"username"`
	Password            string `json:"password"`
	Engine              string `json:"engine"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	DbClusterIdentifier string `json:"dbClusterIdentifier"`
	DbName              string `json:"dbName"`
}

// This structure wraps all user data registered in Cognito and database.
// See users table in database.
type User struct {
	Uuid      string `json:uuid"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// This structure wraps all category's data.
// See category table in database.
type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

// This structure wraps all products' data.
// See products table in database.
type Product struct {
	Id           int     `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	CreatedAt    string  `json:"createdAt"`
	UpdatedAt    string  `json:"updatedAt"`
	Price        float64 `json:"price,omitempty"`
	Stock        int     `json:"stock"`
	Path         string  `json:"path"`
	CategoryId   int     `json:"categoryId"`
	CategoryPath string  `json:"categoryPath,omitempty"`
	Search       string  `json:"search,omitempty"`
}

type Pageable struct {
	TotalElements int       `json:"totalElements"`
	Content       []Product `json:"content"`
}

type PageableUsers struct {
	TotalElements int    `json:"totalElements"`
	Content       []User `json:"content"`
}

type Address struct {
	Id         int    `json:"id"`
	Title      string `json:"title"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postalCode"`
	Phone      string `json:"phone"`
	UserName   string `json:"userName"`
}

type Order struct {
	Id        int            `json:"id"`
	UserUuid  string         `json:"userUuid"`
	AddressId int            `json:"addressId"`
	Date      string         `json:"date"`
	Total     float64        `json:"total"`
	Details   []OrderDetails `json:"details"`
}

type OrderDetails struct {
	Id        int     `json:"id"`
	OrderId   int     `json:"orderId"`
	ProductId int     `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
