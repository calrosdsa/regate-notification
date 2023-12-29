package main

import (
	"context"
	"google.golang.org/api/option"
	"notification/core"

	// "net"
	// "strconv"
	"database/sql"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"fmt"

	"log"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile(`config.json`)
	// viper.SetConfigFile("/home/regate/regate-services/notification/app/config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// if viper.GetBool(`debug`) {
	// 	log.Println("Service RUN on DEBUG mode")
	// }
}

func main() {
	loc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		log.Println(loc)
	}
	host := viper.GetString(`database.host`)
	port := viper.GetInt(`database.port`)
	user := viper.GetString(`database.user`)
	password := viper.GetString(`database.pass`)
	dbname := viper.GetString(`database.name`)
	time.Local = loc
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}

	defer db.Close()
	app, _, _ := SetupFirebase()
	core.Init(db, app)
}


func SetupFirebase() (*firebase.App, context.Context, *messaging.Client) {

	ctx := context.Background()

	// serviceAccountKeyFilePath, err := filepath.Abs(fmt.Sprintf("%s/serviceAccountKey.json",viper.GetString("path")))
	serviceAccountKeyFilePath, err := filepath.Abs("./serviceAccountKey.json")
	if err != nil {
		log.Println("Unable to load serviceAccountKeys.json file")
	}

	opt := option.WithCredentialsFile(serviceAccountKeyFilePath)

	//Firebase admin SDK initialization
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Println("Firebase load error")
	}

	//Messaging client
	client, _ := app.Messaging(ctx)

	return app, ctx, client
}
