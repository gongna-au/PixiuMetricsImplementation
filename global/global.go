package global

//主要完成一些公共使用类代码的封装以一些全局关键字的定义。
import "github.com/prometheus/client_golang/prometheus"

const (
	// Exporter Namespace.
	// Namespace将会经常被scrape和collector甚至main调用
	Namespace = "pixiu"
)

func NewDesc(subsystem, name, help string) *prometheus.Desc {
	return prometheus.NewDesc(
		prometheus.BuildFQName(Namespace, subsystem, name),
		help, nil, nil,
	)
}
