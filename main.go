package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var gBot *tgbotapi.BotAPI
var gToken string
var gChatID int64

var gUsersInChat Users
var gUsefulActivities = Activities{
	// Саморазвитие
	{"yoga", "Yoga (15 minutes)", 1},
	{"meditation", "Meditation (15 minutes)", 1},
	{"language", "Learning a foreign language (15 minutes)", 1},
	{"swimming", "Swimming (15 minutes)", 1},
	{"walk", "Walk (15 minutes)", 1},
	{"chores", "Chores", 1},

	// Work
	{"work_learning", "Studying work materials (15 minutes)", 1},
	{"portfolio_work", "Working on a portfolio project (15 minutes)", 1},
	{"resume_edit", "Resume editing (15 minutes)", 1},

	// Creativity
	{"creative", "Creative creation (15 minutes)", 1},
	{"reading", "Reading fiction literature (15 minutes)", 1},
}

var gRewards = Activities{
	// Entertainment
	{"watch_series", "Watching a series (1 episode)", 10},
	{"watch_movie", "Watching a movie (1 item)", 30},
	{"social_nets", "Browsing social networks (30 minutes)", 10},

	// Food
	{"eat_sweets", "300 kcal of sweets", 60},
}

type User struct {
	id    int64
	name  string
	coins uint16
}
type Users []*User

type Activity struct {
	code, name string
	coins      uint16
}
type Activities []*Activity

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

func isCallbackQuery(update *tgbotapi.Update) bool {
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
func sendStringMessage(msg string) {
	gBot.Send(tgbotapi.NewMessage(gChatID, msg))
}

func askToPrintIntro() {
	msg := tgbotapi.NewMessage(gChatID, "Во вступительных сообщениях ты можешь найти смысл данного бота, и правила игры. Что думаешь?")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_PRINT_INTRO, BUTTON_CODE_PRINT_INTRO),
		getKeyboardRow(BUTTON_TEXT_SKIP_INTRO, BUTTON_CODE_SKIP_INTRO),
	)
	gBot.Send(msg)
}
func showMenu() {
	msg := tgbotapi.NewMessage(gChatID, "Выбери один из вариантов:")
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		getKeyboardRow(BUTTON_TEXT_BALANCE, BUTTON_CODE_BALANCE),
		getKeyboardRow(BUTTON_CODE_USEFUL_ACTIVITIES, BUTTON_CODE_USEFUL_ACTIVITIES),
		getKeyboardRow(BUTTON_TEXT_REWARDS, BUTTON_CODE_REWARDS),
	)
	gBot.Send(msg)
}
func showBalance(user *User) {
	msg := fmt.Sprintf("%s, твой кошелёк пока пуст %s \nЗатрекай полезное действие, чтобы получить монеты", user.name, EMOJI_DONT_KNOW)
	if coins := user.coins; coins > 0 {
		msg = fmt.Sprintf("%s, you have %d %s", user.name, coins, EMOJI_COIN)
	}
	gBot.Send(tgbotapi.NewMessage(gChatID, msg))
	//sendStringMessage(msg)
	showMenu()

}

func callbackQueryIsMissing(update *tgbotapi.Update) bool {
	return update.CallbackQuery == nil || update.CallbackQuery.From == nil
}

func getUserFromUpdate(update *tgbotapi.Update) (user *User, found bool) {
	if callbackQueryIsMissing(update) {
		return
	}
	userID := update.CallbackQuery.From.ID
	for _, userInChat := range gUsersInChat {
		if userID == userInChat.id {
			return userInChat, true
		}
	}
	return
}

func storeUserFromUpdate(update *tgbotapi.Update) (user *User, found bool) {
	if callbackQueryIsMissing(update) {
		return
	}
	from := update.CallbackQuery.From
	user = &User{id: from.ID, name: strings.TrimSpace(from.FirstName + " " + from.LastName), coins: 0}
	gUsersInChat = append(gUsersInChat, user)
	return user, true
}
func showActivities(activities Activities, message string, isUseful bool) {
	activitiesButtonsRows := make([]([]tgbotapi.InlineKeyboardButton), 0, len(gUsefulActivities)+1)
	for _, activity := range activities {
		activityDescription := ""
		if isUseful {
			activityDescription = fmt.Sprintf("+ %d %s: %s", activity.coins, EMOJI_COIN, activity.name)
		} else {
			activityDescription = fmt.Sprintf("- %d %s", activity.coins, EMOJI_COIN, activity.name)
		}
		activitiesButtonsRows = append(activitiesButtonsRows, getKeyboardRow(activityDescription, activity.code))
	}
	activitiesButtonsRows = append(activitiesButtonsRows, getKeyboardRow(BUTTON_TEXT_PRINT_MENU, BUTTON_CODE_PRINT_MENU))

	msg := tgbotapi.NewMessage(gChatID, message)
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(activitiesButtonsRows...)
	gBot.Send(msg)
}
func showUsefulActivities() {
	showActivities(gUsefulActivities, "Трекай полезное действие или возвращайся в главное меню:", false)
}

