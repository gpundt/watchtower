package query

// Global function to start listening on API endpoint
func InitializeQueryEndpoints() {
	initializeHealthCheckAPIEndpoint()
	initializeServerCPUEndpoint()
	initializeServerStorageEndpoint()
	initializeServerMemoryEndpoint()
	initializeServerTempEndpoint()
}
