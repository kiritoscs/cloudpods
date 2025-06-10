package models

import (
	"context"
	"fmt"

	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/jsonutils"
	billing_api "yunion.io/x/onecloud/pkg/apis/billing"
	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/lockman"
	"yunion.io/x/onecloud/pkg/httperrors"
	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/util/stringutils2"
	"yunion.io/x/pkg/errors"
	"yunion.io/x/pkg/util/compare"
	"yunion.io/x/sqlchemy"
)

type SLighthouseManager struct {
	// 由于资源是用户资源，因此定义为Virtual资源
	db.SVirtualResourceBaseManager
	db.SExternalizedResourceBaseManager
	SDeletePreventableResourceBaseManager

	SCloudregionResourceBaseManager
	SManagedResourceBaseManager
}

var LighthouseManager *SLighthouseManager

func init() {
	LighthouseManager = &SLighthouseManager{
		SVirtualResourceBaseManager: db.NewVirtualResourceBaseManager(
			SLighthouse{},
			"lighthouses_tbl",
			"lighthouse",
			"lighthouses",
		),
	}
	LighthouseManager.SetVirtualObject(LighthouseManager)
}

type SLighthouse struct {
	db.SVirtualResourceBase
	db.SExternalizedResourceBase
	SManagedResourceBase
	SBillingResourceBase

	SCloudregionResourceBase
	SDeletePreventableResourceBase
}

func (manager *SLighthouseManager) GetContextManagers() [][]db.IModelManager {
	return [][]db.IModelManager{
		{CloudregionManager},
	}
}

func (man *SLighthouseManager) ListItemFilter(
	ctx context.Context,
	q *sqlchemy.SQuery,
	userCred mcclient.TokenCredential,
	query api.LighthouseListInput,
) (*sqlchemy.SQuery, error) {
	var err error
	q, err = man.SVirtualResourceBaseManager.ListItemFilter(ctx, q, userCred, query.VirtualResourceListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SVirtualResourceBaseManager.ListItemFilter")
	}
	q, err = man.SExternalizedResourceBaseManager.ListItemFilter(ctx, q, userCred, query.ExternalizedResourceBaseListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SExternalizedResourceBaseManager.ListItemFilter")
	}
	q, err = man.SDeletePreventableResourceBaseManager.ListItemFilter(ctx, q, userCred, query.DeletePreventableResourceBaseListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SDeletePreventableResourceBaseManager.ListItemFilter")
	}
	q, err = man.SManagedResourceBaseManager.ListItemFilter(ctx, q, userCred, query.ManagedResourceListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SManagedResourceBaseManager.ListItemFilter")
	}
	q, err = man.SCloudregionResourceBaseManager.ListItemFilter(ctx, q, userCred, query.RegionalFilterListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SCloudregionResourceBaseManager.ListItemFilter")
	}

	return q, nil
}

func (man *SLighthouseManager) OrderByExtraFields(
	ctx context.Context,
	q *sqlchemy.SQuery,
	userCred mcclient.TokenCredential,
	query api.LighthouseListInput,
) (*sqlchemy.SQuery, error) {
	q, err := man.SVirtualResourceBaseManager.OrderByExtraFields(ctx, q, userCred, query.VirtualResourceListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SVirtualResourceBaseManager.OrderByExtraFields")
	}
	q, err = man.SCloudregionResourceBaseManager.OrderByExtraFields(ctx, q, userCred, query.RegionalFilterListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SCloudregionResourceBaseManager.OrderByExtraFields")
	}
	q, err = man.SManagedResourceBaseManager.OrderByExtraFields(ctx, q, userCred, query.ManagedResourceListInput)
	if err != nil {
		return nil, errors.Wrap(err, "SManagedResourceBaseManager.OrderByExtraFields")
	}
	return q, nil
}

func (man *SLighthouseManager) QueryDistinctExtraField(q *sqlchemy.SQuery, field string) (*sqlchemy.SQuery, error) {
	q, err := man.SVirtualResourceBaseManager.QueryDistinctExtraField(q, field)
	if err == nil {
		return q, nil
	}
	q, err = man.SCloudregionResourceBaseManager.QueryDistinctExtraField(q, field)
	if err == nil {
		return q, nil
	}
	q, err = man.SManagedResourceBaseManager.QueryDistinctExtraField(q, field)
	if err == nil {
		return q, nil
	}
	return q, httperrors.ErrNotFound
}

