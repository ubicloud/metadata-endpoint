## How to release?

```
git tag -a vNEW_VERSION -m "First release"
git push origin vNEW_VERSION

goreleaser check
goreleaser release
```