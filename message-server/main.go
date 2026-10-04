package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/kaleb-white/letthemknow/message-server/handlers"
	"github.com/kaleb-white/letthemknow/message-server/log"
	"github.com/kaleb-white/letthemknow/message-server/models"
	"github.com/kaleb-white/letthemknow/message-server/sqlite"
)

// Centralized dependency injection
func initialize(ctx context.Context) (*models.Config, error) {
	c := models.Config{}
	
	contactStore := sqlite.ContactStore{}
	err := contactStore.Init(ctx)
	if err != nil {
		return &c, nil
	}
	c.ContactStore = &contactStore

	contactListStore := sqlite.ContactListStore{}
	err = contactListStore.Init(ctx)
	if err != nil {
		return &c, nil
	}
	c.ContactListStore = &contactListStore

	contactDetailStore := sqlite.ContactDetailStore{}
	c.ContactDetailStore = &contactDetailStore

	return &c, nil
}

func NewServer(
    config *models.Config,
) http.Handler {
	mux := http.NewServeMux()
	addRoutes(
			mux,
			config,
	)
	var handler http.Handler = mux
	return handler
}

func run(
	ctx context.Context,  
	stdout, stderr io.Writer, 
	args []string,
	getenv func(string) string,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// register logging package with given stdout, err
	log.Configure(stdout, stderr, getenv)

	// run initialization
	i, err := initialize(ctx)
	if err != nil {
		return err
	}
	
	// build server	

	// run server and return result

	
	return nil
}


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloWorld)
	fmt.Println("Starting server on 8080...")
	http.ListenAndServe(":8080", mux)
}
