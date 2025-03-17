package helm

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"k8s.io/apimachinery/pkg/util/validation"
)

type ApplicationInfo struct {
	K8sParam     string                 `json:"k8sParams"`
	ApolloParam  string                 `json:"apolloParams"`
	Metadata     *chart.Metadata        `json:"metadata"`
	Templates    []*chart.File          `json:"templates"`
	Values       map[string]interface{} `json:"values"`
	Files        []*chart.File          `json:"files"`
	Dependencies []*chart.Chart         `json:"dependencies"`
	ChartJson    string                 `json:"chartJson"` // 全量的chart数据
}

type RenderOpt struct {
	AppId       uint
	ReleaseName string
	IsUpgrade   bool
	Namespace   string
	KubeVersion string
}

func convertPortToSlice(input interface{}) []interface{} {
	// 将数组拼接的字符串，转换为数组结构
	arrayString := input.(string)
	// 移除字符串中的 "[" 和 "]"
	arrayString = strings.Trim(arrayString, "[]")

	// 将字符串分割为数组
	arrayStringSlice := strings.Split(arrayString, " ")

	// 将字符串数组转换为 interface{} 数组
	var interfaceSlice []interface{}
	for _, strValue := range arrayStringSlice {
		interfaceSlice = append(interfaceSlice, strValue)
	}
	return interfaceSlice
}

func Render3(p *ApplicationInfo, opt *RenderOpt, parameters *map[string]interface{}) (string, error) {
	if msgs := validation.IsDNS1123Label(opt.ReleaseName); opt.ReleaseName != "" && len(msgs) > 0 {
		return "", fmt.Errorf("release name %s is not a valid DNS label: %s", opt.ReleaseName, strings.Join(msgs, ";"))
	}

	if opt.Namespace == "" {
		return "", fmt.Errorf("namespace is not allow empty")
	}

	log.Println(parameters)

	kv, err := chartutil.ParseKubeVersion(opt.KubeVersion)
	if err != nil {
		log.Printf("解析KubeVersion %v失败", err)
	}

	// 创建一个 helm 配置对象
	settings := cli.New()
	client := action.NewInstall(&action.Configuration{})
	// 创建一个 action 对象，用于执行 helm 操作
	client.ReleaseName = opt.ReleaseName // 设置 release 名称
	client.Namespace = opt.Namespace     // 设置 命名空间

	client.DryRun = true     // 设置为 dry-run 模式，不实际安装 chart
	client.Replace = true    // 设置为替换模式，如果 release 已存在则覆盖
	client.ClientOnly = true // 设置为客户端模式，不与集群交互
	client.KubeVersion = kv  // kv版本号
	// client.ChartPathOptions.Verify = true // will verify the chart.

	// 加载 chart 文件
	//log.Infof("加载chart的路径为%s", p.ChartPath)
	//chartPath, err := client.ChartPathOptions.LocateChart(p.ChartPath, settings)
	//if err != nil {
	//	log.Errorf("LocateChart失败。chart:%s chartPath: %s失败，失败详情: %v。", opt.ReleaseName, p.ChartPath, err)
	//}
	//chart, err := loader.Load(chartPath)
	//if err != nil {
	//	log.Errorf("加载chart文件失败。chart:%s chartPath: %s失败，失败详情: %v。", opt.ReleaseName, p.ChartPath, err)
	//}

	// 加载数据库中的chartJson字符串，转换为*chart.Chart类型
	var localChart *chart.Chart
	err = json.Unmarshal([]byte(p.ChartJson), &localChart)
	if err != nil {
		log.Fatal("Failed to unmarshal chart data:", err)
	}

	// dependencies 是私有变量，需要set进去
	localChart.SetDependencies(p.Dependencies...)

	// 创建包含 map 的切片
	var StringValuesRes []string
	var ValuesRes []string
	for key, value := range *parameters {
		// 将键值对格式化为字符串，并添加到切片中
		switch value.(type) {
		case bool, int:
			ValuesRes = append(ValuesRes, fmt.Sprintf(`%s=%v`, key, value))
		//case []uint32:
		//	ValuesRes = append(ValuesRes, fmt.Sprintf(`%s=%v`, key, strings.Join(value, ",")))
		default:
			StringValuesRes = append(StringValuesRes, fmt.Sprintf(`%s=%v`, key, value))
		}
	}

	// 创建一个 values 对象，用于存储 chart 的值
	vals := values.Options{
		StringValues: StringValuesRes, // --set-string
		Values:       ValuesRes,       // --set

	}
	valueOpts := &vals
	opts := getter.All(settings)
	val, err := valueOpts.MergeValues(opts)
	if err != nil {
		log.Fatal(err)
	}

	if allocateNodePortNumVal, ok := val["allocateNodePortNum"]; ok {
		val["allocateNodePortNum"] = convertPortToSlice(allocateNodePortNumVal)
	}

	// 执行 helm template 操作，并获取结果
	rel, err := client.Run(localChart, val)
	if err != nil {
		log.Fatalf("渲染chart:%s 失败详情: %v。", opt.ReleaseName, err)
	}

	log.Printf("[%v]%v的helm3渲染后的yaml文件为%v", opt.Namespace, opt.ReleaseName, rel.Manifest)

	// 将结果输出为 yaml 格式
	return rel.Manifest, nil
}
