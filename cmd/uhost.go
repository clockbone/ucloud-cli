// Copyright © 2018 NAME HERE tony.li@ucloud.cn
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

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/ucloud/ucloud-sdk-go/services/uhost"
	sdk "github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/helpers/waiter"
	"github.com/ucloud/ucloud-sdk-go/ucloud/log"

	. "github.com/ucloud/ucloud-cli/base"
	"github.com/ucloud/ucloud-cli/ux"
)

//NewCmdUHost ucloud uhost
func NewCmdUHost() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uhost",
		Short: "List,create,delete,stop,restart,poweroff or scale UHost instance",
		Long:  `List,create,delete,stop,restart,poweroff or scale UHost instance`,
		Args:  cobra.NoArgs,
	}
	cmd.AddCommand(NewCmdUHostList())
	cmd.AddCommand(NewCmdUHostCreate())
	cmd.AddCommand(NewCmdUHostDelete())
	cmd.AddCommand(NewCmdUHostStop())
	cmd.AddCommand(NewCmdUHostStart())
	cmd.AddCommand(NewCmdUHostReboot())
	cmd.AddCommand(NewCmdUHostPoweroff())
	cmd.AddCommand(NewCmdUHostScale())

	return cmd
}

//UHostRow UHost表格行
type UHostRow struct {
	UHostName      string
	ResourceID     string
	UGroup         string
	ClassicNetwork string
	Config         string
	Type           string
	CreationTime   string
	State          string
}

//NewCmdUHostList [ucloud uhost list]
func NewCmdUHostList() *cobra.Command {
	req := BizClient.NewDescribeUHostInstanceRequest()
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all UHost Instances",
		Long:  `List all UHost Instances`,
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := BizClient.DescribeUHostInstance(req)
			if err != nil {
				HandleError(err)
				return
			}
			if global.json {
				PrintJSON(resp.UHostSet)
			} else {
				list := make([]UHostRow, 0)
				for _, host := range resp.UHostSet {
					row := UHostRow{}
					row.UHostName = host.Name
					row.ResourceID = host.UHostId
					row.UGroup = host.Tag
					for _, ip := range host.IPSet {
						if row.ClassicNetwork != "" {
							row.ClassicNetwork += " | "
						}
						if ip.Type == "Private" {
							row.ClassicNetwork += fmt.Sprintf("%s", ip.IP)
						} else {
							row.ClassicNetwork += fmt.Sprintf("%s %s", ip.IP, ip.Type)
						}
					}
					osName := strings.SplitN(host.OsName, " ", 2)
					cupCore := host.CPU
					memorySize := host.Memory / 1024
					diskSize := 0
					for _, disk := range host.DiskSet {
						if disk.Type == "Data" {
							diskSize += disk.Size
						}
					}
					row.Config = fmt.Sprintf("%s cpu:%d memory:%dG disk:%dG", osName[0], cupCore, memorySize, diskSize)
					row.CreationTime = FormatDate(host.CreateTime)
					row.State = host.State
					row.Type = host.UHostType + "/" + host.HostType
					list = append(list, row)
				}
				PrintTable(list, []string{"UHostName", "ResourceID", "UGroup", "ClassicNetwork", "Config", "Type", "CreationTime", "State"})
			}
		},
	}
	cmd.Flags().SortFlags = false
	req.ProjectId = cmd.Flags().String("project-id", ConfigInstance.ProjectID, "Optional. Assign project-id")
	req.Region = cmd.Flags().String("region", ConfigInstance.Region, "Optional. Assign region")
	req.Zone = cmd.Flags().String("zone", "", "Optional. Assign availability zone")
	cmd.Flags().StringSliceVar(&req.UHostIds, "resource-id", make([]string, 0), "Optional. UHost Instance ID, multiple values separated by comma(without space)")
	req.Tag = cmd.Flags().String("ugroup", "", "Optional. UGroup")
	req.Offset = cmd.Flags().Int("offset", 0, "Optional. Offset default 0")
	req.Limit = cmd.Flags().Int("limit", 20, "Optional. Limit default 20, max value 100")

	return cmd
}

