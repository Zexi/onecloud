package db

import (
	"context"

	"github.com/yunionio/jsonutils"
	"github.com/yunionio/log"
	"github.com/yunionio/mcclient"
)

type SEnabledStatusStandaloneResourceBase struct {
	SStatusStandaloneResourceBase

	Enabled bool `nullable:"false" default:"false" list:"user"` // = Column(Boolean, nullable=False, default=False)
}

type SEnabledStatusStandaloneResourceBaseManager struct {
	SStatusStandaloneResourceBaseManager
}

func NewEnabledStatusStandaloneResourceBaseManager(dt interface{}, tableName string, keyword string, keywordPlural string) SEnabledStatusStandaloneResourceBaseManager {
	return SEnabledStatusStandaloneResourceBaseManager{SStatusStandaloneResourceBaseManager: NewStatusStandaloneResourceBaseManager(dt, tableName, keyword, keywordPlural)}
}

func (self *SEnabledStatusStandaloneResourceBase) AllowPerformEnable(ctx context.Context, userCred mcclient.TokenCredential, query jsonutils.JSONObject, data jsonutils.JSONObject) bool {
	return userCred.IsSystemAdmin()
}

func (self *SEnabledStatusStandaloneResourceBase) PerformEnable(ctx context.Context, userCred mcclient.TokenCredential, query jsonutils.JSONObject, data jsonutils.JSONObject) (jsonutils.JSONObject, error) {
	if !self.Enabled {
		_, err := self.GetModelManager().TableSpec().Update(self, func() error {
			self.Enabled = true
			return nil
		})
		if err != nil {
			log.Errorf("PerformEnable save update fail %s", err)
			return nil, err
		}
		OpsLog.LogEvent(self, ACT_ENABLE, "", userCred)
	}
	return nil, nil
}

func (self *SEnabledStatusStandaloneResourceBase) AllowPerformDisable(ctx context.Context, userCred mcclient.TokenCredential, query jsonutils.JSONObject, data jsonutils.JSONObject) bool {
	return userCred.IsSystemAdmin()
}

func (self *SEnabledStatusStandaloneResourceBase) PerformDisable(ctx context.Context, userCred mcclient.TokenCredential, query jsonutils.JSONObject, data jsonutils.JSONObject) (jsonutils.JSONObject, error) {
	if self.Enabled {
		_, err := self.GetModelManager().TableSpec().Update(self, func() error {
			self.Enabled = false
			return nil
		})
		if err != nil {
			log.Errorf("PerformDisable save update fail %s", err)
			return nil, err
		}
		OpsLog.LogEvent(self, ACT_DISABLE, "", userCred)
	}
	return nil, nil
}
