package main

import (
	"flag"
	"os"
	"runtime"
	"time"

	"service/internal/conf/v1"
	"service/internal/out/broker"
	mylog "service/pkg/logger"

	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/joho/godotenv"

	krlogrus "github.com/go-kratos/kratos/contrib/log/logrus/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	_ "go.uber.org/automaxprocs"
)

// Path to configs
var flagconf string

// Generation instance ID: ENV > hostname > fallback
func instanceID() string {
	if v := os.Getenv("SERVICE_ID"); v != "" {
		return v
	}
	if h, err := os.Hostname(); err == nil && h != "" {
		return h
	}
	return "instance-" + time.Now().UTC().Format("20060102T150405Z")
}

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func newLogger(mode string) klog.Logger {
	lr := mylog.Init(mode)
	base := krlogrus.NewLogger(lr)
	return klog.With(base, "caller", klog.DefaultCaller)
}

func newApp(logger klog.Logger, app *conf.App, gs *grpc.Server, hs *http.Server, broker *broker.Broker, data *conf.Data) *kratos.App {

	// Start MQTT broker
	go broker.StartMQTT(data)

	md := map[string]string{"env": envOr("APP_ENV", "dev"), "go": runtime.Version()}

	return kratos.New(
		kratos.ID(instanceID()),
		kratos.Name(app.Name),       // from config.yaml
		kratos.Version(app.Version), // from config.yaml
		kratos.Metadata(md),
		kratos.Logger(logger),
		kratos.Server(gs, hs),
	)
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func loadDotEnv() {
	candidates := []string{".env", "../.env", "../../.env"}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			_ = godotenv.Load(p)
			return
		}
	}
}

func main() {
	flag.Parse()
	loadDotEnv()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	logger := newLogger(bc.App.Mode)

	app, cleanup, err := wireApp(bc.App, bc.Server, bc.Data, logger)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := app.Run(); err != nil {
		panic(err)
	}
}
