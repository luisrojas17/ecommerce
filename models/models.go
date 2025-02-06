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
	Id   int    `json:"categID"`
	Name string `json:"categName"`
	Path string `json:"categPath"`
}
