# Homebrew Tap Setup

This project is configured to publish to a Homebrew tap via GoReleaser. Follow these steps to set it up.

## 1. Create the Tap Repository

Create a new GitHub repository named `{{ .Scaffold.homebrew_name }}` under the `{{ .Scaffold.homebrew_owner }}` account. This repository will hold the Homebrew cask definitions that GoReleaser publishes.

## 2. Generate a Personal Access Token (PAT)

GoReleaser needs a PAT with write access to the tap repository (separate from `GITHUB_TOKEN` which is scoped to the source repo).

1. Go to **GitHub Settings > Developer Settings > Personal Access Tokens > Fine-grained tokens**
2. Create a new token with:
   - **Repository access**: Select the `{{ .Scaffold.homebrew_name }}` repository
   - **Permissions**: Contents (Read and write)
3. Copy the generated token

## 3. Add the Secret to Your Source Repository

1. Go to your source repository **Settings > Secrets and variables > Actions**
2. Click **New repository secret**
3. Name: `HOMEBREW_TAP_GITHUB_TOKEN`
4. Value: Paste the PAT from the previous step

## 4. Verify

After the next release, GoReleaser will push a cask definition to your tap repository. Users can then install with:

```sh
brew tap {{ .Scaffold.homebrew_owner }}/{{ .Scaffold.homebrew_name | trimPrefix "homebrew-" }}
brew install {{ .Scaffold.gomod | pathBase }}
```

## References

- [GoReleaser Homebrew Casks](https://goreleaser.com/customization/homebrew_casks/)
