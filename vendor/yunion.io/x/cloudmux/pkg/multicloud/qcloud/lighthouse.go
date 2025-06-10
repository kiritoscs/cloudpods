package qcloud

import (
	"fmt"
	"strconv"
	"time"

	billing_api "yunion.io/x/cloudmux/pkg/apis/billing"
	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud"
	"yunion.io/x/jsonutils"
)

// SLighthouse 轻量应用服务器
type SLighthouse struct {
	multicloud.SVirtualResourceBase
	multicloud.SBillingBase
	QcloudTags
	region             *SRegion
	Placement          Placement
	InstanceId         string
	InstanceName       string
	InstanceChargeType InstanceChargeType //PREPAID：表示预付费，即包年包月 POSTPAID_BY_HOUR：表示后付费，即按量计费 CDHPAID：CDH付费，即只对CDH计费，不对CDH上的实例计费。
	InstanceState      string             //	实例状态。取值范围：PENDING：表示创建中 LAUNCH_FAILED：表示创建失败 RUNNING：表示运行中 STOPPED：表示关机 STARTING：表示开机中 STOPPING：表示关机中 REBOOTING：表示重启中 SHUTDOWN：表示停止待销毁 TERMINATING：表示销毁中。
	CreatedTime        time.Time          //	创建时间。按照ISO8601标准表示，并且使用UTC时间。格式为：YYYY-MM-DDThh:mm:ssZ。
}

// 获取 Lighthouse 资源列表
func (self *SRegion) GetLighthouses(ids []string, offset int, limit int) ([]SLighthouse, int, error) {
	params := make(map[string]string)
	if limit < 1 || limit > 50 {
		limit = 50
	}

	params["Region"] = self.Region
	params["Limit"] = fmt.Sprintf("%d", limit)
	params["Offset"] = fmt.Sprintf("%d", offset)
	instances := make([]SLighthouse, 0)
	if ids != nil && len(ids) > 0 {
		for index, id := range ids {
			params[fmt.Sprintf("InstanceIds.%d", index)] = id
		}
	}
	body, err := self.lighthouseRequest("DescribeInstances", params)
	if err != nil {
		return nil, 0, err
	}
	err = body.Unmarshal(&instances, "InstanceSet")
	if err != nil {
		return nil, 0, err
	}
	total, _ := body.Float("TotalCount")
	return instances, int(total), nil
}

// 获取单个 Lighthouse 资源
func (self *SRegion) GetLighthouseById(instanceId string) (*SLighthouse, error) {
	instances, _, err := self.GetLighthouses([]string{instanceId}, 0, 1)
	if err != nil {
		return nil, err
	}
	if len(instances) == 0 {
		return nil, cloudprovider.ErrNotFound
	}
	if len(instances) > 1 {
		return nil, cloudprovider.ErrDuplicateId
	}
	// if instances[0].InstanceState == "LAUNCH_FAILED" {
	// 	return nil, cloudprovider.ErrNotFound
	// }
	return &instances[0], nil
}

func (self *SLighthouse) GetId() string {
	return self.InstanceId
}

func (self *SLighthouse) GetGlobalId() string {
	return self.InstanceId
}

func (self *SLighthouse) GetName() string {
	if len(self.InstanceName) > 0 && self.InstanceName != "未命名" {
		return self.InstanceName
	}
	return self.InstanceId
}

func (self *SLighthouse) Refresh() error {
	new, err := self.region.GetInstance(self.InstanceId)
	if err != nil {
		return err
	}
	return jsonutils.Update(self, new)
}

func (self *SLighthouse) GetCreatedAt() time.Time {
	return self.CreatedTime
}

func (self *SLighthouse) GetBillingType() string {
	switch self.InstanceChargeType {
	case PrePaidInstanceChargeType:
		return billing_api.BILLING_TYPE_PREPAID
	case PostPaidInstanceChargeType:
		return billing_api.BILLING_TYPE_POSTPAID
	default:
		return billing_api.BILLING_TYPE_PREPAID
	}
}

func (self *SLighthouse) GetProjectId() string {
	return strconv.Itoa(self.Placement.ProjectId)
}

// 获取资源状态
func (self *SLighthouse) GetStatus() string {
	switch self.InstanceState {
	case "PENDING":
		return api.VM_DEPLOYING
	case "LAUNCH_FAILED":
		return api.VM_DEPLOY_FAILED
	case "RUNNING":
		return api.VM_RUNNING
	case "STOPPED":
		return api.VM_READY
	case "STARTING", "REBOOTING":
		return api.VM_STARTING
	case "STOPPING":
		return api.VM_STOPPING
	case "SHUTDOWN":
		return api.VM_DEALLOCATED
	case "TERMINATING":
		return api.VM_DELETING
	default:
		return api.VM_UNKNOWN
	}
}

// 实现 GetIElasticSearchs 接口
func (self *SRegion) GetILighthouses() ([]cloudprovider.ICloudLighthouse, error) {
	// 获取当前region的所有elasticsearch实例
	lighthouses, _, err := self.GetLighthouses(nil, 0, 100)
	if err != nil {
		return nil, err
	}
	ret := []cloudprovider.ICloudLighthouse{}
	for i := range lighthouses {
		// 这里需要赋值，例如删除, 就可以使用 lighthouses[i].region.DeleteLighthouse(lighthouses[i].InstanceId)
		lighthouses[i].region = self
		ret = append(ret, &lighthouses[i])
	}
	return ret, nil
}

// 实现 GetIElasticSearchById 接口
func (self *SRegion) GetILighthouseById(id string) (cloudprovider.ICloudLighthouse, error) {
	lighthouse, err := self.GetLighthouseById(id)
	if err != nil {
		return nil, err
	}
	return lighthouse, nil
}
