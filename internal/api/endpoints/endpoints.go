package endpoints

const (
	// Roots
	rootEndpoint       = "/api/v1"
	queryEndpoint      = rootEndpoint + "/query"
	submissionEndpoint = rootEndpoint + "/submit"

	// Queries
	QueryHealthCheck   = queryEndpoint + "/health_check"
	QueryServerCPU     = queryEndpoint + "/server_cpu"
	QueryServerMemory  = queryEndpoint + "/server_memory"
	QueryServerStorage = queryEndpoint + "/server_storage"
	QueryServerTemp    = queryEndpoint + "/server_temp"

	// Submissions
	SubmitHostCPU     = submissionEndpoint + "/host_cpu"
	SubmitHostMemory  = submissionEndpoint + "/host_memory"
	SubmitHostStorage = submissionEndpoint + "/host_storage"
	SubmitHostTemp    = submissionEndpoint + "/host_temp"
)
