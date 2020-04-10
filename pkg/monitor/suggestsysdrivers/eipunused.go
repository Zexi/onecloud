package suggestsysdrivers

import (
	"context"
	"database/sql"
	"fmt"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
	"yunion.io/x/pkg/errors"

	"yunion.io/x/onecloud/pkg/apis/monitor"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/mcclient"
	"yunion.io/x/onecloud/pkg/mcclient/auth"
	"yunion.io/x/onecloud/pkg/mcclient/modules"
	"yunion.io/x/onecloud/pkg/monitor/models"
)

type EIPUnused struct {
	monitor.EIPUnused
}

func (_ *EIPUnused) GetType() string {
	return monitor.EIP_UN_USED
}

func NewEIPUsedDriver() models.ISuggestSysRuleDriver {
	return &EIPUnused{
		monitor.EIPUnused{
			Status: "",
		},
	}
}

func (dri *EIPUnused) ValidateSetting(input *monitor.SSuggestSysAlertSetting) error {
	if input.EIPUnused == nil {
		return errors.ErrNotFound
	}
	return nil
}

func (rule *EIPUnused) Run(instance *monitor.SSuggestSysAlertSetting) {
	oldAlert := make([]models.DSuggestSysAlert, 0)
	q := models.SuggestSysAlertManager.Query()
	q.Equals("type", monitor.EIP_UN_USED)
	err := db.FetchModelObjects(models.SuggestSysAlertManager, q, &oldAlert)
	if err != nil && err != sql.ErrNoRows {
		log.Errorln(errors.Wrap(err, "db.FetchModelObjects"))
		return
	}
	newAlert, err := rule.getEIPUnused(instance)
	if err != nil {
		log.Errorln(errors.Wrap(err, "getEIPUnused error"))
		return
	}

	DealAlertData(oldAlert, newAlert.Value())
}

func (rule *EIPUnused) getEIPUnused(instance *monitor.SSuggestSysAlertSetting) (*jsonutils.JSONArray, error) {
	//处理逻辑
	session := auth.GetAdminSession(context.Background(), "", "")
	query := jsonutils.NewDict()
	query.Add(jsonutils.NewString("scope"), "system")
	rtn, err := modules.Elasticips.List(session, query)
	if err != nil {
		return nil, err
	}
	EIPUnsedArr := jsonutils.NewArray()
	for _, row := range rtn.Data {
		if row.ContainsIgnoreCases("associate_type") || row.ContainsIgnoreCases("associate_id") {
			continue
		}
		alertData := jsonutils.DeepCopy(row).(*jsonutils.JSONDict)
		id, _ := alertData.GetString("id")
		alertData.Add(jsonutils.NewString(id), "res_id")
		alertData.Remove("id")

		input := &monitor.SSuggestSysAlertSetting{
			EIPUnused: &monitor.EIPUnused{Status: rule.Status},
		}
		alertData.Add(jsonutils.Marshal(input), "monitor_config")
		if instance != nil {
			alertData.Add(jsonutils.Marshal(instance), "monitor_config")
		}

		problem := jsonutils.NewDict()
		problem.Add(jsonutils.NewString(rule.GetType()), "eip")
		alertData.Add(problem, "problem")

		alertData.Add(jsonutils.NewString("释放未使用的EIP"), "suggest")
		alertData.Add(jsonutils.NewString(rule.GetType()), "type")

		alertData.Add(jsonutils.NewString(monitor.DRIVER_ACTION), "action")

		alertData.Add(row, "res_meta")
		EIPUnsedArr.Add(alertData)
	}
	return EIPUnsedArr, nil
}

func DealAlertData(oldAlerts []models.DSuggestSysAlert, newAlerts []jsonutils.JSONObject) {
	oldMap := make(map[string]models.DSuggestSysAlert, 0)
	for _, alert := range oldAlerts {
		oldMap[alert.ResId] = alert
	}

	for _, newAlert := range newAlerts {
		res_id, _ := newAlert.GetString("res_id")
		if oldAlert, ok := oldMap[res_id]; ok {
			//更新的alert
			_, err := db.Update(&oldAlert, func() error {
				err := newAlert.Unmarshal(&oldAlert)
				if err != nil {
					errMsg := fmt.Sprintf("unmarshal fail: %s", err)
					log.Errorf(errMsg)
				}
				return nil
			})
			if err != nil {
				log.Errorln("更新alert失败", err)
			}
			delete(oldMap, res_id)
		} else {
			//新增的alert
			adminCredential := auth.AdminCredential()
			_, err := db.DoCreate(models.SuggestSysAlertManager, context.Background(), adminCredential, nil, newAlert,
				adminCredential)
			if err != nil {
				log.Errorln(err)
			}
		}
	}

	for _, oldAlert := range oldMap {
		err := oldAlert.Delete(context.Background(), auth.AdminCredential())
		if err != nil {
			log.Errorln("删除旧alert数据失败", err)
		}
	}
}

//cronjob
func (rule *EIPUnused) DoSuggestSysRule(ctx context.Context, userCred mcclient.TokenCredential,
	isStart bool) {
	suggestSysSettingMap, err := models.SuggestSysRuleManager.FetchSuggestSysAlartSettings(rule.GetType())
	if err != nil {
		log.Errorln("DoSuggestSysRule error :", err)
		return
	}
	rule.Run(suggestSysSettingMap[rule.GetType()])
}
