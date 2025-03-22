package domain

type User struct {
	ID   interface{} `bson:"_id"`
	Name string      `bson:"name"`
	Age  int         `bson:"age"`
}

type UserCreate struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
