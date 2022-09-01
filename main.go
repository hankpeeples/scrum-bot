package main

import (
	"fmt"
	"os"

	"github.com/hankpeeples/scrum-bot/bot"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("DOTENV error: %v\n", err)
	}

	token := os.Getenv("TOKEN")

	bot.Start(token)
}
