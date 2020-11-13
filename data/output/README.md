This folder contains per-state CSV files with discrete changes to vote totals as reported by NYT.

To generate the CSV files locally, simply
```
go run update.go
```
from within `data/output`.

You can find the latest version of these data on Google Sheets [here](https://docs.google.com/spreadsheets/d/1Ez2bupWlmf7V-nE17ScjQzm-Y_Cyi1qjT9ZO9hDDOC4/edit?usp=sharing).

## Columns
Each CSV column is defined as:
* `time` - timestamp of the vote event adjusted to Eastern time
* `pct` - percentage of the total expected vote processed
* `votes` - total number of votes tallied
* `biden pct` - percentage of `votes` belonging to Biden
* `trump pct` - percentage of `votes` belonging to Trump
* `other pct` - percentage of `votes` belonging to third-party candidates
* `biden tot` - total number of `votes` belonging to Biden
* `trump tot` - total number of `votes` belonging to Trump
* `other tot` - total number of `votes` belonging to third-party candidates
* `batch` - number of `votes` added/removed in this event
* `biden bat` - number of `votes` added/removed from Biden
* `trump bat` - number of `votes` added/removed from Trump
* `other bat` - number of `votes` added/removed from third-perty candidates
* `note` - summary text for any anomalies detected for the event

Keep in mind that precision of NYT's data is to three decimal places -- 0.001.  As such, small changes in totals should be considered noise.  This analysis only reports on vote changes that exceed 0.1% of the event `votes`.