//NewCmdUHostCreate [ucloud uhost create]
func NewCmdUHostCreate() *cobra.Command {
	req := BizClient.NewCreateUHostInstanceRequest()
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create UHost instance",
		Long:  "Create UHost instance",
		Run: func(cmd *cobra.Command, args []string) {
			*req.Memory *= 1024
			req.LoginMode = sdk.String("Password")
			resp, err := BizClient.CreateUHostInstance(req)
			if err != nil {
				HandleError(err)
				return
			}
			Cxt.Printf("UHost:%v created successfully!\n", resp.UHostIds)
		},
	}

	n1Zone := map[string]bool{
		"cn-bj2-01": true,
		"cn-bj2-03": true,
		"cn-sh2-01": true,
		"hk-01":     true,
	}
	defaultUhostType := "N2"
	if _, ok := n1Zone[ConfigInstance.Zone]; ok {
		defaultUhostType = "N1"
	}

	req.Disks = make([]uhost.UHostDisk, 2)
	req.Disks[0].IsBoot = sdk.Bool(true)
	req.Disks[1].IsBoot = sdk.Bool(false)

	flags := cmd.Flags()
	flags.SortFlags = false
	req.CPU = flags.Int("cpu", 4, "Required. The count of CPU cores. Optional parameters: {1, 2, 4, 8, 12, 16, 24, 32}")
	req.Memory = flags.Int("memory", 8, "Required. Memory size. Unit: GB. Range: [1, 128], multiple of 2")
	req.Password = flags.String("password", "", "Required. Password of the uhost user(root/ubuntu)")
	req.ImageId = flags.String("image-id", "", "Required. The ID of image. see 'ucloud image list'")
	req.VPCId = flags.String("vpc-id", "", "Optional. VPC ID. This field is required under VPC2.0. See 'ucloud vpc list'")
	req.SubnetId = flags.String("subnet-id", "", "Optional. Subnet ID. This field is required under VPC2.0. See 'ucloud subnet list'")
	req.Name = flags.String("name", "UHost", "Optional. UHost instance name")
	req.ChargeType = flags.String("charge-type", "Month", "Optional.'Year',pay yearly;'Month',pay monthly;'Dynamic', pay hourly(requires access)")
	req.Quantity = flags.Int("quantity", 1, "Optional. The duration of the instance. N years/months.")
	req.ProjectId = flags.String("project-id", ConfigInstance.ProjectID, "Optional. Assign project-id")
	req.Region = flags.String("region", ConfigInstance.Region, "Optional. Assign region")
	req.Zone = flags.String("zone", ConfigInstance.Zone, "Optional. Assign availability zone")
	req.UHostType = flags.String("type", defaultUhostType, "Optional. Default is 'N2' of which cpu is V4 and sata disk. also support 'N1' means V3 cpu and sata disk;'I2' means V4 cpu and ssd disk;'D1' means big data model;'G1' means GPU type, model for K80;'G2' model for P40; 'G3' model for V100")
	req.NetCapability = flags.String("net-capability", "Normal", "Optional. Default is 'Normal', also support 'Super' which will enhance multiple times network capability as before")
	req.Disks[0].Type = flags.String("boot-disk-type", "LOCAL_NORMAL", "Optional. Enumeration value. 'LOCAL_NORMAL', Ordinary local disk; 'CLOUD_NORMAL', Ordinary cloud disk; 'LOCAL_SSD',local ssd disk; 'CLOUD_SSD',cloud ssd disk; 'EXCLUSIVE_LOCAL_DISK',big data. The disk only supports a limited combination.")
	req.Disks[0].Size = flags.String("boot-disk-size", "20", "Optional. Default 20G. Windows should be bigger than 40G Unit GB")
	req.Disks[0].BackupType = flags.String("boot-disk-backup-type", "NONE", "Optional. Enumeration value, 'NONE' or 'DATAARK'. DataArk supports real-time backup, which can restore the disk back to any moment within the last 12 hours. (Normal Local Disk and Normal Cloud Disk Only)")
	req.Disks[1].Type = flags.String("data-disk-type", "LOCAL_NORMAL", "Optional. Enumeration value. 'LOCAL_NORMAL', Ordinary local disk; 'CLOUD_NORMAL', Ordinary cloud disk; 'LOCAL_SSD',local ssd disk; 'CLOUD_SSD',cloud ssd disk; 'EXCLUSIVE_LOCAL_DISK',big data. The disk only supports a limited combination.")
	req.Disks[1].Size = flags.String("data-disk-size", "20", "Optional. Disk size. Unit GB")
	req.Disks[1].BackupType = flags.String("data-disk-backup-type", "NONE", "Optional. Enumeration value, 'NONE' or 'DATAARK'. DataArk supports real-time backup, which can restore the disk back to any moment within the last 12 hours. (Normal Local Disk and Normal Cloud Disk Only)")
	req.NetworkId = flags.String("network-id", "", "Optional. Network ID (no need to fill in the case of VPC2.0). In the case of VPC1.0, if not filled in, we will choose the basic network; if it is filled in, we will choose the subnet. See DescribeSubnet.")
	req.SecurityGroupId = flags.String("firewall-id", "", "Optional. Firewall Id, default: Web recommended firewall. see 'ucloud firewall list'.")
	req.Tag = flags.String("ugroup", "Default", "Optional. Business group")
	req.CouponId = flags.String("coupon-id", "", "Optional. Coupon ID, The Coupon can deduct part of the payment,see DescribeCoupon or https://accountv2.ucloud.cn")

	cmd.Flags().SetFlagValues("charge-type", "Month", "Year", "Dynamic", "Trial")
	cmd.Flags().SetFlagValues("cpu", "1", "2", "4", "8", "12", "16", "24", "32")
	cmd.Flags().SetFlagValues("type", "N2", "N1", "I2", "D1", "G1", "G2", "G3")
	cmd.Flags().SetFlagValues("net-capability", "Normal", "Super")
	cmd.Flags().SetFlagValues("boot-disk-type", "LOCAL_NORMAL", "CLOUD_NORMAL", "LOCAL_SSD", "CLOUD_SSD", "EXCLUSIVE_LOCAL_DISK")
	cmd.Flags().SetFlagValues("boot-disk-backup-type", "NONE", "DATAARK")
	cmd.Flags().SetFlagValues("data-disk-type", "LOCAL_NORMAL", "CLOUD_NORMAL", "LOCAL_SSD", "CLOUD_SSD", "EXCLUSIVE_LOCAL_DISK")
	cmd.Flags().SetFlagValues("data-disk-backup-type", "NONE", "DATAARK")

	cmd.Flags().SetFlagValuesFunc("image-id", func() []string {
		req := BizClient.NewDescribeImageRequest()
		projectID, _ := flags.GetString("project-id")
		if projectID == "" {
			projectID = ConfigInstance.ProjectID
		}
		req.ProjectId = sdk.String(projectID)

		region, _ := flags.GetString("region")
		if region == "" {
			region = ConfigInstance.Region
		}
		req.Region = sdk.String(region)

		zone, _ := flags.GetString("zone")
		if zone == "" {
			zone = ConfigInstance.Zone
		}
		req.Zone = sdk.String(zone)
		req.ImageType = sdk.String("Base")
		result := make([]string, 0)
		resp, err := BizClient.DescribeImage(req)
		if err == nil {
			for _, image := range resp.ImageSet {
				result = append(result, image.ImageId)
			}
		}
		return result
	})

	cmd.MarkFlagRequired("cpu")
	cmd.MarkFlagRequired("memory")
	cmd.MarkFlagRequired("password")
	cmd.MarkFlagRequired("image-id")

	return cmd
}

