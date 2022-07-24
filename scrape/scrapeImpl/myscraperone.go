package scrapeImpl

import (
	"context"
	"github.com/PixiuMetricsImplementation/global"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Subsystem.
	myscraperonesubs = "myscraperonesubs"
)

var (
	myscraperoneDesc = prometheus.NewDesc(
		//myscraperonesubs是我当前scrapeImpl的名字
		//"args"具体一个指标的名称
		//"my scraper one test args ." 这个是HELP说明
		// 描述清楚这个指标是干嘛的即可
		// 在字符数组里只写了args，那么最终的效果是args="argsone1"
		prometheus.BuildFQName(global.Namespace, myscraperonesubs, "args"),
		"my scraper one test args .",
		[]string{"args"}, nil,
	)
)

type MyScraperOne struct{}

// Name of the Scraper. Should be unique.
func (MyScraperOne) Name() string {
	return myscraperonesubs
}

// Help describes the role of the Scraper.
func (MyScraperOne) Help() string {
	return "my scraper one"
}

//dc string这块，根据不同的监控目的组件这里需要换成不同的链接驱动
//如果是mysqld_exporter的话需要换成Scrape(ctx context.Context, db *sql.DB, ch chan<- prometheus.Metric) error
//如果是redis_exporter的话换成 Scrape(ctx context.Context, c *redis.Conn, ch chan<- prometheus.Metric) error
//Impl里切换成相对应的驱动就好
func (MyScraperOne) Scrape(ctx context.Context, dc string, ch chan<- prometheus.Metric) error {
	//get some from datacentor从数据中心获取数据
	//dc.dosomthing...
	//may be return error，will stop this scrape's register可能是返回错误，会停止这个scrape的注册
	//return error
	//因为我在myscraperoneDesc里写的是[]string{"args"}，只有一个label，所以我在注入数据的时候只写"argsone1"一个就行
	ch <- prometheus.MustNewConstMetric(
		myscraperoneDesc, prometheus.CounterValue, 0.01, "argsone1",
	)
	ch <- prometheus.MustNewConstMetric(
		global.NewDesc(myscraperonesubs, "subsystemnameone", "Generic metric"),
		prometheus.UntypedValue,
		0.01,
	)
	return nil
}
