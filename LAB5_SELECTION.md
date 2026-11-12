# Lab #5 - System Selection and Reproducibility Record

## People and collaboration

- Working mode: [ ] Solo  [x] Pair
- Student A (name / GitHub): Alessandro Fantesini / [box-bit](https://github.com/box-bit)
- Student B (name / GitHub; if pair): Ondrej Müller / [Hogigo](https://github.com/Hogigo)
- Fork owner: box-bit
- Collaborator invitation accepted? <!-- TODO: Yes / No -->
- Individual work differentiation: <!-- TODO: links to each student's QUALITY_INVESTIGATION.md (later labs) -->

## First shortlist - catalogue only

| Candidate | Catalogue-based reason | One concern to investigate |
|---|---|---|
| Gorse (#6) | Recommender engine with REST API and dashboard, written in Go, Apache-2.0 | Needs Docker/DB services; is there a held-out benchmark? |
| Ultralytics YOLO (#18) | Detection/segmentation with built-in validation; CPU inference works | AGPL-3.0; COCO8 only checks plumbing |
| FiftyOne (#19) | CV evaluation platform; quickstart has reference boxes + predictions | Source work needs Node/Yarn app; not a model itself |
| MMDetection (#21) | Detection benchmarks with COCO ground truth and mAP | MMCV build constraints; CPU-only builds miss some operators |

## Evidence-based candidate assessment

Use Y = Yes; P = Partial; N = No; U = Unknown. Repository data checked via the GitHub API on 2026-11-12; Gorse also checked in source at v0.5.11.

| Criterion | A: Gorse | B: Ultralytics YOLO | C: FiftyOne | D: MMDetection |
|---|---|---|---|---|
| C1 Public repository / licence / data rights | **Y**: [Apache-2.0](https://github.com/gorse-io/gorse/blob/v0.5.11/LICENSE); datasets have their own terms | **P**: [AGPL-3.0](https://github.com/ultralytics/ultralytics/blob/main/LICENSE) + CLA | **Y**: [Apache-2.0](https://github.com/voxel51/fiftyone/blob/develop/LICENSE); COCO terms | **Y**: [Apache-2.0](https://github.com/open-mmlab/mmdetection/blob/main/LICENSE); weight terms vary |
| C2 Pinned revision / runtime | **Y**: [v0.5.11](https://github.com/gorse-io/gorse/releases/tag/v0.5.11), Go ≥ 1.26 | **Y**: v8.4.211, pip; near-daily releases | **P**: v1.22.1; source needs Node/Yarn + MongoDB | **P**: v3.3.0 (2024); [requires `mmcv<2.2.0`](https://github.com/open-mmlab/mmdetection/blob/main/mmdet/__init__.py) |
| C3 Nontrivial source + refactoring | **Y**: ~45k lines of Go (non-test); [`.golangci.yml`](https://github.com/gorse-io/gorse/blob/v0.5.11/.golangci.yml) | **Y**: large Python package | **Y**: large Python + TypeScript | **Y**: modular `mmdet/` |
| C4 Tests and expected-output **oracle** | **Y**: 89 test files; [`model_test.go:48`](https://github.com/gorse-io/gorse/blob/v0.5.11/model/cf/model_test.go#L48) asserts NDCG ≈ 0.36 | **Y**: `tests/`, run in CI | **Y**: `tests/`, pytest, e2e | **P**: some tests need compiled operators |
| C5 AI/recommender **ground-truth** source and terms | **Y**: pre-split ml-100k / ml-1m / pinterest-20 / frappe ([`built_in.go`](https://github.com/gorse-io/gorse/blob/v0.5.11/model/built_in.go)) | **Y**: COCO val2017 labels | **Y**: 200 quickstart images with boxes | **Y**: COCO-format annotations |
| C6 Evaluation procedure and metrics | **Y**: NDCG, Precision, Recall, HR, MAP, MRR, AUC ([`evaluator.go`](https://github.com/gorse-io/gorse/blob/v0.5.11/model/cf/evaluator.go)) | **Y**: mAP, precision, recall ([docs](https://docs.ultralytics.com/modes/val/)) | **Y**: [`evaluate_detections`](https://docs.voxel51.com/tutorials/evaluate_detections.html) | **Y**: COCO mAP |
| C7 Upstream GitHub history | **Y**: active (push 2026-11-12), ~70 contributors | **Y**: very active, many bot commits | **Y**: very active | **P**: no push since 2024-08 |
| C8 CI/DevOps possibility | **Y**: [GitHub Actions](https://github.com/gorse-io/gorse/tree/v0.5.11/.github/workflows) tests + lint + Docker; Prometheus/OpenTelemetry | **Y**: CI, Docker, docs | **Y**: ~23 workflows | **P**: CircleCI only |
| C9 Laptop setup feasibility | **Y**: built and ran on CPU (see below) | **Y**: CPU / small GPU (not tried) | **P**: heavy source setup (not tried) | **P**: high build risk (not tried) |
| C10 Continuation to Labs 6-9 / credible question | **Y**: see Decision | **Y** | **P**: evaluates, doesn't predict | **P**: setup risk |

### Ground-truth verification notes

- **Gorse:**
  - Functional oracle: [`model_test.go:48`](https://github.com/gorse-io/gorse/blob/v0.5.11/model/cf/model_test.go#L48) `assert.InDelta(t, 0.36, score.NDCG, 0.01)`.
  - Ground truth: held-out implicit feedback in MovieLens ml-1m and the other built-in datasets.
  - Protocol: fixed train/test files; held-out positives are ranked against sampled negatives ([`evaluator.go:35-55`](https://github.com/gorse-io/gorse/blob/v0.5.11/model/cf/evaluator.go#L35-L55)). Sampled negatives may inflate scores, which is a point to check in Lab #8.
  - Terms: MovieLens under [GroupLens terms](https://grouplens.org/datasets/movielens/) (no redistribution); other datasets **U**.
  - Playground API output is not ground truth.
- **Ultralytics YOLO:** oracle = `tests/`; ground truth = COCO val2017; terms = COCO (CC BY 4.0 annotations).
- **FiftyOne:** oracle = `tests/`; ground truth = quickstart boxes; protocol = `evaluate_detections`; terms = COCO.
- **MMDetection:** oracle = `tests/`; ground truth = COCO annotations; terms per dataset/weight.

## Decision

- **Final SUT and short description:** Gorse v0.5.11, an open-source AI recommender system engine in Go (REST API, dashboard, built-in model training and evaluation).
- **Why it is better suited than rejected alternatives:**
  - Only candidate run end to end from source on a laptop CPU.
  - Has built-in ground truth and ranking metrics.
  - Has an active upstream, existing CI and observability.
  - Supports all of Labs 6-9: #6 golangci-lint/gocyclo, #7 GitHub Actions + Prometheus, #8 model comparison with NDCG/Recall, #9 upstream history mining.
- **Rejected alternatives and reasons:**
  - FiftyOne: not a model; heavy source setup.
  - MMDetection: inactive since 2024, old pinned dependencies.
- **Backup SUT and trigger for switching:** Ultralytics YOLO, if Gorse can't be built on our laptops or can't support a credible Lab #8 question.
- **Biggest unresolved risk / mitigation:**
  - Full test suite needs many databases: run package-level tests locally (`go test ./model/...`) and leave the rest to CI.
  - Go instead of Python: use Go tooling.

## Repository provenance and pair access

- **Original upstream repo URL:** https://github.com/gorse-io/gorse
- **Selected original tag (if any):** [`v0.5.11`](https://github.com/gorse-io/gorse/releases/tag/v0.5.11) (2026-07-14), commit `de9141d5d87236289d8d7f102bb9588234ff63cf`. We chose it over `v0.6.0-alpha.1` because that is a pre-release whose [release notes](https://github.com/gorse-io/gorse/releases/tag/v0.6.0-alpha.1) warn that "APIs, configurations, and data storage might change between pre-release versions".
- **Fork URL:** https://github.com/box-bit/gorse
- **Working branch:** [`semester-investigation`](https://github.com/box-bit/gorse/tree/semester-investigation)
- **Git and runtime versions:** git 2.53.0, Go 1.27.3 linux/amd64 (`go.mod` requires 1.26)
- **Applicable code licence (link + note):** [Apache-2.0](https://github.com/gorse-io/gorse/blob/v0.5.11/LICENSE), permissive; keep notices.
- **External dataset licence/limits:** MovieLens under GroupLens terms; other datasets unknown. None committed to the fork.
- **Partner collaborator invitation accepted? (pair only):** <!-- TODO: Yes / No -->

Setup commands:

```bash
git clone https://github.com/box-bit/gorse.git && cd gorse
git remote add upstream https://github.com/gorse-io/gorse.git
git fetch upstream --tags
git switch -c semester-investigation v0.5.11
git rev-parse HEAD   # de9141d5d87236289d8d7f102bb9588234ff63cf
git push -u origin semester-investigation
```

Issues, PRs and reviews are not copied into the fork; the Lab #9 history study uses the upstream repository.

## Actual execution - not just installation

| Check | Official documentation link | Exact local command / path + input | Observed output / short log + evidence link | Status |
|---|---|---|---|---|
| Install/build source-aware baseline | [README](https://github.com/gorse-io/gorse/blob/v0.5.11/README.md) | `go build ./cmd/gorse-in-one` on `semester-investigation` | Build OK; reports API v0.2.7, Go 1.27.3 | Success |
| Functionality 1: train + evaluate BPR on ml-1m | [`model_test.go`](https://github.com/gorse-io/gorse/blob/v0.5.11/model/cf/model_test.go) | `go test ./model/cf/ -run 'TestBPR_MovieLens$' -count=1 -v` (ml-1m auto-downloaded) | `NDCG@10=0.361, Precision@10=0.063, Recall@10=0.631`; `--- PASS (33.38s)` | Success |
| Functionality 2: personalised recommendations via REST API | [README Quick Start](https://github.com/gorse-io/gorse/blob/v0.5.11/README.md#quick-start) | Commands below: start the playground, send 5 LLM-repo "star" events for user `bob`, request 10 recommendations | Feedback: `{"RowAffected": 5}`. After ~2 min training (`fit bpr complete`): litgpt, unsloth, onnxruntime, deepseek-v3, web-llm, ktransformers… (LLM-related, as documented) | Success (workaround) |

Commands for Functionality 2 (run from the repo root):

```bash
# Terminal 1: build and start the playground (downloads ~114 MB of GitHub-repo data)
go build ./cmd/gorse-in-one
mkdir -p "$PWD/gorse-cache" "$PWD/gorse-blob"
GORSE_BLOB_URI="$PWD/gorse-blob" ./gorse-in-one --playground --cache-path "$PWD/gorse-cache"

# Terminal 2: insert 5 "star" events for user "bob"
curl -X POST http://127.0.0.1:8088/api/feedback -H 'Content-Type: application/json' -d '[
 {"FeedbackType":"star","UserId":"bob","ItemId":"ollama:ollama","Value":1.0,"Timestamp":"2022-02-24"},
 {"FeedbackType":"star","UserId":"bob","ItemId":"huggingface:transformers","Value":1.0,"Timestamp":"2022-02-25"},
 {"FeedbackType":"star","UserId":"bob","ItemId":"rasbt:llms-from-scratch","Value":1.0,"Timestamp":"2022-02-26"},
 {"FeedbackType":"star","UserId":"bob","ItemId":"vllm-project:vllm","Value":1.0,"Timestamp":"2022-02-27"},
 {"FeedbackType":"star","UserId":"bob","ItemId":"hiyouga:llama-factory","Value":1.0,"Timestamp":"2022-02-28"}
]'
# Output:
# {
#  "RowAffected": 5
# }

# After "fit bpr complete" appears in terminal 1 (~2 min), fetch recommendations
curl "http://127.0.0.1:8088/api/recommend/bob?n=10"
```

- **Is the installed code from the local source fork, or only a separately installed wheel/image?** Local source from the fork's `semester-investigation` branch (commit `de9141d`), built with `go build`, not the Docker image.
- **Blockers and fallback date/action:**
  - `go: download go1.26 for linux/amd64: toolchain not available`: the apt `golang-go` package is older than the `go 1.26` required by `go.mod`. Fix: `go env -w GOTOOLCHAIN=go1.26.0` or install Go from https://go.dev/dl/.
  - Default blob path `/var/lib/gorse/blob` isn't writable ([`config.toml:365`](https://github.com/gorse-io/gorse/blob/v0.5.11/config/config.toml#L365)): set `GORSE_BLOB_URI`.
  - `failed to init meta database: unable to open database file`: `--cache-path` must be an absolute, existing folder.
  - The log shows `cannot fetch embedding from <nil>` for some items.


**Date of record:** 2026-11-12
