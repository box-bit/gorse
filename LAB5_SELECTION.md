# Lab #5 — System Selection

## 1. Student details

Alessandro Fantesini - Ondrej Müller
box-bit - Hogigo
<!-- TODO: name(s), GitHub handle(s), individual/pair -->

## 2. Shortlist and comparison

<!-- TODO: Tasks 1–2 -->

## 3. Final selection

<!-- TODO -->

## 4. Repository information

| Item | Value |
|---|---|
| Original repository | <https://github.com/gorse-io/gorse> |
| Fork | <https://github.com/box-bit/gorse> |
| Working branch | [`semester-investigation`](https://github.com/box-bit/gorse/tree/semester-investigation) |
| Pinned version | [`v0.5.11`](https://github.com/gorse-io/gorse/releases/tag/v0.5.11) (latest stable release, 2026-07-14) |
| Pinned commit SHA | `de9141d5d87236289d8d7f102bb9588234ff63cf` |
| Go version required | `go 1.26` (from `go.mod`) |

Setup commands used:

```bash
git clone https://github.com/box-bit/gorse.git
cd gorse
git remote add upstream https://github.com/gorse-io/gorse.git
git fetch upstream --tags
git switch -c semester-investigation v0.5.11
git rev-parse HEAD   # de9141d5d87236289d8d7f102bb9588234ff63cf
git push -u origin semester-investigation
```

Why `v0.5.11` rather than `v0.6.0-alpha.1`: `v0.6.0-alpha.1` is marked as a pre-release, and its
[release notes](https://github.com/gorse-io/gorse/releases/tag/v0.6.0-alpha.1) warn that
"APIs, configurations, and data storage might change between pre-release versions". A stable
release gives a reproducible baseline for Labs #6–#9.

Note: issues, pull requests and reviews are not copied into the fork; the development history
(Lab #9) is investigated on the upstream repository.

## 5. Software exploration

<!-- TODO: Task 4 -->

## 6. Future-lab feasibility

<!-- TODO: Task 5 -->
