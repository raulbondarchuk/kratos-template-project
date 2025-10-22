package traffic

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// Configuration (simplified)
type KeyBy string

const (
	KeyGlobal KeyBy = "global"
	KeyIP     KeyBy = "ip"
	KeyUser   KeyBy = "user"
)

type Config struct {
	// maximum concurrent requests
	InflightMax int
	// token bucket (RPS/Burst)
	RateRPS   float64
	RateBurst int
	// grouping key (global/ip/user)
	KeyBy KeyBy
	// Adaptive protection by CPU (BBR)
	EnableCPU    bool
	CPUWindow    time.Duration
	CPUBuckets   int
	CPUThreshold int64 // 800 = 80%
	CPUQuota     float64

	LogHelper *log.Helper
}
