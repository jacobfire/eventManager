package main

import (
	"calendar/internal/app/apiserver"
	"calendar/pkg/api"
	"calendar/pkg/operations"
	"flag"
	"google.golang.org/grpc"
	"log"
	"net"
)

const MIGRATION_FILE_EVENTS = "event"

var (
	configPath string
	needMigration bool
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.toml", "config path")
	flag.BoolVar(&needMigration, "migrate", false, "Start migrations if we need")
}

func main() {
	flag.Parse()

	go func() {
		server := apiserver.New()

		if needMigration {
			err := server.CreateMigrationFiles(MIGRATION_FILE_EVENTS)
			if err != nil {
				log.Fatal(err)

				return
			}
			err = server.Migrate()
			if err != nil {
				log.Println(err)

				return
			}

			return
		}

		err := server.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	s := grpc.NewServer()
	srv := &operations.GRPCServer{}
	api.RegisterCalendarServer(s, srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	err = s.Serve(l)
	if err != nil {
		log.Fatal(err)
	}
}
