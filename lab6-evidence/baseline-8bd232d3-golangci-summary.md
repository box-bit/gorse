## Lab 6 - golangci-lint results

- Commit: `8bd232d3bb391d3e831217085f39f8efddee45bb`
- Branch: `semester-investigation`
- Scope: `./master/...` (non-test code), config `.golangci-lab6.yml`
- Total findings: **61**

| Linter | Findings |
|---|---|
| dupl | 8 |
| funlen | 16 |
| gocognit | 20 |
| gocyclo | 13 |
| nestif | 3 |
| revive | 1 |

<details><summary>All findings</summary>

| Linter | Location | Message |
|---|---|---|
| funlen | `master/master.go:152` | Function 'Serve' has too many statements (57 > 40) |
| gocognit | `master/master.go:152` | cognitive complexity 40 of func `(*Master).Serve` is high (> 15) |
| gocyclo | `master/master.go:152` | cyclomatic complexity 28 of func `(*Master).Serve` is high (> 15) |
| dupl | `master/master.go:207` | 207-218 lines are duplicate of `master/master.go:222-233` |
| dupl | `master/master.go:222` | 222-233 lines are duplicate of `master/master.go:207-218` |
| gocognit | `master/metrics.go:242` | cognitive complexity 16 of func `(*OnlineEvaluator).Add` is high (> 15) |
| gocognit | `master/metrics.go:277` | cognitive complexity 21 of func `(*OnlineEvaluator).Evaluate` is high (> 15) |
| funlen | `master/rest.go:81` | Function 'CreateWebService' is too long (165 > 60) |
| dupl | `master/rest.go:184` | 184-192 lines are duplicate of `master/rest.go:201-209` |
| dupl | `master/rest.go:201` | 201-209 lines are duplicate of `master/rest.go:184-192` |
| gocognit | `master/rest.go:435` | cognitive complexity 24 of func `(*Master).checkLogin` is high (> 15) |
| nestif | `master/rest.go:439` | `if m.Config.OIDC.Enable` has complex nested blocks (complexity: 13) |
| nestif | `master/rest.go:466` | `if m.Config.OIDC.Enable` has complex nested blocks (complexity: 9) |
| gocognit | `master/rest.go:645` | cognitive complexity 17 of func `(*Master).getStats` is high (> 15) |
| gocyclo | `master/rest.go:645` | cyclomatic complexity 18 of func `(*Master).getStats` is high (> 15) |
| funlen | `master/rest.go:857` | Function 'getRecommend' is too long (61 > 60) |
| gocognit | `master/rest.go:857` | cognitive complexity 16 of func `(*Master).getRecommend` is high (> 15) |
| funlen | `master/rest.go:934` | Function 'getUserFeedback' has too many statements (55 > 40) |
| gocognit | `master/rest.go:934` | cognitive complexity 20 of func `(*Master).getUserFeedback` is high (> 15) |
| gocyclo | `master/rest.go:934` | cyclomatic complexity 16 of func `(*Master).getUserFeedback` is high (> 15) |
| funlen | `master/rest.go:1148` | Function 'getRankerPrompt' has too many statements (75 > 40) |
| gocognit | `master/rest.go:1148` | cognitive complexity 29 of func `(*Master).getRankerPrompt` is high (> 15) |
| gocyclo | `master/rest.go:1148` | cyclomatic complexity 24 of func `(*Master).getRankerPrompt` is high (> 15) |
| funlen | `master/rest.go:1271` | Function 'importExportUsers' has too many statements (67 > 40) |
| gocognit | `master/rest.go:1271` | cognitive complexity 44 of func `(*Master).importExportUsers` is high (> 15) |
| gocyclo | `master/rest.go:1271` | cyclomatic complexity 20 of func `(*Master).importExportUsers` is high (> 15) |
| funlen | `master/rest.go:1369` | Function 'importExportItems' has too many statements (77 > 40) |
| gocognit | `master/rest.go:1369` | cognitive complexity 58 of func `(*Master).importExportItems` is high (> 15) |
| gocyclo | `master/rest.go:1369` | cyclomatic complexity 24 of func `(*Master).importExportItems` is high (> 15) |
| funlen | `master/rest.go:1492` | Function 'importExportFeedback' has too many statements (75 > 40) |
| gocognit | `master/rest.go:1492` | cognitive complexity 55 of func `(*Master).importExportFeedback` is high (> 15) |
| gocyclo | `master/rest.go:1492` | cyclomatic complexity 23 of func `(*Master).importExportFeedback` is high (> 15) |
| funlen | `master/rest.go:1715` | Function 'dump' has too many statements (62 > 40) |
| gocognit | `master/rest.go:1715` | cognitive complexity 33 of func `(*Master).dump` is high (> 15) |
| gocyclo | `master/rest.go:1715` | cyclomatic complexity 21 of func `(*Master).dump` is high (> 15) |
| funlen | `master/rest.go:1826` | Function 'Restore' has too many statements (69 > 40) |
| gocognit | `master/rest.go:1826` | cognitive complexity 97 of func `(*Master).Restore` is high (> 15) |
| gocyclo | `master/rest.go:1826` | cyclomatic complexity 30 of func `(*Master).Restore` is high (> 15) |
| funlen | `master/tasks.go:51` | Function 'loadDataset' has too many statements (88 > 40) |
| gocognit | `master/tasks.go:51` | cognitive complexity 34 of func `(*Master).loadDataset` is high (> 15) |
| gocyclo | `master/tasks.go:51` | cyclomatic complexity 25 of func `(*Master).loadDataset` is high (> 15) |
| gocognit | `master/tasks.go:229` | cognitive complexity 21 of func `(*Master).runLoadDatasetTask` is high (> 15) |
| gocyclo | `master/tasks.go:229` | cyclomatic complexity 17 of func `(*Master).runLoadDatasetTask` is high (> 15) |
| funlen | `master/tasks.go:274` | Function 'LoadDataFromDatabase' has too many statements (171 > 40) |
| gocognit | `master/tasks.go:274` | cognitive complexity 224 of func `(*Master).LoadDataFromDatabase` is high (> 15) |
| gocyclo | `master/tasks.go:274` | cyclomatic complexity 84 of func `(*Master).LoadDataFromDatabase` is high (> 15) |
| revive | `master/tasks.go:274` | argument-limit: maximum number of arguments per function exceeded; max 5 but got 9 |
| dupl | `master/tasks.go:329` | 329-351 lines are duplicate of `master/tasks.go:389-411` |
| dupl | `master/tasks.go:389` | 389-411 lines are duplicate of `master/tasks.go:329-351` |
| nestif | `master/tasks.go:463` | `if len(negFeedbackTypes) > 0` has complex nested blocks (complexity: 6) |
| funlen | `master/tasks.go:742` | Function 'updateItemToItem' is too long (72 > 60) |
| gocognit | `master/tasks.go:742` | cognitive complexity 30 of func `(*Master).updateItemToItem` is high (> 15) |
| funlen | `master/tasks.go:859` | Function 'updateUserToUser' is too long (63 > 60) |
| gocognit | `master/tasks.go:859` | cognitive complexity 29 of func `(*Master).updateUserToUser` is high (> 15) |
| funlen | `master/tasks.go:963` | Function 'trainCollaborativeFiltering' has too many statements (68 > 40) |
| gocognit | `master/tasks.go:963` | cognitive complexity 22 of func `(*Master).trainCollaborativeFiltering` is high (> 15) |
| gocyclo | `master/tasks.go:963` | cyclomatic complexity 21 of func `(*Master).trainCollaborativeFiltering` is high (> 15) |
| funlen | `master/tasks.go:1111` | Function 'trainClickThroughRatePrediction' has too many statements (52 > 40) |
| gocognit | `master/tasks.go:1249` | cognitive complexity 22 of func `(*Master).collectGarbage` is high (> 15) |
| dupl | `master/tasks.go:1262` | 1262-1275 lines are duplicate of `master/tasks.go:1276-1289` |
| dupl | `master/tasks.go:1276` | 1276-1289 lines are duplicate of `master/tasks.go:1262-1275` |

</details>
