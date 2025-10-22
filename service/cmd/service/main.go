package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"service/internal/conf/v1"
	"service/internal/out/broker"
	mylog "service/pkg/logger"

	krlogrus "github.com/go-kratos/kratos/contrib/log/logrus/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"
	_ "go.uber.org/automaxprocs"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func instanceID() string {
	if v := os.Getenv("SERVICE_ID"); v != "" {
		return v
	}
	if h, err := os.Hostname(); err == nil && h != "" {
		return h
	}
	return "instance-" + time.Now().UTC().Format("20060102T150405Z")
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func newLogger(mode string) klog.Logger {
	lr := mylog.Init(mode)
	base := krlogrus.NewLogger(lr)
	return klog.With(base, "caller", klog.DefaultCaller)
}

func newApp(logger klog.Logger, app *conf.App, gs *grpc.Server, hs *http.Server, b *broker.Broker, data *conf.Data) *kratos.App {
	// safe start broker
	if b != nil && data != nil {
		go b.Start(data)
	}

	md := map[string]string{"env": envOr("APP_ENV", "dev"), "go": runtime.Version()}

	// add only non-nil servers
	var servers []transport.Server
	if gs != nil {
		servers = append(servers, gs)
	}
	if hs != nil {
		servers = append(servers, hs)
	}

	return kratos.New(
		kratos.ID(instanceID()),
		kratos.Name(app.GetName()),
		kratos.Version(app.GetVersion()),
		kratos.Metadata(md),
		kratos.Logger(logger),
		kratos.Server(servers...),
	)
}

func loadDotEnv() {
	for _, p := range []string{".env", "../.env", "../../.env"} {
		if _, err := os.Stat(p); err == nil {
			_ = godotenv.Load(p)
			return
		}
	}
}

func main() {
	flag.Parse()
	loadDotEnv()

	// Context for graceful shutdown (Ctrl+C / SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	c := config.New(config.WithSource(file.NewSource(flagconf)))
	defer c.Close()

	if err := c.Load(); err != nil {
		log.Fatalf("config load: %v", err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		log.Fatalf("config scan: %v", err)
	}

	logger := newLogger(bc.App.GetMode())

	app, cleanup, err := wireApp(&bc, logger)
	if err != nil {
		log.Fatalf("bootstrap: %v", err)
	}

	// safe cleanup (nil-safe + protected from panic)
	defer func() {
		if cleanup != nil {
			defer func() { _ = recover() }()
			cleanup()
		}
	}()

	// Start
	go func() {
		if err := app.Run(); err != nil {
			// log and request shutdown
			logger.Log(klog.LevelError, "msg", "app.Run failed", "err", err)
			stop()
		}
	}()

	<-ctx.Done()

	// (optional) graceful shutdown timeout, if app has Stop(ctx)
	// shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	// _ = app.Stop(shutCtx)
}
