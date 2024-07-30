package processor

import (
	"math"
	"testing"
	"time"

	"github.com/openshift/cluster-health-analyzer/pkg/prom"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)

// TODOs:
// - write basic tests
// - rename GroupsCollection to IncidentsMapper

func TestGroupsCollectionProcessAlertsBatch(t *testing.T) {
	start := model.TimeFromUnixNano(
		time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC).UnixNano())

	gc := GroupsCollection{}

	// Time-based matcher with infinite distance: it should match only when
	// the alert is within the [processor.timeMatchTimeDelta] time range.
	timeMatcher := GroupMatcher{
		RootGroupID: "time-matcher",
		GroupID:     "time-matcher",
		Start:       start.Add(1 * time.Hour),
		Modified:    start.Add(1 * time.Hour),
		End:         start.Add(3 * time.Hour),
		Distance:    math.Inf(1),
	}
	gc.AddGroup(&timeMatcher)

	// Case 1: Alert is within the time range of time-based matcher.
	//
	// It should match the original group.
	alerts := []prom.Alert{
		{Name: "Alert1", Labels: map[string]string{"alertname": "Alert1"}},
	}
	case1 := gc.ProcessAlertsBatch(alerts, start.Add(1*time.Hour+10*time.Minute).Time())

	assert.Equal(t, "time-matcher", case1[0].Labels["group_id"])

	// Case 2: 2 alerts outside of the time range of time-based matcher,
	//
	// They should not match the original group, but they should both become part
	// of a new group.
	alerts = []prom.Alert{
		{Name: "Alert2", Labels: map[string]string{"alertname": "Alert2"}},
		{Name: "Alert3", Labels: map[string]string{"alertname": "Alert3"}},
	}
	case2 := gc.ProcessAlertsBatch(alerts, start.Add(3*time.Hour).Time())
	assert.NotEqual(t, "time-matcher", case2[0].Labels["group_id"])
	assert.Equal(t, case2[0].Labels["group_id"], case2[1].Labels["group_id"])

	// Case 3: Alert with same alertname as one from case 2 fires within
	// [processor.fuzzyMatchTimeDelta] time range.
	//
	// It should match the group created in case 2.
	alerts = []prom.Alert{
		{Name: "Alert2", Labels: map[string]string{"alertname": "Alert2"}},
	}
	case3 := gc.ProcessAlertsBatch(alerts, start.Add(7*time.Hour).Time())
	assert.Equal(t, case2[0].Labels["group_id"], case3[0].Labels["group_id"])
}

func TestGroupsCollectionPruneGroups(t *testing.T) {
}

func TestGroupsCollectionProcessHistoricalAlerts(t *testing.T) {
	// rv := prom.RangeVector{
	// 	prom.Range{
	// 		Metric: prom.LabelSet{
	// 			Labels: map[string]string{
	// 				"alertname": "SomeAlert",
	// 			},
	// 		},
	// 		Samples: []model.SamplePair{{start.Add(1 * time.Hour), 1}},
	// 	},
	// }
}
