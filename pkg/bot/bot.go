package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"movieSpider/pkg/config"
	"movieSpider/pkg/download"
	"movieSpider/pkg/httpClient"
	"movieSpider/pkg/log"
	"movieSpider/pkg/types"
	"os"
	"strconv"
	"strings"
)

var (
	pageNum    *int
	Resolution types.Resolution
)

const (
	//MovieListCMD    = "/movie_list"
	MoveDownloadCMD = "/movie_download"
	//DownloadCMD     = "/download"
	//StartMovieCMD   = "/star_movie"
	//StarListCMD     = "/star_list"
)

type tgBot struct {
	botToken string
	IDs      []int
	bot      *tgbotapi.BotAPI
}

func NewTgBot(BotToken string, TgIDs []int) *tgBot {
	client := httpClient.NewHttpClient()
	bot, err := tgbotapi.NewBotAPIWithClient(config.TG.BotToken, "https://api.telegram.org/bot%s/%s", client)
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
	return &tgBot{
		BotToken, TgIDs, bot,
	}
}

func (t *tgBot) StartBot() {
	//var err error
	//
	//var (
	//	isMovie bool
	//	//isStar  bool
	//)

	//bot.Debug = true

	log.Infof("Authorized on account %s", t.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {

		index := 1
		pageNum = &index
		if update.Message != nil {

			switch {
			// movie_list 指令
			//case update.Message.Text == MovieListCMD:
			//	isMovie = true
			//	isStar = false
			//	dataStr := GetMovesData(1)
			//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, dataStr)
			//	msg.ParseMode = "HTML"
			//	msg.DisableWebPagePreview = true
			//	msg.ReplyMarkup = getMovieInlineKeyboardMarkup()
			//	if _, err = t.bot.Send(msg); err != nil {
			//		log.Error(err)
			//		continue
			//	}
			// movie_download 指令
			case strings.Contains(update.Message.Text, MoveDownloadCMD):
				// 如果参数长度不够直接continue 防止地址越界
				pars, ok := t.checkPars(update.Message.Text, update.Message.Chat.ID, update, MoveDownloadCMD)
				if !ok {
					continue
				}

				downloader := download.NewDownloader(config.Downloader.Scheduling)
				downloadMsg := downloader.DownloadByName(pars[1], pars[2])
				fmt.Println(downloadMsg)

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, downloadMsg)
				msg.ReplyToMessageID = update.Message.MessageID
				if _, err := t.bot.Send(msg); err != nil {
					log.Error(err)
				}

				//for _, v := range config.Aria2cList {
				//	fmt.Println(pars[1])
				//	fmt.Println(v)
				//	//if v.Label == strings.Split(update.Message.Text, " ")[1] {
				//	//	ok := t.checkUser(update.Message.Chat.ID, update)
				//	//	if !ok {
				//	//		continue
				//	//	}
				//	//	movieID, err := getMovieID(update.Message.Text)
				//	//	if err != nil {
				//	//		log.Error(err)
				//	//		continue
				//	//	}
				//	//
				//	//	movie, err := model.GetMovieByID(movieID)
				//	//	if err != nil {
				//	//		log.Error(err)
				//	//		continue
				//	//	}
				//	//
				//	//	git, err := aria2.DownloadByUrl(movie.Magnet, v.Label)
				//	//	if err != nil {
				//	//		log.Error(err)
				//	//		continue
				//	//	}
				//	//	log.Infof("%s start download.", movie.Name)
				//	//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("git_id: %s", git))
				//	//	msg.ReplyToMessageID = update.Message.MessageID
				//	//	if _, err := t.bot.Send(msg); err != nil {
				//	//		log.Error(err)
				//	//		continue
				//	//	}
				//	//}
				//}

				// download 指令
				//case strings.Contains(update.Message.Text, DownloadCMD):
				//	// 如果参数长度不够直接continue 防止地址越界
				//	if !t.checkPars(update.Message.Text, update.Message.Chat.ID, update) {
				//		continue
				//	}
				//	for _, v := range config.Aria2cList {
				//		if v.Label == strings.Split(update.Message.Text, " ")[1] {
				//			ok := t.checkUser(update.Message.Chat.ID, update)
				//			if !ok {
				//				continue
				//			}
				//			url := strings.Split(update.Message.Text, " ")[2]
				//			git, err := aria2.DownloadByUrl(url, v.Label)
				//			if err != nil {
				//				log.Error(err)
				//				continue
				//			}
				//			log.Infof("url %#s start download.", url)
				//			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("git_id: %s", git))
				//			msg.ReplyToMessageID = update.Message.MessageID
				//			if _, err := t.bot.Send(msg); err != nil {
				//				log.Error(err)
				//				continue
				//			}
				//		}
				//	}
				//case strings.Contains(update.Message.Text, StartMovieCMD):
				//	name := strings.Split(update.Message.Text, " ")[1]
				//	if name == "" {
				//		continue
				//	}
				//	model.SaveStar(name)
				//case strings.Contains(update.Message.Text, StarListCMD):
				//	isStar = true
				//	isMovie = false
				//	dataStr := GetStarsData(1)
				//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, dataStr)
				//	msg.ParseMode = "HTML"
				//	msg.DisableWebPagePreview = true
				//	msg.ReplyMarkup = getMovieInlineKeyboardMarkup()
				//	if _, err = t.bot.Send(msg); err != nil {
				//		log.Error(err)
				//		continue
				//	}
			}

		}
		//else if update.CallbackQuery != nil && isMovie {
		//	err = t.MovesCallbackQuery(update)
		//	if err != nil {
		//		log.Error(err)
		//		continue
		//	}
		//}
		//else if update.CallbackQuery != nil && isStar {
		//	err = t.StarCallbackQuery(update)
		//	if err != nil {
		//		log.Error(err)
		//		continue
		//	}
		//}

	}
}

