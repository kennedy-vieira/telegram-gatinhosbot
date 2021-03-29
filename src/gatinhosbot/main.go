package main

import (
	//bot

	"fmt"
	"log"
	"os"

	tb "gopkg.in/tucnak/telebot.v2"

	//db
	"database/sql"

	_ "github.com/lib/pq"
	//http request
)

func sendScheduledPics(bt *tb.Bot) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("erro ao abrir db")
		return
	}
	defer db.Close()
	//pinga a db e verifica se a conexão foi feita
	err = db.Ping()
	if err != nil {
		fmt.Println("erro ao pringar db,func random_image_opt")
		return
	}

	query := `SELECT chat_id FROM users WHERE random_image_opt = true `

	rows, err := db.Query(query)
	if err == nil {
		r := get_catApi_pic_Url()
		var id int64
		for rows.Next() {
			rows.Scan(&id)
			aa := tb.ChatID(id)
			bt.Send(aa, r)
		}
	} else {
		fmt.Println("erro a retornar rows rand ,func sendScheduledPics")
	}

}
func sendScheduledReddit(bt *tb.Bot) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("erro ao abrir db")
		return
	}
	defer db.Close()
	//pinga a db e verifica se a conexão foi feita
	err = db.Ping()
	if err != nil {
		fmt.Println("erro ao pringar db,func random_image_opt")
		return
	}

	query := `SELECT chat_id FROM users WHERE aww_update_opt = true `

	rows, err := db.Query(query)
	if err == nil {
		a, b, c := get_raww_top3_pics_url()
		var id int64
		for rows.Next() {
			rows.Scan(&id)
			aa := tb.ChatID(id)
			bt.Send(aa, a)
			bt.Send(aa, b)
			bt.Send(aa, c)
		}
	} else {
		fmt.Println("erro a retornar rows aww,func sendScheduledReddit")
	}

}

func main() {

	fmt.Println(get_catApi_pic_Url())
	//bot
	var (
		port           = os.Getenv("PORT")
		publicURL      = os.Getenv("PUBLIC_URL") // you must add it to your config vars
		token          = os.Getenv("TOKEN")      // you must add it to your config vars
		admin_username = os.Getenv("ADMINUSERNAMETELEGRAM")
	)

	webhook := &tb.Webhook{
		Listen:   ":" + port,
		Endpoint: &tb.WebhookEndpoint{PublicURL: publicURL},
	}

	pref := tb.Settings{
		Token:  token,
		Poller: webhook,
	}

	bot, err := tb.NewBot(pref)
	if err != nil {

		fmt.Println("erro ao iniciar o bot,func main")
		log.Fatal(err)
	}

	// msg com schedule
	if len(os.Args) >= 2 && os.Args[1] == "gatinhos" {
		sendScheduledPics(bot)
		os.Exit(0)
		return
	}
	if len(os.Args) >= 2 && os.Args[1] == "reddit" {
		sendScheduledReddit(bot)
		os.Exit(0)
		return
	}
	bot.Handle("/criadb", func(m *tb.Message) {
		if m.Sender.Username == admin_username {
			err := startdb()
			if err != nil {
				bot.Send(m.Sender, "erro")
			} else {
				bot.Send(m.Sender, "sucesso")
			}
		} else {
			bot.Send(m.Sender, "Usuario não autorizado para esta ação")
		}
	})

	bot.Handle("/start", func(m *tb.Message) {
		a, err := findInDB(m.Chat.ID)
		if err != nil {
			if !a {
				err := insertUserDB(m.Sender, m.Chat.ID)
				if err != nil {
					log.Fatal(err)
				}
			}
			bot.Send(m.Sender, "Bem-vindo(a) "+m.Sender.FirstName)
		} else {
			bot.Send(m.Sender, "Erro interno, tente novamente mais tarde")
		}
	})

	bot.Handle("/miau", func(m *tb.Message) {
		aa := tb.ChatID(m.Chat.ID)
		bot.Send(aa, "meow meow")

	})

	bot.Handle("/get_gatinho", func(m *tb.Message) {
		sb := get_catApi_pic_Url()
		bot.Send(m.Sender, sb)

	})
	bot.Handle("/get_top_aww", func(m *tb.Message) {
		a, b, c := get_raww_top3_pics_url()
		bot.Send(m.Sender, a)
		bot.Send(m.Sender, b)
		bot.Send(m.Sender, c)
	})

	//opções de configurações

	inlineBtnAww := tb.InlineButton{
		Unique: "awwconf",
		Text:   "r/aww",
	}

	inlineBtnRand := tb.InlineButton{
		Unique: "riconf",
		Text:   "Gatinhos Aleatorios",
	}

	inlineBtns := [][]tb.InlineButton{
		[]tb.InlineButton{inlineBtnAww, inlineBtnRand},
	}
	bot.Handle(&inlineBtnRand, func(c *tb.Callback) {
		state, err := changeRandomImageOpt(c.Message.Chat.ID)
		fmt.Println(state)
		bot.Respond(c, &tb.CallbackResponse{
			ShowAlert: false,
		})
		if err == nil {
			bot.Send(c.Sender, "Recebimento de imagens aleatorias de gatinhos foi "+state)
		} else {
			bot.Send(c.Sender, "Aconteceu algum erro inesperado, tente usar o comando /start")
		}
	})
	bot.Handle(&inlineBtnAww, func(c *tb.Callback) {
		state, err := changeAwwOpt(c.Message.Chat.ID)
		bot.Respond(c, &tb.CallbackResponse{
			ShowAlert: false,
		})
		if err == nil {
			bot.Send(c.Sender, "Recebimento automatico do conteudo de r/aww foi "+state)
		} else {
			bot.Send(c.Sender, "Aconteceu algum erro inesperado, tente usar o comando /start")
		}
	})
	bot.Handle("/config", func(m *tb.Message) {
		bot.Send(m.Sender,
			`Aqui você pode ativar/desativar o recebimento automatico de conteudo,
			 escolha a opção que deseja mudar as configurações`,
			&tb.ReplyMarkup{InlineKeyboard: inlineBtns})
	})

	bot.Start()

}
