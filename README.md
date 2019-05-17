# sqm
A Query Builder with Benefits

## In Development
While the public API is stable, there are still internal tasks left

## Quick example

```go
package main

import (
	"database/sql"

	_ "github.com/lib/pq"

	"fmt"

	"github.com/g-ferreira-dev/sqm"
)

type User struct {
	UUID string `db:"uuid"`

	Name  string `db:"name"`
	Email string `db:"email"`

	CreatedAt  int `db:"created_at"`
	ModifiedAt int `db:"modified_at"`
}

func connect() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s user=%s password=123456 port=5433 "+
			"dbname=%s sslmode=disable",
		"localhost", "postgres", "postgres",
	)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

var isMichael = sqm.Group(
	sqm.C("name", sqm.Equal, "Michael"),
	sqm.C("email", sqm.Like, "%michael%"),
)

func main() {
	db, err := connect()
	if err != nil {
		panic(err)
	}

	isJorge := sqm.Group(
		sqm.C("name", sqm.Equal, "Jorge"),
		sqm.C("email", sqm.Like, "%jorge%"),
	)

	isMichaelOrJorge := sqm.Group(isMichael, sqm.Or, isJorge)

	isWoman := sqm.Group(
		sqm.C("created_at", sqm.Equal, "2"),
	)

	cQuery := sqm.Using(db, "users").
		Where(isMichaelOrJorge, sqm.Or, isWoman)

	var counter int
	err = cQuery.Count(&counter)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Count: ", counter)

	var users []User
	err = cQuery.Select(&users)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Entries: ", users)

	newUser := User{
		UUID:       "ZXY987",
		Name:       "Jorge",
		Email:      "jorge@gmail.com",
		CreatedAt:  123,
		ModifiedAt: 456,
	}

	_, err = sqm.Using(db, "users").Insert(newUser)
	if err != nil {
		fmt.Println(err)
	}

	_, err = sqm.Using(db, "users").Update(newUser)
	if err != nil {
		fmt.Println(err)
	}

	_, err = cQuery.Delete()
	if err != nil {
		fmt.Println(err)
	}
}
```
