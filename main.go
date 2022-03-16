package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/m1ome/vsratoslav/drawer"
	"github.com/syndtr/goleveldb/leveldb"
	"gopkg.in/telebot.v3"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	filepath string
	token    string
	db       string
	percent  int64
)

func randomPhrase(phrases []string) string {
	return phrases[rand.Intn(len(phrases))]
}

func dbChatKey(chatId int64) []byte {
	return []byte(fmt.Sprintf("chat:%d", chatId))
}

func getPercent(db *leveldb.DB, chatId int64) int64 {
	key := dbChatKey(chatId)
	exists, err := db.Has(key, nil)
	if !exists || err != nil {
		return percent
	}

	val, err := db.Get(key, nil)
	if err != nil {
		return percent
	}

	v, err := strconv.ParseInt(string(val), 10, 64)
	if err != nil {
		db.Delete([]byte(key), nil)
		return percent
	}

	return v
}

func setPercent(db *leveldb.DB, chatId int64, percent int64) error {
	key := dbChatKey(chatId)
	return db.Put(key, []byte(fmt.Sprintf("%d", percent)), nil)
}

func main() {
	rand.Seed(time.Now().Unix())

	flag.StringVar(&filepath, "file", "phrases.json", "")
	flag.StringVar(&token, "token", "", "")
	flag.StringVar(&db, "db", "db", "")
	flag.Int64Var(&percent, "percent", 25, "percentage of action")
	flag.Parse()

	db, err := leveldb.OpenFile(db, nil)
	if err != nil {
		log.Fatalf("can't open database file")
	}

	file, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatalf("can't read file %s: %v", filepath, err)
	}

	var phrases []string
	if err := json.Unmarshal(file, &phrases); err != nil {
		log.Fatalf("can't unmarshal phrases: %v", err)
	}

	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	bot.Handle("/setrandom", func(c telebot.Context) error {
		admins, err := bot.AdminsOf(c.Message().Chat)
		if err != nil {
			return nil
		}

		isAdmin := false
		for _, admin := range admins {
			if admin.User.ID == c.Message().Sender.ID {
				isAdmin = true
				break
			}
		}

		if !isAdmin {
			return nil
		}

		percent, err := strconv.ParseInt(c.Message().Payload, 10, 64)
		if err != nil {
			return nil
		}

		if err := setPercent(db, c.Chat().ID, percent); err != nil {
			log.Printf("error settings percent: %v", err)
			return nil
		}

		log.Printf("set percent to %d for chat %d", percent, c.Chat().ID)

		return nil
	})

	bot.Handle(telebot.OnPhoto, func(ctx telebot.Context) error {
		percent := getPercent(db, ctx.Chat().ID)
		if rand.Int63n(100) > percent {
			return nil
		}

		photo := ctx.Message().Photo

		reader, err := bot.File(&photo.File)
		if err != nil {
			log.Printf("error reaading photo: %v", err)
			return nil
		}
		defer reader.Close()

		phrase := randomPhrase(phrases)
		buf, err := drawer.DrawText(reader, "./Lobster-Regular.ttf", phrase)
		if err != nil {
			log.Printf("error drawing text on image: %v", err)
			return nil
		}

		if err := ctx.Send(&telebot.Photo{
			File: telebot.FromReader(buf),
		}); err != nil {
			log.Printf("error sending image: %v", err)
			return nil
		}

		return nil
	})

	bot.Start()
}
