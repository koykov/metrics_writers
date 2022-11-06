package queue

import (
	"time"

	q "github.com/koykov/queue"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusMetrics is a Prometheus implementation of queue.MetricsWriter.
type PrometheusMetrics struct {
	name string
	prec time.Duration
}

var (
	promQueueSize, promWorkerIdle, promWorkerActive, promWorkerSleep        *prometheus.GaugeVec
	promQueueIn, promQueueOut, promQueueRetry, promQueueLeak, promQueueLost *prometheus.CounterVec

	promWorkerWait *prometheus.HistogramVec

	_ = NewPrometheusMetrics
)

func init() {
	promWorkerIdle = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "queue_workers_idle",
		Help: "Indicates how many workers idle.",
	}, []string{"queue"})
	promWorkerActive = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "queue_workers_active",
		Help: "Indicates how many workers active.",
	}, []string{"queue"})
	promWorkerSleep = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "queue_workers_sleep",
		Help: "Indicates how many workers sleep.",
	}, []string{"queue"})

	promQueueSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "queue_size",
		Help: "Actual queue size.",
	}, []string{"queue"})

	promQueueIn = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "queue_in",
		Help: "How many items comes to the queue.",
	}, []string{"queue"})
	promQueueOut = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "queue_out",
		Help: "How many items leaves queue.",
	}, []string{"queue"})
	promQueueRetry = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "queue_retry",
		Help: "How many retries occurs.",
	}, []string{"queue"})
	promQueueLeak = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "queue_leak",
		Help: "How many items dropped on the floor due to queue is full.",
	}, []string{"queue"})
	promQueueLost = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "queue_lost",
		Help: "How many items throw to the trash due to force close.",
	}, []string{"queue"})

	buckets := append(prometheus.DefBuckets, []float64{15, 20, 30, 40, 50, 100, 150, 200, 250, 500, 1000, 1500, 2000, 3000, 5000}...)
	promWorkerWait = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "queue_wait",
		Help:    "How many worker waits due to delayed execution.",
		Buckets: buckets,
	}, []string{"queue"})

	prometheus.MustRegister(promWorkerIdle, promWorkerActive, promWorkerSleep, promQueueSize,
		promQueueIn, promQueueOut, promQueueRetry, promQueueLeak, promQueueLost,
		promWorkerWait)
}

func NewPrometheusMetrics(name string) *PrometheusMetrics {
	return NewPrometheusMetricsWP(name, time.Nanosecond)
}

func NewPrometheusMetricsWP(name string, precision time.Duration) *PrometheusMetrics {
	if precision == 0 {
		precision = time.Nanosecond
	}
	m := &PrometheusMetrics{
		name: name,
		prec: precision,
	}
	return m
}

func (m PrometheusMetrics) WorkerSetup(active, sleep, stop uint) {
	promWorkerActive.DeleteLabelValues(m.name)
	promWorkerSleep.DeleteLabelValues(m.name)
	promWorkerIdle.DeleteLabelValues(m.name)

	promWorkerActive.WithLabelValues(m.name).Add(float64(active))
	promWorkerSleep.WithLabelValues(m.name).Add(float64(sleep))
	promWorkerIdle.WithLabelValues(m.name).Add(float64(stop))
}

func (m PrometheusMetrics) WorkerInit(_ uint32) {
	promWorkerActive.WithLabelValues(m.name).Inc()
	promWorkerIdle.WithLabelValues(m.name).Add(-1)
}

func (m PrometheusMetrics) WorkerSleep(_ uint32) {
	promWorkerSleep.WithLabelValues(m.name).Inc()
	promWorkerActive.WithLabelValues(m.name).Add(-1)
}

func (m PrometheusMetrics) WorkerWakeup(_ uint32) {
	promWorkerActive.WithLabelValues(m.name).Inc()
	promWorkerSleep.WithLabelValues(m.name).Add(-1)
}

func (m PrometheusMetrics) WorkerWait(_ uint32, delay time.Duration) {
	promWorkerWait.WithLabelValues(m.name).Observe(float64(delay.Nanoseconds() / int64(m.prec)))
}

func (m PrometheusMetrics) WorkerStop(_ uint32, force bool, status q.WorkerStatus) {
	promWorkerIdle.WithLabelValues(m.name).Inc()
	if force {
		switch status {
		case q.WorkerStatusActive:
			promWorkerActive.WithLabelValues(m.name).Add(-1)
		case q.WorkerStatusSleep:
			promWorkerSleep.WithLabelValues(m.name).Add(-1)
		}
	} else {
		promWorkerSleep.WithLabelValues(m.name).Add(-1)
	}
}

func (m PrometheusMetrics) QueuePut() {
	promQueueIn.WithLabelValues(m.name).Inc()
	promQueueSize.WithLabelValues(m.name).Inc()
}

func (m PrometheusMetrics) QueuePull() {
	promQueueOut.WithLabelValues(m.name).Inc()
	promQueueSize.WithLabelValues(m.name).Dec()
}

func (m PrometheusMetrics) QueueRetry() {
	promQueueRetry.WithLabelValues(m.name).Inc()
}

func (m PrometheusMetrics) QueueLeak() {
	promQueueLeak.WithLabelValues(m.name).Inc()
	promQueueSize.WithLabelValues(m.name).Dec()
}

func (m PrometheusMetrics) QueueLost() {
	promQueueLost.WithLabelValues(m.name).Inc()
	promQueueSize.WithLabelValues(m.name).Dec()
}