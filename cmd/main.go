package main

import (
	"log"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/ghazlabs/idn-remote-scheduler/internal/core"
	wa "github.com/ghazlabs/idn-remote-scheduler/internal/driven/publisher"
	"github.com/ghazlabs/idn-remote-scheduler/internal/driven/scheduler"
	"github.com/ghazlabs/idn-remote-scheduler/internal/driver"
	"github.com/go-co-op/gocron/v2"
	"github.com/go-resty/resty/v2"
)

func main() {
	var cfg config
	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	waPublisher, err := wa.NewWaPublisher(wa.WaPublisherConfig{
		HttpClient:   resty.New(),
		Username:     cfg.WAPublisherUsername,
		Password:     cfg.WAPublisherPassword,
		WaApiBaseUrl: cfg.WAPublisherApiBaseUrl,
	})
	if err != nil {
		log.Fatalf("failed to create wa publisher: %v", err)
	}

	gocronClient, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to create gocron client: %v", err)
	}

	_, err = scheduler.NewGoCronScheduler(scheduler.GoCronSchedulerConfig{
		Client:    gocronClient,
		Publisher: waPublisher,
		Storage:   nil,
	})
	if err != nil {
		log.Fatalf("failed to create gocron scheduler: %v", err)
	}

	service, err := core.NewService(core.ServiceConfig{})
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	api, err := driver.NewAPI(driver.APIConfig{
		Service:        service,
		ClientUsername: cfg.ClientUsername,
		ClientPassword: cfg.ClientPassword,
	})
	if err != nil {
		log.Fatalf("failed to create api: %v", err)
	}

	// initialize server
	s := &http.Server{
		Addr:        cfg.Port,
		Handler:     api.GetHandler(),
		ReadTimeout: time.Second * 30,
	}
	// run server
	log.Printf("server is listening on %v", cfg.Port)
	err = s.ListenAndServe()
	if err != nil {
		log.Fatalf("unable to run server due: %v", err)
	}
}

type config struct {
	Port string `env:"PORT,required" envDefault:"9865"`

	ClientUsername string `env:"CLIENT_USERNAME,required" envDefault:"admin"`
	ClientPassword string `env:"CLIENT_PASSWORD,required" envDefault:"admin"`
	DefaultNumbers string `env:"DEFAULT_NUMBERS,required" envDefault:"120363026176938692@g.us"`

	WAPublisherApiBaseUrl string `env:"WA_PUBLISHER_API_BASE_URL,required" envDefault:"http://localhost:8080"`
	WAPublisherUsername   string `env:"WA_PUBLISHER_USERNAME,required" envDefault:"admin"`
	WAPublisherPassword   string `env:"WA_PUBLISHER_PASSWORD,required" envDefault:"admin"`
}
