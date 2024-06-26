package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var gBot *tgbotapi.BotAPI
var gToken string
var gChatID int64

func init() {
	//_ = os.Setenv(TOKEN_NAME_IN_OS, "INSERT_YOUR_TOKEN")
	_ = os.Setenv(TOKEN_NAME_IN_OS, "7491991680:AAFTnpQhTW-Qs6yGzljbERyogQHmH4db4W4")
	gToken = os.Getenv(TOKEN_NAME_IN_OS)

	if gToken = os.Getenv(TOKEN_NAME_IN_OS); gToken == "" {
		panic(fmt.Errorf(`failed to load environment variable "%s"`, TOKEN_NAME_IN_OS))
	}

	var err error
	if gBot, err = tgbotapi.NewBotAPI(gToken); err != nil {
		log.Panic(err)
	}
	gBot.Debug = true
}

func isStartMessage(update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Text == "/start"
}

func isCallbackQuey(update *tgbotapi.Update) bool {
	return update.CallbackQuery != nil && update.CallbackQuery.Data != ""
}

func delay(seconds uint8) {
	time.Sleep(time.Second * time.Duration(seconds))

}

func sendMessageWithDelay(delayInSec uint8, message string) {
	gBot.Send(tgbotapi.NewMessage(gChatID, message))
	delay(delayInSec)
}

func printIntro(update *tgbotapi.Update) {
	sendMessageWithDelay(2, "Привет! "+EMOJI_SUNGLASSES)
	sendMessageWithDelay(7, "Расскажи пожалуйста немного о себе")
	sendMessageWithDelay(7, EMOJI_SMILE)
	sendMessageWithDelay(1, "Здорово! Приятно с тобой познакомиться!")
	sendMessageWithDelay(10, "Наверное, каждый хоть раз играл в игру, где нужно прокачивать персонажа, делая его сильнее, умнее или красивее. Это приятно, потому что каждое действие приносит результаты. В жизни, правда, систематические действия со временем начинают приносить заметный результат.")
	sendMessageWithDelay(14, `Перед вами две таблицы: "Полезные дела" и "Награды". В первой таблице перечислены несложные короткие занятия, за каждое из которых вы получите указанное количество монет. Во второй таблице вы увидите список занятий, которые вы можете делать только после того, как оплатите их монетами, заработанными на предыдущем этапе. Готовы начать?`)
	sendMessageWithDelay(1, EMOJI_COIN)
	sendMessageWithDelay(10, `Рассмотрим пример! Вы занимаетесь йогой полчаса, за что получаете 2 монеты. Затем вас ждут 2 часа изучения программирования, за которые вы получаете 8 монет. Теперь вы можете посмотреть 1 серию "Интернов" и остаться при своих. Все просто! Готовы составить свои таблицы?`)
	sendMessageWithDelay(6, `Чтобы не потерять свои монеты, отмечайте выполненные задания в таблице "Полезные дела". А чтобы получить награду, не забудьте её "купить", прежде чем приступать к ней. Готовы начать заполнять таблицы? Предлагаю начать с "Полезных дел".`)
}
func getKeyboardRow(buttonText, buttonCode string) []tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonCode))
}

func askToPrintIntro() {
	msg := tgbotapi.NewMessage(gChatID, "Во вступительных сообщениях ты можешь найти смысл данного бота, и правила игры. Что думаешь?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_PRINT_INTRO, BUTTON_CODE_PRINT_INTRO),
		getKeyboardRow(BUTTON_TEXT_SKIP_INTRO, BUTTON_CODE_SKIP_INTRO),
	)
	gBot.Send(msg)
}
func showMenu(update *tgbotapi.Update) {
	msg := tgbotapi.NewMessage(gChatID, "Выбери один из вариантов:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_BALANCE, BUTTON_CODE_BALANCE),
		getKeyboardRow(BUTTON_CODE_USEFUL_ACTIVITIES, BUTTON_CODE_USEFUL_ACTIVITIES),
		getKeyboardRow(BUTTON_TEXT_REWARDS, BUTTON_CODE_REWARDS),
	)
	gBot.Send(msg)
}
func updateProcessing(update *tgbotapi.Update) {
	choiceCode := update.CallbackQuery.Data
	log.Println("[%T] %s", time.Now(), choiceCode)
	switch choiceCode {
	case BUTTON_CODE_PRINT_INTRO:
		printIntro(update)
		showMenu(update)
	case BUTTON_CODE_SKIP_INTRO:
		showMenu(update)
	}
}

func main() {
	log.Printf("Authorized on account %s", gBot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = UPDATE_CONFIG_TIMEOUT

	for update := range gBot.GetUpdatesChan(updateConfig) {
		if isCallbackQuey(&update) {
			updateProcessing(&update)
		} else if isStartMessage(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			gChatID = update.Message.Chat.ID
			askToPrintIntro()
		}

	}

}
