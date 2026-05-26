package endpoints

const (
	// Roots
	RootEndpoint  = "/api/v1"
	QueryEndpoint = RootEndpoint + "/query"

	// Queries
	HealthCheckEndpoint   = QueryEndpoint + "/health_check"
	ServerCPUEndpoint     = QueryEndpoint + "/server_cpu"
	ServerMemoryEndpoint  = QueryEndpoint + "/server_memory"
	ServerStorageEndpoint = QueryEndpoint + "/server_storage"
	ServerTempEndpoint    = QueryEndpoint + "/server_temp"
)
