package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/kaleb-white/letthemknow/message-server/log"
	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/sqlite"
)

// Centralized dependency injection and initialization
func initialize(ctx context.Context) (*models.Config, error) {
	c := models.Config{}
	
	contactStore := sqlite.ContactStore{}
	err := contactStore.Init(ctx)
	if err != nil {
		return &c, err
	}
	c.ContactStore = contactStore

	contactListStore := sqlite.ContactListStore{}
	err = contactListStore.Init(ctx)
	if err != nil {
		return &c, err
	}
	c.ContactListStore = contactListStore

	contactDetailStore := sqlite.ContactDetailStore{}
	c.ContactDetailStore = contactDetailStore

	return &c, nil
}

func NewServer(
    config *models.Config,
) http.Handler {
	mux := http.NewServeMux()
	AddRoutes(
			mux,
			config,
	)
	return mux
}

func run(
	ctx context.Context,  
	stdout, stderr io.Writer, 
	args []string,
	getenv func(string) string,
) error {
	sqlite.OpenDB()
	defer sqlite.CloseDB()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Register logging package with given stdout, err
	log.Configure(stdout, stderr, getenv)

	// Run initialization
	c, err := initialize(ctx)
	if err != nil {
		return err
	}
	
	// Build server	
	s := NewServer(c)

	// Build uri
	url := getenv("API_URL")
	if url == "" {
		return errors.New("Could not start server, env missing API_URL")
	}
	port := getenv("API_PORT")
	if port == "" {
		return errors.New("Could not start server, env missing API_PORT")
	}
	uri := fmt.Sprintf("%s:%s", url, port)	

	// Run server and return result
	return http.ListenAndServe(uri, s)
}


func main() {
	ctx := context.Background()
	err := run(
		ctx,
		os.Stdout,
		os.Stdin,
		os.Args,
		os.Getenv,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
