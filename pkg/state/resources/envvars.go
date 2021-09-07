package resources

import (
	"fmt"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"

	"github.com/riser-platform/riser-server/pkg/core"
	"github.com/riser-platform/riser-server/pkg/util"
)

func k8sEnvVars(ctx *core.DeploymentContext) []corev1.EnvVar {
	envVars := []corev1.EnvVar{}
	// User defined  vars
	for key, val := range ctx.DeploymentConfig.App.Environment {
		envVars = append(envVars, corev1.EnvVar{
			Name:  strings.ToUpper(key),
			Value: val.String(),
		})
	}

	// Secret vars
	for _, secret := range ctx.Secrets {
		secretEnv := corev1.EnvVar{
			Name: strings.ToUpper(secret.Name),
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key:      "data",
					Optional: util.PtrBool(false),
					LocalObjectReference: corev1.LocalObjectReference{
						Name: fmt.Sprintf("%s-%s-%d", ctx.DeploymentConfig.App.Name, secret.Name, secret.Revision),
					},
				},
			},
		}

		envVars = append(envVars, secretEnv)
	}

	// Platform vars
	envVars = append(envVars, platformVarsFromContext(ctx)...)
	sort.Slice(envVars, func(i, j int) bool { return envVars[i].Name < envVars[j].Name })
	return envVars
}

func platformVarsFromContext(ctx *core.DeploymentContext) []corev1.EnvVar {
	return []corev1.EnvVar{
		{Name: "RISER_APP", Value: string(ctx.DeploymentConfig.App.Name)},
		{Name: "RISER_DEPLOYMENT", Value: string(ctx.DeploymentConfig.Name)},
		{Name: "RISER_DEPLOYMENT_REVISION", Value: fmt.Sprintf("%d", ctx.RiserRevision)},
		{Name: "RISER_ENVIRONMENT", Value: string(ctx.DeploymentConfig.EnvironmentName)},
		{Name: "RISER_NAMESPACE", Value: string(ctx.DeploymentConfig.Namespace)},
	}
}
