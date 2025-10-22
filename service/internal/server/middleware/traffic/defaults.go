package traffic

import (
	"runtime"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var DefaultConfig = Config{
	InflightMax:  400,                            // Valor maximo de concurrentes requests (concurrencia) // Максимальное количество одновременных запросов (конкурентность)
	RateRPS:      150,                            // Valor objetivo de requests por segundo (token bucket) // Целевой лимит запросов в секунду (token bucket)
	RateBurst:    300,                            // Valor maximo de tokens (capacidad del bucket) // Максимальное количество токенов (емкость bucket)
	KeyBy:        KeyIP,                          // Como agrupar el limite: por IP (ip)/usuario (user)/global (global) // Как группировать лимит: по IP (ip)/пользователю (user)/глобально (global)
	EnableCPU:    true,                           // Activar proteccion adaptativa por CPU (BBR)
	CPUWindow:    10 * time.Second,               // Ventana de observacion BBR para calcular metricas // Окно наблюдения BBR для расчета метрик
	CPUBuckets:   100,                            // Numero de buckets dentro de la ventana (precision de mediciones) // Число бакетов внутри окна (точность измерений)
	CPUThreshold: 800,                            // Umbral de carga de CPU en milesimas (800 = 80%) // Порог загрузки CPU в тысячных (800 = 80%)
	CPUQuota:     float64(runtime.GOMAXPROCS(0)), // Quota efectiva de CPU para BBR (cantidad de CPU disponibles) // Эффективная квота CPU для BBR (количество доступных CPU)
}

var DefaultConfigTest = Config{
	InflightMax:  10,                             // Poca concurrencia para pruebas rápidas // Маленькая конкурентность для быстрых тестов
	RateRPS:      5,                              // Bajo RPS para pruebas rápidas // Низкий RPS для быстрой проверки ограничения
	RateBurst:    10,                             // Pequeño pico // Небольшой всплеск
	KeyBy:        KeyIP,                          // Agrupar por IP // Группировка по IP
	EnableCPU:    true,                           // BBR activado // BBR включен
	CPUWindow:    1 * time.Second,                // Ventana corta para respuestas rápidas en pruebas // Короткое окно для быстрых реакций в тестах
	CPUBuckets:   10,                             // Menos buckets — más rápido cálculo // Меньше бакетов — быстрее расчет
	CPUThreshold: 800,                            // 80%
	CPUQuota:     float64(runtime.GOMAXPROCS(0)), // Usar CPU efectivas // Использовать эффективные CPU
}

// backward-compatible helpers used by server constructors
func HTTPConfig(logger log.Logger) Config     { return DefaultConfigWithLog(DefaultConfig, logger) }
func HTTPConfigTest(logger log.Logger) Config { return DefaultConfigWithLog(DefaultConfigTest, logger) } // tiny defaults for fast tests

func GRPCConfig(logger log.Logger) Config     { return DefaultConfigWithLog(DefaultConfig, logger) }
func GRPCConfigTest(logger log.Logger) Config { return DefaultConfigWithLog(DefaultConfigTest, logger) } // tiny defaults for fast tests

func DefaultConfigWithLog(cfg Config, logger log.Logger) Config {
	cfg.LogHelper = log.NewHelper(logger)
	return cfg
}
