package main

import (
	"github.com/LeO0999/exercise4/database"
)

//  Viết hàm: sau khi tạo user thì insert user_id vào user_point với số điểm 10.
func InsertUsertoPoint(user database.User) error {
	err := db.InsertUser(user)
	if err != nil {
		return err
	}

	point := database.Point{UserID: user.ID, Points: 10}
	err = db.InsertPoint(point)
	if err != nil {
		return err
	}
	return nil
}
