package models

import (
	"fmt"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"regexp"
)

type User struct {
	Username string
	Password string
}

func Collection() (*mgo.Collection, error) {
	db, err := DB()
	if err != nil {
		return nil, err
	}
	collection := db.C("users")
	index := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = collection.EnsureIndex(index)
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var userRegex = regexp.MustCompile("^\\w*$")

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Username,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{4},
		revel.Match{userRegex},
	)
	ValidatePassword(v, user.Password).Key("user.Password")
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{5},
	)
}

func (user *User) Save() error {
	collection, err := Collection()
	if err != nil {
		return err
	}
	err = collection.Insert(user)
	if err != nil {
		return err
	}
	return nil
}
