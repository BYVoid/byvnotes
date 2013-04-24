package user

import (
	"fmt"
	"github.com/robfig/revel"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"regexp"
	"crypto/sha256"
	"io"
	"encoding/base64"
	"byvnotes/app/models/db"
)

type User struct {
	Username string
	Password string
}

func collection() (*mgo.Collection, error) {
	mongo, err := db.DB()
	if err != nil {
		return nil, err
	}
	collection := mongo.C("users")
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

func Get(username string) (*User, error) {
	collection, err := collection()
	if err != nil {
		return nil, err
	}
	user := User{}
	err = collection.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func Count() (int, error) {
	collection, err := collection()
	if err != nil {
		return 0, err
	}
	return collection.Count()
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Username)
}

var userRegex = regexp.MustCompile("^\\w*$")

func (user *User) Validate(v *revel.Validation) {
	v.Check(user.Username,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{3},
		revel.Match{userRegex},
	).Key("username")
	v.Check(user.Password,
		revel.Required{},
		revel.MaxSize{15},
		revel.MinSize{3},
	).Key("password")
}

func (user *User) CheckPassword(password string) bool {
	password = getPasswordHash(password)
	return user.Password == password
}

func (user *User) Save() error {
	collection, err := collection()
	if err != nil {
		return err
	}
	// Generate hash checksum of password
	user.Password = getPasswordHash(user.Password)
	err = collection.Insert(user)
	if err != nil {
		return err
	}
	return nil
}

func getPasswordHash(password string) string {
	hash := sha256.New()
	io.WriteString(hash, password)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}
