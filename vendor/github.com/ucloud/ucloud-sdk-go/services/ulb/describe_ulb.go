//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api ULB DescribeULB

package ulb

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// DescribeULBRequest is request schema for DescribeULB action
type DescribeULBRequest struct {
	request.CommonBase

	// 数据偏移量，默认为0
	Offset *int `required:"false"`

	// 数据分页值，默认为20
	Limit *int `required:"false"`

	// 负载均衡实例的Id。 若指定则返回指定的负载均衡实例的信息； 若不指定则返回当前数据中心中所有的负载均衡实例的信息
	ULBId *string `required:"false"`

	// ULB所属的VPC
	VPCId *string `required:"false"`

	// ULB所属的子网ID
	SubnetId *string `required:"false"`

	// ULB所属的业务组ID
	BusinessId *string `required:"false"`
}

// DescribeULBResponse is response schema for DescribeULB action
type DescribeULBResponse struct {
	response.CommonBase

	// 满足条件的ULB总数
	TotalCount int

	// ULB列表，每项参数详见 ULBSet
	DataSet []ULBSet
}

// NewDescribeULBRequest will create request of DescribeULB action.
func (c *ULBClient) NewDescribeULBRequest() *DescribeULBRequest {
	req := &DescribeULBRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// DescribeULB - 获取ULB详细信息
func (c *ULBClient) DescribeULB(req *DescribeULBRequest) (*DescribeULBResponse, error) {
	var err error
	var res DescribeULBResponse

	err = c.client.InvokeAction("DescribeULB", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
