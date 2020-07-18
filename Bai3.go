package main

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/LeO0999/exercise4/database"
)

//insert 100 bản ghi vào user:
func InsertUser() {
	user := database.User{}
	for i := 0; i < 100; i++ {
		user.ID = strconv.FormatInt(int64(i+3), 10)
		user.Name = "Person" + user.ID
		err := db.InsertUser(user)
		if err != nil {
			log.Println(err)
		}
	}
}

func PrintUser(buffchan chan *database.Data, wg *sync.WaitGroup) {
	for {
		select {
		case data := <-buffchan:
			fmt.Printf("Line %v - %v - %v\n", data.Iden, data.User.ID, data.User.Name)
			wg.Done()
		}
	}
}

func GetNameOfUser() error {
	buffchan := make(chan *database.Data, 100)
	defer close(buffchan)
	var wg sync.WaitGroup

	for i := 0; i < 2; i++ {
		go PrintUser(buffchan, &wg)
	}
	err := db.ScanforRow(buffchan, &wg)
	if err != nil {
		return err
	}
	wg.Wait()
	return nil

}
