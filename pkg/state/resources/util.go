package resources

import (
	"fmt"
	"strconv"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/util"
)

// deploymentLabels are labels common to Riser deployment resources
func deploymentLabels(ctx *core.DeploymentContext) map[string]string {
	return map[string]string{
		riserLabel("deployment"):  ctx.DeploymentConfig.Name,
		riserLabel("environment"): ctx.DeploymentConfig.EnvironmentName,
		riserLabel("app"):         string(ctx.DeploymentConfig.App.Name),
	}
}

// deploymentAnnotations are annotations common to Riser deployment resources
func deploymentAnnotations(ctx *core.DeploymentContext) map[string]string {
	return map[string]string{
		riserLabel("revision"):       strconv.FormatInt(ctx.RiserRevision, 10),
		riserLabel("server-version"): util.VersionString,
	}
}

// riserLabel returns a fully qualified riser label or annotation (e.g. riser.dev/your-label)
func riserLabel(labelName string) string {
	return fmt.Sprintf("riser.dev/%s", labelName)
}
