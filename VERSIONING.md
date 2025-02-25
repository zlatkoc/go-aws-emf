# Versioning Policy

This project follows [Semantic Versioning](https://semver.org/) (SemVer).

## Version Format

Versions are in the format `vMAJOR.MINOR.PATCH` where:

- **MAJOR** version increases for incompatible API changes
- **MINOR** version increases for new functionality in a backward-compatible manner
- **PATCH** version increases for backward-compatible bug fixes

## Go Module Versioning

For Go projects, we follow the [Go Modules versioning conventions](https://go.dev/ref/mod#versions):

1. Versions v0.x.x are considered development versions and may have breaking changes between minor releases
2. Starting with v1.0.0, we guarantee API compatibility within the same major version
3. Import paths will use major version suffixes for v2 and beyond (e.g., `/v2`, `/v3`)

## Creating Releases

To create a new release:

1. Update code and ensure all tests pass
2. Tag the commit with the appropriate version:
   ```
   git tag -a v1.0.0 -m "Release v1.0.0"
   git push origin v1.0.0
   ```
3. The GitHub Actions workflow will automatically:
   - Run tests
   - Generate a changelog
   - Create a GitHub release

## Pre-releases

For pre-release versions, append a suffix like `-beta.1`, `-rc.1`:

```
git tag -a v1.0.0-beta.1 -m "Beta release v1.0.0-beta.1"
git push origin v1.0.0-beta.1
```

These tags will be marked as pre-releases in GitHub.