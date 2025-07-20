package schedulers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"pinmarker/services"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type AuditScheduler struct {
	TrackService services.TrackService
}

func NewAuditScheduler(
	trackService services.TrackService,
) *AuditScheduler {
	return &AuditScheduler{
		TrackService: trackService,
	}
}

func (s *AuditScheduler) SchedulerAuditAppsUserTotal() {
	// Open the JSON
	file, err := os.Open("configs/admin_telegram.json")
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// Decode JSON
	var admins []Admin
	if err := json.NewDecoder(file).Decode(&admins); err != nil {
		log.Fatalf("failed to decode json: %v", err)
	}

	// Service : Get All Error Audit
	res, err := s.TrackService.GetAppsUserTotal()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Send to Telegram
	if len(admins) > 0 && len(res) > 0 {
		for _, dt := range admins {
			bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
			if err != nil {
				log.Println("Failed to connect to Telegram bot")
				return
			}

			telegramID, err := strconv.ParseInt(dt.TelegramUserID, 10, 64)
			if err != nil {
				log.Println("Invalid Telegram User Id")
				return
			}

			var summary string
			for _, stats := range res {
				summary += fmt.Sprintf("- %s (%d Users)", stats.AppName, stats.Total)
			}

			msgText := fmt.Sprintf("[ADMIN] Hello %s, the system just checked the apps summary. Here's the result :\n%s", dt.Username, summary)
			msg := tgbotapi.NewMessage(telegramID, msgText)

			_, err = bot.Send(msg)
			if err != nil {
				log.Println("Failed to send message to Telegram")
				return
			}
		}
	}
}
