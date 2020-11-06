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

package tasks

import (
	"context"

	"yunion.io/x/jsonutils"

	"yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/cloudcommon/db"
	"yunion.io/x/onecloud/pkg/cloudcommon/db/taskman"
	"yunion.io/x/onecloud/pkg/cloudcommon/notifyclient"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/util/logclient"
)

type InstanceSnapshotResetTask struct {
	taskman.STask
}

func init() {
	taskman.RegisterTask(InstanceSnapshotResetTask{})
}

func (self *InstanceSnapshotResetTask) taskFail(
	ctx context.Context, isp *models.SInstanceSnapshot, guest *models.SGuest, reason jsonutils.JSONObject) {

	if guest == nil {
		guest = models.GuestManager.FetchGuestById(isp.GuestId)
	}
	guest.SetStatus(self.UserCred, compute.VM_SNAPSHOT_RESET_FAILED, reason.String())

	db.OpsLog.LogEvent(guest, db.ACT_VM_RESET_SNAPSHOT_FAILED, reason, self.UserCred)
	logclient.AddActionLogWithStartable(self, guest, logclient.ACT_VM_RESET, reason, self.UserCred, false)
	notifyclient.NotifySystemErrorWithCtx(ctx, guest.GetId(), isp.Name, compute.VM_SNAPSHOT_RESET_FAILED, reason.String())
	self.SetStageFailed(ctx, reason)
}

func (self *InstanceSnapshotResetTask) taskComplete(
	ctx context.Context, isp *models.SInstanceSnapshot, guest *models.SGuest, data jsonutils.JSONObject) {

	if guest == nil {
		guest = models.GuestManager.FetchGuestById(isp.GuestId)
	}
	guest.StartSyncstatus(ctx, self.UserCred, "")

	db.OpsLog.LogEvent(isp, db.ACT_VM_RESET_SNAPSHOT, "instance snapshot reset success", self.UserCred)
	logclient.AddActionLogWithStartable(self, guest, logclient.ACT_VM_RESET, false, self.UserCred, true)
	self.SetStageComplete(ctx, nil)
}

func (self *InstanceSnapshotResetTask) OnInit(
	ctx context.Context, obj db.IStandaloneModel, data jsonutils.JSONObject) {

	isp := obj.(*models.SInstanceSnapshot)
	guest := models.GuestManager.FetchGuestById(isp.GuestId)

	self.SetStage("OnInstanceSnapshotReset", nil)
	params := jsonutils.NewDict()
	params.Set("disk_index", jsonutils.NewInt(0))
	if err := isp.GetRegionDriver().RequestResetToInstanceSnapshot(ctx, guest, isp, self, params); err != nil {
		self.taskFail(ctx, isp, guest, jsonutils.NewString(err.Error()))
		return
	}
}

func (self *InstanceSnapshotResetTask) OnKvmDiskReset(
	ctx context.Context, isp *models.SInstanceSnapshot, data jsonutils.JSONObject) {

	guest := models.GuestManager.FetchGuestById(isp.GuestId)

	diskIndex, err := self.Params.Int("disk_index")
	if err != nil {
		self.taskFail(ctx, isp, guest, jsonutils.NewString(err.Error()))
		return
	}
	params := jsonutils.NewDict()
	params.Set("disk_index", jsonutils.NewInt(diskIndex+1))
	if err := isp.GetRegionDriver().RequestResetToInstanceSnapshot(ctx, guest, isp, self, params); err != nil {
		self.taskFail(ctx, isp, guest, jsonutils.NewString(err.Error()))
		return
	}
}

func (self *InstanceSnapshotResetTask) OnKvmDiskResetFailed(
	ctx context.Context, isp *models.SInstanceSnapshot, data jsonutils.JSONObject) {
	self.taskFail(ctx, isp, nil, data)
}

func (self *InstanceSnapshotResetTask) OnInstanceSnapshotReset(ctx context.Context, isp *models.SInstanceSnapshot, data jsonutils.JSONObject) {
	guest, _ := isp.GetGuest()
	self.taskComplete(ctx, isp, guest, data)
}

func (self *InstanceSnapshotResetTask) OnInstanceSnapshotResetFailed(ctx context.Context, isp *models.SInstanceSnapshot, data jsonutils.JSONObject) {
	self.taskFail(ctx, isp, nil, data)
}
