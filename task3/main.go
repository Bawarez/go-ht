package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"io/ioutil"
	"os"
)

type User struct {
	Id string `json:"id"`
	Email string `json:"email"`
	Age int `json:"age"`
}

type Arguments map[string]string

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}


func Perform(args Arguments, writer io.Writer) error {
	operation := args["operation"]
	fileName := args["fileName"]

	if operation == "" {
		return errors.New("-operation flag has to be specified")
	}
	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	var err error

	switch operation {
	case "add":
		item := args["item"]
		if item == "" {
			return errors.New("-item flag has to be specified")
		}
		err = add(item, fileName)
	case "list":
		var users []User
		users, err = list(fileName)
		bytes, _ := json.Marshal(users)
		writer.Write(bytes)
	case "remove":
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		err = remove(id, fileName)
	case "findById":
		var user User
		id := args["id"]
		if id == "" {
			return errors.New("-id flag has to be specified")
		}
		user, err = findById(id, fileName)

		var bytes []byte
		if user.Id != "" {
			bytes, _ = json.Marshal(user)
		} else {
			bytes = []byte("")
		}
		writer.Write(bytes)
	default:
		return errors.New("Operation " + operation + " not allowed!")
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