func showRewards() {
	showActivities(gRewards, "Купи вознаграждение или возвращайся в главное меню:", false)
}
func findActivity(activities Activities, choiceCode string) (activity *Activity, found bool) {
	for _, activity := range activities {
		if choiceCode == activity.code {
			return activity, true
		}
	}
	return
}

func processUsefulActivity(activity *Activity, user *User) {
	errorMsg := ""
	if activity.coins == 0 {
		errorMsg = fmt.Sprintf(`у активности "%s" не указана стоимость`, activity.name)
	} else if user.coins+activity.coins > MAX_USER_COINS {
		errorMsg = fmt.Sprintf("у тебя не может быть больше %d %s", MAX_USER_COINS, EMOJI_COIN)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, прости, но %s %s Твой баланс остался без изменений.", user.name, errorMsg, EMOJI_SAD)
	} else {
		user.coins += activity.coins
	}
	gBot.Send(tgbotapi.NewMessage(gChatID, resultMessage))
}

func processReward(activity *Activity, user *User) {
	errorMsg := ""
	if activity.coins == 0 {
		errorMsg = fmt.Sprintf(`у вознаграждения "%s" не указана стоимость`, activity.name)
	} else if user.coins < activity.coins {
		errorMsg = fmt.Sprintf(`у тебя сейчас  %d %s. ты на можешь себе позволить "%s" за %d %s`, user.coins, EMOJI_COIN, activity.name, activity.coins, EMOJI_COIN)
	}

	resultMessage := ""
	if errorMsg != "" {
		resultMessage = fmt.Sprintf("%s, извиняюсь, но %s %s Твой баланс остался без изменений, вознаграждение не доступно %s", user.name, errorMsg, EMOJI_SAD, EMOJI_DONT_KNOW)
	} else {
		user.coins -= activity.coins
		resultMessage = fmt.Sprintf(`%s, вознаграждение "%s" оплачено, приступай! %d %s было сеяир с ивренр счёта. Теперь у тебя %d %s`, user.name, activity.name, activity.coins, EMOJI_COIN, user.coins, EMOJI_COIN)
	}
	gBot.Send(tgbotapi.NewMessage(gChatID, resultMessage))
}

func updateProcessing(update *tgbotapi.Update) {
	user, found := getUserFromUpdate(update)
	if !found {
		if user, found = storeUserFromUpdate(update); !found {
			gBot.Send(tgbotapi.NewMessage(gChatID, "Не получается идентифицировать пользователя"))
			return
		}
	}

	choiceCode := update.CallbackQuery.Data
	log.Println("[%T] %s", time.Now(), choiceCode)
	switch choiceCode {
	case BUTTON_CODE_BALANCE:
		showBalance(user)
	case BUTTON_CODE_USEFUL_ACTIVITIES:
		showUsefulActivities()
	case BUTTON_CODE_REWARDS:
		printIntro(update)
		showRewards()
	case BUTTON_CODE_SKIP_INTRO:
		showMenu()
	case BUTTON_CODE_PRINT_MENU:
		showMenu()
	default:
		if usefulActivity, found := findActivity(gUsefulActivities, choiceCode); found {
			processUsefulActivity(usefulActivity, user)
			delay(2)

			showUsefulActivities()
			return
		}
		if reward, found := findActivity(gRewards, choiceCode); found {
			processReward(reward, user)
			delay(2)
			showRewards()
			return
		}
		log.Printf(`[%T] !!!!!!!!! ERROR: Unknown code "%s"`, time.Now(), choiceCode)
		msg := fmt.Sprintf("%s, I'm sorry, I don't recognize code '%s' %s Please report this error to my creator.", user.name, choiceCode, EMOJI_SAD)
		gBot.Send(tgbotapi.NewMessage(gChatID, msg))
	}
}

func main() {
	log.Printf("Authorized on account %s", gBot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = UPDATE_CONFIG_TIMEOUT

	for update := range gBot.GetUpdatesChan(updateConfig) {
		if isCallbackQuery(&update) {
			updateProcessing(&update)
		} else if isStartMessage(&update) {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			gChatID = update.Message.Chat.ID
			askToPrintIntro()
		}
	}
}
