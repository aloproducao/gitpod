// Copyright (c) 2020 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package workspace_test

import (
	"context"
	"testing"
	"time"

	"github.com/gitpod-io/gitpod/test/pkg/integration"
	agent "github.com/gitpod-io/gitpod/test/tests/workspace/workspace_agent/api"
)

func TestWorkspaceInstrumentation(t *testing.T) {
	it := integration.NewTest(t)
	defer it.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	nfo, stopWs := integration.LaunchWorkspaceFromContextURL(it, ctx, "github.com/gitpod-io/gitpod")
	defer stopWs(true)

	rsa, err := it.Instrument(integration.ComponentWorkspace, "workspace", integration.WithInstanceID(nfo.LatestInstance.ID))
	if err != nil {
		t.Fatal(err)
	}
	defer rsa.Close()

	var ls agent.ListDirResponse
	err = rsa.Call("WorkspaceAgent.ListDir", &agent.ListDirRequest{
		Dir: "/workspace/gitpod",
	}, &ls)
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range ls.Files {
		t.Log(f)
	}
}

func TestLaunchWorkspaceDirectly(t *testing.T) {
	it := integration.NewTest(t)
	defer it.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	nfo := integration.LaunchWorkspaceDirectly(it, ctx)
	defer integration.DeleteWorkspace(it, ctx, nfo.Req.Id)
}
