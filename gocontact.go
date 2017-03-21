package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/urfave/cli"
)

// App Name and Version
const (
	AppName = "GoContact"
	AppVer  = "0.0.1"
)

const (
	page = `<!DOCTYPE HTML><html><body>
	<script src='https://www.google.com/recaptcha/api.js'></script>
	<form method="post">
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
  <div class="g-recaptcha" data-sitekey="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"></div>
	<input type="submit" value="Send">
	</form>
	</body></html>`
)

func homePage(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fmt.Fprint(writer, page)
}

func run(ctx *cli.Context) error {
	http.HandleFunc("/", homePage)
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("failed to start server", err)
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = AppName
	app.Usage = "Go Contact Form Collector"
	app.Version = AppVer
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Action = run
	app.Run(os.Args)
}
