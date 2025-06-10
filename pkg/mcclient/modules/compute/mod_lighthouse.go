package compute

import (
	"yunion.io/x/onecloud/pkg/mcclient/modulebase"
	"yunion.io/x/onecloud/pkg/mcclient/modules"
)

type LighthouseManager struct {
	modulebase.ResourceManager
}

var (
	Lighthouses LighthouseManager
)

func init() {
	Lighthouses = LighthouseManager{modules.NewComputeManager("lighthouse", "lighthouses",
		[]string{},
		[]string{})}

	modules.RegisterCompute(&Lighthouses)
}
