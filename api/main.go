package main

import (
	"errors"
	"fmt"
	authMiddleware "github.com/itelman/forum/internal/handler/users/middleware"
	"github.com/itelman/forum/internal/service/categories"
	"github.com/itelman/forum/internal/service/filters"
	"github.com/itelman/forum/internal/service/posts"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/itelman/forum/internal/exception"
	"github.com/itelman/forum/internal/handler"
	"github.com/itelman/forum/internal/handler/home"
	usersHandlers "github.com/itelman/forum/internal/handler/users"
	"github.com/itelman/forum/internal/middleware/dynamic"
	"github.com/itelman/forum/internal/middleware/standard"
	"github.com/itelman/forum/internal/service/users"
	"github.com/itelman/forum/pkg/templates"

	_ "github.com/mattn/go-sqlite3"
)

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

	tmplRender := templates.NewTemplateRender(deps.templateCache)
	exceptionHandlers := exception.NewExceptions(errorLog, tmplRender)

	usersSvc := users.NewService(
		users.WithAPI(conf.ApiLink),
	)

	authMid := authMiddleware.NewMiddleware(usersSvc, exceptionHandlers)
	dynamicMiddleware := dynamic.NewMiddleware(authMid, exceptionHandlers)
	defaultHandlers := handler.NewHandlers(dynamicMiddleware, exceptionHandlers, tmplRender)

	postsSvc := posts.NewService(
		posts.WithAPI(conf.ApiLink),
	)

	categoriesSvc := categories.NewService(
		categories.WithAPI(conf.ApiLink),
	)

	filtersSvc := filters.NewService(
		filters.WithAPI(conf.ApiLink),
	)

	/*
		commentsSvc := comments.NewService(
			comments.WithSqlite(deps.sqlite),
		)

		postReactionsSvc := post_reactions.NewService(
			post_reactions.WithSqlite(deps.sqlite),
		)

		commentReactionsSvc := comment_reactions.NewService(
			comment_reactions.WithSqlite(deps.sqlite),
		)

		activitySvc := activity.NewService(
			activity.WithSqlite(deps.sqlite),
		)
	*/

	mux := http.NewServeMux()

	home.NewHandlers(defaultHandlers, postsSvc, categoriesSvc, filtersSvc).RegisterMux(mux)
	usersHandlers.NewHandlers(defaultHandlers, usersSvc).RegisterMux(mux)
	/*
		postsHandlers.NewHandlers(defaultHandlers, postsSvc, categoriesSvc).RegisterMux(mux)
		commentsHandlers.NewHandlers(defaultHandlers, commentsSvc).RegisterMux(mux)
		postReactionsHandlers.NewHandlers(defaultHandlers, postReactionsSvc).RegisterMux(mux)
		commentReactionsHandlers.NewHandlers(defaultHandlers, commentReactionsSvc).RegisterMux(mux)
		activityHandlers.NewHandlers(defaultHandlers, activitySvc).RegisterMux(mux)
	*/

	fileServer := http.FileServer(http.Dir(conf.UI.CSSDir))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", conf.Port),
		ErrorLog:     errorLog,
		Handler:      standard.NewMiddleware(exceptionHandlers, infoLog).Chain(mux),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		infoLog.Printf("Starting server on %s", conf.Host)

		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errorLog.Fatalf("Server error: %v", err)
		}
	}()

	sig := <-stop
	infoLog.Printf("Received shutdown signal: %v", sig)

	infoLog.Println("Server shutting down...")
}
