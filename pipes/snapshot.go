package pipes

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"strings"

	"github.com/spf13/viper"
	"github.com/turbot/pipe-fittings/constants"
	"github.com/turbot/pipe-fittings/export"
	"github.com/turbot/pipe-fittings/sperr"
	"github.com/turbot/pipe-fittings/steampipeconfig"
	steampipecloud "github.com/turbot/pipes-sdk-go"
)

func PublishSnapshot(ctx context.Context, snapshot *steampipeconfig.SteampipeSnapshot, share bool) (string, error) {
	snapshotLocation := viper.GetString(constants.ArgSnapshotLocation)
	// snapshotLocation must be set (validation should ensure this)
	if snapshotLocation == "" {
		return "", sperr.New("to share a snapshot, snapshot-location must be set")
	}

	// if snapshot location is a workspace handle, upload it
	if steampipeconfig.IsPipesWorkspaceIdentifier(snapshotLocation) {
		url, err := uploadSnapshot(ctx, snapshot, share)
		if err != nil {
			return "", sperr.Wrap(err)
		}
		return fmt.Sprintf("\nSnapshot uploaded to %s\n", url), nil
	}

	// otherwise assume snapshot location is a file path
	filePath, err := exportSnapshot(snapshot)
	if err != nil {
		return "", sperr.Wrap(err)
	}
	return fmt.Sprintf("\nSnapshot saved to %s\n", filePath), nil
}

func exportSnapshot(snapshot *steampipeconfig.SteampipeSnapshot) (string, error) {
	exporter := &export.SnapshotExporter{}

	fileName := export.GenerateDefaultExportFileName(snapshot.FileNameRoot, exporter.FileExtension())
	dirName := viper.GetString(constants.ArgSnapshotLocation)
	filePath := path.Join(dirName, fileName)

	err := exporter.Export(context.Background(), snapshot, filePath)
	if err != nil {
		return "", sperr.Wrap(err)
	}
	return filePath, nil
}

func uploadSnapshot(ctx context.Context, snapshot *steampipeconfig.SteampipeSnapshot, share bool) (string, error) {
	client := newPipesClient(viper.GetString(constants.ArgPipesToken))

	cloudWorkspace := viper.GetString(constants.ArgSnapshotLocation)
	parts := strings.Split(cloudWorkspace, "/")
	if len(parts) != 2 {
		return "", sperr.New("failed to resolve username and workspace handle from workspace %s", cloudWorkspace)
	}
	identityHandle := parts[0]
	workspaceHandle := parts[1]

	// no determine whether this is a user or org workspace
	// get the identity
	identity, _, err := client.Identities.Get(ctx, identityHandle).Execute()
	if err != nil {
		return "", sperr.Wrap(err)
	}

	workspaceType := identity.Type

	// set the visibility
	visibility := "workspace"
	if share {
		visibility = "anyone_with_link"
	}

	// resolve the snapshot title
	title := resolveSnapshotTitle(snapshot)
	slog.Debug("Uploading snapshot", "title", title)
	// populate map of tags tags been set?
	tags := getTags()

	cloudSnapshot, err := snapshot.AsCloudSnapshot()
	if err != nil {
		return "", sperr.Wrap(err)
	}

	// strip verbose/sensitive fields
	err = steampipeconfig.StripSnapshot(cloudSnapshot)
	if err != nil {
		return "", sperr.Wrap(err)
	}

	req := steampipecloud.CreateWorkspaceSnapshotRequest{Data: *cloudSnapshot, Tags: tags, Visibility: &visibility}
	req.SetTitle(title)

	var uploadedSnapshot steampipecloud.WorkspaceSnapshot
	if identity.Type == "user" {
		uploadedSnapshot, _, err = client.UserWorkspaceSnapshots.Create(ctx, identityHandle, workspaceHandle).Request(req).Execute()
	} else {
		uploadedSnapshot, _, err = client.OrgWorkspaceSnapshots.Create(ctx, identityHandle, workspaceHandle).Request(req).Execute()
	}
	if err != nil {
		return "", sperr.Wrap(err)
	}

	snapshotId := uploadedSnapshot.Id
	snapshotUrl := fmt.Sprintf("https://%s/%s/%s/workspace/%s/snapshot/%s",
		viper.GetString(constants.ArgPipesHost),
		workspaceType,
		identityHandle,
		workspaceHandle,
		snapshotId)

	return snapshotUrl, nil
}

func resolveSnapshotTitle(snapshot *steampipeconfig.SteampipeSnapshot) string {
	if titleArg := viper.GetString(constants.ArgSnapshotTitle); titleArg != "" {
		return titleArg
	}
	// is there a title property set on the snapshot
	if snapshotTitle := snapshot.Title; snapshotTitle != "" {
		return snapshotTitle
	}
	// fall back to the fully qualified name of the root resource (which is also the FileNameRoot)
	return snapshot.FileNameRoot
}

func getTags() map[string]any {
	tags := viper.GetStringSlice(constants.ArgSnapshotTag)
	res := map[string]any{}

	for _, tagStr := range tags {
		parts := strings.Split(tagStr, "=")
		if len(parts) != 2 {
			continue
		}
		res[parts[0]] = parts[1]
	}
	return res
}
