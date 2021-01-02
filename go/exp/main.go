package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "your-password"
	dbname   = "lenslocked_exp" //
)

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
	Orders []Order
}

type Order struct{
	gorm.Model
	UserId      uint
	Amount      int
	Description string
}

func main()  {
	//psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	//	"password=%s dbname=%s sslmode=disable",
	//	host, port, user, password, dbname)
	//// db, err := gorm.Open("postgres", psqlInfo)
	//us, err := models.NewUserService(psqlInfo)
	//if err != nil {
	//	panic(err)
	//}
	//defer us.Close()
	//fmt.Println("Successfully connected!")
	//us.DestructiveReset()

	//db.LogMode(true)
	//db.AutoMigrate(&User{}, &Order{})

	//name, email := getInfo()
	//u := &User{
	//	Name:  name,
	//	Email: email,
	//}
	//if err = db.Create(u).Error; err != nil {
	//	panic(err)
	//}

	//var user User
	//db.First(&user)
	//if db.Error != nil {
	//	panic(db.Error)
	//}
	//
	//createOrder(db, user, 1001,"Fake Description #1")
	//createOrder(db, user, 9999,"Fake Description #2")
	//createOrder(db, user, 8800,"Fake Description #3")

	//var user User
	//db.Preload("Orders").First(&user)
	//if db.Error != nil {
	//	panic(db.Error)
	//}
	//fmt.Print("Email:", user.Email)
	//fmt.Print("Number of orders:", len(user.Orders))
	//fmt.Print("Orders:", user.Orders)

	// with actual gorm model

	//user := models.User{
	//	Name:     "Michael Scott",
	//	Email:    "michaelscott@dindermufflin.com",
	//	Password: "bestboss",
	//}
	//err = us.Create(&user)
	//if err != nil {
	//	panic(err)
	//}

	//us := models.NewUserService(psqlInfo)
	//if err := us.

	//foundUser, err := us.

	// Update a user
	//user.Name = "Updated Name"
	//if err := us.Update(&user); err != nil {
	//	panic(err)
	//}
	//
	//foundUser, err := us.ByEmail("michael@dundermifflin.com")
	//if err != nil {
	//	panic(err)
	//}
	//// Because of an update, the name should now
	//// be "Updated Name"
	//fmt.Println(foundUser)

	// Delete a user
	//if err := us.Delete(foundUser.ID); err != nil {
	//	panic(err)
	//}
	//// Verify the user is deleted
	//_, err = us.ByID(foundUser.ID)
	//if err != models.ErrNotFound {
	//	panic("user was not deleted!")
	//}

	//fmt.Println( rand.String(10))
	//fmt.Println( rand.RememberToken())

	//hmac := hash.NewHMAC("my-secret-key")
	//// This should print out:
	////   4waUFc1cnuxoM2oUOJfpGZLGP1asj35y7teuweSFgPY=
	//fmt.Println(hmac.Hash("this is my string to hash"))


	// Verify that the user has a Remember and RememberHash
	//fmt.Printf("%+v\n", user)
	//if user.Remember == "" {
	//	panic("Invalid remember token")
	//}
	//
	//// Now verify that we can lookup a user with that remember
	//// token
	//user2, err := us.ByRemember(user.Remember)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v\n", *user2)

}

func getInfo() (name, email string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What is your name?")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	fmt.Println("What is your email?")
	email, _ = reader.ReadString('\n')
	email = strings.TrimSpace(email)
	return name, email
}

func createOrder(db *gorm.DB, user User, amount int, desc string) {
	db.Create(&Order{
		UserId:      user.ID,
		Amount:      amount,
		Description: desc,
	})
	if db.Error != nil {
		panic(db.Error)
	}
}