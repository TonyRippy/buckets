#!/bin/sh
set -eu

# Invoked by make release VERSION=v0.1.2. See docs/runbooks/release.md.
version=${RELEASE_VERSION:-}
if ! printf '%s\n' "$version" | LC_ALL=C grep -Eq '^v[01]\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$'; then
    echo 'Usage: make release VERSION=v0.1.2 (stable v0/v1 versions only)' >&2
    exit 1
fi
tag="go/buckets/$version"

command -v gh >/dev/null
gh auth status
if [ -n "$(git status --porcelain)" ]; then
    echo 'Commit or stash all changes first.' >&2
    exit 1
fi

git fetch origin --tags
if [ "$(git rev-parse HEAD)" != "$(git rev-parse refs/remotes/origin/main)" ]; then
    echo 'HEAD must match origin/main; update your checkout first.' >&2
    exit 1
fi
if git show-ref --verify --quiet "refs/tags/$tag"; then
    echo "Tag $tag already exists; see the runbook for recovery." >&2
    exit 1
fi

"${RELEASE_MAKE:-make}" -C go build test
if [ -n "$(git status --porcelain)" ]; then
    echo 'Build/tests changed the checkout; inspect changes first.' >&2
    exit 1
fi

git tag -a "$tag" -m "Go buckets $version"
git push origin "refs/tags/$tag"
gh release create "$tag" --repo TonyRippy/buckets --verify-tag \
    --title "Go buckets $version" --generate-notes
