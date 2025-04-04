package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	authMiddleware "github.com/itelman/forum/internal/handler/users/middleware"
	"github.com/itelman/forum/internal/service/activity"
	"github.com/itelman/forum/internal/service/categories"
	"github.com/itelman/forum/internal/service/comment_reactions"
	"github.com/itelman/forum/internal/service/comments"
	"github.com/itelman/forum/internal/service/filters"
	"github.com/itelman/forum/internal/service/post_reactions"
	"github.com/itelman/forum/internal/service/posts"

	"github.com/itelman/forum/internal/exception"
	"github.com/itelman/forum/internal/handler"
	activityHandlers "github.com/itelman/forum/internal/handler/activity"
	commentsHandlers "github.com/itelman/forum/internal/handler/comments"
	"github.com/itelman/forum/internal/handler/home"
	postsHandlers "github.com/itelman/forum/internal/handler/posts"
	commentReactionsHandlers "github.com/itelman/forum/internal/handler/reactions/comment_reactions"
	postReactionsHandlers "github.com/itelman/forum/internal/handler/reactions/post_reactions"
	usersHandlers "github.com/itelman/forum/internal/handler/users"
	"github.com/itelman/forum/internal/middleware/dynamic"
	"github.com/itelman/forum/internal/middleware/standard"
	"github.com/itelman/forum/internal/service/users"
	"github.com/itelman/forum/pkg/templates"
)

func checkDependencyHealth(url string, timeout time.Duration) error {
	client := &http.Client{Timeout: timeout}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("dependency API returned status: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	f, err := os.OpenFile("/tmp/info.log", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer f.Close()

	conf := newConfig()
	deps, err := NewDependencies(
		WithTemplateCache(conf.UI.TmplDir),
	)
	if err != nil {
		errorLog.Fatal(err)
	}

	tmplRender := templates.NewTemplateRender(deps.templateCache, deps.flashManager)
	exceptionHandlers := exception.NewExceptions(errorLog, tmplRender)

	usersSvc := users.NewService(
		users.WithAPI(conf.ApiLink),
	)

	authMid := authMiddleware.NewMiddleware(usersSvc, exceptionHandlers, deps.flashManager)
	dynamicMiddleware := dynamic.NewMiddleware(authMid, exceptionHandlers)
	defaultHandlers := handler.NewHandlers(dynamicMiddleware, exceptionHandlers, tmplRender, deps.flashManager)

	postsSvc := posts.NewService(
		posts.WithAPI(conf.ApiLink),
	)

	categoriesSvc := categories.NewService(
		categories.WithAPI(conf.ApiLink),
	)

	filtersSvc := filters.NewService(
		filters.WithAPI(conf.ApiLink),
	)

	commentsSvc := comments.NewService(
		comments.WithAPI(conf.ApiLink),
	)

	postReactionsSvc := post_reactions.NewService(
		post_reactions.WithAPI(conf.ApiLink),
	)

	commentReactionsSvc := comment_reactions.NewService(
		comment_reactions.WithAPI(conf.ApiLink),
	)

	activitySvc := activity.NewService(
		activity.WithAPI(conf.ApiLink),
	)

	mux := http.NewServeMux()

	home.NewHandlers(defaultHandlers, postsSvc, categoriesSvc, filtersSvc).RegisterMux(mux)
	usersHandlers.NewHandlers(defaultHandlers, usersSvc).RegisterMux(mux)
	postsHandlers.NewHandlers(defaultHandlers, postsSvc, categoriesSvc).RegisterMux(mux)
	commentsHandlers.NewHandlers(defaultHandlers, commentsSvc).RegisterMux(mux)
	postReactionsHandlers.NewHandlers(defaultHandlers, postReactionsSvc).RegisterMux(mux)
	commentReactionsHandlers.NewHandlers(defaultHandlers, commentReactionsSvc).RegisterMux(mux)
	activityHandlers.NewHandlers(defaultHandlers, activitySvc).RegisterMux(mux)

	fileServer := http.FileServer(http.Dir(conf.UI.CSSDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Step 1: Check if the required REST API is healthy
	infoLog.Println("Checking dependency API health...")
	if err := checkDependencyHealth(conf.ApiLink+"/health", 10*time.Second); err != nil {
		errorLog.Fatal(err)
	}

	// Step 2: Start the main server if dependency is healthy
	infoLog.Println("Dependency API is healthy, starting server...")

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", conf.Host, conf.Port),
		ErrorLog:     errorLog,
		Handler:      standard.NewMiddleware(exceptionHandlers, infoLog).Chain(mux),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		infoLog.Printf("Starting server at %s", fmt.Sprintf("http://%s", srv.Addr))

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errorLog.Fatal(err)
		}
	}()

	sig := <-stop
	infoLog.Printf("Received shutdown signal: %v", sig)

	infoLog.Println("Server shutting down...")
}
