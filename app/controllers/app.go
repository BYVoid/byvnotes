package controllers

import (
	"github.com/robfig/revel"
	"byvnotes/app/models"
)

type Application struct {
	*revel.Controller
}

func init() {

}

func (c Application) Index() revel.Result {
	var user models.User
	user.Username = "byvoid"
	user.Password = "byvoid"
	user.Save()
	return c.Render()
}
