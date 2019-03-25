package baremetal

import (
	"fmt"

	computeapi "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/compute/baremetal"
	"yunion.io/x/onecloud/pkg/scheduler/algorithm/predicates"
	"yunion.io/x/onecloud/pkg/scheduler/core"
)

type StoragePredicate struct {
	BasePredicate
}

func (p *StoragePredicate) Name() string {
	return "baremetal_storage"
}

func (p *StoragePredicate) Clone() core.FitPredicate {
	return &StoragePredicate{}
}

func toBaremetalDisks(disks []*computeapi.DiskConfig) []*baremetal.Disk {
	ret := make([]*baremetal.Disk, len(disks))
	for i, disk := range disks {
		ret[i] = &baremetal.Disk{
			Backend:    disk.Backend,
			ImageID:    disk.ImageId,
			Fs:         &disk.Fs,
			Format:     disk.Format,
			MountPoint: &disk.Mountpoint,
			Driver:     &disk.Driver,
			Cache:      &disk.Cache,
			Size:       int64(disk.SizeMb),
			Storage:    &disk.Storage,
		}
	}
	return ret
}

func (p *StoragePredicate) Execute(u *core.Unit, c core.Candidater) (bool, []core.PredicateFailureReason, error) {
	h := predicates.NewPredicateHelper(p, u, c)
	schedData := u.SchedData()

	candidate, err := h.BaremetalCandidate()
	if err != nil {
		return false, nil, err
	}

	layouts, err := baremetal.CalculateLayout(
		schedData.BaremetalDiskConfigs,
		candidate.StorageInfo,
	)

	if err == nil && baremetal.CheckDisksAllocable(layouts, toBaremetalDisks(schedData.Disks)) {
		h.SetCapacity(int64(1))
	} else {
		h.SetCapacity(int64(0))
		h.AppendPredicateFailMsg(fmt.Sprintf("%s err: %v", predicates.ErrNoEnoughStorage, err))
	}

	return h.GetResult()
}