// ValidateCreateData 验证创建参数
func (man *SLighthouseManager) ValidateCreateData(ctx context.Context, userCred mcclient.TokenCredential, ownerId mcclient.IIdentityProvider, query jsonutils.JSONObject, input api.LighthouseCreateInput) (api.LighthouseCreateInput, error) {
	return input, httperrors.NewNotImplementedError("Not Implemented")
}

// FetchCustomizeColumns 获取自定义字段
func (man *SLighthouseManager) FetchCustomizeColumns(
	ctx context.Context,
	userCred mcclient.TokenCredential,
	query jsonutils.JSONObject,
	objs []interface{},
	fields stringutils2.SSortedStrings,
	isList bool,
) []api.LighthouseDetails {
	rows := make([]api.LighthouseDetails, len(objs))
	virtRows := man.SVirtualResourceBaseManager.FetchCustomizeColumns(ctx, userCred, query, objs, fields, isList)
	manRows := man.SManagedResourceBaseManager.FetchCustomizeColumns(ctx, userCred, query, objs, fields, isList)
	regRows := man.SCloudregionResourceBaseManager.FetchCustomizeColumns(ctx, userCred, query, objs, fields, isList)
	for i := range rows {
		rows[i] = api.LighthouseDetails{
			VirtualResourceDetails:  virtRows[i],
			ManagedResourceInfo:     manRows[i],
			CloudregionResourceInfo: regRows[i],
		}
	}
	return rows
}

func (self *SCloudregion) GetLighthouses(managerId string) ([]SLighthouse, error) {
	q := LighthouseManager.Query().Equals("cloudregion_id", self.Id)
	if len(managerId) > 0 {
		q = q.Equals("manager_id", managerId)
	}
	ret := []SLighthouse{}
	err := db.FetchModelObjects(LighthouseManager, q, &ret)
	if err != nil {
		return nil, errors.Wrapf(err, "db.FetchModelObjects")
	}
	return ret, nil
}

func (self *SCloudregion) SyncLighthouses(ctx context.Context, userCred mcclient.TokenCredential, provider *SCloudprovider, exts []cloudprovider.ICloudLighthouse) compare.SyncResult {
	// 加锁防止重入
	lockman.LockRawObject(ctx, LighthouseManager.KeywordPlural(), fmt.Sprintf("%s-%s", provider.Id, self.Id))
	defer lockman.ReleaseRawObject(ctx, LighthouseManager.KeywordPlural(), fmt.Sprintf("%s-%s", provider.Id, self.Id))

	result := compare.SyncResult{}

	dbLighthouse, err := self.GetLighthouses(provider.Id)
	if err != nil {
		result.Error(err)
		return result
	}

	removed := make([]SLighthouse, 0)
	commondb := make([]SLighthouse, 0)
	commonext := make([]cloudprovider.ICloudLighthouse, 0)
	added := make([]cloudprovider.ICloudLighthouse, 0)
	// 本地和云上资源列表进行比对
	err = compare.CompareSets(dbLighthouse, exts, &removed, &commondb, &commonext, &added)
	if err != nil {
		result.Error(err)
		return result
	}

	// 删除云上没有的资源
	for i := 0; i < len(removed); i++ {
		err := removed[i].syncRemoveCloudLighthouse(ctx, userCred)
		if err != nil {
			result.DeleteError(err)
			continue
		}
		result.Delete()
	}

	// 和云上资源属性进行同步
	for i := 0; i < len(commondb); i++ {
		err := commondb[i].SyncWithCloudLighthouse(ctx, userCred, commonext[i])
		if err != nil {
			result.UpdateError(err)
			continue
		}
		result.Update()
	}

	// 创建本地没有的云上资源
	for i := 0; i < len(added); i++ {
		_, err := self.newFromCloudLighthouse(ctx, userCred, provider, added[i])
		if err != nil {
			result.AddError(err)
			continue
		}
		result.Add()
	}
	return result
}

// 判断资源是否可以删除
func (self *SLighthouse) ValidateDeleteCondition(ctx context.Context, info jsonutils.JSONObject) error {
	if self.DisableDelete.IsTrue() {
		return httperrors.NewInvalidStatusError("Lighthouse is locked, cannot delete")
	}
	return self.SStatusStandaloneResourceBase.ValidateDeleteCondition(ctx, info)
}

func (self *SLighthouse) syncRemoveCloudLighthouse(ctx context.Context, userCred mcclient.TokenCredential) error {
	return self.Delete(ctx, userCred)
}

