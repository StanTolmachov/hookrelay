# Contributing

## Branches

- `main`: Main branch for production code

### Branch prefixes:
- `feat`: for new features
- `fix`: for bug fixes
- `chore`: for maintenance tasks
- `refactor`: for code refactoring
- `docs`: for documentation updates

Format: `<type>/HR-<number>-<short description>`

Workflow:

`git switch main && git pull`

`git switch -c feat/HR-1-service-skeleton`


## Commits

- `feat`: Feature implementation
- `fix`: Bug fix
- `chore`: Maintenance task
- `refactor`: Code refactoring
- `docs`: Documentation update
- `test`: Test implementation

Format: `<type>(<scope>): <short description>`

## Before opening a PR

```bash
go vet ./... && go test ./... -race && gofmt -l .
```

## Pull Requests

`git push -u origin <type>/HR-<number>-<short description>`

`gh pr create`

Name for pr:
`<type>(<scope>): <short description>`

### Body:

- What was done
- How to test
- Resolutions


`Closes #1` - #N is the number of issue

## Rules
- squash merge only; the branch is deleted afterwards
- push into main is closed, only through PR
- main is always green and deployable
- decisions with a rejected alternative go into DECISIONS.md in the same PR