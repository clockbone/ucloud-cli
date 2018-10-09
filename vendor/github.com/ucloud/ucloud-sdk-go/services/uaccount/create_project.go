//Code is generated by ucloud code generator, don't modify it by hand, it will cause undefined behaviors.
//go:generate ucloud-gen-go-api UAccount CreateProject

package uaccount

import (
	"github.com/ucloud/ucloud-sdk-go/sdk"
	"github.com/ucloud/ucloud-sdk-go/sdk/request"
	"github.com/ucloud/ucloud-sdk-go/sdk/response"
)

// CreateProjectRequest is request schema for CreateProject action
type CreateProjectRequest struct {
	request.CommonBase

	// 项目名称
	ProjectName *string `required:"true"`

	// 项目父节点Id, 不填写创建顶层项目
	ParentId *string `required:"false"`
}

// CreateProjectResponse is response schema for CreateProject action
type CreateProjectResponse struct {
	response.CommonBase

	// 所创建项目的Id
	ProjectId string
}

// NewCreateProjectRequest will create request of CreateProject action.
func (c *UAccountClient) NewCreateProjectRequest() *CreateProjectRequest {
	cfg := c.client.GetConfig()

	return &CreateProjectRequest{
		CommonBase: request.CommonBase{
			Region:    sdk.String(cfg.Region),
			ProjectId: sdk.String(cfg.ProjectId),
		},
	}
}

// CreateProject - 创建项目
func (c *UAccountClient) CreateProject(req *CreateProjectRequest) (*CreateProjectResponse, error) {
	var err error
	var res CreateProjectResponse

	err = c.client.InvokeAction("CreateProject", req, &res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}