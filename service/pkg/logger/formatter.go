package logger

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05.000")

	var levelColor *color.Color
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = color.New(color.FgMagenta).Add(color.Bold)
	case logrus.InfoLevel:
		levelColor = color.New(color.FgBlue).Add(color.Bold)
	case logrus.WarnLevel:
		levelColor = color.New(color.FgYellow).Add(color.Bold)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = color.New(color.FgRed).Add(color.Bold)
	default:
		levelColor = color.New(color.FgWhite).Add(color.Bold)
	}

	level := strings.ToUpper(entry.Level.String())
	prefix := levelColor.Sprintf("%s", padLevel(level))

	if logType, exists := entry.Data["type"]; exists {
		switch logType {
		case "route":
			var routeColor *color.Color
			method := ""
			path := ""
			if m, ok := entry.Data["method"].(string); ok {
				method = strings.ToUpper(m)
			}
			if p, ok := entry.Data["path"].(string); ok {
				path = p
			}
			var methodLabel string
			switch method {
			case "GET":
				routeColor = color.New(color.FgHiWhite).Add(color.Bold).Add(color.BgHiBlue)
				methodLabel = routeColor.Sprintf(" %-10s", method)
			case "POST":
				routeColor = color.New(color.FgBlack).Add(color.Bold).Add(color.BgHiCyan)
				methodLabel = routeColor.Sprintf(" %-10s", method)
			case "PUT":
				routeColor = color.New(color.FgBlack).Add(color.Bold).Add(color.BgHiYellow)
				methodLabel = routeColor.Sprintf(" %-10s", method)
			case "DELETE":
				routeColor = color.New(color.FgWhite).Add(color.Bold).Add(color.BgHiRed)
				methodLabel = routeColor.Sprintf(" %-10s", method)
			default:
				routeColor = color.New(color.BgHiBlack).Add(color.Bold).Add(color.BgHiWhite)
				methodLabel = routeColor.Sprintf(" %-10s", method)
			}
			purple := color.New(color.FgMagenta).Add(color.Bold)

			routePrefix := purple.Sprintf("%s %s Route Accessed", ">>>", timestamp)
			msg := routePrefix

			ip := ""
			if ipVal, ok := entry.Data["ip"].(string); ok && ipVal != "" {
				ip = ipVal
			}

			if method != "" && path != "" {
				if ip != "" {
					msg += fmt.Sprintf(" | %s | %-16s | \"%s\"", methodLabel, ip, path)
				} else {
					msg += fmt.Sprintf(" | %s | %-16s | \"%s\"", methodLabel, "", path)
				}
			} else if method != "" {
				if ip != "" {
					msg += fmt.Sprintf(" | %s | %-16s", methodLabel, ip)
				} else {
					msg += fmt.Sprintf(" | %s | %-16s", methodLabel, "")
				}
			} else if path != "" {
				if ip != "" {
					msg += fmt.Sprintf(" | %-16s | \"%s\"", ip, path)
				} else {
					msg += fmt.Sprintf(" | %-16s | \"%s\"", "", path)
				}
			}

			for k, v := range entry.Data {
				if k == "type" || k == "method" || k == "path" || k == "ip" {
					continue
				}
				msg += fmt.Sprintf(" [%s:%v]",
					purple.Sprint(k),
					color.New(color.FgWhite).Add(color.Bold).Sprint(v))
			}

			if entry.Message != "" {
				msg += " " + purple.Sprint(entry.Message)
			}
			return []byte(msg + "\n"), nil
		case "service":
			prefix = color.New(color.FgYellow).Add(color.Bold).Sprintf("[SERVICE]")
		default:
			prefix = levelColor.Sprintf("%s", padLevel(strings.ToUpper(fmt.Sprint(logType))))
		}
		delete(entry.Data, "type")
	}

	msg := fmt.Sprintf("%s %s %s",
		color.New(color.FgWhite).Add(color.Bold).Sprint(timestamp),
		prefix,
		color.New(color.FgWhite).Add(color.Bold).Sprint(entry.Message),
	)

	for k, v := range entry.Data {
		msg += fmt.Sprintf(" [%s:%v]",
			color.New(color.FgCyan).Add(color.Bold).Sprint(k),
			color.New(color.FgWhite).Add(color.Bold).Sprint(v))
	}

	return []byte(msg + "\n"), nil
}

func padLevel(level string) string {
	return fmt.Sprintf("[%s]%s", level, strings.Repeat(" ", 7-len(level)))
}
