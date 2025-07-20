package schedulers

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"pinmarker/utils"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Admin struct {
	TelegramUserID string `json:"telegram_user_id"`
	Username       string `json:"username"`
}

type HouseKeepingScheduler struct {
}

func NewHouseKeepingScheduler() *HouseKeepingScheduler {
	return &HouseKeepingScheduler{}
}

func (s *HouseKeepingScheduler) SchedulerMonthlyLog() {
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

	// Helpers : Clean Logs
	logPath, err := utils.GetLastMonthLogFilePath()
	if err != nil {
		log.Println("Log file not found:", err)
		return
	}

	// Open the log file
	fileBytes, err := os.Open(logPath)
	if err != nil {
		log.Println("Failed to open log file:", err)
		return
	}
	defer fileBytes.Close()

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

			file, err := os.Open(logPath)
			if err != nil {
				log.Println("Failed to open log file:", err)
				return
			}
			defer file.Close()

			fileInfo, err := file.Stat()
			if err != nil {
				log.Println("Failed to stat log file:", err)
				return
			}

			fileReader := tgbotapi.FileReader{
				Name:   fileInfo.Name(),
				Reader: file,
				Size:   fileInfo.Size(),
			}

			doc := tgbotapi.NewDocumentUpload(telegramID, fileReader)
			doc.ParseMode = "html"
			doc.Caption = fmt.Sprintf("[ADMIN] Hello %s, here is housekeeping log for %s %d",
				dt.Username, time.Now().AddDate(0, -1, 0).Format("January"), time.Now().AddDate(0, -1, 0).Year())

			_, err = bot.Send(doc)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}

		if err := utils.DeleteFileByPath(logPath); err != nil {
			log.Println("Failed to delete log file:", err)
		}
	}
}
