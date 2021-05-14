// Copyright 2021 Yunion
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

package regions

import (
	"github.com/spf13/cobra"

	"yunion.io/x/onecloud/pkg/cloudmux/cloudprovider"
	"yunion.io/x/onecloud/pkg/cloudmux/cmcli/util/cmd"
	"yunion.io/x/onecloud/pkg/util/printutils"
)

func NewCmdRegionsList(f cmd.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List virtual regions",
		Run: f.RunWithProvider(func(provider cloudprovider.IProvider, _ []string) error {
			regions, err := provider.Compute().Regions().List()
			if err != nil {
				return err
			}
			printutils.PrintInterfaceList(regions, 0, 0, 0, nil)
			return nil
		}),
	}
	return cmd
}
