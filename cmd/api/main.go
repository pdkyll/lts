package main

import (
	"context"
	"strings"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jackc/pgx"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	
	"gitlab.com/m0ta/lts/app/config"
	"gitlab.com/m0ta/lts/app/controller"
	"gitlab.com/m0ta/lts/app/router"
	"gitlab.com/m0ta/lts/app/service"
	"gitlab.com/m0ta/lts/app/store"
	"gitlab.com/m0ta/lts/app/utils"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	// config
	cfg := config.Get()

	// logger
	
	// Init repository store (with postgresql inside)
	store, err := store.New(ctx)
	if err != nil {
		return utils.ErrorWrap(err, "store.New failed")
	}

	//init DB
	// TODO do migration
	if cfg.InitDB {
		var a string
		fmt.Print("== IMPORTANT == Start database initialization?(y)")
		fmt.Fscanln(os.Stdin, &a)
		if a == "y" {
			err := initDB()
			if err != nil {
				return utils.ErrorWrap(err, "InitDB failed")
			}
		} else {
			return utils.ErrorNew("Database initialization canceled")
		}
		fmt.Println("Database initialization completed")
		return nil
	}

	// Init service manager
	serviceManager, err := service.NewManager(ctx, store)
	if err != nil {
		return utils.ErrorWrap(err, "manager.New failed")
	}

	// Init controllers
	cUser 		:= controller.NewUsers(ctx, serviceManager)
	
	// Initialize Fiber instance
	app := fiber.New()
	app.Use(cors.New())

	api := router.SetupRoutes(app, cUser)
	router.SetupRoutesForUser(api, cUser)

	s := app.Stack()
	for _, v := range s {
		for _, w := range v {
			if (w.Method == "GET") || (w.Method == "POST") || (w.Method == "DELETE") || (w.Method == "PATCH") {
				serviceManager.Logger.Info.Println(w.Method, w.Path)
			}
		}
	}

	// start api server
	log.Fatal(app.Listen(cfg.HTTPAddr))

	return nil
}

//initDB ...something
func initDB () error {
	fmt.Println("== Database initialization started..")
	cfg := config.Get()
	if cfg.PgURL == "" {
		return utils.ErrorNew("No URL to connect Postgre")
	}
	
	pgxConfig, err := pgx.ParseConnectionString(cfg.PgURL)
	if err != nil {
		return err
	}

	conn, err := pgx.Connect(pgxConfig)
	if err != nil {
		return err
	}
	
	_, err = conn.Exec("SELECT 1")
	if err != nil {
		return err
	}

	defer conn.Close()

	// reading sql file to Exec
	dirData := "./data/init"
	lst, err := ioutil.ReadDir(dirData)
	if err != nil {
		return err
	}
	for _, val := range lst {
		if !val.IsDir() {
			if (strings.HasSuffix(val.Name(), ".sql")) {
				data, err := ioutil.ReadFile(fmt.Sprintf("%v/%v", dirData, val.Name()))
				if err != nil {
					return err
				}

				strSQL := string(data)

				_, err = conn.Exec(strSQL)
				if err != nil {
					return err
				}
				fmt.Println("-- Initialized DB from sql:", val.Name())
			}
		}
	}
	
	fmt.Println("== Database initialization completed.. successfully!")
	return nil
}