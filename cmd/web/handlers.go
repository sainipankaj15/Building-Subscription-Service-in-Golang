package main

import (
	"building-subscritpion-service/data"
	"fmt"
	"html/template"
	"net/http"
)

func (app *Config) HomePage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "home.page.gohtml", nil)
}

func (app *Config) LoginPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.gohtml", nil)
}

func (app *Config) PostLoginPage(w http.ResponseWriter, r *http.Request) {

	_ = app.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// fetch the email and password from the request

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := app.Models.User.GetByEmail(email)

	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid Credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// if code is there it means user is in our database : Now will verify for password
	validPassword, err := user.PasswordMatches(password)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Invalid Credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !validPassword {
		msg := Message{
			To:      email,
			Subject: "Failed log in attempt",
			Data:    "Invalid Login Attempt",
		}
		app.sendEmail(msg)

		app.Session.Put(r.Context(), "error", "Invalid Credentials")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// if all okay till there : it means user log in : Saving the information in Session
	app.Session.Put(r.Context(), "userID", user.ID)
	app.Session.Put(r.Context(), "user", user)

	app.Session.Put(r.Context(), "flash", "Successfully Login")

	//Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (app *Config) Logout(w http.ResponseWriter, r *http.Request) {
	// clean up session
	// Two steps 1. Destory entire session and renew session
	_ = app.Session.Destroy(r.Context())
	_ = app.Session.RenewToken(r.Context())

	//Redirect to home page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *Config) RegisterPage(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "register.page.gohtml", nil)
}

func (app *Config) PostRegisterPage(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.ErrorLog.Println(err)
	}

	// TODO : Validate - All Data which is very necessary will do later

	//Create a user
	user := data.User{
		Email:     r.Form.Get("email"),
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Password:  r.Form.Get("password"),
		Active:    0,
		IsAdmin:   0,
	}

	_, err = user.Insert(user)

	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to create this user")

		//Redirect to home page
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// Sending an activation Email
	url := fmt.Sprintf("http://localhost/activation?email=%s", user.Email)

	// This is our url but now we will use only Signed URL to avoid hackign and tempering the URL : Signer.go
	signedURL := GenerateTokenFromString(url)
	app.InfoLog.Println(signedURL)

	msg := Message{
		To:       user.Email,
		Subject:  "Activate your account",
		Template: "confirmation-email",
		Data:     template.HTML(signedURL),
	}

	app.sendEmail(msg)
	app.Session.Put(r.Context(), "flash", "Confirmation Mail Sent from Our backend, Please check your Mail")

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}

func (app *Config) ActivateAccount(w http.ResponseWriter, r *http.Request) {

}
