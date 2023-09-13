package admin

import depinjection "backend-service/pkg/common/dep_injection"

var Module = depinjection.BulkProvide(
	[]any{},
	"admin-controller",
)
