package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	gocontact "bitbucket.org/nicluo/gocontact/src"

	"github.com/gorilla/schema"
	"github.com/urfave/cli"
)

type Form struct {
	Name               string `schema:"name"`
	Email              string `schema:"email"`
	Subject            string `schema:"subject"`
	Message            string `schema:"message"`
	GRecaptchaResponse string `schema:"g-recaptcha-response"`
}

type SubmitResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

// App Name and Version
const (
	AppName = "GoContact"
	AppVer  = "0.0.1"
)

const (
	page = `<!DOCTYPE HTML><html><body>
	<script src='https://www.google.com/recaptcha/api.js'></script>
	<form method="post" action="/">
	<div>
	<label>Your Name (required)</label>
	<input type="text" name="name" required />
	</div>
	<div>
	<label>Your Email (required)</label>
	<input type="email" name="email" required />
	</div>
	<div>
	<label>Subject</label>
	<input type="text" name="subject" required />
	</div>
	<div>
	<label>Message</label>
	<textarea name="message" cols="40" rows="10"></textarea>
	</div>
	<div class="g-recaptcha" data-sitekey="%s"></div>
	<input type="submit" value="Send">
	</form>
	</body></html>`
)

var decoder = schema.NewDecoder()
var siteKey string

func demo(writer http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(writer, page, siteKey)
}

func submit(writer http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("error parsing form: %v\n", err)
		response, _ := json.Marshal(&SubmitResponse{Success: false, Error: "invalid-form"})
		fmt.Fprint(writer, string(response))
		return
	}

	var form Form
	err = decoder.Decode(&form, r.PostForm)
	if err != nil {
		log.Printf("error decoding message: %v\n", err)
		response, _ := json.Marshal(&SubmitResponse{Success: false, Error: "invalid-decoded-form"})
		fmt.Fprint(writer, string(response))
		return
	}

	ip, _ := gocontact.RequestIP(r)
	result := gocontact.VerifyRecaptcha(ip, form.GRecaptchaResponse)
	log.Printf("reCAPTCHA: %v\n", result)

	if !result {
		response, _ := json.Marshal(&SubmitResponse{Success: false, Error: "recaptcha-failed"})
		fmt.Fprint(writer, string(response))
		return
	}

	err = gocontact.SendContactMail(form.Subject, form.Name, form.Email, form.Message)
	if err != nil {
		log.Printf("error sending mail: %v\n", err)
		response, _ := json.Marshal(&SubmitResponse{Success: false, Error: "smtp-failed"})
		fmt.Fprint(writer, string(response))
	} else {
		log.Println("smtp send successful")
		response, _ := json.Marshal(&SubmitResponse{Success: true})
		fmt.Fprint(writer, string(response))
	}
}

func run(ctx *cli.Context) error {
	gocontact.InitRecaptcha(ctx.String("private-key"))
	gocontact.InitMail(ctx.String("smtp-sender"), ctx.String("smtp-password"), ctx.String("smtp-host"), ctx.Int("smtp-port"), ctx.String("smtp-to"))

	http.HandleFunc("/", submit)
	if ctx.BoolT("demo") {
		http.HandleFunc("/demo", demo)
		siteKey = ctx.String("site-key")
	}

	if err := http.ListenAndServe(fmt.Sprintf(":%d", ctx.Int("port")), nil); err != nil {
		log.Fatal("failed to start server", err)
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Usage = "Go Contact Form Collector"
	app.Version = AppVer
	app.Flags = append(app.Flags, []cli.Flag{
		cli.BoolTFlag{
			Name:  "demo",
			Usage: "show demo form for testing",
		},
		cli.IntFlag{
			Name:   "port",
			Usage:  "server port",
			Value:  3000,
			EnvVar: "PORT",
		},
		cli.StringFlag{
			Name:   "site-key",
			Usage:  "site key for reCaptcha (only required for demo)",
			EnvVar: "SITE_KEY",
		},
		cli.StringFlag{
			Name:   "private-key",
			Usage:  "private key for reCaptcha",
			EnvVar: "PRIVATE_KEY",
		},
		cli.StringFlag{
			Name:   "smtp-sender",
			Usage:  "sender email",
			EnvVar: "SMTP_SENDER",
		},
		cli.StringFlag{
			Name:   "smtp-to",
			Usage:  "contact recipient email",
			EnvVar: "SMTP_TO",
		},
		cli.StringFlag{
			Name:   "smtp-password",
			Usage:  "sender password",
			EnvVar: "SMTP_PASSWORD",
		},
		cli.StringFlag{
			Name:   "smtp-host",
			Usage:  "smtp host",
			EnvVar: "SMTP_HOST",
		},
		cli.IntFlag{
			Name:   "smtp-port",
			Usage:  "smtp port",
			Value:  587,
			EnvVar: "SMTP_PORT",
		},
	}...)
	app.Action = run
	app.Run(os.Args)
}
