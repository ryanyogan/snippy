package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.yogan.dev/internal/models"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

const (
	username = "doadmin"
	password = "AVNS_F6Tom7HO4JzlkeRvDNn"
	hostname = "snippetbox-do-user-6565302-0.b.db.ondigitalocean.com:25060"
	dbname   = "snippetbox"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// t := "ryan:hummer@/snippetbox?parseTime=true"
	// s := "mysql://doadmin:AVNS_F6Tom7HO4JzlkeRvDNn@snippetbox-do-user-6565302-0.b.db.ondigitalocean.com:25060/snippetbox?ssl-mode=REQUIRED&?parseTime=true"
	dsn := flag.String("dsn", dsn("snippetbox?parseTime=true"), "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	logger.Info("Listening: http://localhost"+port, "addr", port)
	err = http.ListenAndServe(":"+port, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
