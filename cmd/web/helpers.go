package main 

func (app *Config) sendEmail(msg Message){
	app.WaitGroup.Add(1)
	app.Mailer.MailerChan <- msg
}