//NewCmdUHostDelete ucloud uhost delete
func NewCmdUHostDelete() *cobra.Command {
	isDestory := sdk.Bool(false)
	isEipReleased := sdk.Bool(false)

	req := BizClient.NewTerminateUHostInstanceRequest()
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete Uhost instance",
		Long:  "Delete Uhost instance",
		Run: func(cmd *cobra.Command, args []string) {
			if *isDestory {
				req.Destroy = sdk.Int(1)
			} else {
				req.Destroy = sdk.Int(0)
			}
			if *isEipReleased {
				req.EIPReleased = sdk.String("yes")
			} else {
				req.EIPReleased = sdk.String("no")
			}
			hostIns, err := describeUHostByID(*req.UHostId, *req.ProjectId, *req.Region, *req.Zone)
			if err != nil {
				HandleError(err)
			} else if hostIns != nil {
				if hostIns.State == "Running" {
					_req := BizClient.NewStopUHostInstanceRequest()
					_req.ProjectId = req.ProjectId
					_req.Region = req.Region
					_req.Zone = req.Zone
					_req.UHostId = req.UHostId
					stopUhostIns(_req)
				}
			}
			resp, err := BizClient.TerminateUHostInstance(req)
			if err != nil {
				HandleError(err)
			} else {
				Cxt.Printf("UHost:[%v] deleted successfully!\n", resp.UHostId)
			}
		},
	}

	req.ProjectId = cmd.Flags().String("project-id", ConfigInstance.ProjectID, "Assign project-id")
	req.Region = cmd.Flags().String("region", ConfigInstance.Region, "Assign region")
	req.Zone = cmd.Flags().String("zone", "", "availability zone")
	req.UHostId = cmd.Flags().String("resource-id", "", "ResourceID of the uhost instance( or UHostId)")
	isDestory = cmd.Flags().Bool("destory", false, "false,the uhost instance will be thrown to UHost recycle If you have permission; true,the uhost instance will be deleted directly")
	isEipReleased = cmd.Flags().Bool("eip-released", false, "false,Unbind EIP only; true, Unbind EIP and release it")
	cmd.Flags().SetFlagValues("destory", "true", "false")
	cmd.Flags().SetFlagValues("eip-released", "true", "false")
	cmd.MarkFlagRequired("resource-id")

	return cmd
}

