// Copyright 2019 Yunion
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package suggestsysdrivers

import (
	"database/sql"
	"fmt"
	"time"

	"yunion.io/x/log"

	"yunion.io/x/onecloud/pkg/cloudcommon/cronman"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/monitor/models"
)

func init() {
	models.RegisterSuggestSysRuleDrivers(NewEIPUsedDriver())
}

func InitSuggestSysRuleCronjob() {
	objs := make([]models.SSuggestSysRule, 0)
	q := models.SuggestSysRuleManager.Query()
	if q == nil {
		fmt.Println(" query is nil")
	}
	err := db.FetchModelObjects(models.SuggestSysRuleManager, q, &objs)
	if err != nil && err != sql.ErrNoRows {
		log.Errorln("InitSuggestSysRuleCronjob db.FetchModelObjects error")
	}

	for _, suggestSysRuleConfig := range objs {
		if suggestSysRuleConfig.Enabled.Bool() {
			dur, _ := time.ParseDuration(suggestSysRuleConfig.Period)
			cronman.GetCronJobManager().AddJobAtIntervalsWithStartRun(suggestSysRuleConfig.Type, dur,
				models.GetSuggestSysRuleDrivers()[suggestSysRuleConfig.Type].DoSuggestSysRule, true)
		}
	}
	if len(objs) == 0 {
		for _, driver := range models.GetSuggestSysRuleDrivers() {
			cronman.GetCronJobManager().AddJobAtIntervalsWithStartRun(driver.GetType(), time.Duration(30)*time.Second,
				driver.DoSuggestSysRule, true)
		}
	}
}
