package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Message struct {
	Name string
	Body string
	Time int64
}

func main() {
	if len(os.Args) < 2 || os.Args[1] == "help" {
		help_txt, err := os.ReadFile("./help.txt")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(help_txt))
		return
	}
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	db := setup_db()
	defer db.Close()

	create_cmd := flag.NewFlagSet("create_cmd", flag.ExitOnError)
	playlist_name := create_cmd.String("create", "", "create")
	delete_flag := create_cmd.String("delete", "", "delete")

	// TODO: add the Charm library for all CLI functions

	switch os.Args[1] {
	case "add":
		if len(os.Args) != 3 {
			log.Fatal("no channel tag provided")
		}
		tag := os.Args[2]
		_, exists := find_row(db, tag, "./sql/read_row.sql")
		if exists {
			log.Println("Channel is already tracked ;)")
			os.Exit(0)
		}
		key := os.Getenv("API_KEY")
		item, err := get_channel_ID(tag, key)

		if err != nil {
			log.Fatal(err)
		}
		id, title, real_tag := item[0], item[1], item[2]
		createChannelRow(db, id, real_tag, title)
	case "remove":
		if len(os.Args) != 3 {
			log.Fatal("no channel tag provided")
		}
		tag := os.Args[2]
		tag, exists := find_row(db, tag, "./sql/read_row.sql")
		if !exists {
			log.Println("no channel found")
			os.Exit(0)
		}
		deleteRow(db, tag)
	case "cli":

	case "playlist":
		if len(os.Args) == 2 {
			playlists := read_playlists(db)
			fmt.Println(playlists)
			return
		}
		api_key, access_token := os.Getenv("API_KEY"), os.Getenv("ACCESS_TOKEN")
		create_cmd.Parse(os.Args[2:])

		if len(*delete_flag) != 0 {
			err := delete_playlist(db, *delete_flag)
			if err != nil {
				os.Exit(1)
			}
			os.Exit(0)
		}
		if len(*playlist_name) == 0 {
			log.Fatal("playlist name not provided")
		}
		query := get_user_input("Enter search terms: ")
		filter := get_user_input("Filter: ")

		q := csv_string(query)
		f := csv_string(filter)

		playlist_resp := create_playlist(db, *playlist_name, q, f, api_key, access_token)
		populate_playlist(db, query, filter, playlist_resp.Id)

	case "create_table":
		createTable(db, "./sql/create_playlist_table.sql")
	case "delete_table":
		deleteTable(db, "./sql/delete_playlist_table.sql")
	case "refresh":
		refresh_quota(db)
	case "quota":
		quota := read_quota(db)
		fmt.Println(time.Now().Unix() - quota.timestamp.Unix())
	case "read":
		channels := readChannels(db)
		fmt.Println(channels)
	case "insert":
		insert_row(db)
	default:
		log.Fatal("Invalid subcommand. To see usable commands, use 'cli help'")
	}
}
