package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gofrs/uuid"
	"gopkg.in/yaml.v3"
	"log"
	"medods/cmd/login/internal/adapters/notify"
	"medods/cmd/login/internal/adapters/repo"
	http_api "medods/cmd/login/internal/api/http"
	"medods/cmd/login/internal/app"
	"medods/cmd/login/internal/auth"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// TODO
type (
	config struct {
		Server serverConfig `yaml:"server"`
		DB     dbConfig     `yaml:"db"`
		Keys   keys         `yaml:"keys"`
	}
	serverConfig struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	}

	dbConfig struct {
		MigrateDir string         `yaml:"migrate_dir"`
		Driver     string         `yaml:"driver"`
		Postgres   repo.Connector `yaml:"postgres"`
	}
	keys struct {
		AccessKey  string `yaml:"access_key"`
		RefreshKey string `yaml:"refresh_key"`
	}
)

// ./cmd/compost/config.yml
// /build/config.yml
var (
	cfgFile = flag.String("cfg", "/build/config.yml", "path to config file")
)

func main() {
	fmt.Println(uuid.Must(uuid.NewV4()).String())
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	defer cancel()

	cfg := configRead(*cfgFile)

	postgresDB, err := repo.New(ctx, repo.Config{
		Postgres:   cfg.DB.Postgres,
		MigrateDir: cfg.DB.MigrateDir,
		Driver:     cfg.DB.Driver,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer postgresDB.Close()

	authModule := auth.New(cfg.Keys.AccessKey)

	appModule := app.New(postgresDB, authModule, &notify.Notify{})

	httpApi := http_api.New(appModule)

	log.Fatal(run(ctx, httpApi, cfg.Server))

}

func run(ctx context.Context, handler http.Handler, cfg serverConfig) error {
	srv := &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, cfg.Port),
		Handler: handler,
	}

	errc := make(chan error, 1)
	go func() {
		errc <- srv.ListenAndServe()
	}()

	log.Printf("started %s", net.JoinHostPort(cfg.Host, cfg.Port))
	defer log.Println("shutdown")

	var err error
	select {
	case err = <-errc:
	case <-ctx.Done():
		err = srv.Shutdown(context.Background())
	}

	if err != nil {
		return fmt.Errorf("srv.ListenAndServe: %w", err)
	}

	return nil
}

func configRead(cfgPath string) config {
	file, err := os.Open(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	cfg := config{}
	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
