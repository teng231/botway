package initx

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/abdfnx/botway/constants"
	"github.com/abdfnx/tran/dfs"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

var (
	c = config.New(".")
)

func GetBotName() string {
	c.AddDriver(yaml.Driver)
	c.LoadFiles(".botway.yaml")

	return c.String("bot.name")
}

func GetBotType() string {
	c.AddDriver(yaml.Driver)
	c.LoadFiles(".botway.yaml")

	return c.String("bot.type")
}

func DockerInit() {
	err := dfs.CreateDirectory(filepath.Join(constants.HomeDir, ".botway"))

	if err != nil {
		log.Fatal(err)
	}

	viper.AddConfigPath(constants.BotwayDirPath())
	viper.SetConfigName("botway")
	viper.SetConfigType("json")

	t := GetBotType()
	bot_token := ""
	app_token := ""
	cid := ""

	if t == "discord" {
		bot_token = "DISCORD_TOKEN"
		app_token = "DISCORD_CLIENT_ID"
		cid = "bot_app_id"
	} else if t == "slack" {
		bot_token = "SLACK_TOKEN"
		app_token = "SLACK_APP_TOKEN"
		cid = "bot_app_token"
	} else if t == "telegram" {
		bot_token = "TELEGRAM_TOKEN"
	}

	env := viper.New()
	env.SetConfigType("env")

	secretsFile, serr := ioutil.ReadFile(filepath.Join("config", "botway-deploy-tokens.env"))

	if serr != nil {
		panic(serr)
	}

	env.ReadConfig(bytes.NewBuffer(secretsFile))

	viper.SetDefault("botway.bots." + GetBotName() + ".bot_token", env.Get(bot_token))
	viper.SetDefault("botway.bots_names", []string{GetBotName()})

	if t != "telegram" {
		viper.SetDefault("botway.bots." + GetBotName() + "." + cid, env.Get(app_token))
	}

	if t == "discord" {
		if constants.Gerr != nil {
			panic(constants.Gerr)
		} else {
			guilds := gjson.Get(string(constants.Guilds), "guilds.#")
			
			for x := 0; x < int(guilds.Int()); x++ {
				server := gjson.Get(string(constants.Guilds), "guilds." + fmt.Sprint(x)).String()

				sgi := strings.ToUpper(server) + "_GUILD_ID"

				viper.Set("botway.bots." + GetBotName() + ".guilds." + server + ".server_id", env.Get(sgi))
			}
		}
	}

	if err := viper.SafeWriteConfig(); err != nil {
		if os.IsNotExist(err) {
			err = viper.WriteConfig()

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatal(err)
		}
	}

	fmt.Println("🐋 Done")
}
