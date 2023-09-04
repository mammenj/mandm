package messages

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mailjet/mailjet-apiv3-go/v4"
)

func SendMail(fromEmail, emailTo, subject, body string) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	mailjetClient := mailjet.NewMailjetClient(os.Getenv("MJ_APIKEY_PUBLIC"), os.Getenv("MJ_APIKEY_PRIVATE"))
	messagesInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: fromEmail,
				Name:  "Admin",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: emailTo,
					Name:  emailTo,
				},
			},
			Subject:  subject,
			TextPart: body,
		},
	}
	messages := mailjet.MessagesV31{Info: messagesInfo}
	res, err := mailjetClient.SendMailV31(&messages)
	if res != nil {
		log.Printf("Data: %+v\n", res)
	} else {
		log.Println("NO DATA FROM EMAIL Provider")
	}
	if err != nil {
		log.Println(err)
		return err
	}
	return nil

}
