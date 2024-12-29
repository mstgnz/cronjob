package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// CronjobExecutionTotal toplam cronjob çalıştırma sayısı
	CronjobExecutionTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cronjob_execution_total",
			Help: "Total number of cronjob executions",
		},
		[]string{"job_id", "status"},
	)

	// CronjobExecutionDuration cronjob çalışma süresi
	CronjobExecutionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cronjob_execution_duration_seconds",
			Help:    "Duration of cronjob execution in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"job_id"},
	)

	// CronjobLastExecutionTime son çalışma zamanı
	CronjobLastExecutionTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cronjob_last_execution_timestamp",
			Help: "Timestamp of last cronjob execution",
		},
		[]string{"job_id"},
	)

	// CronjobErrors hata sayısı
	CronjobErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cronjob_errors_total",
			Help: "Total number of cronjob errors",
		},
		[]string{"job_id", "error_type"},
	)

	// ActiveCronjobs aktif cronjob sayısı
	ActiveCronjobs = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cronjob_active_jobs",
			Help: "Number of currently active cronjobs",
		},
	)

	// SchedulerLeaderInfo scheduler leader bilgisi
	SchedulerLeaderInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cronjob_scheduler_leader_info",
			Help: "Information about scheduler leadership",
		},
		[]string{"instance_id", "job_id"},
	)

	// Yeni detaylı metrikler
	CronjobMemoryUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cronjob_memory_usage_bytes",
			Help: "Memory usage of cronjob in bytes",
		},
		[]string{"job_id"},
	)

	CronjobCPUUsage = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cronjob_cpu_usage_seconds",
			Help: "CPU usage of cronjob in seconds",
		},
		[]string{"job_id"},
	)

	CronjobQueueSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cronjob_queue_size",
			Help: "Number of jobs waiting in the queue",
		},
	)

	CronjobRetryCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cronjob_retry_total",
			Help: "Total number of job retries",
		},
		[]string{"job_id"},
	)

	CronjobHTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cronjob_http_request_duration_seconds",
			Help:    "Duration of HTTP requests made by cronjobs",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"job_id", "method", "status_code"},
	)

	CronjobDatabaseOperations = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cronjob_database_operations_total",
			Help: "Total number of database operations",
		},
		[]string{"job_id", "operation", "status"},
	)

	CronjobNotificationsSent = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cronjob_notifications_total",
			Help: "Total number of notifications sent",
		},
		[]string{"job_id", "type", "status"},
	)

	CronjobResourceSaturation = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cronjob_resource_saturation",
			Help: "Resource saturation levels (0-1)",
		},
		[]string{"resource_type"},
	)

	CronjobLastSuccess = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cronjob_last_success_timestamp",
			Help: "Timestamp of last successful execution",
		},
		[]string{"job_id"},
	)

	CronjobScheduleLag = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cronjob_schedule_lag_seconds",
			Help: "Difference between scheduled and actual execution time",
		},
		[]string{"job_id"},
	)
)
