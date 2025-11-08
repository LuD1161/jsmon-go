# Release Process

This document describes how to create a new release of JSMon-Go.

## Prerequisites

1. GitHub repository set up with proper permissions
2. GoReleaser installed locally (for testing): `brew install goreleaser`
3. Docker installed (for testing Docker builds)
4. GitHub Container Registry (GHCR) access configured

## Automated Release (GitHub Actions)

The project uses GoReleaser with GitHub Actions for fully automated releases.

### Steps

1. **Ensure everything is committed and pushed**
   ```bash
   git status
   git push origin main
   ```

2. **Create and push a new tag**
   ```bash
   # Create a new tag (semver format: v1.0.0)
   git tag -a v1.0.0 -m "Release v1.0.0: Initial release"

   # Push the tag to trigger the release workflow
   git push origin v1.0.0
   ```

3. **Monitor the release**
   - Go to: https://github.com/LuD1161/jsmon-go/actions
   - Watch the "Release" workflow run
   - It will automatically:
     - Build binaries for all platforms (Linux, macOS, Windows)
     - Create GitHub release with binaries and changelog
     - Build Docker images for amd64 and arm64
     - Push multi-arch images to ghcr.io
     - Create version tags (latest, v1, v1.0, v1.0.0)

4. **Verify the release**
   - Check GitHub releases: https://github.com/LuD1161/jsmon-go/releases
   - Check Docker images: https://github.com/LuD1161/jsmon-go/pkgs/container/jsmon-go
   - Test a binary download
   - Test Docker image pull

## Manual Release (Local Testing)

For testing before pushing a tag:

```bash
# Test the release process locally (dry-run)
goreleaser release --snapshot --clean

# Check the output in dist/
ls -la dist/

# Test a binary
./dist/jsmon_linux_amd64_v1/jsmon --help
```

## Release Checklist

Before creating a release:

- [ ] All tests pass (`make test`)
- [ ] Documentation is up to date
- [ ] CHANGELOG is updated (GoReleaser will auto-generate, but manual updates are better)
- [ ] Version follows semver (MAJOR.MINOR.PATCH)
- [ ] No uncommitted changes
- [ ] Main branch is up to date

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (v2.0.0): Incompatible API changes
- **MINOR** (v1.1.0): New features, backwards compatible
- **PATCH** (v1.0.1): Bug fixes, backwards compatible

## What Gets Released

### Binaries
- Linux: amd64, arm64
- macOS: amd64 (Intel), arm64 (Apple Silicon)
- Windows: amd64

### Docker Images
All images pushed to `ghcr.io/LuD1161/jsmon-go`:

- `latest` - Latest release (multi-arch)
- `v1` - Latest v1.x.x (multi-arch)
- `v1.0` - Latest v1.0.x (multi-arch)
- `v1.0.0` - Specific version (multi-arch)
- Architecture-specific tags: `-amd64`, `-arm64`

### Archives
- `.tar.gz` for Linux/macOS
- `.zip` for Windows
- Contains: binary, README.md, LICENSE, .env.example

### Checksums
- SHA256 checksums for all artifacts

## Troubleshooting

### Release fails with "permission denied"
- Ensure `GITHUB_TOKEN` has `contents: write` and `packages: write` permissions
- Check GitHub Actions settings

### Docker push fails
- Ensure GitHub Container Registry is enabled
- Check package permissions: https://github.com/LuD1161/jsmon-go/settings/packages

### Tag already exists
```bash
# Delete local tag
git tag -d v1.0.0

# Delete remote tag
git push origin :refs/tags/v1.0.0

# Recreate and push
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

## Rollback

If a release has issues:

1. **Delete the GitHub release**
   - Go to Releases page
   - Delete the problematic release

2. **Delete the tag**
   ```bash
   git tag -d v1.0.0
   git push origin :refs/tags/v1.0.0
   ```

3. **Delete Docker images** (if needed)
   - Go to Packages
   - Delete the version tag

4. **Fix the issues and re-release**

## Post-Release

After a successful release:

1. Announce on social media (if applicable)
2. Update homebrew tap (if you create one)
3. Update documentation links
4. Monitor issues for bug reports

## Emergency Hotfix Release

For urgent bug fixes:

1. Create a hotfix branch from the tag
   ```bash
   git checkout -b hotfix/v1.0.1 v1.0.0
   ```

2. Fix the issue and commit

3. Tag the hotfix
   ```bash
   git tag -a v1.0.1 -m "Hotfix: Critical bug fix"
   ```

4. Push and merge
   ```bash
   git push origin v1.0.1
   git checkout main
   git merge hotfix/v1.0.1
   git push origin main
   ```
