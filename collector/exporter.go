package collector

import (
	"context"
	"github.com/PixiuMetricsImplementation/global"
	"github.com/PixiuMetricsImplementation/scrape"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"sync"
	"time"
)

//将写好的一个个scrapeImpl接管过来
//记录监控目的组件是否可达
//记录下每个scrapeImpl的执行时间
const (
	// Subsystem(s).
	exporter = "exporter"
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(global.Namespace, exporter, "collector_duration_seconds"),
		"Collector time duration.",
		[]string{"collector"}, nil,
	)
)

type Metrics struct {
	ExporterUp prometheus.Gauge
}

//生成形如mysql_up,memcached_up,redis_up的metrics，
func NewMetrics() Metrics {
	return Metrics{
		ExporterUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: global.Namespace,
			Name:      "up",
			Help:      "Whether the datacenter is up.",
		}),
	}
}

//Exporter就是数据回收
type Exporter struct {
	ctx      context.Context
	scrapers []scrape.Scraper
	metrics  Metrics
	dsn      string
	//监控自己本身
	duration     prometheus.Gauge
	scrapeError  prometheus.Gauge
	totalScrapes prometheus.Counter
}

//New 给main函数调用的，无需考虑其他组件。
//ctx是http请求那里传过来的，需要使用ctx将本exporter抓取到的所有数据返回到response里，必带
//dsn是exporter刚启动时从配置文件或者启动参数拿来的数据，用这个获取链接,然后通过；建立连接，然后把具体的连接传递给一个个的Impl，每个数据采集器（scrapeImpl）采集周期都是15s
//dsn可以是结构体，看构造情况方便修改即可
//所有的scrapeImpl写好后先在main里声明到数组里，main函数调用collector的New时告诉collector有哪些scrapeImpl需要采集数据。至此完成了exporter的核心功能
func New(ctx context.Context, dsn string, metrics Metrics, scrapers []scrape.Scraper) *Exporter {
	return &Exporter{
		ctx:      ctx,
		dsn:      dsn,
		scrapers: scrapers,
		metrics:  metrics,
	}
}
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.metrics.ExporterUp.Desc()
}
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(e.ctx, ch)
	ch <- e.metrics.ExporterUp
}

//在这里开始启动所有传过来的scrapeImpl
func (e *Exporter) scrape(ctx context.Context, ch chan<- prometheus.Metric) {

	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, 0.01, "version")
	e.metrics.ExporterUp.Set(1)
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, scraper := range e.scrapers {
		wg.Add(1)
		go func(scraper scrape.Scraper) {
			defer wg.Done()
			label := "collect." + scraper.Name()
			scrapeTime := time.Now()
			if err := scraper.Scrape(ctx, "dc", ch); err != nil {
				log.Println(err)
			}
			ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, time.Since(scrapeTime).Seconds(), label)
		}(scraper)
	}
}
