package main

import (
	"building-subscritpion-service/data"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/phpdave11/gofpdf"
	"github.com/phpdave11/gofpdf/contrib/gofpdi"
)

var pathToManual = "./pdf"
var tmpPath = "./tmp"

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
	validPassword, err := app.Models.User.PasswordMatches(password)
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

	_, err = app.Models.User.Insert(user)

	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to create this user")

		//Redirect to home page
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	// Sending an activation Email
	url := fmt.Sprintf("http://localhost:8080/activate?email=%s", user.Email)

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

	// validate URL
	url := r.RequestURI

	app.InfoLog.Println(url)

	testurl := fmt.Sprintf("http://localhost:8080%s", url)

	app.InfoLog.Println(testurl)

	okay := VerifyToken(testurl)

	if !okay {
		app.Session.Put(r.Context(), "error", "Invalid token")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Till now , URL is validated : Will activate the account
	user, err := app.Models.User.GetByEmail(r.URL.Query().Get("email"))

	if err != nil {
		app.Session.Put(r.Context(), "error", "No user found")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// Became user active
	user.Active = 1

	// Updating the user in database as well
	err = app.Models.User.Update(*user)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to update the user")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.Session.Put(r.Context(), "flash", "Account activated, You can login now")
	http.Redirect(w, r, "/login ", http.StatusSeeOther)
}

func (app *Config) ChooseSubscription(w http.ResponseWriter, r *http.Request) {

	// We have already wrote a middleware for Auth so no need to worry
	// This page will be avalible only for Logged in users

	plans, err := app.Models.Plan.GetAll()
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to find the plans")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		app.ErrorLog.Println(err)
		return
	}

	dataMap := make(map[string]any)
	dataMap["plans"] = plans

	app.render(w, r, "plans.page.gohtml", &TemplateData{
		Data: dataMap,
	})
}

func (app *Config) SubscribeToPlan(w http.ResponseWriter, r *http.Request) {

	// get the id of the plan which user chose
	id := r.URL.Query().Get("id")

	// Covert string to int
	planID, _ := strconv.Atoi(id)

	// Get the plans from the DB
	plan, err := app.Models.Plan.GetOne(planID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Unable to find the plan")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		app.ErrorLog.Println(err)
		return
	}

	// get the user from the session
	user, ok := app.Session.Get(r.Context(), "user").(data.User)
	if !ok {
		app.Session.Put(r.Context(), "error", "Login in First.")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		app.ErrorLog.Println(err)
		return
	}

	// Will generate a invoice and send to the user via mail
	app.WaitGroup.Add(1)
	go func() {
		defer app.WaitGroup.Done()

		invoice, err := app.GetInvoice(user, plan)
		if err != nil {
			app.ErrorChan <- err
		}

		msg := Message{
			To:       user.Email,
			Subject:  "Your Invoice Mail",
			Data:     invoice,
			Template: "invoice",
		}

		app.sendEmail(msg)
	}()

	// Will generate a manual
	app.WaitGroup.Add(1)
	go func() {
		defer app.WaitGroup.Done()

		pdf := app.GenerateManual(user, plan)
		err := pdf.OutputFileAndClose(fmt.Sprintf("%s/%d_manual.pdf", tmpPath, user.ID))
		if err != nil {
			app.ErrorChan <- err
			return
		}

		msg := Message{
			To:       user.Email,
			Subject:  "Your Manual",
			Data:     "Your user manual is atached",
			Template: "invoice",
			AttachmentMap: map[string]string{
				"Manual.pdf": fmt.Sprintf("%s/%d_manual.pdf", tmpPath, user.ID),
			},
		}

		app.sendEmail(msg)
	}()

	// Subscribe the user to a new plan
	err = app.Models.Plan.SubscribeUserToPlan(user, *plan)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Error while subscribing the plans")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		app.ErrorLog.Println(err)
	}

	u, err := app.Models.User.GetOne(user.ID)
	if err != nil {
		app.Session.Put(r.Context(), "error", "Error while getting user from the database")
		http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
		app.ErrorLog.Println(err)
	}

	app.Session.Put(r.Context(), "user", u)

	//redirect now
	app.Session.Put(r.Context(), "flash", "Subscribed")
	http.Redirect(w, r, "/members/plans", http.StatusSeeOther)
}

func (app *Config) GetInvoice(user data.User, plan *data.Plan) (string, error) {
	app.InfoLog.Println("Ammount is", plan.PlanAmountFormatted)
	return plan.PlanAmountFormatted, nil
}

func (app *Config) GenerateManual(user data.User, plan *data.Plan) *gofpdf.Fpdf {

	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetMargins(10, 13, 10)

	importer := gofpdi.NewImporter()

	time.Sleep(6 * time.Second)

	t := importer.ImportPage(pdf, fmt.Sprintf("%s/manual.pdf", pathToManual), 1, "/MediaBox")
	pdf.AddPage()

	importer.UseImportedTemplate(pdf, t, 0, 0, 215.9, 0)

	pdf.SetX(75)
	pdf.SetY(150)
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s %s", user.FirstName, user.LastName), "", "C", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 4, fmt.Sprintf("%s User Guide", plan.PlanName), "", "C", false)

	return pdf
}
