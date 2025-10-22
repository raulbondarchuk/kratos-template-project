package individual_quotas

import "time"

// Configuración por defecto.
const (
	defaultServiceURL   = "http://10.70.20.80:10000"
	defaultRefreshEvery = 24 * time.Hour
	defaultBurstFactor  = 2.0
	defaultStrictMatch  = true
)
