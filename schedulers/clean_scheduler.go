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

type CleanScheduler struct {
	TrackService services.TrackService
}

func NewCleanScheduler(
	trackService services.TrackService,
) *CleanScheduler {
	return &CleanScheduler{
		TrackService: trackService,
	}
}

func (s *CleanScheduler) SchedulerCleanAllTracksCreatedByDays() {
	days := 30

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

	// Service : Delete All Tracks By Days Created
	deletedRow, err := s.TrackService.DeleteAllTracksByDaysCreated(days)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Send to Telegram
	if len(admins) > 0 {
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

			msgText := fmt.Sprintf("[ADMIN] Hello %s, the system just clean track history that have passed %d days with total %d item deleted", dt.Username, days, deletedRow)
			msg := tgbotapi.NewMessage(telegramID, msgText)

			_, err = bot.Send(msg)
			if err != nil {
				log.Println("Failed to send message to Telegram")
				return
			}
		}
	}
}
