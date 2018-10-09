//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UNet DescribeShareBandwidth

package unet

import (
	"github.com/ucloud/ucloud-sdk-go/sdk"
	"github.com/ucloud/ucloud-sdk-go/sdk/request"
	"github.com/ucloud/ucloud-sdk-go/sdk/response"
)

// DescribeShareBandwidthRequest is request schema for DescribeShareBandwidth action
type DescribeShareBandwidthRequest struct {
	request.CommonBase

	// 需要返回的共享带宽Id
	ShareBandwidthIds []string `required:"false"`
}

// DescribeShareBandwidthResponse is response schema for DescribeShareBandwidth action
type DescribeShareBandwidthResponse struct {
	response.CommonBase

	// 共享带宽信息组 参见 UnetShareBandwidthSet
	DataSet []UnetShareBandwidthSet
}

// NewDescribeShareBandwidthRequest will create request of DescribeShareBandwidth action.
func (c *UNetClient) NewDescribeShareBandwidthRequest() *DescribeShareBandwidthRequest {
	cfg := c.client.GetConfig()

	return &DescribeShareBandwidthRequest{
		CommonBase: request.CommonBase{
			Region:    sdk.String(cfg.Region),
			ProjectId: sdk.String(cfg.ProjectId),
		},
	}
}

// DescribeShareBandwidth - 获取共享带宽信息
func (c *UNetClient) DescribeShareBandwidth(req *DescribeShareBandwidthRequest) (*DescribeShareBandwidthResponse, error) {
	var err error
	var res DescribeShareBandwidthResponse

	err = c.client.InvokeAction("DescribeShareBandwidth", req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}