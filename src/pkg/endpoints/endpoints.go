package endpoints

const (
	// Roots
	rootEndpoint       = "/api/v1"
	queryEndpoint      = rootEndpoint + "/query"
	RegistrationEndpoint = rootEndpoint + "/registration"
	SubmissionEndpoint = rootEndpoint + "/submit"

	// Queries
	QueryHealthCheck   = queryEndpoint + "/health_check"
	QueryServerCPU     = queryEndpoint + "/server_cpu"
	QueryServerMemory  = queryEndpoint + "/server_memory"
	QueryServerStorage = queryEndpoint + "/server_storage"
	QueryServerTemp    = queryEndpoint + "/server_temp"
	
	// Registrations
	RegisterAgent = RegistrationEndpoint + "/agent"

	// Submissions
	SubmitHostCheckIn = SubmissionEndpoint + "/host_check_in"
	SubmitHostCPU     = SubmissionEndpoint + "/host_cpu"
	SubmitHostMemory  = SubmissionEndpoint + "/host_memory"
	SubmitHostStorage = SubmissionEndpoint + "/host_storage"
	SubmitHostTemp    = SubmissionEndpoint + "/host_temp"
)
