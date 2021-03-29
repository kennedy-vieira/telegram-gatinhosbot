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

	"encoding/json"
	"net/http"
)

func startdb() error {
	//db
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {

		fmt.Println("erro ao abrir a  db,func startDB")
		log.Fatal(err)
	}

	//pinga a db e verifica se a conexão foi feita
	err = db.Ping()
	if err != nil {

		fmt.Println("erro ao pingar db, func startdb")
		log.Fatal(err)
	}

	defer db.Close()
	_, err = db.Exec(`CREATE TABLE users(
		user_id SERIAL PRIMARY KEY,
		t_username TEXT UNIQUE,
		first_name TEXT,
		chat_id BIGINT UNIQUE NOT NULL,
		aww_update_opt BOOL,
		random_image_opt BOOL 
		);`)
	if err != nil {

		fmt.Println("erro ao iniciar a  db,func startDB")
		return err
	}

	return nil
}
func insertUserDB(u *tb.User, chatid int64) error {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {

		fmt.Println("erro ao inserir usuario na  db, func insertUserDB")
		return err
	}

	//pinga a db e verifica se a conexão foi feita
	err = db.Ping()
	if err != nil {

		fmt.Println("erro ao pingar db, func insertUSerDB")
		return err
	}
	defer db.Close()
	sqlStatement := `
	INSERT INTO users(user_id,t_username,first_name,chat_id,aww_update_opt,random_image_opt)
	VALUES($1,$2,$3,$4,$5,$6)`

	_, err = db.Exec(sqlStatement, u.ID, u.Username, u.FirstName, chatid, false, false)
	if err != nil {

		fmt.Println("erro ao executar comando na db.func insertUserDb")
		return err
	}
	return nil
}
func findInDB(xxx int64) (bool, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {

		fmt.Println("erro ao abrir  db, func findINdb")
		return false, err
	}

	//pinga a db e verifica se a conexão foi feita
	err = db.Ping()
	if err != nil {

		fmt.Println("erro ao pingar db,func findIndDB")
		return false, err
	}
	defer db.Close()
	query := `SELECT user_id
	FROM users
	WHERE
	chat_id = $1
	`
	var result int
	row := db.QueryRow(query, xxx)

	switch err := row.Scan(&result); err {
	case sql.ErrNoRows:
		fmt.Println("erro: no rows returned, func find in db")
		return false, err

	case nil:
		fmt.Println("sem erro: rows returned result =, func find in db")
		fmt.Println(result)
		return true, nil
	default:
		fmt.Println("erro: func find in db")
		fmt.Println(err)
		return false, err
	}
}

