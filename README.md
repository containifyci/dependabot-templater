# dependabot-templater

This is a little go application to render dependabot.yml for multi folder terraform projects.

It solves the following missing feature of Dependabot for now.
[Support Nested Terraform Code](https://github.com/dependabot/dependabot-core/issues/649)

## Usage

Just specify the path and it will look for terraform folders that define a backend and generate a full Dependabot configuration for all found folders.
```bash
./dependabot-templater [type/package-ecosystem] [path]
```

### Terraform

```bash
dependabot-templater terraform iac-example/projects/staging
```
Output (Snippet)

```yaml
---
# https://docs.github.com/github/administering-a-repository/configuration-options-for-dependency-updates
version: 2
registries:
  git-terraform-modules:
    type: git
    url: https://github.com
    username: x-access-token
    password: ${{ secrets.REGISTRIES_PAT_TOKEN }}
updates:
  - package-ecosystem: "terraform"
    directory: "projects/staging/"
    schedule:
      interval: "weekly"
      day: "sunday"
    labels:
      - dependencies
    open-pull-requests-limit: 1
    commit-message:
      include: "scope"
    ignore:
      - dependency-name: "*"
        update-types: ["version-update:semver-patch", "version-update:semver-minor"]
    registries:
      - git-terraform-modules
```


### Github Actions

```bash
dependabot-templater gha github-actions
```

Output (Snippet)
```yaml
---
# https://docs.github.com/github/administering-a-repository/configuration-options-for-dependency-updates
version: 2
updates:
  - package-ecosystem: "github-actions"
    directory: "github-actions/golangci-lint"
    schedule:
      interval: "weekly"
      day: "sunday"
    labels:
      - dependencies
    commit-message:
      prefix: "chore:"
      include: "scope"
    groups:
      minor:
        patterns:
        - "*"
        update-types:
        - "minor"
        - "patch"
  - package-ecosystem: "github-actions"
    directory: "github-actions/remote-access-ssh"
    schedule:
      interval: "weekly"
      day: "sunday"
    labels:
      - dependencies
    commit-message:
      prefix: "chore:"
      include: "scope"
    groups:
      minor:
        patterns:
        - "*"
        update-types:
        - "minor"
        - "patch"
```

### NPM

```bash
dependabot-templater npm github-actions
```

Output (Snippet)
```yaml
---
# https://docs.github.com/github/administering-a-repository/configuration-options-for-dependency-updates
version: 2
updates:
  - package-ecosystem: "npm"
    directory: "argocd-app-diff"
    schedule:
      interval: "weekly"
      day: "sunday"
    labels:
      - dependencies
    commit-message:
      include: "scope"
    groups:
      minor:
        patterns:
        - "*"
        update-types:
        - "minor"
        - "patch"
  - package-ecosystem: "npm"
    directory: "link-from-comment"
    schedule:
      interval: "weekly"
      day: "sunday"
    labels:
      - dependencies
    commit-message:
      include: "scope"
    groups:
      minor:
        patterns:
        - "*"
        update-types:
        - "minor"
        - "patch"
```

### Templates

The templates are at the moment embedded in the binary and can't be customized.
In order to change them you have to adjust the template file and run `make build`
that will store the adjusted template version into the new binary.

## Build

```bash
make build
```

## Test

```bash
make test
```