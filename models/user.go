package models

type User struct {
    ID       string    `bson:"_id" json:"id"`
    Email    string `bson:"email" json:"email"`
    Password string `bson:"password" json:"password"`
    FirstName string `bson:"firstName" json:"firstName"`
    LastName string `bson:"lastName" json:"lastName"`
}
