# Workshop repo justfile - docs + common chores.

default:
    @just --list

# --- Docs (MkDocs + Material) ---

# Install the docs toolchain into the project venv.
docs-install:
    uv sync --group docs

# Serve the docs site locally with hot reload.
docs-serve:
    uv run --group docs mkdocs serve

# Serve on a specific port (useful when 8000 is busy).
docs-serve-port PORT="8000":
    uv run --group docs mkdocs serve --dev-addr=127.0.0.1:{{PORT}}

# Build the static site into ./site.
docs-build:
    uv run --group docs mkdocs build

# Build with strict validation - fails on broken links / warnings.
docs-validate:
    uv run --group docs mkdocs build --strict

# Deploy to GitHub Pages manually (CI does this on push to main).
docs-deploy:
    uv run --group docs mkdocs gh-deploy --force

# Remove generated site output.
docs-clean:
    rm -rf site/
