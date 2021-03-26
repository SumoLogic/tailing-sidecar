# Tailing Sidecar Release Instruction

1. Prepare Release pull request with Helm Chart version change in [Chart.yaml](../helm/tailing-sidecar-operator/Chart.yaml).

1. Create the release tag for commit with Helm Chart version change, e.g.

   ```bash
   git tag -a v0.1.0 -m "Release v0.1.0"
   ```

1. Push the release tag, e.g.

   ```bash
   git push origin v0.1.0
   ```

1. For major and minor version change prepare release branch, e.g.

    ```bash
    git checkout -b release-v0.1
    git push origin release-v0.1
    ```

1. Cut the release

   - Go to https://github.com/SumoLogic/tailing-sidecar/releases and click "Draft a new release".
   - Compare changes since the last release.
   - Prepare release notes.
