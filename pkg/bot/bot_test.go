package bot

import (
	"movieSpider/pkg/config"
	"movieSpider/pkg/model"
	"testing"
)

func TestNewTgBot(t *testing.T) {
	config.InitConfig("/home/ycd/Data/Daddylab/source_code/src/go-source/tools-cmd/movieSpider/bin/movieSpider/config.yaml")
	model.NewMovieDB()
	bot := NewTgBot(config.TG.BotToken, config.TG.TgIDs)

	bot.StartBot()

}
