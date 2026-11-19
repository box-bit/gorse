# Lab #6 - Investigation Record

**Team / student(s):** Alessandro Fantesini (`box-bit`)  
**Mode:** individual  
**Approved upstream project:** <https://github.com/gorse-io/gorse>  
**Fork / working branch:** [`box-bit/gorse@semester-investigation`](https://github.com/box-bit/gorse/tree/semester-investigation)  
**Original Lab #5 version:** `de9141d5d87236289d8d7f102bb9588234ff63cf` (upstream tag `v0.5.11`)  
**Lab #6 pre-refactoring baseline:** source revision [`140cfb1bc72876a51356e8446e3b311687814cd0`](https://github.com/box-bit/gorse/commit/140cfb1bc72876a51356e8446e3b311687814cd0) (Go source identical to `de9141d5`; this commit only adds `LAB5_SELECTION.md`). Baseline scan commit: [to fill after the Lab 6 workflow is pushed]

> Keep this file short and evidence-based. Replace the placeholders as you go. Every analysis/test link must point to your **fork's GitHub-visible results**. Do not submit unsupported local-only tool screenshots.

## Task 1 - Stakeholder, GQM, and Code Smell Selection

**Selected stakeholder and exact fictional concern from handout:** Scenario A - Alex, Software Maintainer: *"Our software keeps evolving. Future contributors need to understand and modify its components, but certain design choices might make that work unnecessarily difficult. Please investigate possible code smells, identify maintainability concerns worth addressing, and demonstrate a focused improvement without changing the software's intended behavior."*

**GQM goal:** Analyze **the `master` package of Gorse** for the purpose of **finding code smells that make future changes harder, and fixing one of them without changing behavior**, with respect to **maintainability (how easily a future contributor can understand and modify the code)**, from the viewpoint of **Alex, the software maintainer**, in the context of **Gorse v0.5.11 on the `semester-investigation` branch of the `box-bit/gorse` fork, analyzed with static analysis tools run in GitHub Actions**.

| Initial question | Question wording | Relevant metric(s) / observation(s) | Why useful? |
|---|---|---|---|
| **Q1** | Which functions in the `master` package are the hardest to understand and change, and why? | Cyclomatic and cognitive complexity per function; function length (lines/statements); nesting depth; duplicated code blocks; number of commits that changed each file in the last 12 months (churn); manual reading of the flagged code | Shows *where* the problems are. Churn shows if the code really changes often, so a problem there costs more for future contributors. Manual reading checks if a warning is a real design problem. |
| **Q2** | After the refactoring, is the selected code easier to understand and change, and does it still behave the same? | The same metrics as Q1, before and after, on the same function(s) and with the same tool settings; number of places a contributor must edit for a typical change; result of the `master` package tests run in GitHub Actions | Checks if the change really helps a future contributor, not only if a warning disappears. The tests give evidence (not proof) that the behavior did not change. |

**Code-smell types chosen from the lecture:** [individual: 3-4; pair: 5-6; selection does NOT guarantee they exist]

| Smell type | Lecture slide(s) | Why relevant to our GQM? | How we plan to investigate it (tool/rule/metric/inspection) |
|---|---|---|---|
| 1. **Long Method** — *In Go:* a long function or method (`func (m *Master) ...`), with many statements, branches and nested `if`/`for`/`switch` blocks. | [slide] | A long function does many things at once. A future contributor must read all of it to change one small part, and is more likely to break something. This affects Q1 (where is the code hard to change?). | **Automatic detection:** `golangci-lint` with `funlen` (lines/statements), `gocyclo` (cyclomatic complexity), `gocognit` (cognitive complexity) and `nestif` (nesting depth), run in GitHub Actions. Then read the flagged functions in the source. |
| 2. **Duplicated Code** — *In Go:* the same or almost the same block of statements in two or more functions or files (e.g. the same error handling or request parsing copied in several REST handlers). | [slide] | When the same logic exists in several places, one change must be repeated in every copy. If one copy is forgotten, the copies behave differently. This affects Q1 and Q2 (number of places to edit). | **Automatic detection:** `golangci-lint` with `dupl` (token-based clone detection), run in GitHub Actions. Then compare the copies by hand, because some similar code is intentional or too small to matter. |
| 3. **Large Class** — *In Go:* there are no classes. The equivalent is a **struct type and all its methods** (methods have a receiver, e.g. `(m *Master)`). The methods of one type can be in different files, so we measure the **type**, not the file. Candidate: the `Master` struct (27 fields, 59 methods in `master.go`, `rest.go`, `tasks.go`, `rpc.go`). | [slide] | A type with too many responsibilities (REST API, scheduled tasks, gRPC, configuration) is hard to understand as a whole, and changes for different reasons touch the same type. This affects Q1 (why is the code hard to change?). | **Metric-based indication + manual inspection:** count fields and methods of the type and group the methods by responsibility. No standard Go linter reports this smell directly, so the result is a reasoned candidate, not a proven design problem. The counts come from the source at the baseline commit (GitHub permalinks). |
| 4. **Long Parameter List** — *In Go:* a function with many inputs in its signature. Names grouped with one type count separately (`a, b uint` = 2 parameters). The receiver `(m *Master)` is not counted (like `this` in Java). Go functions can also return several values; these are results, not parameters. Candidate: `LoadDataFromDatabase` in `master/tasks.go` (9 parameters, 3 results). | [slide] | Many parameters make a function hard to call correctly (wrong order, many similar `[]...` arguments) and every new option changes all callers. This affects Q1 and Q2. | **Automatic detection:** `golangci-lint` with `revive` rule `argument-limit` (and `function-result-limit` for results), run in GitHub Actions. Then check the callers to see which parameters always travel together. |

**Primary GitHub-integrated tool:** [golangci-lint](https://golangci-lint.run/) v2.15.0, run in GitHub Actions with the official [golangci-lint-action](https://github.com/golangci/golangci-lint-action) (workflow [`.github/workflows/lab6-golangci.yml`](.github/workflows/lab6-golangci.yml), config [`.golangci-lab6.yml`](.golangci-lab6.yml))  
**Complementary tool (only if needed):** N/A (duplication is covered by the `dupl` linter inside golangci-lint). Regression check: `go test ./master/...` in GitHub Actions (workflow [`.github/workflows/lab6-tests.yml`](.github/workflows/lab6-tests.yml)).  
**Why these tools match our questions and selected smells:** golangci-lint is the standard Go linter runner and has one linter for the symptoms of three of the four smells: `funlen`, `gocyclo`, `gocognit`, `nestif` (Long Method), `dupl` (Duplicated Code) and `revive` `argument-limit` (Long Parameter List). It runs on GitHub for each pushed commit, so the before/after results are tied to a SHA and use the same committed config. Large Class has no Go linter rule, so it is investigated with counts and manual inspection.

## Task 2 - Detect and Investigate Code Smells

**Analysis configuration:**

- Fork branch and **baseline SHA**: [full SHA]
- Source directory/directories analyzed: [paths]
- Exclusions (generated/test/vendor/weight data, etc.): [paths or none]
- Tool(s), versions and relevant enabled rules/settings: [details]
- **Baseline GitHub scan / Actions run link(s):** [direct links]

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
| C1 | | | |
| C2 | | | |
| C3 (optional) | | | |

**One warning we challenged, or one potentially missed issue:** [brief explanation of tool limitations, possible false alarm or incomplete coverage]

## Task 3 - Prioritize Technical Debt and Plan Refactoring

| Candidate | Expected future impact / relevance | Estimated effort and behavior risk | Prioritize now or defer? Why? |
|---|---|---|---|
| C1 | | | |
| C2 | | | |
| Additional (if needed) | | | |

**Selected documented issue(s) and individual ownership:** [1-2 per student; pair total 2-4]

| Contributor / GitHub handle | Candidate ID | Planned refactoring | Expected benefit | GitHub regression check planned |
|---|---|---|---|---|
| | | | | |
| [second member if pair] | | | | |

## Task 4 - Refactor, Verify, and Commit

**Baseline test/check run(s) on GitHub:** [URL(s) and outcomes]  
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
