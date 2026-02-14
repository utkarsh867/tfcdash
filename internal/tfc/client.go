package tfc

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-tfe"
)

type Client struct {
	tfeClient *tfe.Client
	Org       string
}

func NewClient(token string) (*Client, error) {
	config := &tfe.Config{
		Token: token,
	}
	tfeClient, err := tfe.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Client{
		tfeClient: tfeClient,
	}, nil
}

func (c *Client) SetOrg(org string) {
	c.Org = org
}

func (c *Client) ListWorkspaces(ctx context.Context) ([]*tfe.Workspace, error) {
	options := &tfe.WorkspaceListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	}
	var allWorkspaces []*tfe.Workspace
	for {
		wl, err := c.tfeClient.Workspaces.List(ctx, c.Org, options)
		if err != nil {
			return nil, err
		}
		allWorkspaces = append(allWorkspaces, wl.Items...)
		if wl.CurrentPage >= wl.TotalPages {
			break
		}
		options.PageNumber = wl.NextPage
	}
	return allWorkspaces, nil
}

func (c *Client) ListRuns(ctx context.Context, workspaceID string) ([]*tfe.Run, error) {
	options := &tfe.RunListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 20,
		},
		Include: []tfe.RunIncludeOpt{tfe.RunPlan},
	}
	rl, err := c.tfeClient.Runs.List(ctx, workspaceID, options)
	if err != nil {
		return nil, err
	}
	return rl.Items, nil
}

func (c *Client) ApplyRun(ctx context.Context, runID string, comment string) error {
	options := tfe.RunApplyOptions{
		Comment: tfe.String(comment),
	}
	return c.tfeClient.Runs.Apply(ctx, runID, options)
}

func (c *Client) GetRun(ctx context.Context, runID string) (*tfe.Run, error) {
	return c.tfeClient.Runs.Read(ctx, runID)
}

func (c *Client) GetCurrentUser(ctx context.Context) (*tfe.User, error) {
	return c.tfeClient.Users.ReadCurrent(ctx)
}

func (c *Client) GetPlanJSONOutput(ctx context.Context, planID string) ([]byte, error) {
	urlBytes, err := c.tfeClient.Plans.ReadJSONOutput(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan JSON URL: %w", err)
	}
	return urlBytes, nil
}

func (c *Client) GetApplyLogs(ctx context.Context, applyID string) (string, error) {
	apply, err := c.tfeClient.Applies.Read(ctx, applyID)
	if err != nil {
		return "", err
	}
	return string(apply.Status), nil
}

func (c *Client) ListOrganizations(ctx context.Context) ([]*tfe.Organization, error) {
	options := &tfe.OrganizationListOptions{
		ListOptions: tfe.ListOptions{
			PageSize: 100,
		},
	}
	var allOrgs []*tfe.Organization
	for {
		orgList, err := c.tfeClient.Organizations.List(ctx, options)
		if err != nil {
			return nil, err
		}
		allOrgs = append(allOrgs, orgList.Items...)
		if orgList.CurrentPage >= orgList.TotalPages {
			break
		}
		options.PageNumber = orgList.NextPage
	}
	return allOrgs, nil
}
