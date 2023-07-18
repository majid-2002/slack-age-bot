package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
	"log"
	"os"
	"strconv"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent) {
	for event := range analyticsChannel {
		fmt.Println("Command Events:")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	bot := slacker.NewClient(os.Getenv("SLACK_BOT_TOKEN"), os.Getenv("SLACK_APP_TOKEN"))

	go printCommandEvents(bot.CommandEvents())

	//greet the user when the bot is tagged in the channel
	bot.Init(func() {
		fmt.Println("Connected!")
	})

	bot.Command("my name is <name>", &slacker.CommandDefinition{
		Description: "Set your name",
		Examples: []string{
			"my name is Jhon",
			"iam Jhon",
		},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			name := request.StringParam("name", "")
			r := fmt.Sprintf("Nice to meet you, %s!", name)
			response.Reply(r)
		},
	})

	bot.Command("my yob is <year>", &slacker.CommandDefinition{
		Description: "Calculate your age based on the year of birth",
		Examples: []string{
			"my yob is 1990",
		},
		Handler: func(botCtx slacker.BotContext, request slacker.Request, response slacker.ResponseWriter) {
			year := request.Param("year")
			yob, err := strconv.Atoi(year)
			if err != nil {
				fmt.Println("Error converting year to int")
				return
			}
			age := 2023 - yob
			r := fmt.Sprintf("Your age is %d", age)
			response.Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = bot.Listen(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
