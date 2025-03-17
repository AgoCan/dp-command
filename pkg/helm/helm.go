package helm

import (
	"helm.sh/helm/v3/pkg/action"
)

type Helm struct {
	actionConfig *action.Configuration
}
