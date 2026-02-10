# GitHub Actions Setup

This project uses GitHub Actions for CI/CD pipelines.

## Workflow Files

### `ci.yml` - Continuous Integration

Runs on every push to `main` and on pull requests to `main`.

**Jobs:**
- **test** - Runs on Ubuntu latest
  - Checks out code
  - Sets up Go 1.24.4
  - Runs tests in `./test/...`

### `cd.yml` - Continuous Deployment

Runs automatically after `ci.yml` completes successfully.

**Jobs:**
- **docker** - Builds and publishes Docker image to GitHub Container Registry
  - Extracts version from git tags (or uses `dev` if no tags)
  - Logs in to GHCR
  - Builds and pushes Docker image with tags:
    - `latest`
    - Version tag (e.g., `v1.0.0`)
    - Git SHA tag

## Docker Image Tags

| Tag | Description |
|-----|-------------|
| `latest` | Most recent build on main branch |
| `v*.*.*` | Version tag from git tags |
| `SHA` | Specific commit SHA |

## Setup Requirements

1. **GitHub Container Registry (GHCR)** permissions must be enabled
2. Docker buildx should be available in the runner
3. No additional secrets required (uses `GITHUB_TOKEN`)

## Running Locally

```bash
# Run CI tests
go test -v ./test/...

# Build Docker image locally
docker build -t ghcr.io/owner/repo:dev .
```

