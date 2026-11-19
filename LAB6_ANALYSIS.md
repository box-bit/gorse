# Lab #6 - Investigation Record

**Team / student(s):** Alessandro Fantesini (`box-bit`)  
**Mode:** individual  
**Approved upstream project:** <https://github.com/gorse-io/gorse>  
**Fork / working branch:** [`box-bit/gorse@semester-investigation`](https://github.com/box-bit/gorse/tree/semester-investigation)  
**Original Lab #5 version:** `de9141d5d87236289d8d7f102bb9588234ff63cf` (upstream tag `v0.5.11`)  
**Lab #6 pre-refactoring baseline:** source revision [`140cfb1bc72876a51356e8446e3b311687814cd0`](https://github.com/box-bit/gorse/commit/140cfb1bc72876a51356e8446e3b311687814cd0) (Go source identical to `de9141d5`; this commit only adds `LAB5_SELECTION.md`). Baseline scan commit: [`8bd232d3bb391d3e831217085f39f8efddee45bb`](https://github.com/box-bit/gorse/commit/8bd232d3bb391d3e831217085f39f8efddee45bb) (adds only the Lab 6 workflows, config and this file; no Go file differs from `de9141d5`)

> Keep this file short and evidence-based. Replace the placeholders as you go. Every analysis/test link must point to your **fork's GitHub-visible results**. Do not submit unsupported local-only tool screenshots.

## Task 1 - Stakeholder, GQM, and Code Smell Selection

**Selected stakeholder and exact fictional concern from handout:** Scenario A - Alex, Software Maintainer: *"Our software keeps evolving. Future contributors need to understand and modify its components, but certain design choices might make that work unnecessarily difficult. Please investigate possible code smells, identify maintainability concerns worth addressing, and demonstrate a focused improvement without changing the software's intended behavior."*

**GQM goal:** Analyze **the `master` package of Gorse** for the purpose of **finding code smells that make future changes harder, and fixing one of them without changing behavior**, with respect to **maintainability (how easily a future contributor can understand and modify the code)**, from the viewpoint of **Alex, the software maintainer**, in the context of **Gorse v0.5.11 on the `semester-investigation` branch of the `box-bit/gorse` fork, analyzed with static analysis tools run in GitHub Actions**.

| Initial question | Question wording | Relevant metric(s) / observation(s) | Why useful? |
|---|---|---|---|
| **Q1** | Which functions in the `master` package are the hardest to understand and change, and why? | Cyclomatic and cognitive complexity per function; function length (lines/statements); nesting depth; duplicated code blocks; number of commits that changed each file in the last 12 months (churn); manual reading of the flagged code | Shows *where* the problems are. Churn shows if the code really changes often, so a problem there costs more for future contributors. Manual reading checks if a warning is a real design problem. |
| **Q2** | After the refactoring, is the selected code easier to understand and change, and does it still behave the same? | The same metrics as Q1, before and after, on the same function(s) and with the same tool settings; number of places a contributor must edit for a typical change; result of the `master` package tests run in GitHub Actions | Checks if the change really helps a future contributor, not only if a warning disappears. The tests give evidence (not proof) that the behavior did not change. |

**Code-smell types chosen from the lecture:** [individual: 3-4; pair: 5-6; selection does NOT guarantee they exist]

| Smell type | Why relevant to our GQM? | How we plan to investigate it (tool/rule/metric/inspection) |
|---|---|---|
| 1. **Long Method** — *In Go:* a long function or method (`func (m *Master) ...`), with many statements, branches and nested `if`/`for`/`switch` blocks. | A long function does many things at once. A future contributor must read all of it to change one small part, and is more likely to break something. This affects Q1 (where is the code hard to change?). | **Automatic detection:** `golangci-lint` with `funlen` (lines/statements), `gocyclo` (cyclomatic complexity), `gocognit` (cognitive complexity) and `nestif` (nesting depth), run in GitHub Actions. Then read the flagged functions in the source. |
| 2. **Duplicated Code** — *In Go:* the same or almost the same block of statements in two or more functions or files (e.g. the same error handling or request parsing copied in several REST handlers). | When the same logic exists in several places, one change must be repeated in every copy. If one copy is forgotten, the copies behave differently. This affects Q1 and Q2 (number of places to edit). | **Automatic detection:** `golangci-lint` with `dupl` (token-based clone detection), run in GitHub Actions. Then compare the copies by hand, because some similar code is intentional or too small to matter. |
| 3. **Large Class** — *In Go:* there are no classes. The equivalent is a **struct type and all its methods** (methods have a receiver, e.g. `(m *Master)`). The methods of one type can be in different files, so we measure the **type**, not the file. Candidate: the `Master` struct (27 fields, 59 methods in `master.go`, `rest.go`, `tasks.go`, `rpc.go`). | A type with too many responsibilities (REST API, scheduled tasks, gRPC, configuration) is hard to understand as a whole, and changes for different reasons touch the same type. This affects Q1 (why is the code hard to change?). | **Metric-based indication + manual inspection:** count fields and methods of the type and group the methods by responsibility. No standard Go linter reports this smell directly, so the result is a reasoned candidate, not a proven design problem. The counts come from the source at the baseline commit (GitHub permalinks). |
| 4. **Long Parameter List** — *In Go:* a function with many inputs in its signature. Names grouped with one type count separately (`a, b uint` = 2 parameters). The receiver `(m *Master)` is not counted (like `this` in Java). Go functions can also return several values; these are results, not parameters. Candidate: `LoadDataFromDatabase` in `master/tasks.go` (9 parameters, 3 results). | Many parameters make a function hard to call correctly (wrong order, many similar `[]...` arguments) and every new option changes all callers. This affects Q1 and Q2. | **Automatic detection:** `golangci-lint` with `revive` rule `argument-limit` (and `function-result-limit` for results), run in GitHub Actions. Then check the callers to see which parameters always travel together. |

**Primary GitHub-integrated tool:** [golangci-lint](https://golangci-lint.run/) v2.15.0, run in GitHub Actions with the official [golangci-lint-action](https://github.com/golangci/golangci-lint-action) (workflow [`.github/workflows/lab6-golangci.yml`](.github/workflows/lab6-golangci.yml), config [`.golangci-lab6.yml`](.golangci-lab6.yml))  
**Complementary tool (only if needed):** N/A (duplication is covered by the `dupl` linter inside golangci-lint). Regression check: `go test ./master/...` in GitHub Actions (workflow [`.github/workflows/lab6-tests.yml`](.github/workflows/lab6-tests.yml)).  
**Why these tools match our questions and selected smells:** golangci-lint is the standard Go linter runner and has one linter for the symptoms of three of the four smells: `funlen`, `gocyclo`, `gocognit`, `nestif` (Long Method), `dupl` (Duplicated Code) and `revive` `argument-limit` (Long Parameter List). It runs on GitHub for each pushed commit, so the before/after results are tied to a SHA and use the same committed config. Large Class has no Go linter rule, so it is investigated with counts and manual inspection.

## Task 2 - Detect and Investigate Code Smells

**Analysis configuration:**

- Fork branch and **baseline SHA**: `semester-investigation` @ [`8bd232d3bb391d3e831217085f39f8efddee45bb`](https://github.com/box-bit/gorse/commit/8bd232d3bb391d3e831217085f39f8efddee45bb)
- Source directory/directories analyzed: `./master/...` (package `github.com/gorse-io/gorse/master`)
- Exclusions (generated/test/vendor/weight data, etc.): `*_test.go` (`run.tests: false`); no generated code in `master/`
- Tool(s), versions and relevant enabled rules/settings: golangci-lint v2.15.0 via `golangci/golangci-lint-action@v9`, config [`.golangci-lab6.yml`](.golangci-lab6.yml): `funlen` (lines > 60 or statements > 40), `gocyclo` (> 15), `gocognit` (> 15), `nestif` (complexity >= 5), `dupl` (threshold 100 tokens), `revive` `argument-limit` (> 5) and `function-result-limit` (> 3); all other linters disabled; `uniq-by-line: false` so findings of different linters on the same function line are all kept
- **Baseline GitHub scan / Actions run link(s):** [Lab6 - golangci-lint run 44535359098](https://github.com/box-bit/gorse/actions/runs/44535359098) (success). **Baseline totals: 61 findings** = `gocognit` 20, `funlen` 16, `gocyclo` 13, `dupl` 8, `nestif` 3, `revive` 1 (full list in the run summary and artifact `lab6-golangci-report-8bd232d3bb391d3e831217085f39f8efddee45bb`)

| Selected smell type | Found / not found / not assessable | Tool result or source link | Interpretation / limitations |
|---|---|---|---|
| | | | |
| | | | |
| | | | |
| | | | |
| | | | |
| | | | |

**Candidate issues inspected in source (at least two):**

| Candidate ID | File/entity permalink | What was actually observed? | Why it might matter / why uncertain |
|---|---|---|---|
| C1 | [`(*Master).LoadDataFromDatabase` - `master/tasks.go` L274-L740](https://github.com/box-bit/gorse/blob/8bd232d3bb391d3e831217085f39f8efddee45bb/master/tasks.go#L274-L740) | One function of ~470 lines doing 6 separate steps, each marked by a `// STEP n` comment (load users, load items, load explicit negative feedback, load positive feedback, load read feedback, build the CTR dataset). Tool results on this function: `funlen` 171 statements (> 40), `gocyclo` 84 (> 15), `gocognit` 224 (> 15), `revive` argument-limit 9 parameters (> 5), and `dupl` lines [329-351](https://github.com/box-bit/gorse/blob/8bd232d3bb391d3e831217085f39f8efddee45bb/master/tasks.go#L329-L351) = [389-411](https://github.com/box-bit/gorse/blob/8bd232d3bb391d3e831217085f39f8efddee45bb/master/tasks.go#L389-L411): the same label-indexing algorithm ("remember first occurrence, add to index at second occurrence") is written once for user labels and once for item labels. Manual reading also shows near-duplicate blocks in steps 3, 4 and 5 (same feedback-stream options, same "convert ids to indexes, skip unknown" code) that `dupl` did not report. 5 call sites in 3 files (`master/tasks.go`, `master/tasks_test.go`, `cmd/gorse-bench/main.go`). | **Plausible maintainability concern.** Several independent symptoms agree (Long Method, Duplicated Code, Long Parameter List) and the code changes often: `master/tasks.go` was changed in 48 commits in the 12 months before v0.5.11, and the function signature changed several times (e.g. #1215 added `negFeedbackTypes`, so every caller had to be edited). A fix in the label logic must be made twice. Uncertainty: the 6 steps share many local variables (`dataSet`, `itemGroups`, `negativeSet`, `positiveFeedback`, `span`, `mu`), so splitting the function is not free, and part of the length is concurrency code that is intentionally inline. |
| C2 | [`(*Master).importExportUsers` L1271-L1367](https://github.com/box-bit/gorse/blob/8bd232d3bb391d3e831217085f39f8efddee45bb/master/rest.go#L1271-L1367), [`importExportItems` L1369-L1490](https://github.com/box-bit/gorse/blob/8bd232d3bb391d3e831217085f39f8efddee45bb/master/rest.go#L1369-L1490), [`importExportFeedback` L1492-L1617](https://github.com/box-bit/gorse/blob/8bd232d3bb391d3e831217085f39f8efddee45bb/master/rest.go#L1492-L1617) - `master/rest.go` | Tool results: Long Method symptoms on all three (`funlen` 67 / 77 / 75 statements, `gocyclo` 20 / 24 / 23, `gocognit` 44 / 58 / 55). No `dupl` finding. Manual reading: the three handlers have the same structure: `checkLogin` block, then a `switch` on the HTTP method; GET streams all rows as JSON lines, POST decodes the file line by line, validates, inserts in batches of `batchSize`, and triggers a new scheduled run. | **Plausible but weaker concern.** Each function is long but linear and easy to read top-down, so the complexity numbers overstate how hard it is to understand. The parallel structure means a change to the import/export protocol (e.g. a new validation or response format) must be made three times. Uncertainty: the duplication is a **manual observation, not a tool result**; the copies differ in types and validation details, so this needs judgement about how much is really shared. |

**One warning we challenged, or one potentially missed issue:** *Missed issue:* `dupl` (token-based, threshold 100 tokens) reports only 2 of the duplicated blocks in C1 and none in C2. The near-duplicate feedback-loading blocks in steps 3-5 of C1 and the parallel structure of the three C2 handlers differ in identifiers, types and small details, so token-based clone detection does not match them. A "0 duplicates" result for a function is therefore not evidence that it has no duplication. *Challenged warning:* for C2, the very high `gocognit` values (44-58) come mostly from repeated, flat `if err != nil { ...; return }` error handling, which is the normal Go style. They overstate how hard these handlers are to understand compared with C1, where the complexity comes from nested loops, shared state and concurrency.

## Task 3 - Prioritize Technical Debt and Plan Refactoring

| Candidate | Expected future impact / relevance | Estimated effort and behavior risk | Prioritize now or defer? Why? |
|---|---|---|---|
| C1a - duplicated label indexing in `LoadDataFromDatabase` (users vs items, `dupl` [L329-351](https://github.com/box-bit/gorse/blob/8bd232d3bb391d3e831217085f39f8efddee45bb/master/tasks.go#L329-L351) / [L389-411](https://github.com/box-bit/gorse/blob/8bd232d3bb391d3e831217085f39f8efddee45bb/master/tasks.go#L389-L411)) | Any change to the label rule ("keep a label only if >= 2 users/items have it, add it back to the first one") must be made twice; if one copy is forgotten, users and items silently follow different rules, with no error and (before our new test) no failing test. `master/tasks.go` is the most frequently changed file of the package (48 commits in 12 months), so this code is likely to be touched again. | **Low effort** (one helper, ~40 lines, inside one file; no signature or caller change). **Risk:** subtle - the "add back to the first user/item at the second occurrence" step is easy to break, and existing tests only check the *number* of labels. | **Now.** Small, local and behavior-preserving; removes verified duplication (tool + manual); measurable before/after. |
| C1b - Long Parameter List of `LoadDataFromDatabase` (9 parameters, 5 call sites in 3 files) | Every new data-source option changes the signature and all callers (e.g. #1215 added `negFeedbackTypes`). | **Medium effort** (Introduce Parameter Object across `master/tasks.go`, `master/tasks_test.go`, `cmd/gorse-bench/main.go`); low behavior risk (compiler-checked). | **Defer.** Real cost, but the change is spread across packages and adds less understanding benefit than C1a for a first step. |
| C1c - splitting `LoadDataFromDatabase` into its 6 steps | Largest understandability benefit (the function is ~470 lines). | **High effort and risk:** steps share many local variables; steps 3-5 run in parallel goroutines with a shared mutex. | **Defer.** Not a small, safe change; C1a is a first step in this direction. |
| C2 - three parallel `importExport*` handlers in `master/rest.go` | A protocol change must be made three times, but each handler is linear and readable; `rest.go` changes a bit less (36 commits in 12 months). | **Medium effort:** a shared helper needs generics or callbacks and touches HTTP error paths; tests exist in `rest_test.go`. | **Defer.** Lower impact; duplication only from manual inspection. |
| Large Class - `Master` struct (27 fields, 59 methods) | Root cause of many of the above (one type for REST, tasks, gRPC, config). | **Very high effort:** redesign of the central type. | **Defer.** Out of scope for a small behavior-preserving change. |

**Selected documented issue(s) and individual ownership:** C1a, owned by Alessandro Fantesini (`box-bit`).

| Contributor / GitHub handle | Candidate ID | Planned refactoring | Expected benefit | GitHub regression check planned |
|---|---|---|---|---|
| Alessandro Fantesini (`box-bit`) | C1a | **Extract Method** (removing verified duplication): move the label-indexing logic into one small helper type with an `add` method, used once for user labels and once for item labels in `LoadDataFromDatabase`. | One place (with a name and a comment) for the label rule instead of two copies; a future change is made once; `LoadDataFromDatabase` becomes shorter and less complex. Expected measurement: the 2 `dupl` findings for L329-351 / L389-411 disappear, `funlen`/`gocyclo`/`gocognit` values of `LoadDataFromDatabase` decrease (it will still be flagged - C1c is deferred). | `Lab6 - master tests` workflow (`go test ./master/...`). A new **characterization test** `TestLoadDataFromDatabaseLabels` is committed **before** the refactoring: it checks which labels each user/item receives, including the "added back to the first user/item" step. Mutation check (local): removing that step makes the new test fail while the two existing `LoadDataFromDatabase` tests still pass, which shows the gap the new test closes. |

## Task 4 - Refactor, Verify, and Commit

**Baseline test/check run(s) on GitHub:** [Lab6 - master tests run 44535359180](https://github.com/box-bit/gorse/actions/runs/44535359180): `go test ./master/...` -> `ok`, 45 tests passed, 0 failed, 0 skipped  
**After-refactoring SHA(s):** [full SHA(s)]  
**After-refactoring scan / workflow run(s):** [URL(s)]  
**After-refactoring test/check run(s):** [URL(s) and outcomes]

| Issue / contributor | Before (same scope and measurement) | After | Test / verification result | What does this actually establish? |
|---|---|---|---|---|
| | | | | |
| | | | | |

**Answers to the original GQM questions:**

- **Q1:** [answer based on evidence, not only a warning count]
- **Q2:** [answer based on evidence, not only a warning count]

**Conclusion and limits:** [What was improved? Did any measurements worsen or remain unchanged? What cannot be claimed about maintainability or preserved behavior?]

### Individual implementation commits - permanent GitHub links

> Each student creates and pushes **their own** changes from **their own machine/account**. Include direct `/commit/FULL_SHA` links, **not just repository or PR URLs**. Do not squash individual source commits into a single partner-authored commit.

| Contributor / GitHub account | Documented issue (C1, C2, ...) | What was refactored? | **Permanent commit link** | PR or GitHub Actions push-run link |
|---|---|---|---|---|
| Student 1 | | | `https://github.com/OWNER/REPO/commit/FULL_SHA` | |
| Student 1 (second issue, if applicable) | | | | |
| Student 2 (pair only) | | | | |
| Student 2 (second issue, if applicable) | | | | |

**Completed template location on `semester-investigation`:** [GitHub file URL]  
**Template/report commit link:** [GitHub `/commit/FULL_SHA` URL]

**Evidence preservation note:** [If any GitHub Actions artifact is relevant, provide a run link and preserve its essential figures here because workflow artifacts can expire.]
