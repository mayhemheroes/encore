package custommetrics

import (
	"github.com/rs/zerolog"
)

// CloudWatch logs-based metrics support up to three dimensions:
//
// https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html#logs-metric-filters-dimensions
//
// For this reason, we check that the number of tags passed in by the caller
// isn't greater than three. This doesn't cover the case where the same metric is
// published with different sets of tags over multiple calls to the functions
// defined in this file.

type awsMetricsManager struct {
	logger zerolog.Logger
}

func (m *awsMetricsManager) Counter(name string, tags map[string]string) {
	// See comment above.
	if len(tags) > 3 {
		m.logger.Trace().Str("dropped_metric_name", name).Msg("dropping metric")
		return
	}

	loggerCtx := m.logger.With().Str("metric_name", name)
	for k, v := range tags {
		loggerCtx = loggerCtx.Str(k, v)
	}
	logger := loggerCtx.Logger()
	logger.Trace().Send()
}