func changeAwwOpt(chat_id int64) (string, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("erro ao abrir db")
		return "", err
	}
	defer db.Close()
	//pinga a db e verifica se a conexão foi feita
	err = db.Ping()
	if err != nil {
		fmt.Println("erro ao pringar db")
		return "", err
	}

	querytrue := "UPDATE users SET aww_update_opt = $1 where chat_id = $2"
	queryfalse := "UPDATE users SET aww_update_opt = $1 where chat_id = $2"

	querySearch := `
	SELECT aww_update_opt
	FROM users
	WHERE chat_id = $1;
	`
	var a interface{}
	err = db.QueryRow(querySearch, chat_id).Scan(&a)
	if err != nil {
		fmt.Println("erro : erro ao buscar aww_update_opt ,func changeawwopt")
		return "", err
	} else {

		switch a.(type) {
		case bool:

			if a.(bool) {
				_, err := db.Exec(queryfalse, !a.(bool), chat_id)
				if err != nil {
					fmt.Println("erro: erro ao exec queryfalse, func changeawwopt")
					fmt.Println(err)
					return "", err
				} else {
					return "desativado", nil
				}
			} else {
				_, err := db.Exec(querytrue, !a.(bool), chat_id)
				if err != nil {
					fmt.Println(err)
					fmt.Println("erro: erro ao exec querytrue, func changeawwopt")
					return "", err
				} else {

					return "ativado", nil
				}
			}
		default:
			{
				fmt.Println("erro: não é bool caiu em default,func changeawwopt")
				return "", fmt.Errorf("erro: tipo inesperado,deveria ser bool")
			}

		}
	}
}
func changeRandomImageOpt(chat_id int64) (string, error) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("erro ao abrir db")
		return "", err
	}
	defer db.Close()
	//pinga a db e verifica se a conexão foi feita
	err = db.Ping()
	if err != nil {
		fmt.Println("erro ao pringar db,func random_image_opt")
		return "", err
	}

	querytrue := "UPDATE users SET random_image_opt = $1 where chat_id = $2"
	queryfalse := "UPDATE users SET random_image_opt = $1 where chat_id = $2"

	querySearch := `
	SELECT random_image_opt
	FROM users
	WHERE chat_id = $1;
	`
	var a interface{}
	err = db.QueryRow(querySearch, chat_id).Scan(&a)
	if err != nil {
		fmt.Println("erro : erro ao buscar random_image_opt ,func changerandom_image_opt")
		return "", err
	} else {

		switch a.(type) {
		case bool:

			if a.(bool) {
				_, err := db.Exec(queryfalse, !a.(bool), chat_id)
				if err != nil {
					fmt.Println("erro: erro ao exec queryfalse, func changerandom_image_opt")
					fmt.Println(err)
					return "", err
				} else {
					return "desativado", nil
				}
			} else {
				_, err := db.Exec(querytrue, !a.(bool), chat_id)
				if err != nil {
					fmt.Println(err)
					fmt.Println("erro: erro ao exec querytrue, func changerandom_image_opt")
					return "", err
				} else {

					return "ativado", nil
				}
			}
		default:
			{
				fmt.Println("erro: não é bool caiu em default,func changerandom_image_opt")
				return "", fmt.Errorf("erro: tipo inesperado,deveria ser bool,func random_image_opt")
			}

		}
	}
}

//struct para armazenar informações retornadas pelo catApi
type CatResponse []struct {
	Url string `json:"url"`
}

//conjunto de structs para armazenar informações retornadas pela api do reddit
type redditResponse struct {
	Kind string      `json:"kind"`
	Data dataReponse `json:"data"`
}
type dataReponse struct {
	ModHash  string           `json:"modhash"`
	Dist     int              `json:"dist"`
	Children []childrenReddit `json:"children"`
	After    string           `json:"after"`
	Before   string           `json:"before"`
}
type childrenReddit struct {
	Kind string             `json:"kind"`
	Data dataChildrenReddit `json:"data"`
}
type dataChildrenReddit struct {
	Title string `json:"title"`
	Url   string `json:"url"`
	Type  string `json:"post_hint"`
}

func get_catApi_pic_Url() string {
	url := "https://api.thecatapi.com/v1/images/search"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var urlcat CatResponse

	if err := json.NewDecoder(resp.Body).Decode(&urlcat); err != nil {
		//log.Fatal("ooopsss! an error occurred, please try again")
		log.Fatal(err)
	}
	return urlcat[0].Url

}
func get_raww_top3_pics_url() (string, string, string) {
	url := "https://www.reddit.com/r/aww/hot/.json?limit=25"
	cliente := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("User-Agent", "tbot.ld")

	resp, err := cliente.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	var rr redditResponse

	if err := json.NewDecoder(resp.Body).Decode(&rr); err != nil {
		//log.Fatal("ooopsss! an error occurred, please try again")
		log.Fatal(err)
	}
	var str1, str2, str3 string = "", "", ""
	for i, j := 0, 0; i < len(rr.Data.Children) && j < 3; i++ {
		if rr.Data.Children[i].Data.Type == "image" {
			switch j {
			case 0:
				j++
				str1 = rr.Data.Children[i].Data.Url
			case 1:
				j++
				str2 = rr.Data.Children[i].Data.Url
			case 2:
				j++
				str3 = rr.Data.Children[i].Data.Url
			}
		}
	}
	return str1, str2, str3
}
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
