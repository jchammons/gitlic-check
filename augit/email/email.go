package email

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/solarwinds/gitlic-check/augit/models"
)

var (
	emailHost = "smtp.gmail.com"
	emailPort = "587"
	emailFrom = os.Getenv("EMAIL_USERNAME")
	emailPass = os.Getenv("EMAIL_PASSWORD")
)

func SendOwnerListEmail(toEmail, serviceAccount string, owners []*models.GithubOwner) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", emailHost, emailPort),
		smtp.PlainAuth("", emailFrom, emailPass, emailHost),
		emailFrom, []string{toEmail}, getEmailBody(emailFrom, toEmail, serviceAccount, owners))
}

func getEmailBody(sender, recipient, serviceAccount string, owners []*models.GithubOwner) []byte {
	from := fmt.Sprintf("From: \"SolarWinds.io - GitHub Management\" <%s>\n", sender)
	to := fmt.Sprintf("To: %s\n", recipient)
	subj := fmt.Sprintf("Subject: Next Steps: You registered %s as a service account\n", serviceAccount)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n"
	y, m, d := time.Now().Date()

	var body string
	buf := new(bytes.Buffer)
	parseData := struct {
		ServiceAccount string
		Date           string
		Link           string
		Owners         []*models.GithubOwner
	}{
		serviceAccount,
		fmt.Sprint(m, " ", d, ", ", y),
		fmt.Sprint(os.Getenv("ROOT_URL"), "/blog/github-data-audit"),
		owners,
	}
	t, err := template.ParseFiles("templates/sa_confirmation.html")
	if err != nil {
		log.Printf("Failed to find notification template file: %v\n", err)
		// body = getPlainText()
	} else {
		err = t.Execute(buf, parseData)
		if err != nil {
			log.Printf("Failed to parse HTML notification template: %v\n", err)
			// body = getPlainText(txn, parseData.Link)
		} else {
			body = buf.String()
		}
	}
	msg := fmt.Sprint(from, to, subj, mime, "\n", body)
	return []byte(msg)
}
