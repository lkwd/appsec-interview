// this file isn't really part of the interview exercise - it just provides a simple stub for
// user database. Things like having a hardcoded admin password are therefore not in scope

package users

import "sync"

type User struct {
	Name string
	Password []byte
}

// UserDatabase - simulate a datastore
type UserDatabase interface {
	Get(name string) *User
}

type userDatabase struct {
	lock sync.RWMutex
	users []User
}

func (us *userDatabase) Get(name string) *User {
	us.lock.RLock()
	defer us.lock.RUnlock()
	for _, u := range us.users {
		if u.Name == name {
			return &u
		}
	}
	return nil
}

var db = userDatabase{
	lock:  sync.RWMutex{},
	users: []User{
		// password below is bcrypt hash of "up6WvcAsb6iodPgVYG3SofRYoEw2ALYq"
		{Name: "admin", Password: []byte("$2a$13$6bkw.CNTd4ZjAm8PEaIeleni6QfjXGULYtQco9ZTo2APiv31V8h8y")},
	},
}

func Database() UserDatabase {
	return &db
}