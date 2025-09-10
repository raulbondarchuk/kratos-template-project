package logger

import (
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	instance *logrus.Logger
	once     sync.Once
)

func Init(mode string) *logrus.Logger {
	once.Do(func() {
		instance = logrus.New()
		instance.SetFormatter(&CustomFormatter{})
		instance.SetOutput(os.Stdout)

		color.NoColor = false
		mode := strings.ToLower(mode)
		var logLevel string
		switch mode {
		case "dev", "local", "test":
			logLevel = "debug"
		case "prod":
			logLevel = "info"
		default:
			logLevel = "debug"
		}
		lvl, err := logrus.ParseLevel(logLevel)
		if err != nil {
			lvl = logrus.InfoLevel
		}
		instance.SetLevel(lvl)
		if mode == "" {
			mode = "unknown"
		}
		Info("Logger initialized", map[string]interface{}{"mode": mode})
	})
	return instance
}

func getLogger() *logrus.Logger {
	if instance == nil {
		Init("dev")
	}
	return instance
}

// --- Primary methods ---

func Info(msg string, fields ...map[string]interface{})    { logWithType("info", msg, fields...) }
func Error(msg string, fields ...map[string]interface{})   { logWithType("error", msg, fields...) }
func Warn(msg string, fields ...map[string]interface{})    { logWithType("warn", msg, fields...) }
func Debug(msg string, fields ...map[string]interface{})   { logWithType("debug", msg, fields...) }
func Service(msg string, fields ...map[string]interface{}) { logWithType("service", msg, fields...) }
func Gorm(msg string, fields ...map[string]interface{})    { logWithType("gorm", msg, fields...) }

func Route(method, path string, fields ...map[string]interface{}) {
	data := map[string]interface{}{"method": method, "path": path}

	if len(fields) > 0 && fields[0] != nil {
		for k, v := range fields[0] {
			data[k] = v
		}
	}
	logWithType("route", "", data)
}

func logWithType(logType, msg string, fields ...map[string]interface{}) {
	data := map[string]interface{}{"type": logType}
	if len(fields) > 0 && fields[0] != nil {
		for k, v := range fields[0] {
			data[k] = v
		}
	}
	entry := getLogger().WithFields(data)
	switch logType {
	case "error":
		entry.Error(msg)
	case "warn":
		entry.Warn(msg)
	case "debug":
		entry.Debug(msg)
	default:
		entry.Info(msg)
	}
}

/*
	time.Sleep(3 * time.Second)
	fmt.Printf("\n")
	logger.Route("GET", "/", map[string]interface{}{"ip": "127.0.0.1"})
	fmt.Printf("\n")
	logger.Route("POST", "/", map[string]interface{}{"ip": "127.0.0.1"})
	fmt.Printf("\n")
	logger.Route("PUT", "/", map[string]interface{}{"ip": "127.0.0.1"})
	fmt.Printf("\n")
	logger.Route("DELETE", "/", map[string]interface{}{"ip": "127.0.0.1"})
	fmt.Printf("\n")
	logger.Route("OPTIONS", "/", map[string]interface{}{"ip": "127.0.0.1"})
	fmt.Printf("\n")
	logger.Info("Application initialized")
	fmt.Printf("\n")
	logger.Info("Otro ejemplo de info", map[string]interface{}{"name": "John", "age": 30})
	fmt.Printf("\n")
	logger.Error("Error occurred", map[string]interface{}{"error": "test error"})
	fmt.Printf("\n")
	logger.Warn("Warning occurred", map[string]interface{}{"warning": "test warning"})
	fmt.Printf("\n")
	logger.Debug("Debug message", map[string]interface{}{"debug": "test debug"})
*/

// import (
// 	"fmt"

// 	"github.com/gin-gonic/gin"
// )

// func RouteRequestLogger() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		fmt.Print("\n")
// 		Route(c.Request.Method, c.Request.RequestURI, map[string]interface{}{"ip": c.ClientIP()})
// 		c.Next()
// 	}
// }