//func (t *tgBot) StarCallbackQuery(update tgbotapi.Update) error {
//	pg, err := strconv.Atoi(update.CallbackQuery.Data)
//	if err != nil {
//		return err
//	}
//
//	pageNum = &pg
//
//	dataStr := GetStarsData(*pageNum)
//	msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, dataStr)
//	msg.ParseMode = "HTML"
//	msg.DisableWebPagePreview = true
//	msg.ReplyMarkup = getMovieInlineKeyboardMarkup()
//
//	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
//	if _, err := t.bot.Request(callback); err != nil {
//		return err
//	}
//
//	if _, err := t.bot.Send(msg); err != nil {
//		return err
//	}
//	return nil
//}

//func (t *tgBot) MovesCallbackQuery(update tgbotapi.Update) error {
//	pg, err := strconv.Atoi(update.CallbackQuery.Data)
//	if err != nil {
//		return err
//	}
//
//	pageNum = &pg
//
//	dataStr := GetMovesData(*pageNum)
//	msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, dataStr)
//	msg.ParseMode = "HTML"
//	msg.DisableWebPagePreview = true
//	msg.ReplyMarkup = getMovieInlineKeyboardMarkup()
//
//	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
//	if _, err := t.bot.Request(callback); err != nil {
//		return err
//
//	}
//
//	if _, err := t.bot.Send(msg); err != nil {
//		return err
//
//	}
//	return nil
//}

func getMovieID(str string) (int, error) {
	sile := strings.Split(str, " ")
	if len(sile) < 2 {
		return 0, errors.New("getMovieID id is 0")
	} else {
		movieID, err := strconv.Atoi(sile[2])
		if err != nil {
			return 0, err
		}
		return movieID, nil
	}

}

func getMovieInlineKeyboardMarkup() *tgbotapi.InlineKeyboardMarkup {
	if *pageNum <= 1 {
		if *pageNum == 1 {
			Markup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("上一页", "0"),
					tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("下一页(第%d页)", *pageNum+1), fmt.Sprintf("%d", *pageNum+1)),
				),
			)
			return &Markup
		}
		if *pageNum == 0 {
			Markup := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("上一页", "1"),
					tgbotapi.NewInlineKeyboardButtonData("下一页(第2页)", "2"),
				),
			)
			return &Markup
		}
	} else if *pageNum > 1 {
		Markup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("上一页(第%d页)", *pageNum-1), fmt.Sprintf("%d", *pageNum-1)),
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("下一页(第%d页)", *pageNum+1), fmt.Sprintf("%d", *pageNum+1)),
			),
		)
		return &Markup
	}

	return nil
}

//func GetMovesData(pangIndex int) string {
//	movies, err := model.GetMovies(20, pangIndex)
//	if err != nil {
//		log.Panic(err)
//	}
//	var dataStr string
//	for _, v := range movies {
//		dataStr += fmt.Sprintf("%-8d | %s,%s,%s,<a href='%s' > %s </a>\n", v.ID, v.Year, v.Area, v.Type, v.URL, v.Name)
//	}
//	return dataStr
//}

func inArray(val int, array []int) (ok bool, i int) {
	for i = range array {
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}

func (t *tgBot) checkUser(ChatID int64, update tgbotapi.Update) bool {
	ok, _ := inArray(int(ChatID), config.TG.TgIDs)
	if !ok {
		msg := tgbotapi.NewMessage(ChatID, "您没有权限")
		msg.ReplyToMessageID = update.Message.MessageID
		if _, err := t.bot.Send(msg); err != nil {
			log.Error(err)
			return false
		}
		return false
	}
	return ok
}

func (t *tgBot) checkPars(pars string, ChatID int64, update tgbotapi.Update, cmd string) ([]string, bool) {
	log.Infof("Msg: %s", update.Message.Text)
	cmdAndargs := removeSpaceItem(strings.Split(pars, " "))
	switch cmd {
	case MoveDownloadCMD:
		if len(strings.Split(pars, " ")) < 2 {
			msg := tgbotapi.NewMessage(ChatID, "参数长度不够")
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := t.bot.Send(msg); err != nil {
				log.Error(err)
				return []string{}, false
			}
			log.Warnf("参数长度不够")
		}
		return cmdAndargs, true
	default:
		return cmdAndargs, false
	}

}

func removeSpaceItem(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

//func GetStarsData(pageIndex int) string {
//	movies, err := model.GetStarLimit(20, pageIndex)
//	if err != nil {
//		log.Panic(err)
//	}
//	var dataStr string
//	for _, v := range movies {
//		var isDownload string
//		if v.IsDownload == 1 {
//			isDownload = "已下载"
//		} else {
//			isDownload = "未下载"
//		}
//		dataStr += fmt.Sprintf("%-8d |%-s,<a> %s </a>\n", v.ID, isDownload, v.Name)
//	}
//	return dataStr
//}
