package submission

import (
	Endpoints "watchtower/internal/api/endpoints"
	Query "watchtower/internal/api/query"
)

func SubmitHostMetricsMaster() {
	submitHostMetrics(Endpoints.SubmitHostCPU, Query.GenerateHostCPUJSON())
	submitHostMetrics(Endpoints.SubmitHostMemory, Query.GenerateHostMemoryJSON())
	submitHostMetrics(Endpoints.SubmitHostStorage, Query.GenerateHostStorageJSON())
	submitHostMetrics(Endpoints.SubmitHostTemp, Query.GenerateHostTempJSON())
}

func InitializeHostMetricsAPIEndpoints() {
	initializeHostCPUSubmissionEndpoint()
	initializeHostMemorySubmissionEndpoint()
	initializeHostStorageSubmissionEndpoint()
	initializeHostTempSubmissionEndpoint()
}
