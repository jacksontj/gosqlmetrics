package sqlmetrics

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

// Options for the Collector
type Options struct {
	Prefix      string
	Labels      []string
	LabelValues []string
}

type metrics struct {
	maxConnsDesc *prometheus.Desc
	// Pool Status
	openConns *prometheus.Desc
	inUse     *prometheus.Desc
	idle      *prometheus.Desc

	// Counters
	waitCount         *prometheus.Desc
	waitDuration      *prometheus.Desc
	maxIdleClosed     *prometheus.Desc
	maxLifetimeClosed *prometheus.Desc
}

// NewCollector returns a collector for the given db
func NewCollector(o Options, db *sql.DB) *Collector {
	return &Collector{
		o:  o,
		db: db,
		m: metrics{
			maxConnsDesc: prometheus.NewDesc(
				o.Prefix+"connections_max",
				"Max number of open connections to the DB",
				o.Labels, nil,
			),
			openConns: prometheus.NewDesc(
				o.Prefix+"connections_open",
				"Current number of established connections bith inuse and idle",
				o.Labels, nil,
			),
			inUse: prometheus.NewDesc(
				o.Prefix+"connections_in_use",
				"The number of connections currently in use",
				o.Labels, nil,
			),
			idle: prometheus.NewDesc(
				o.Prefix+"connections_idle",
				"The number of idle connections",
				o.Labels, nil,
			),
			waitCount: prometheus.NewDesc(
				o.Prefix+"connections_wait_count_total",
				"The total number of connections waited for",
				o.Labels, nil,
			),
			waitDuration: prometheus.NewDesc(
				o.Prefix+"connections_wait_duration_seconds_total",
				"The total time blocked waiting for a new connection in seconds",
				o.Labels, nil,
			),
			maxIdleClosed: prometheus.NewDesc(
				o.Prefix+"connections_max_idle_closed_total",
				"The total number of connections closed due to SetMaxIdleConns",
				o.Labels, nil,
			),
			maxLifetimeClosed: prometheus.NewDesc(
				o.Prefix+"connections_max_lifetime_closed_total",
				"The total number of connections closed due to SetConnMaxLifetime",
				o.Labels, nil,
			),
		},
	}
}

// Collector is a prometheus Collector which collects metrics from a sql.DB
type Collector struct {
	o  Options
	db *sql.DB
	m  metrics
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	stats := c.db.Stats()

	ch <- prometheus.MustNewConstMetric(
		c.m.maxConnsDesc,
		prometheus.GaugeValue,
		float64(stats.MaxOpenConnections),
		c.o.LabelValues...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.m.openConns,
		prometheus.GaugeValue,
		float64(stats.OpenConnections),
		c.o.LabelValues...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.m.inUse,
		prometheus.GaugeValue,
		float64(stats.InUse),
		c.o.LabelValues...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.m.idle,
		prometheus.GaugeValue,
		float64(stats.Idle),
		c.o.LabelValues...,
	)

	// Counters
	ch <- prometheus.MustNewConstMetric(
		c.m.waitCount,
		prometheus.CounterValue,
		float64(stats.WaitCount),
		c.o.LabelValues...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.m.waitDuration,
		prometheus.CounterValue,
		float64(stats.WaitDuration.Seconds()),
		c.o.LabelValues...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.m.maxIdleClosed,
		prometheus.CounterValue,
		float64(stats.MaxIdleClosed),
		c.o.LabelValues...,
	)
	ch <- prometheus.MustNewConstMetric(
		c.m.maxLifetimeClosed,
		prometheus.CounterValue,
		float64(stats.MaxLifetimeClosed),
		c.o.LabelValues...,
	)
}