//NewCmdUHostStop ucloud uhost stop
func NewCmdUHostStop() *cobra.Command {
	req := BizClient.NewStopUHostInstanceRequest()
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Shut down uhost instance",
		Long:  "Shut down uhost instance",
		Run: func(cmd *cobra.Command, args []string) {
			stopUhostIns(req)
		},
	}
	req.ProjectId = cmd.Flags().String("project-id", ConfigInstance.ProjectID, "Assign project-id")
	req.Region = cmd.Flags().String("region", ConfigInstance.Region, "Assign region")
	req.Zone = cmd.Flags().String("zone", "", "Assign availability zone")
	req.UHostId = cmd.Flags().String("resource-id", "", "ResourceID of the uhost instance( or UHostId)")
	cmd.MarkFlagRequired("resource-id")

	return cmd
}

func stopUhostIns(req *uhost.StopUHostInstanceRequest) {
	resp, err := BizClient.StopUHostInstance(req)
	if err != nil {
		HandleError(err)
	} else {
		text := fmt.Sprintf("UHost:[%v] is shuting down", resp.UhostId)
		done := poll(resp.UhostId, *req.ProjectId, *req.Region, *req.Zone, "Stopped")
		ux.DotSpinner.Start(text)
		<-done
		ux.DotSpinner.Stop()
	}
}

//NewCmdUHostStart ucloud uhost start
func NewCmdUHostStart() *cobra.Command {
	req := BizClient.NewStartUHostInstanceRequest()
	cmd := &cobra.Command{
		Use:     "start",
		Short:   "Start Uhost instance",
		Long:    "Start Uhost instance",
		Example: "ucloud uhost start --resource-id uhost-xxx",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := BizClient.StartUHostInstance(req)
			if err != nil {
				HandleError(err)
			} else {
				text := fmt.Sprintf("UHost:[%v] is starting", resp.UhostId)
				done := poll(resp.UhostId, *req.ProjectId, *req.Region, *req.Zone, "Running")
				ux.DotSpinner.Start(text)
				<-done
				ux.DotSpinner.Stop()
			}
		},
	}
	req.ProjectId = cmd.Flags().String("project-id", ConfigInstance.ProjectID, "Assign project-id")
	req.Region = cmd.Flags().String("region", ConfigInstance.Region, "Assign region")
	req.Zone = cmd.Flags().String("zone", "", "Assign availability zone")
	req.UHostId = cmd.Flags().String("resource-id", "", "ResourceID of the uhost instance( or UHostId)")
	req.DiskPassword = cmd.Flags().String("disk-password", "", "Encrypted disk password")
	cmd.MarkFlagRequired("resource-id")
	return cmd
}

//NewCmdUHostReboot ucloud uhost restart
func NewCmdUHostReboot() *cobra.Command {
	req := BizClient.NewRebootUHostInstanceRequest()
	cmd := &cobra.Command{
		Use:     "restart",
		Short:   "Restart/reboot Uhost instance",
		Long:    "Restart/reboot Uhost instance",
		Example: "ucloud uhost restart --resource-id uhost-xxx",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := BizClient.RebootUHostInstance(req)
			if err != nil {
				HandleError(err)
			} else {
				text := fmt.Sprintf("UHost:[%v] is restarting", resp.UhostId)
				done := poll(resp.UhostId, *req.ProjectId, *req.Region, *req.Zone, "Running")
				ux.DotSpinner.Start(text)
				<-done
				ux.DotSpinner.Stop()
			}
		},
	}
	req.ProjectId = cmd.Flags().String("project-id", ConfigInstance.ProjectID, "Assign project-id")
	req.Region = cmd.Flags().String("region", ConfigInstance.Region, "Assign region")
	req.Zone = cmd.Flags().String("zone", "", "Assign availability zone")
	req.UHostId = cmd.Flags().String("resource-id", "", "ResourceID of the uhost instance( or UHostId)")
	req.DiskPassword = cmd.Flags().String("disk-password", "", "Encrypted disk password")
	cmd.MarkFlagRequired("resource-id")
	return cmd
}

