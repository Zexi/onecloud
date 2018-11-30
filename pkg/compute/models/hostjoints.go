package models

import (
	"context"

	"yunion.io/x/jsonutils"
	"yunion.io/x/onecloud/pkg/mcclient"

	"yunion.io/x/onecloud/pkg/cloudcommon/consts"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/policy"
)

type SHostJointsManager struct {
	db.SJointResourceBaseManager
}

func NewHostJointsManager(dt interface{}, tableName string, keyword string, keywordPlural string, slave db.IStandaloneModelManager) SHostJointsManager {
	return SHostJointsManager{
		SJointResourceBaseManager: db.NewJointResourceBaseManager(
			dt,
			tableName,
			keyword,
			keywordPlural,
			HostManager,
			slave,
		),
	}
}

type SHostJointsBase struct {
	db.SJointResourceBase
}

func (manager *SHostJointsManager) AllowListItems(ctx context.Context, userCred mcclient.TokenCredential, query jsonutils.JSONObject) bool {
	return userCred.IsAdminAllow(consts.GetServiceType(), manager.KeywordPlural(), policy.PolicyActionList)
}

func (manager *SHostJointsManager) AllowCreateItem(ctx context.Context, userCred mcclient.TokenCredential, query jsonutils.JSONObject, data jsonutils.JSONObject) bool {
	return userCred.IsAdminAllow(consts.GetServiceType(), manager.KeywordPlural(), policy.PolicyActionCreate)
}

func (manager *SHostJointsManager) AllowListDescendent(ctx context.Context, userCred mcclient.TokenCredential, model db.IStandaloneModel, query jsonutils.JSONObject) bool {
	return userCred.IsAdminAllow(consts.GetServiceType(), manager.KeywordPlural(), policy.PolicyActionList)
}

func (manager *SHostJointsManager) AllowAttach(ctx context.Context, userCred mcclient.TokenCredential, master db.IStandaloneModel, slave db.IStandaloneModel) bool {
	return userCred.IsAdminAllow(consts.GetServiceType(), manager.KeywordPlural(), policy.PolicyActionDelete)
}

func (self *SHostJointsBase) AllowGetDetails(ctx context.Context, userCred mcclient.TokenCredential, query jsonutils.JSONObject) bool {
	return userCred.IsAdminAllow(consts.GetServiceType(), self.GetModelManager().KeywordPlural(), policy.PolicyActionGet)
}

func (self *SHostJointsBase) AllowUpdateItem(ctx context.Context, userCred mcclient.TokenCredential) bool {
	return userCred.IsAdminAllow(consts.GetServiceType(), self.GetModelManager().KeywordPlural(), policy.PolicyActionUpdate)
}

func (self *SHostJointsBase) AllowDeleteItem(ctx context.Context, userCred mcclient.TokenCredential, query jsonutils.JSONObject, data jsonutils.JSONObject) bool {
	return userCred.IsAdminAllow(consts.GetServiceType(), self.GetModelManager().KeywordPlural(), policy.PolicyActionDelete)
}
