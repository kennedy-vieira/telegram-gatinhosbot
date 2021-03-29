package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"

	tb "gopkg.in/tucnak/telebot.v2"
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
