package firebase

type Item struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type User struct {
	Name  string `firestore:"name"`
	Email string `firestore:"email"`
	Age   int    `firestore:"age"`
}