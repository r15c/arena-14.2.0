// Copyright 2024 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topnode

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/kubeflow/arena/pkg/apis/types"
	"gopkg.in/yaml.v2"
)

/*
format like:

NAME                       IPADDRESS      ROLE    STATUS  GPU(Total)  GPU(Allocated)  GPU_MODE
cn-shanghai.192.168.7.178  192.168.7.178  master  Ready   0           0               none
cn-shanghai.192.168.7.179  192.168.7.179  master  Ready   0           0               none
cn-shanghai.192.168.7.180  192.168.7.180  master  Ready   0           0               none
cn-shanghai.192.168.7.181  192.168.7.181  <none>  Ready   0           0               none
cn-shanghai.192.168.7.182  192.168.7.182  <none>  Ready   1           0               exclusive
cn-shanghai.192.168.7.186  192.168.7.186  <none>  Ready   4           0               topology
cn-shanghai.192.168.7.183  192.168.7.183  <none>  Ready   4           2.1             share
*/
func DisplayNodeSummary(nodeNames []string, targetNodeType types.NodeType, format types.FormatStyle, showMetric bool) error {
	totalGPUs := float64(0)
	allocatedGPUs := float64(0)
	unhealthyGPUs := float64(0)
	nodes, err := BuildNodes(nodeNames, targetNodeType, showMetric)
	if err != nil {
		return err
	}
	allNodeInfos := types.AllNodeInfo{}
	for _, processer := range GetSupportedNodePorcessers() {
		allNodeInfos = processer.Convert2NodeInfos(nodes, allNodeInfos)
	}
	switch format {
	case types.JsonFormat:
		data, _ := json.MarshalIndent(allNodeInfos, "", "    ")
		fmt.Printf("%v", string(data))
		return nil
	case types.YamlFormat:
		data, _ := yaml.Marshal(allNodeInfos)
		fmt.Printf("%v", string(data))
		return nil
	}
	var showNodeType bool
	var isUnhealthy bool
	nodeTypes := map[types.NodeType]bool{}
	for _, node := range nodes {
		nodeTypes[node.Type()] = true
		if !node.AllDevicesAreHealthy() {
			isUnhealthy = true
		}
	}
	if len(nodeTypes) == 1 {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		processers := GetSupportedNodePorcessers()
		for i := len(processers) - 1; i >= 0; i-- {
			processer := processers[i]
			processer.DisplayNodesCustomSummary(w, nodes)
		}
		_ = w.Flush()
		return nil
	}

	delete(nodeTypes, types.NormalNode)

	if len(nodeTypes) > 1 {
		showNodeType = true
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	header := []string{"NAME", "IPADDRESS", "ROLE", "STATUS", "GPU(Total)", "GPU(Allocated)"}
	if showNodeType {
		header = append(header, "GPU(Mode)")
	}
	if isUnhealthy {
		header = append(header, "UNHEALTHY")
	}
	PrintLine(w, header...)
	processers := GetSupportedNodePorcessers()
	for i := len(processers) - 1; i >= 0; i-- {
		processer := processers[i]
		t, a, u := processer.DisplayNodesSummary(w, nodes, showNodeType, isUnhealthy)
		totalGPUs += t
		allocatedGPUs += a
		unhealthyGPUs += u
	}
	if len(nodeNames) != 0 {
		_ = w.Flush()
		return nil
	}
	PrintLine(w, "---------------------------------------------------------------------------------------------------")
	PrintLine(w, "Allocated/Total GPUs In Cluster:")
	allocatedPercent := float64(0)
	if totalGPUs != 0 {
		allocatedPercent = float64(allocatedGPUs) / float64(totalGPUs) * 100
	}
	unhealthyPercent := float64(0)
	if totalGPUs != 0 {
		unhealthyPercent = float64(unhealthyGPUs) / float64(totalGPUs) * 100
	}
	PrintLine(w, fmt.Sprintf("%v/%v (%.1f%%)", allocatedGPUs, totalGPUs, allocatedPercent))
	if unhealthyGPUs != 0 {
		PrintLine(w, "Unhealthy/Total GPUs In Cluster:")
		PrintLine(w, fmt.Sprintf("%v/%v (%.1f%%)", unhealthyGPUs, totalGPUs, unhealthyPercent))
	}
	_ = w.Flush()
	return nil
}
