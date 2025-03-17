package k8s

import (
	"context"

	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Client) GetNodes() (nodes *v1.NodeList, err error) {
	nodes, err = c.ClientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		zap.S().Error("get node err: ", err)
		return nil, err
	}
	return nodes, err
}

func (c *Client) GetNodeIps() (ips map[string]string, err error) {
	nodes, err := c.GetNodes()
	if err != nil {
		return nil, err
	}
	// 这里应该给map开辟一个空间，用items的长度
	ips = make(map[string]string, len(nodes.Items))

	for _, node := range nodes.Items {
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" {
				ips[node.Name] = addr.Address
				break
			}
		}
	}
	return ips, nil
}

func (c *Client) GetNode(nodeName string) (node *v1.Node, err error) {
	node, err = c.ClientSet.CoreV1().Nodes().Get(context.Background(), nodeName, metav1.GetOptions{})
	if err != nil {
		zap.S().Error("get node err: ", err)
		return nil, err
	}
	return node, err
}

func (c *Client) DeleteNode(nodes *v1.NodeList, addresses []string) (err error) {
	nodeNames := []string{}

	for _, address := range addresses {
		for _, node := range nodes.Items {
			for _, addr := range node.Status.Addresses {
				if addr.Type == "InternalIP" {
					if addr.Address == address {
						nodeNames = append(nodeNames, node.Name)
					}
				}
			}
		}
	}

	zap.S().Info("delete nodes: ", nodeNames)
	for _, name := range nodeNames {
		err = c.ClientSet.CoreV1().Nodes().Delete(context.TODO(), name, metav1.DeleteOptions{})
		if err != nil {
			zap.S().Error("delete node:", err)
		}
	}

	return
}

type NodeLabelInfo struct {
	IP     string
	Labels map[string]string
}

/**
 * @description: 获取节点的信息，返回节点的标签信息
 * @return {*}
 */
func (c *Client) GetNodeLabels() (nodeLabels map[string]NodeLabelInfo, err error) {
	nodes, err := c.GetNodes()
	if err != nil {
		return nil, err
	}
	// 这里应该给map开辟一个空间，用items的长度
	nodeLabels = make(map[string]NodeLabelInfo, len(nodes.Items))

	for _, node := range nodes.Items {
		nodeIp := ""
		for _, addr := range node.Status.Addresses {
			if addr.Type == "InternalIP" { // 选择你需要的地址类型，比如 InternalIP 或 ExternalIP
				nodeIp = addr.Address
				break
			}
		}
		nodeLabels[node.Name] = NodeLabelInfo{
			IP:     nodeIp,
			Labels: node.Labels,
		}
	}
	return nodeLabels, nil
}
