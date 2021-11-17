package main

import (
	"calendar/pkg/api"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"strconv"
)

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("not enough arguments")
	}

	x, err := strconv.Atoi(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	//y, err := strconv.Atoi(flag.Arg(1))
	//if err != nil {
	//	log.Fatal(err)
	//}

	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	client := api.NewCalendarClient(conn)


	result, err := client.GetById(context.Background(), &api.GetByIdRequest{Id: int32(x)})

	fmt.Println(result)
}
