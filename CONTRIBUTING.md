# Contributing to Opencode Agent Switcher

Thank you for your interest in contributing to opencode-agent-switcher! This document provides guidelines and instructions for contributing.

## Development Setup

### Prerequisites

- Go 1.23 or later
- Make (optional, for using Makefile targets)

### Getting Started

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/opencode-agent-switcher.git
   cd opencode-agent-switcher
   ```
3. Create a feature branch:
   ```bash
   git checkout -b feature/my-feature
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```

## Development Workflow

We use **GitFlow** branching model:

| Branch | Purpose |
|--------|---------|
| `main` | Production releases |
| `develop` | Integration branch (default) |
| `feature/*` | New features → PR to develop |
| `release/*` | Release preparation |
| `hotfix/*` | Emergency fixes |

### Making Changes

1. Create a branch from `develop`:
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feature/my-feature
   ```

2. Make your changes

3. Run checks:
   ```bash
   make check
   ```
   Or individually:
   ```bash
   make fmt        # Format code
   make vet        # Run go vet
   make lint       # Run golangci-lint
   make test       # Run tests
   ```

4. Commit with clear messages:
   ```bash
   git commit -m "feat: add new feature description"
   ```

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `chore:` - Maintenance tasks
- `refactor:` - Code refactoring
- `test:` - Adding/updating tests

## Pull Request Process

1. Push your branch:
   ```bash
   git push origin feature/my-feature
   ```

2. Open a Pull Request to `develop`

3. Ensure CI checks pass

4. Request review from maintainers

5. Address review feedback

6. Once approved, a maintainer will merge your PR

## Code Style

- Follow standard Go conventions
- Run `gofmt -s -w .` before committing
- Add comments for exported functions
- Keep functions focused and small
- Handle errors properly (don't ignore them)

## Testing

Run tests:
```bash
make test
```

Generate coverage report:
```bash
make test-coverage
```

## Questions or Issues?

- Open a [Discussion](https://github.com/mario-gc/opencode-agent-switcher/discussions) for questions
- Open an [Issue](https://github.com/mario-gc/opencode-agent-switcher/issues) for bugs or feature requests

## License

By contributing, you agree that your contributions will be licensed under the MIT License.