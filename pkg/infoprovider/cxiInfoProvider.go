// Copyright 2018 Intel Corp. All Rights Reserved.
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

package infoprovider

import (
	"os"
	"path/filepath"

	"github.com/golang/glog"
	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"

	"github.com/k8snetworkplumbingwg/sriov-network-device-plugin/pkg/types"
	"github.com/k8snetworkplumbingwg/sriov-network-device-plugin/pkg/utils"
)

// CxiDevDir is the base directory for CXI char devices.
// It is a variable to allow overriding the path in unit tests.
var CxiDevDir = "/dev"

/*
cxiInfoProvider provides the Cassini (CXI) char device information.
It mounts the /dev/cxiN char device associated with a CXI VF into the container.
*/
type cxiInfoProvider struct {
	pciAddr string
}

// NewCxiInfoProvider returns a new CXI Information Provider
func NewCxiInfoProvider(pciAddr string) types.DeviceInfoProvider {
	return &cxiInfoProvider{
		pciAddr: pciAddr,
	}
}

// *****************************************************************
/* DeviceInfoProvider Interface */

func (ip *cxiInfoProvider) GetName() string {
	return "cxi"
}

func (ip *cxiInfoProvider) GetDeviceSpecs() []*pluginapi.DeviceSpec {
	devSpecs := make([]*pluginapi.DeviceSpec, 0)

	cxiDev, err := utils.GetCxiDeviceFile(ip.pciAddr)
	if err != nil {
		glog.Errorf("GetDeviceSpecs(): error getting cxi device file for device: %s, %s", ip.pciAddr, err.Error())
		return devSpecs
	}

	// Confirm the corresponding /dev char device actually exists before mounting it.
	devPath := filepath.Join(CxiDevDir, filepath.Base(cxiDev))
	if _, err := os.Stat(devPath); err != nil {
		glog.Errorf("GetDeviceSpecs(): cxi device file %s does not exist for device: %s, %s", devPath, ip.pciAddr, err.Error())
		return devSpecs
	}

	devSpecs = append(devSpecs, &pluginapi.DeviceSpec{
		HostPath:      cxiDev,
		ContainerPath: cxiDev,
		Permissions:   "rw",
	})

	return devSpecs
}

func (ip *cxiInfoProvider) GetEnvVal() types.AdditionalInfo {
	envs := make(map[string]string, 0)

	cxiDev, err := utils.GetCxiDeviceFile(ip.pciAddr)
	if err != nil {
		glog.Errorf("GetEnvVal(): error getting cxi device file for device: %s, %s", ip.pciAddr, err.Error())
	} else {
		envs["dev-mount"] = cxiDev
	}

	return envs
}

func (ip *cxiInfoProvider) GetMounts() []*pluginapi.Mount {
	return nil
}

// *****************************************************************
