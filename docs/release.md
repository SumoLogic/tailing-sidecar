# Tailing Sidecar Release Instruction

## Automated Release Process

Releases are automated using GitHub Actions. Only two manual steps are required:

1. Trigger the **Prepare Release** workflow
2. Update the CHANGELOG, review, and merge the resulting PR

### Prerequisites

- A repository secret `GH_RELEASE_PAT` containing a PAT with `repo` scope.
  This is required because tags pushed with `GITHUB_TOKEN` do not trigger
  downstream workflows.
- The `release` label is auto-created by the prepare workflow if missing.

### Steps

#### 1. Trigger the Prepare Release Workflow

Go to **Actions > Prepare Release** and click **Run workflow**.
Enter the version in semver format without a `v` prefix (e.g. `0.21.0`).

The workflow will:

- Validate the version format
- Verify the tag does not already exist
- Create branch `prepare-release-{VERSION}`
- Update `version` and `appVersion` in
  `helm/tailing-sidecar-operator/Chart.yaml`
- Auto-generate a CHANGELOG entry from merged PRs since the last tag
- Open a PR to `main` titled `feat: prepare release v{VERSION}`
  with the `release` label

#### 2. Review and Merge

The PR includes an auto-generated `CHANGELOG.md` entry built from merged
PRs since the last tag. Review and edit it if needed, then merge the PR.

#### 3. Automatic Steps After Merge

After the release PR is merged, the **Finalize Release** workflow:

1. Reads the version from `Chart.yaml`
2. Creates an annotated tag `v{VERSION}` and pushes it
3. Creates a `release-v{MAJOR}.{MINOR}` branch (for new minor versions)
4. Creates a **draft GitHub Release** with auto-generated notes

The tag push triggers the **Releases Otelcol** workflow which:

- Builds and pushes Docker images to GHCR, ECR, and Docker Hub
- Publishes the Helm chart to GitHub Pages

#### 4. Publish the Release

Once builds pass, review the draft release at
<https://github.com/SumoLogic/tailing-sidecar/releases>, edit the notes
if needed, and click **Publish release**.

## Release Branch Policy

A long-lived `release-v{MAJOR}.{MINOR}` branch is created automatically
for each new minor version. Patch releases on an existing minor version
do not create a new branch.

## Manual Release (Fallback)

If the automated process is unavailable, follow these manual steps:

1. Create and merge a PR updating `version` and `appVersion` in
   [Chart.yaml](../helm/tailing-sidecar-operator/Chart.yaml).

2. Create and push the release tag:

   ```bash
   export VERSION=0.21.0
   git checkout main
   git pull
   git tag -a "v${VERSION}" -m "Release v${VERSION}"
   git push origin "v${VERSION}"
   ```

3. For major or minor version changes, create a release branch:

   ```bash
   git checkout -b "release-v${VERSION%.*}"
   git push origin "release-v${VERSION%.*}"
   ```

4. If the GitHub Release was not created automatically, draft one at
   <https://github.com/SumoLogic/tailing-sidecar/releases>.
