package controllers

import (
	"github.com/robfig/revel"
	"byvnotes/app/models/user"
)

type Application struct {
	*revel.Controller
}

func (ctl Application) checkUser() revel.Result {
	if _, ok := ctl.Session["user"]; !ok {
		if ctl.Action != "Application.Login" && ctl.Action != "Application.DoLogin" {
			ctl.Flash.Error("Please log in first")
			return ctl.Redirect(Application.Login)
		}
	}
	return nil
}

func init() {
	revel.InterceptMethod(Application.checkUser, revel.BEFORE)
}

func (ctl Application) Index() revel.Result {
	return ctl.Render()
}

func (ctl Application) Login() revel.Result {
	return ctl.Render()
}

func (ctl Application) DoLogin(username string, password string) revel.Result {
	count, err := user.Count()
	if err != nil {
		return ctl.RenderError(err)
	}
	var userInst *user.User = nil
	if count > 0 {
		userInst, err = user.Get(username)
		ctl.Validation.Required(userInst != nil).Message("User does not exist").Key("username")
		if userInst != nil {
			ctl.Validation.Required(userInst.CheckPassword(password)).Message("Password is wrong").Key("password")
		}
	} else {
		// If this is the first user, create it
		userInst = &user.User{}
		userInst.Username = username
		userInst.Password = password
		userInst.Validate(ctl.Validation)
	}
	if ctl.Validation.HasErrors() {
		ctl.Validation.Keep()
		ctl.FlashParams()
		return ctl.Redirect(Application.Login)
	}
	if count == 0 {
		// Save first created user
		err := userInst.Save()
		if err != nil {
			return ctl.RenderError(err)
		}
	}
	ctl.Session["user"] = userInst.Username
	ctl.Flash.Success("Welcome, " + userInst.Username)
	return ctl.Redirect(Application.Index)
}

func (ctl Application) Logout() revel.Result {
	for key := range ctl.Session {
		delete(ctl.Session, key)
	}
	return ctl.Redirect(Application.Index)
}

func (ctl Application) Settings() revel.Result {
	return ctl.Render()
}

func (ctl Application) EditSettings() revel.Result {
	// to be implemented
	ctl.Flash.Success("Settings saved")
	return ctl.Redirect(Application.Settings)
}
