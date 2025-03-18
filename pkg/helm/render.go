package helm

import (
	"fmt"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
)

// ChartRenderer 用于渲染 Helm chart 的结构体
type ChartRenderer struct {
	settings  *cli.EnvSettings
	namespace string
}

func NewChartRenderer(namespace string) *ChartRenderer {
	settings := cli.New()

	if namespace == "" {
		namespace = "default"
	}

	return &ChartRenderer{
		settings:  settings,
		namespace: namespace,
	}
}

func (r *ChartRenderer) RenderChart(chartPath string, releaseName string, valueFiles []string, valuesParm map[string]interface{}) (string, error) {
	// 加载 chart
	chartRequested, err := loader.Load(chartPath)
	if err != nil {
		return "", fmt.Errorf("加载 chart 失败: %w", err)
	}

	// 创建 action 配置
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(r.settings.RESTClientGetter(), r.namespace, os.Getenv("HELM_DRIVER"), func(format string, v ...interface{}) {
		fmt.Printf(format, v...)
	}); err != nil {
		return "", fmt.Errorf("初始化 Helm 配置失败: %w", err)
	}

	// 创建 install 客户端
	client := action.NewInstall(actionConfig)
	client.DryRun = true
	client.ReleaseName = releaseName
	client.Replace = true
	client.ClientOnly = true
	client.IncludeCRDs = true
	client.Namespace = r.namespace

	// 合并 values
	valueOpts := &values.Options{
		ValueFiles: valueFiles,
		Values:     []string{},
	}

	// 获取 values
	vals, err := r.getValues(valueOpts)
	if err != nil {
		return "", fmt.Errorf("获取 values 失败: %w", err)
	}

	// 合并自定义 values
	for k, v := range valuesParm {
		vals[k] = v
	}

	// 执行安装（干运行模式）
	release, err := client.Run(chartRequested, vals)
	if err != nil {
		return "", fmt.Errorf("渲染 chart 失败: %w", err)
	}

	// 返回渲染后的 YAML
	return release.Manifest, nil
}

// getValues 合并所有的 values 文件和命令行指定的 values
func (r *ChartRenderer) getValues(valueOpts *values.Options) (map[string]interface{}, error) {
	providers := getter.All(r.settings)

	vals, err := valueOpts.MergeValues(providers)
	if err != nil {
		return nil, err
	}

	return vals, nil
}
