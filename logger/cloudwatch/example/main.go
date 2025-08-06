package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/RedHatInsights/insights-operator-utils/logger/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/google/uuid"
)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := cloudwatchlogs.NewFromConfig(cfg)
	g := cloudwatch.NewGroup("test", client)

	stream := uuid.New().String()

	w, err := g.Create(stream)
	if err != nil {
		log.Fatal(err)
	}

	r, err := g.Open(stream)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		var i int
		for {
			i++
			<-time.After(time.Second / 30)
			_, err := fmt.Fprintf(w, "Line %d\n", i)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	io.Copy(os.Stdout, r)
}
