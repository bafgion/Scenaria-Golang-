# Desktop smoke pass for remaining ROADMAP manual QA items

This runbook covers the 16 remaining manual QA items that are expected to be
validated in a real desktop session, after the automated E2E/Go coverage has
already passed.

Scope:
- `§1 Monaco Tabs`: 2 items
- `§2 Batch`: 6 items
- `§5 Backend storage`: 3 items
- `§6 Recorder`: 1 item
- `§7 Reports`: 3 items
- `§8 Runner lifecycle`: 1 item

Out of scope for this pass:
- `§9 Performance` items

Preflight:
- `./scripts/desktop-smoke.ps1` should already be green.
- Use a real desktop session, not mock Wails.
- Keep the app visible and avoid browser devtools unless you need them for verification.

Recommended order:
1. Monaco tabs
2. Batch
3. Backend storage
4. Recorder
5. Reports
6. Runner lifecycle

## Checklist

| # | ROADMAP item | What to verify | Result |
|---|---|---|---|
| 1 | Open 5 files | Five editor tabs can be opened and remain distinct. | [ ] |
| 2 | Close first and last tab | Closing the outer tabs leaves the middle tabs intact and active tab selection remains sane. | [ ] |
| 3 | Add 3 more tests and run | Batch selection can be expanded to 5 items and the run starts with the updated selection. | [ ] |
| 4 | Verify all 5 execute | All selected cases are included in the run output. | [ ] |
| 5 | Remove part of the tests and run | Batch selection can be reduced and the next run respects the new subset. | [ ] |
| 6 | Hotkey right after batch mode | The batch-run hotkey works immediately after enabling batch mode. | [ ] |
| 7 | Folder selection | Folder picker/open project flow still selects the intended folder. | [ ] |
| 8 | Refresh after file delete | Project refresh after deleting a file updates the UI without stale entries. | [ ] |
| 9 | Run with unsaved/temp feature | A run with an unsaved or temp feature behaves correctly and does not lose the target. | [ ] |
| 10 | Settings save with recents/HTTP auth | Saving settings while recents/HTTP auth state is active does not corrupt or lose settings. | [ ] |
| 11 | Project refresh during save/rename/delete | Refreshing the project while save/rename/delete is in progress does not break the workspace. | [ ] |
| 12 | Close browser during picker | Closing the browser while picker is active does not leave the app hanging. | [ ] |
| 13 | Latest pointer + old report | A second report run updates latest pointer, while the older report stays readable. | [ ] |
| 14 | Canceled run partial report | A canceled run still produces a usable partial report view. | [ ] |
| 15 | Allure after failed write simulation | Failed write handling keeps the failure visible and does not silently corrupt Allure output. | [ ] |
| 16 | Close app during run | Closing the app while a run is active shuts down cleanly. | [ ] |

Note:
- Item numbering above intentionally keeps the roadmap wording, but the
  pass itself is the 16-item desktop set. Performance items are tracked
  separately and should not be mixed into this run.

## Evidence to capture

- For each item, record:
  - pass/fail
  - short note if behavior is surprising
  - screenshot or log snippet when the outcome is not obvious
- If an item fails, leave the roadmap checkbox open and add the observed
  repro detail before continuing.
