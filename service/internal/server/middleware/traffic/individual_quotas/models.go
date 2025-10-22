package individual_quotas

// Modelos de red (JSON) usados al descargar las cuotas.
type quotaItem struct {
	Project  string `json:"project"`
	Route    string `json:"route"`
	Quota    int    `json:"quota"`    // RPS
	Interval int    `json:"interval"` // Interval in seconds (0 = 1 second)
}

type quotaResponse struct {
	Items []quotaItem `json:"items"`
}
