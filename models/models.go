package models

type Secret struct {
	Username            string `json:"username"`
	Password            string `json:"password"`
	Engine              string `json:"engine"`
	Host                string `json:"host"`
	Port                int    `json:"port"`
	DbClusterIdentifier string `json:"dbClusterIdentifier"`
}

type User struct {
	Email string `json:"UserEmail"`
	Uuid  string `json:"UserUUID"`
}

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

type Product struct {
	Id           int64   `json:"id"`
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
