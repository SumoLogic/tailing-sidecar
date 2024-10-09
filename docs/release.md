# Tailing Sidecar Release Instruction

1. Create and Merge release pull request with Helm Chart version change in [Chart.yaml](../helm/tailing-sidecar-operator/Chart.yaml)

1. Create the release tag for commit with Helm Chart version change, e.g. when release tag = 0.1.0

   ```bash
   export VERSION=0.1.0
   git checkout main
   git pull
   git tag -a "v${VERSION}" -m "Release v${VERSION}"
   ```

1. Push the release tag, e.g.

   ```bash
   git push origin "v${VERSION}"
   ```

1. For major and minor version change prepare release branch, e.g.

    ```bash
    git checkout -b "release-v${VERSION%.*}"
    git push origin "release-v${VERSION%.*}"
    ```

1. Cut the release

   - Go to https://github.com/SumoLogic/tailing-sidecar/releases and click "Draft a new release".
   - Compare changes since the last release.
   - Prepare release notes.