//NewCmdUHostPoweroff ucloud uhost poweroff
func NewCmdUHostPoweroff() *cobra.Command {
	req := BizClient.NewPoweroffUHostInstanceRequest()
	cmd := &cobra.Command{
		Use:   "poweroff",
		Short: "Analog power off Uhost instnace. Danger! this operation may affect data integrity or cause file system corruption",
		Long:  "Analog power off Uhost instnace. Danger! this operation may affect data integrity or cause file system corruption",
		Run: func(cmd *cobra.Command, args []string) {
			resp, err := BizClient.PoweroffUHostInstance(req)
			if err != nil {
				HandleError(err)
			} else {
				Cxt.Printf("UHost:[%v] is power off\n", resp.UhostId)
			}
		},
	}
	req.ProjectId = cmd.Flags().String("project-id", ConfigInstance.ProjectID, "Assign project-id")
	req.Region = cmd.Flags().String("region", ConfigInstance.Region, "Assign region")
	req.Zone = cmd.Flags().String("zone", "", "Assign availability zone")
	req.UHostId = cmd.Flags().String("resource-id", "", "ResourceID of the uhost instance( or UHostId)")
	return cmd
}

//NewCmdUHostScale ucloud uhost scale
func NewCmdUHostScale() *cobra.Command {
	req := BizClient.NewResizeUHostInstanceRequest()
	cmd := &cobra.Command{
		Use:   "scale",
		Short: "Scale uhost instance,such as cpu core count, memory size, system disk size and data disk size",
		Long:  "Scale uhost instance,such as cpu core count, memory size, system disk size and data disk size",
		Run: func(cmd *cobra.Command, args []string) {
			if *req.CPU == 0 {
				req.CPU = nil
			}
			if *req.Memory == 0 {
				req.Memory = nil
			} else {
				*req.Memory *= 1024
			}
			if *req.DiskSpace == 0 {
				req.DiskSpace = nil
			}
			if *req.BootDiskSpace == 0 {
				req.BootDiskSpace = nil
			}
			resp, err := BizClient.ResizeUHostInstance(req)
			if err != nil {
				HandleError(err)
			} else {
				Cxt.Printf("UHost:[%v] scaled\n", resp.UhostId)
			}
		},
	}
	req.ProjectId = cmd.Flags().String("project-id", ConfigInstance.ProjectID, "Assign project-id")
	req.Region = cmd.Flags().String("region", ConfigInstance.Region, "Assign region")
	req.Zone = cmd.Flags().String("zone", "", "Assign availability zone")
	req.UHostId = cmd.Flags().String("resource-id", "", "ResourceID of the uhost instance( or UHostId)")
	req.CPU = cmd.Flags().Int("cpu", 0, "The number of virtual CPU cores. Series1 {1, 2, 4, 8, 12, 16, 24, 32}. Series2 {1,2,4,8,16}")
	req.Memory = cmd.Flags().Int("memory", 0, "memory size. Unit: GB. Range: [1, 128], multiple of 2")
	req.DiskSpace = cmd.Flags().Int("data-disk-size", 0, "Data disk size,unit GB. Range[10,1000], SSD disk range[100,500]. Step 10")
	req.BootDiskSpace = cmd.Flags().Int("system-disk-size", 0, "System disk size, unit GB. Range[20,100]. Step 10. System disk does not support shrinkage")
	req.NetCapValue = cmd.Flags().Int("net-cap", 0, "NIC scale. 1,upgrade; 2,downgrade; 0,unchanged")
	cmd.MarkFlagRequired("resource-id")
	return cmd
}

func poll(uhostID, projectID, region, zone string, targetState string) chan bool {
	w := waiter.StateWaiter{
		Pending: []string{"pending"},
		Target:  []string{"avaliable"},
		Refresh: func() (interface{}, string, error) {
			inst, err := describeUHostByID(uhostID, projectID, region, zone)
			if err != nil {
				return nil, "", err
			}

			if inst == nil || inst.State != targetState {
				return nil, "pending", nil
			}

			return inst, "avaliable", nil
		},
		Timeout: 5 * time.Minute,
	}

	done := make(chan bool)
	go func() {
		if resp, err := w.Wait(); err != nil {
			log.Error(err)
		} else {
			log.Infof("%#v", resp)
		}
		done <- true
	}()
	return done
}

func describeUHostByID(uhostID, projectID, region, zone string) (*uhost.UHostInstanceSet, error) {
	req := BizClient.NewDescribeUHostInstanceRequest()
	req.UHostIds = []string{uhostID}
	req.ProjectId = &projectID
	req.Region = &region
	req.Zone = &zone

	resp, err := BizClient.DescribeUHostInstance(req)
	if err != nil {
		return nil, err
	}
	if len(resp.UHostSet) < 1 {
		return nil, nil
	}

	return &resp.UHostSet[0], nil
}
