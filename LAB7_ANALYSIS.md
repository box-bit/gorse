# Lab #7 — From Code to Operation: Investigation Record

**Team:** individual  
**Student(s) / GitHub handle(s):** Ondrej Müller — [@Hogigo](https://github.com/Hogigo)  
**Fork / working branch:** [`box-bit/gorse@semester-investigation`](https://github.com/box-bit/gorse/tree/semester-investigation) (Lab #7 work on [`lab7-deploy-verify`](https://github.com/box-bit/gorse/tree/lab7-deploy-verify), merged by PR)  
**Upstream project and pinned version:** <https://github.com/gorse-io/gorse>, tag `v0.5.11` = [`de9141d5`](https://github.com/gorse-io/gorse/commit/de9141d5d87236289d8d7f102bb9588234ff63cf) (Lab #5 version), plus the Lab #6 work already on `semester-investigation` ([`32714b42`](https://github.com/box-bit/gorse/commit/32714b4286be17ea3986c8202c50fe239eabe1fd)). Lab #6 did not touch `server/`.

## 1. Stakeholder, goal and selected functionality

**Stakeholder concern:** Sam, Release Engineer, made specific to Gorse: *"Client applications get their recommendation lists from the Gorse REST API. I need to see that a server built from our fork can be installed and started on a clean machine and returns the right leaderboard for a fixed data set. When we change how the API handles the request limit `n`, the pipeline must show that the new rule is in place and that normal requests still behave exactly as before."*

**GQM goal:** Analyze **the leaderboard/list endpoint of the Gorse REST API (`RestServer.SearchDocuments`) and its delivery pipeline** for the purpose of **evaluation**, with respect to **deployability in a clean environment and detection of an intended behaviour change**, from the viewpoint of **Sam, the release engineer**, in the context of **the `box-bit/gorse` fork built and run on a GitHub-hosted `ubuntu-24.04` runner**.

| Question | Question in your own words | Evidence / measurement and why it helps |
|---|---|---|
| Q1 | Can a server built from the exact pushed commit be installed, started and used on a fresh Ubuntu runner, and does it give the expected leaderboard for a fixed input? | (a) Commit match: `GITHUB_SHA` = checkout `HEAD` = `vcs.revision` stamped in the built `gorse-in-one` binary, with `vcs.modified=false`. This shows that the binary was built from the commit under test, not from a published package. (b) Setup outcome: build time, time until `/api/health/ready`, time until the leaderboard is computed. (c) Functional results: HTTP status and item ids for 5 fixed requests to the deployed server, checked by assertions. |
| Q2 | When the behaviour for `n=-1` is changed, do the checks show the change, and do they show that the unchanged behaviour stayed the same? | (a) Status/body of the changed case `n=-1` in the baseline run and in the changed run. (b) "Counter-check": in each run the changed case is also checked against the *other* version's expectation, and that check must fail. This shows that the assertions can tell the two versions apart. (c) The 4 unchanged checks (same assertions in both runs) and the `server` + `master` package test suites pass in both runs. |

**Selected function/component:** `RestServer.SearchDocuments`, which serves `GET /api/non-personalized/{name}` — [`server/rest.go` at the baseline](https://github.com/box-bit/gorse/blob/def09661903811100f281554fa9753106f18127d/server/rest.go#L613-L689). It is also used for item-to-item, user-to-user, item/user neighbours and collaborative-filtering lists, and by the master dashboard ([`master/rest.go`](https://github.com/box-bit/gorse/blob/def09661903811100f281554fa9753106f18127d/master/rest.go#L1096-L1111)).  
**Why this represents the software's purpose:** Gorse is a recommender system, and its main output is a ranked list of items served to applications over REST. This function reads a ranked list from the cache store, takes the request limit `n` and `offset` into account, removes items the user has already read (`user-id`) and returns the result. That is the core "serve recommendations" path.  
**Representative input and baseline expectation:** 15 items `item-00`…`item-14`; `item-i` gets `15 − i` "star" ratings, so the `most_starred` leaderboard computed by Gorse is `item-00` (15) … `item-14` (1); user `alice` has a "read" event for `item-01`. Requests and expected results:

| Case | Request | Expected (baseline) |
|---|---|---|
| `top_n` | `?n=3` | 200, `item-00,item-01,item-02` |
| `default_n` | (no `n`) | 200, `item-00`…`item-09` (`server.default_n = 10`) |
| `exclude_read_items` | `?n=3&user-id=alice` | 200, `item-00,item-02,item-03` (read `item-01` removed) |
| `non_numeric_n_rejected` | `?n=abc` | 400 |
| `negative_n` | `?n=-1` | 200, **all 15 items** |

**Planned small behavior change:** reject a negative `n` with **400 Bad Request** ("invalid n parameter"). Today a negative `n` goes into an unbounded cache query (`end = offset + n = -1` means "no limit" in the cache layer), so `n=-1` returns the whole list. The API documents `n` only as "Number of returned items", so a negative number is not meaningful. A release engineer wants every response to be bounded by the request: a client bug or a crafted request should not be able to download a full leaderboard (up to `recommend.cache_size` items, for every list endpoint).  
**What must still work unchanged?** Every request with `n ≥ 0`: top-n order (`top_n`), default limit (`default_n`), offset handling (Go test `offset_n`), read-item exclusion (`exclude_read_items`), rejection of a non-numeric `n` (`non_numeric_n_rejected`), and the existing `server` and `master` package tests.

## 2. Baseline: deploy and execute on GitHub

Workflow: [`.github/workflows/lab7-deploy-verify.yml`](https://github.com/box-bit/gorse/blob/def09661903811100f281554fa9753106f18127d/.github/workflows/lab7-deploy-verify.yml) — push-triggered on `semester-investigation` and `lab7-*`. Supporting files: [`lab7/e2e.sh`](https://github.com/box-bit/gorse/blob/def09661903811100f281554fa9753106f18127d/lab7/e2e.sh) (HTTP checks against the deployed binary), [`lab7/gorse-lab7.toml`](https://github.com/box-bit/gorse/blob/def09661903811100f281554fa9753106f18127d/lab7/gorse-lab7.toml) (SQLite, one leaderboard, 10 s task loop), [`lab7/claimed-contract.txt`](https://github.com/box-bit/gorse/blob/def09661903811100f281554fa9753106f18127d/lab7/claimed-contract.txt) (`legacy` / `bounded`), [`server/lab7_request_limit_test.go`](https://github.com/box-bit/gorse/blob/def09661903811100f281554fa9753106f18127d/server/lab7_request_limit_test.go) (same checks as a Go test on the real REST handler behind a local HTTP listener). No inherited workflow runs on these branches (they trigger only on `master`/`release-**`, releases or issues). The Lab #6 workflows run on `semester-investigation` only.

| Evidence | Baseline |
|---|---|
| Baseline commit SHA and link | [`def09661903811100f281554fa9753106f18127d`](https://github.com/box-bit/gorse/commit/def09661903811100f281554fa9753106f18127d) — adds only the Lab #7 workflow, `lab7/` and the Lab #7 test; `server/rest.go` is identical to `semester-investigation` |
| Actions workflow file link | [`lab7-deploy-verify.yml`](https://github.com/box-bit/gorse/blob/def09661903811100f281554fa9753106f18127d/.github/workflows/lab7-deploy-verify.yml) |
| Ubuntu/runtime and relevant dependency versions | GitHub-hosted `ubuntu-24.04`: Ubuntu 24.04.5 LTS, kernel 7.0.0-1017-azure x86_64, runner image `ubuntu24 20261116.376.1`; Go 1.26.11 linux/amd64 (from `go.mod`: `go 1.26`, installed by `actions/setup-go@v7`); Go dependencies pinned by `go.sum`; SQLite through the pure-Go driver (`modernc.org/sqlite`), no external services or secrets |
| How the fork's code is installed/built and verified as the executed code | `actions/checkout@v7` of the pushed commit → `go build ./cmd/gorse-in-one` → the step fails unless `git rev-parse HEAD` = `GITHUB_SHA`, the tree is clean, and the binary's `vcs.revision` = `GITHUB_SHA` with `vcs.modified=false`. The binary is then started on the runner (`127.0.0.1:8088`) and must answer `/api/health/ready` |
| Functional case executed (command/script/test or local endpoint) | `lab7/e2e.sh seed` posts 15 items + 121 feedback rows through `POST /api/items` and `POST /api/feedback`, then waits until the Gorse master has computed `most_starred`; `lab7/e2e.sh check` calls `GET /api/non-personalized/most_starred` with the 5 requests in §1. Then `go test -run TestLab7 ./server/` and the regression `go test ./server/... ./master/...` |
| Expected vs observed behavior | All 5 requests gave the expected result: `top_n` 200 `item-00,item-01,item-02`; `default_n` 200 `item-00`…`item-09`; `exclude_read_items` 200 `item-00,item-02,item-03`; `non_numeric_n_rejected` 400 `strconv.Atoi: parsing "abc": invalid syntax`; **`negative_n` 200 with all 15 items** `item-00`…`item-14` ([e2e-results.jsonl](lab7-evidence/baseline-def09661/e2e-results.jsonl)) |
| Check result and meaningful measurement(s), if any | Job **passed** in 3 min 02 s. Commit match: `GITHUB_SHA` = checkout `HEAD` = binary `vcs.revision` = `def09661…`, `vcs.modified=false`, 0 modified files. Build 107 s (cold Go module cache), server ready after 1 s, leaderboard computed 10 s after seeding (15 items + 121 feedback rows inserted). E2E checks against `legacy`: 9/9 assertions pass, `failures: 0`. E2E counter-check against `bounded`: **failed as required** (`expected '400', got '200'`). Go Lab #7 tests: `ok` (2 tests, 5 subtests); Go counter-check: `FAIL` as required (`expected: 400, actual: 200`). Regression: `server ok`, `master ok` |
| Exact GitHub Actions run URL | <https://github.com/box-bit/gorse/actions/runs/45848750148> — evidence copied to [`lab7-evidence/baseline-def09661/`](lab7-evidence/baseline-def09661/) |

**What does this establish, and what does it not establish?** It shows that the pushed commit (and no other code) can be built, deployed and used on a clean Ubuntu runner without external services, and that it returns the expected leaderboard, honours `n`/`offset`/default `n`, removes read items, and returns all items for `n=-1`. It does not show behaviour with the production stores (Redis/MySQL/…), with many users or concurrent load, or in a multi-node master/server/worker setup.

## 3. Controlled change and verification

**Source/function changed:** [`RestServer.SearchDocuments` in `server/rest.go`](https://github.com/box-bit/gorse/blob/6b17ddc3504a61549dee9f03b6dc1ada9714f1d4/server/rest.go#L634-L639) — 6 added lines:

```go
// A negative n used to fall through to an unbounded query and return the whole
// list. Reject it so that the size of a response is always bounded by the request.
if n < 0 {
    BadRequest(response, errors.Errorf("invalid n parameter: must be non-negative, got %d", n))
    return
}
```

**Reason for change:** keep every response bounded by the request. Before, one request with a negative `n` returned a whole list (up to `recommend.cache_size` items) from every list endpoint, and this was not documented. It also makes the API consistent: a non-numeric `n` was already rejected with 400.  
**Code-change commit(s):** [`6b17ddc3504a61549dee9f03b6dc1ada9714f1d4`](https://github.com/box-bit/gorse/commit/6b17ddc3504a61549dee9f03b6dc1ada9714f1d4) — `server/rest.go` (+6) and `lab7/claimed-contract.txt` (`legacy` → `bounded`). No test assertion, test data or workflow line was changed.

| Item | Before change | After change |
|---|---|---|
| Same input/case | `GET /api/non-personalized/most_starred?n=-1` on the deployed `gorse-in-one`, same 15-item data | same |
| Expected behavior under each version's stated contract | `legacy`: 200, all 15 items | `bounded`: 400, body contains `invalid n parameter`, no items |
| Observed behavior on GitHub runner | **200**, 15 items (`item-00`…`item-14`) | **400**, body `invalid n parameter: must be non-negative, got -1`, 0 items |
| Automated case check/result | `legacy` (claimed): pass, in e2e and Go. `bounded` counter-check: fails as required — e2e `expected '400', got '200'`, Go `expected: 400, actual: 200` | `bounded` (claimed): pass, in e2e and Go. `legacy` counter-check: fails as required — e2e `expected '200', got '400'` and `expected '15', got '0'`, Go `expected: 200, actual: 400` |
| Unchanged regression/invariant check | `top_n`, `default_n`, `exclude_read_items`, `non_numeric_n_rejected` pass in e2e (identical observed ids/status), plus Go `offset_n`; `go test ./server/... ./master/...`: `server ok`, `master ok` | the same assertions pass with identical observed ids/status; `server ok`, `master ok` |
| Exact Actions run URL | <https://github.com/box-bit/gorse/actions/runs/45848750148> (3 min 02 s) | <https://github.com/box-bit/gorse/actions/runs/45849378269> (1 min 05 s) |
| Commit SHA | `def09661903811100f281554fa9753106f18127d` | `6b17ddc3504a61549dee9f03b6dc1ada9714f1d4` |

**How do your checks distinguish the two versions?** The changed case has two expectations in the code, `legacy` (200 + 15 items) and `bounded` (400 + "invalid n parameter" + no items). Each version states which one it claims in `lab7/claimed-contract.txt`. Both check layers (curl against the deployed binary, and the Go test) assert the claimed contract, which must pass, and then the other contract, which must fail. A "must fail" step turns red if the other expectation unexpectedly passes. So in the baseline run the new expectation fails on the old code, and in the changed run the old expectation fails on the new code. The 4 unchanged checks use the same assertions in both commits, so the new version is not accepted just because all expectations were rewritten.

**Did anything fail unexpectedly? What did you learn or fix?** Before the first push, I ran the workflow steps locally with the same shells GitHub uses. This found two problems in the workflow, which I fixed before the baseline commit: (1) a step without `shell: bash` runs as `bash -e` *without* `pipefail`, so `if check | tee log` tested `tee`'s exit status and a "must fail" counter-check could never fail — I added `set -o pipefail` to both counter-check steps; (2) the summary step stopped early when an evidence file was missing — it now runs with `set +e`. I also had to keep the evidence folder outside the checkout: an untracked folder in the source tree makes Go stamp `vcs.modified=true`, and then the commit-match check fails. On GitHub, both runs passed at the first attempt and no check failed unexpectedly. Two observations: (a) both runs show a red annotation `sha256 rest.go: <hash>` although the step passed. This is a false positive: the Go problem matcher installed by `actions/setup-go` reads the provenance line `… rest.go: …` as a compiler error message (`file.go: message`). I did not rename the line between the two runs, so that the pipeline stays identical for both versions; renaming it would remove the annotation. (b) The build took 107 s in the baseline run and 6 s in the changed run, because `actions/setup-go` restored the Go module/build cache saved by the first run. So build time measures cache state, not the change.

## 4. Engineering judgment and submission

**Answer Q1:** Yes, for this setup. In both runs the binary that served the requests was built on a fresh `ubuntu-24.04` runner from the pushed commit: `GITHUB_SHA`, checkout `HEAD` and the `vcs.revision` stamped into `gorse-in-one` were equal (`def09661…` / `6b17ddc3…`), with `vcs.modified=false` and a clean tree, and the `rest.go` hash differs between the runs (`6c789f37…` → `9286f5db…`). The server became ready in 1–2 s with only SQLite and computed the leaderboard 8–10 s after the fixed data was inserted through the public API. All four unchanged requests returned exactly the expected ids/status in both runs.  
**Answer Q2:** Yes. For the same request `n=-1`, the baseline returned 200 with all 15 items and the changed version returned 400 `invalid n parameter`. Each run also checked the *other* version's expectation, and that check failed in both runs (baseline: `expected '400', got '200'`; changed: `expected '200', got '400'`). So the checks do not pass for any behaviour — they tell the two versions apart, in both directions, at two levels (deployed binary over HTTP and in-process Go test). The unchanged checks (identical assertions in both commits) and the existing `server`/`master` suites passed in both runs, so the new version was not accepted because expectations were rewritten.

**Release decision:** **Accept the changed version**, with one release note. The intended behaviour (negative `n` → 400) is verified on the deployed binary, all checked `n ≥ 0` behaviour is unchanged, and the regression suites pass. Because a request that used to succeed (200) now fails (400), this is an API contract change: the release notes should say that clients sending `n=-1` to mean "everything" must switch to an explicit positive `n`.  
**One important limitation or missing check:** The functional checks cover only `GET /api/non-personalized/{name}` on a single-node `gorse-in-one` with SQLite. `SearchDocuments` also serves item-to-item, user-to-user, neighbours, collaborative-filtering lists and the master dashboard, which now also reject a negative `n`. For those paths only the existing upstream tests ran, and they do not send a negative `n`. Also, no real client was checked for whether it relies on `n=-1`, and the production cache stores (Redis, MySQL, Postgres, MongoDB) were not run — the "negative end = unbounded" behaviour belongs to the cache layer and could differ per store.

| GitHub evidence | Link |
|---|---|
| Lab #7 pull request (within your fork) | PR_TBD |
| **Merged commit** on `semester-investigation` | MERGE_TBD |
| Final `LAB7_ANALYSIS.md` on merged branch | <https://github.com/box-bit/gorse/blob/semester-investigation/LAB7_ANALYSIS.md> |
| Student 1 meaningful commit(s) | Ondrej Müller (@Hogigo): [`def09661`](https://github.com/box-bit/gorse/commit/def09661903811100f281554fa9753106f18127d) (workflow + checks, baseline), [`6b17ddc3`](https://github.com/box-bit/gorse/commit/6b17ddc3504a61549dee9f03b6dc1ada9714f1d4) (behaviour change) |
| Student 2 meaningful commit(s), if a pair | N/A (individual) |

**Evidence preservation:** Actions artifacts expire, so the full evidence of both runs (provenance, binary build info, e2e and Go check outputs, JSON results, server log), downloaded from the runs, is committed in [`lab7-evidence/baseline-def09661/`](lab7-evidence/baseline-def09661/) (run [45848750148](https://github.com/box-bit/gorse/actions/runs/45848750148)) and [`lab7-evidence/changed-6b17ddc3/`](lab7-evidence/changed-6b17ddc3/) (run [45849378269](https://github.com/box-bit/gorse/actions/runs/45849378269)). The run pages (logs, step summaries) remain viewable after the artifacts expire.

**LLM assistance (if used):** Claude (Claude Code) was used to explore the Gorse source, to write the workflow, `lab7/e2e.sh`, the Go test and the code change, and to draft this record. I checked independently: the observed behaviour for `n=-1`/`0`/`abc` by probing the handler before choosing the change; that the curl checks and the Go test pass on the claimed contract and fail on the other one, for both commits; the workflow with `actionlint` (with shellcheck) and with a local dry run in a clean clone; and every figure in this file against the GitHub run logs and step summaries.
