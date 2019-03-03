package main

import (
	"io"
	"io/ioutil"
	"os"
	"fmt"
	"flag"
	"encoding/json"
	"errors"
)

type User struct {
	Id string `json:"id"`
	Email string `json:"email"`
	Age int `json:"age"`
}

type Arguments struct {
	id string
	operation string
	item  string
	fileName string
}


func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}


func Perform(args Arguments, writer io.Writer) error {
	if args.operation == "" {
		return errors.New("-operation flag has to be specified")
	}
	if args.fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	var err error

	switch args.operation {
	case "add":
		if args.item == "" {
			return errors.New("-item flag has to be specified")
		}
		err = add(args.item, args.fileName)
	case "list":
		var users []User
		users, err = list(args.fileName)
		for _, user := range users {
			fmt.Println(user)
		}
	case "remove":
		err = remove(args.id, args.fileName)
	case "findById":
		var user User
		user, err = findById(args.id, args.fileName)
		if user.Id != "" {
			_, err = writer.Write([]byte(fmt.Sprint(user)))
		} else {
			_, err = writer.Write([]byte(""))
		}
	default:
		return errors.New("Operation " + args.operation + " not allowed!")
	}

	return err
}


func parseArgs() Arguments {
	var id string
	var operation string
	var item string
	var fileName string

	flag.StringVar(&id, "id", "", "user id")
	flag.StringVar(&operation, "operation", "", "operation")
	flag.StringVar(&item, "item", "", "item")
	flag.StringVar(&fileName, "fileName", "", "file name")
	flag.Parse()

	return Arguments {
		id: 		id,
		operation: 	operation,
		item:      	item,
		fileName:  	fileName,
	}
}


func list(fileName string) ([]User, error) {
	var users []User

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return users, err
	}
	bytes, err :=  ioutil.ReadAll(file)
	if err != nil {
		return users, err
	}

	if len(bytes) > 0 {
		err = json.Unmarshal(bytes, &users)
		if err != nil {
			return users, err
		}
	}
	return users, nil
}

func findById(id string, fileName string) (User, error) {
	var user User

	users, err := list(fileName)
	if err != nil {
		return user, err
	}
	return findUser(id, users)
}

func add(item string, fileName string) error {
	var user User
	err := json.Unmarshal([]byte(item), &user)
	if err != nil {
		return err
	}

	users, err := list(fileName)
	if err != nil {
		return err
	}

	existingUser, err := findUser(user.Id, users)
	if err != nil {
		return nil
	}
	if user == existingUser {
		return errors.New("Item with id " +  string(user.Id) + " already exists")
	}
	users = append(users, user)

	return saveUsers(users, fileName)
}

func remove(id string, fileName string) error {
	users, err := list(fileName);
	if err != nil {
		return err
	}

	user, err := findUser(id, users)
	if err != nil {
		return err
	}

	if user.Id != "" {
		var newUsers []User
		for _, user := range users {
			if user.Id != id {
				newUsers = append(newUsers, user)
			}
		}

		return saveUsers(newUsers, fileName)
	}
	return errors.New("Item with id " + id + " not found")
}

func findUser(id string, users []User) (User, error) {
	var user User

	for _, usr := range users {
		if usr.Id == id {
			user = usr
			break
		}
	}
	return user, nil
}

func saveUsers(users []User, fileName string) error {
	b, _ := json.Marshal(users)
	err := ioutil.WriteFile(fileName, b, 0644)
	return err
}
