//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UDisk DeleteUDisk

package udisk

import (
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"github.com/ucloud/ucloud-sdk-go/ucloud/response"
)

// DeleteUDiskRequest is request schema for DeleteUDisk action
type DeleteUDiskRequest struct {
	request.CommonBase

	// 可用区。参见 [可用区列表](../summary/regionlist.html)
	Zone *string `required:"true"`

	// 要删除的UDisk的Id
	UDiskId *string `required:"true"`
}

// DeleteUDiskResponse is response schema for DeleteUDisk action
type DeleteUDiskResponse struct {
	response.CommonBase
}

// NewDeleteUDiskRequest will create request of DeleteUDisk action.
func (c *UDiskClient) NewDeleteUDiskRequest() *DeleteUDiskRequest {
	req := &DeleteUDiskRequest{}

	// setup request with client config
	c.client.SetupRequest(req)

	// setup retryable with default retry policy (retry for non-create action and common error)
	req.SetRetryable(true)
	return req
}

// DeleteUDisk - 删除UDisk
func (c *UDiskClient) DeleteUDisk(req *DeleteUDiskRequest) (*DeleteUDiskResponse, error) {
	var err error
	var res DeleteUDiskResponse

	err = c.client.InvokeAction("DeleteUDisk", req, &res)
	if err != nil {
		return &res, err
	}

	return &res, nil
}
