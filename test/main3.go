package main

import (
	//"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	cpuTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_temperature_celsius",
		Help: "Current temperature of the CPU.",
	})
	hdFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "hd_errors_total",
			Help: "Number of hard-disk errors.",
		},
		[]string{"device"},
	)
)

func main() {
	// 创建一个自定义的注册表
	registry := prometheus.NewRegistry()

	// 可选: 添加 process 和 Go 运行时指标到我们自定义的注册表中
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())

	queueLength := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "queue_length",
		Help: "The number of items in the queue.",
	})

	totalRequests := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "The total number of handled HTTP requests.",
	})

	requestSummaryDurations := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "http_request_duration_seconds",
		Help: "A summary of the HTTP request durations in seconds.",
		Objectives: map[float64]float64{
			0.5:  0.05,  // 第50个百分位数，最大绝对误差为0.05。
			0.9:  0.01,  // 第90个百分位数，最大绝对误差为0.01。
			0.99: 0.001, // 第90个百分位数，最大绝对误差为0.001。
		},
	},
	)
	temp2 := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "home_temperature_celsius",
			Help: "The current temperature in degrees Celsius.",
		},
		// 指定标签名称
		[]string{"house", "room"},
	)

	registry.MustRegister(queueLength)
	registry.MustRegister(totalRequests)
	//registry.MustRegister(requestDurations)
	registry.MustRegister(requestSummaryDurations)
	registry.MustRegister(temp2)

	// 设置 gague 的值为 39

	//设置 queueLength 的值为 39
	queueLength.Inc()   // +1：Increment the gauge by 1.
	queueLength.Dec()   // -1：Decrement the gauge by 1.
	queueLength.Add(23) // Increment by 23.
	queueLength.Sub(42) // Decrement by 42.

	//设置 totalRequests 的值为 39
	totalRequests.Inc()

	//设置 requestDurations 的值
	//requestDurations.Observe(0.42)

	//设置 requestSummaryDurations 的值
	requestSummaryDurations.Observe(0.42)

	temp2.WithLabelValues("xiaoming", "ben-roorm").Set(27)

	// 暴露自定义指标
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	http.ListenAndServe(":8080", nil)
}
