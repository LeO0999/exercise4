package database

import (
	"errors"
	"log"
	"sync"
	"time"
)

type User struct {
	ID       string
	Name     string
	Birth    int64
	Created  int64
	UpdateAt int64
}
type Data struct {
	Iden int
	User User
}

func (db *Db) InsertUser(u User) error {
	insert, err := db.engine.Insert(&u)
	if err != nil {
		return err
	}
	if insert == 0 {
		return errors.New(" Insert fail")
	}
	return nil
}

// - b2: tạo 1 transaction khi update `birth` thành công thì cộng 10 điểm vào `point` sau đó sửa lại `name ` thành `$name + "updated "`
// 		  nếu 1 quá trình fail thì rollback, xong commit (CreateSesson)
func (db *Db) UpdateUser(user, condiUser *User) error {
	update, err := db.engine.Update(user, condiUser)
	if err != nil {
		return err
	}
	if update == 0 {
		return errors.New("Update fail")
	}
	return nil
}

func (db *Db) ListUser() ([]*User, error) {
	var list []*User
	err := db.engine.Find(&list)
	if err != nil {
		return nil, errors.New("List User fail")
	}
	return list, nil
}

func (db *Db) GetUser(id string) (*User, error) {
	user := &User{ID: id}
	Find, err := db.engine.Get(user)
	if err != nil {
		log.Println("Fail")
		return nil, err
	}
	if !Find {
		return nil, errors.New("Not Found")
	}
	return user, err
}
func (db *Db) UpdateBirthUser(id string, birth int64) error {
	session := db.engine.NewSession()
	defer session.Close()

	session.Begin()

	user := &User{ID: id}
	find, err := session.Get(user)
	if err != nil {
		session.Rollback()
		return err
	}
	if !find {
		session.Rollback()
		return errors.New("Not found user!")
	}

	user.Birth = birth
	_, err1 := session.Update(user, &User{ID: id})
	if err1 != nil {
		session.Rollback()
		return err1
	}

	point := &Point{UserID: user.ID}
	_, err2 := session.Get(point)
	if err2 != nil {
		session.Rollback()
		return err2
	}
	point.Points += 10
	_, err = session.Update(point, &Point{UserID: user.ID})
	if err != nil {
		session.Rollback()
		return err
	}

	user.Name = user.Name + "updated"
	user.UpdateAt = time.Now().UnixNano()
	_, err1 = session.Update(user, &User{ID: id})
	if err1 != nil {
		session.Rollback()
		return err1
	}

	session.Commit()
	return nil
}

func (db *Db) ScanforRow(buffchan chan *Data, wg *sync.WaitGroup) error {
	rows, err := db.engine.Rows(&User{})
	defer rows.Close()
	if err != nil {
		return err
	}
	user := new(User)
	i := 1
	for rows.Next() {
		err2 := rows.Scan(user)
		if err2 == nil {

			dataUser := &Data{Iden: i, User: *user}
			i++
			buffchan <- dataUser
			wg.Add(1)
		}

	}
	return nil

}