// SyncWithCloudLighthouse 同步资源属性
func (self *SLighthouse) SyncWithCloudLighthouse(ctx context.Context, userCred mcclient.TokenCredential, ext cloudprovider.ICloudLighthouse) error {
	diff, err := db.UpdateWithLock(ctx, self, func() error {
		self.ExternalId = ext.GetGlobalId()
		self.Status = ext.GetStatus()

		self.BillingType = ext.GetBillingType()
		if self.BillingType == billing_api.BILLING_TYPE_PREPAID {
			if expiredAt := ext.GetExpiredAt(); !expiredAt.IsZero() {
				self.ExpiredAt = expiredAt
			}
			self.AutoRenew = ext.IsAutoRenew()
		}
		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "db.Update")
	}

	syncVirtualResourceMetadata(ctx, userCred, self, ext)
	if provider := self.GetCloudprovider(); provider != nil {
		SyncCloudProject(ctx, userCred, self, provider.GetOwnerId(), ext, provider.Id)
	}
	db.OpsLog.LogSyncUpdate(self, diff, userCred)
	return nil
}

// newFromCloudLighthouse 同步单个Lighthouse
func (self *SCloudregion) newFromCloudLighthouse(ctx context.Context, userCred mcclient.TokenCredential, provider *SCloudprovider, ext cloudprovider.ICloudLighthouse) (*SLighthouse, error) {
	es := SLighthouse{}
	es.SetModelManager(LighthouseManager, &es)

	es.ExternalId = ext.GetGlobalId()
	es.CloudregionId = self.Id
	es.ManagerId = provider.Id
	es.IsEmulated = ext.IsEmulated()
	es.Status = ext.GetStatus()

	if createdAt := ext.GetCreatedAt(); !createdAt.IsZero() {
		es.CreatedAt = createdAt
	}

	es.BillingType = ext.GetBillingType()
	if es.BillingType == billing_api.BILLING_TYPE_PREPAID {
		if expired := ext.GetExpiredAt(); !expired.IsZero() {
			es.ExpiredAt = expired
		}
		es.AutoRenew = ext.IsAutoRenew()
	}

	var err error
	err = func() error {
		// 这里加锁是为了防止名称重复
		lockman.LockRawObject(ctx, LighthouseManager.Keyword(), "name")
		defer lockman.ReleaseRawObject(ctx, LighthouseManager.Keyword(), "name")

		es.Name, err = db.GenerateName(ctx, LighthouseManager, provider.GetOwnerId(), ext.GetName())
		if err != nil {
			return errors.Wrapf(err, "db.GenerateName")
		}
		return LighthouseManager.TableSpec().Insert(ctx, &es)
	}()
	if err != nil {
		return nil, errors.Wrapf(err, "newFromCloudLighthouse.Insert")
	}

	// 同步标签
	syncVirtualResourceMetadata(ctx, userCred, &es, ext)
	// 同步项目归属
	SyncCloudProject(ctx, userCred, &es, provider.GetOwnerId(), ext, provider.Id)

	db.OpsLog.LogEvent(&es, db.ACT_CREATE, es.GetShortDesc(ctx), userCred)

	return &es, nil
}

func (manager *SLighthouseManager) ListItemExportKeys(ctx context.Context,
	q *sqlchemy.SQuery,
	userCred mcclient.TokenCredential,
	keys stringutils2.SSortedStrings,
) (*sqlchemy.SQuery, error) {
	var err error

	q, err = manager.SVirtualResourceBaseManager.ListItemExportKeys(ctx, q, userCred, keys)
	if err != nil {
		return nil, errors.Wrap(err, "SVirtualResourceBaseManager.ListItemExportKeys")
	}

	if keys.ContainsAny(manager.SManagedResourceBaseManager.GetExportKeys()...) {
		q, err = manager.SManagedResourceBaseManager.ListItemExportKeys(ctx, q, userCred, keys)
		if err != nil {
			return nil, errors.Wrap(err, "SManagedResourceBaseManager.ListItemExportKeys")
		}
	}

	if keys.ContainsAny(manager.SCloudregionResourceBaseManager.GetExportKeys()...) {
		q, err = manager.SCloudregionResourceBaseManager.ListItemExportKeys(ctx, q, userCred, keys)
		if err != nil {
			return nil, errors.Wrap(err, "SCloudregionResourceBaseManager.ListItemExportKeys")
		}
	}

	return q, nil
}
