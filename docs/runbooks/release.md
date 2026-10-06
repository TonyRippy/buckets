# Releasing Go buckets

The Go module is `github.com/TonyRippy/buckets/go/buckets`, rooted at
`go/buckets/`. Its Git tags must include that directory prefix:

| Module version | Git tag |
| --- | --- |
| `v0.1.2` | `go/buckets/v0.1.2` |
| `v1.0.0` | `go/buckets/v1.0.0` |
| `v1.0.0-rc.1` | `go/buckets/v1.0.0-rc.1` |

A repository-root tag such as `v0.1.2` does not release this module. Consumers
use the version alone, for example
`go get github.com/TonyRippy/buckets/go/buckets@v0.1.2`.
These are Go's [module subdirectory tagging rules](https://go.dev/ref/mod#vcs-version).

Choose a new semantic version: patch for compatible fixes, minor for compatible
features, and major for breaking changes once v1 is stable. During v0 development,
minor releases may break compatibility. The current module path supports v0 and
v1. Before releasing v2 or later, migrate the module and imports to a `/v2` (or
later) suffix and update the automation; changing only the tag is insufficient.
See [Go module version numbering](https://go.dev/doc/modules/version-numbers).
The `go 1.26` directive specifies the minimum Go version, not the library version.

## Prerequisites

- Install Git, Make, Go 1.26 or newer, and the GitHub CLI (`gh`).
- Run `gh auth login` and ensure you can push tags and create releases in
  `TonyRippy/buckets`. Git push authentication must also be configured.
- Use a clone whose `origin` points to `TonyRippy/buckets`.
- Merge the intended changes into `main`. Start with a clean checkout, including
  untracked files, and review the changes since the previous module tag.
- Choose an unused version. Published versions are immutable; fixes get a new
  version instead of moving or replacing an existing tag.

## Automated stable release

Run from the repository root, replacing the example with the intended version:

```sh
git switch main
git pull --ff-only origin main
make release VERSION=v0.1.2
```

This command runs `scripts/go-release.sh` and publishes the release. It accepts canonical stable v0/v1 versions
with a leading `v`, authenticates `gh`, fetches tags, requires `HEAD` to match
`origin/main`, rejects an existing tag, and runs the Go build and tests. It then
creates an annotated tag, pushes that tag alone, and runs:

```sh
gh release create go/buckets/v0.1.2 \
  --repo TonyRippy/buckets --verify-tag \
  --title 'Go buckets v0.1.2' --generate-notes
```

`--verify-tag` requires the remote tag to exist. `--generate-notes` uses GitHub's
release notes generator. Review the resulting release and edit its notes as
needed; generated notes may include repository changes outside the Go module.
See the [GitHub CLI release reference](https://cli.github.com/manual/gh_release_create).

## Manual release and prereleases

Use this sequence for a prerelease or when controlling the notes range. Start
with the same prerequisites and updated `main` checkout as above. For a stable
release, set a stable version and omit `--prerelease`.

```sh
VERSION=v1.0.0-rc.1
TAG="go/buckets/$VERSION"

gh auth status
git fetch origin --tags
git status --short                    # Must be empty.
git rev-parse HEAD origin/main        # Both hashes must match.
git tag --list "$TAG"                 # Must be empty.
make -C go build test                 # Stop if either fails.
git status --short                    # Must still be empty.

git tag -a "$TAG" -m "Go buckets $VERSION"
git push origin "refs/tags/$TAG"
gh release create "$TAG" --repo TonyRippy/buckets --verify-tag \
  --title "Go buckets $VERSION" --generate-notes --prerelease
```

Run these commands in order and stop on any failure. Optionally add
`--notes-start-tag go/buckets/v0.1.1` to `gh release create`, substituting the
previous release tag for this module. This is useful when other implementations
also have releases in the repository. For the first release, omit that option.
To review notes before publishing the GitHub release, add `--draft`, review it
with `gh release view "$TAG" --repo TonyRippy/buckets --web`, then publish with
`gh release edit "$TAG" --repo TonyRippy/buckets --draft=false`.
Pushing the tag already makes the module version available to Go, even while the
GitHub release is a draft.

## Verify the published release

```sh
gh release view go/buckets/v0.1.2 --repo TonyRippy/buckets
GOWORK=off go list -m github.com/TonyRippy/buckets/go/buckets@v0.1.2
```

Replace both example versions with the version released. The Go command should
report the module path and requested version. A module proxy may take time to
observe a newly pushed tag; retry before treating this as a release failure.

## Recover from a partial failure

The Make target deliberately refuses to reuse an existing tag. If publication
stops after creating the tag, inspect its commit before continuing:

```sh
git show --no-patch go/buckets/v0.1.2
git ls-remote origin 'refs/tags/go/buckets/v0.1.2*'
gh release view go/buckets/v0.1.2 --repo TonyRippy/buckets
```

If the local tag points to the intended release commit but the push failed, retry
`git push origin refs/tags/go/buckets/v0.1.2`. If the remote tag exists and the
GitHub release is absent, rerun only the `gh release create` command above (with
`--prerelease` for a prerelease). If the release already exists, review or edit it
instead of creating another. Never force-push or retag a published version.
