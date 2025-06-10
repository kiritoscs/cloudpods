package compute

import "yunion.io/x/onecloud/pkg/apis"

// LighthouseCreateInput 资源创建参数, 目前仅站位
type LighthouseCreateInput struct {
}

// LighthouseDetails 资源返回详情
type LighthouseDetails struct {
	apis.VirtualResourceDetails
	ManagedResourceInfo
	CloudregionResourceInfo
}

// LighthouseListInput 资源列表请求参数
type LighthouseListInput struct {
	apis.VirtualResourceListInput
	apis.ExternalizedResourceBaseListInput
	apis.DeletePreventableResourceBaseListInput

	RegionalFilterListInput
	ManagedResourceListInput
}
