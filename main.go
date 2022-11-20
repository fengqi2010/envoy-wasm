package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"regexp"
	"strconv"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// NewPluginContext Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &metricPluginContext{}
}

type metricPluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

// NewHttpContext Override types.DefaultPluginContext.
func (ctx *metricPluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	counters := map[string]proxywasm.MetricCounter{}
	gauges := map[string]proxywasm.MetricGauge{}
	histograms := map[string]proxywasm.MetricHistogram{}
	return &metricHttpContext{
		counters:   counters,
		gauges:     gauges,
		histograms: histograms,
	}
}

type metricHttpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	// counters
	counters map[string]proxywasm.MetricCounter
	// gauges
	gauges map[string]proxywasm.MetricGauge
	// histograms
	histograms map[string]proxywasm.MetricHistogram
}

// # HELP datacanvas_server_request_seconds_max
// # TYPE datacanvas_server_request_seconds_max gauge
// datacanvas_server_request_seconds_max{KIND="KIND_SHADOW",method="PREDICTIONS",modelName="1111",status="ok",} NaN
// datacanvas_server_request_seconds_max{KIND="KIND_M",method="PREDICTIONS",modelName="1111",status="ok",} NaN
// datacanvas_server_request_seconds_max{KIND="KIND_S",method="PREDICTIONS",modelName="1111",status="ok",} NaN
// # HELP datacanvas_server_request_seconds_min
// # TYPE datacanvas_server_request_seconds_min gauge
// datacanvas_server_request_seconds_min{KIND="KIND_SHADOW",method="PREDICTIONS",modelName="1111",status="ok",} NaN
// datacanvas_server_request_seconds_min{KIND="KIND_M",method="PREDICTIONS",modelName="1111",status="ok",} NaN
// datacanvas_server_request_seconds_min{KIND="KIND_S",method="PREDICTIONS",modelName="1111",status="ok",} NaN
// # HELP datacanvas_server_request_seconds_sum
// # TYPE datacanvas_server_request_seconds_sum gauge
// datacanvas_server_request_seconds_sum{KIND="KIND_SHADOW",method="PREDICTIONS",modelName="1111",status="ok",} 0.0
// datacanvas_server_request_seconds_sum{KIND="KIND_M",method="PREDICTIONS",modelName="1111",status="ok",} 0.0
// datacanvas_server_request_seconds_sum{KIND="KIND_S",method="PREDICTIONS",modelName="1111",status="ok",} 0.0
// # HELP datacanvas_server_request_accept_total
// # TYPE datacanvas_server_request_accept_total counter
// datacanvas_server_request_accept_total{method="PREDICTIONS",modelName="1111",} 1.0
// # HELP datacanvas_server_request_total
// # TYPE datacanvas_server_request_total counter
// datacanvas_server_request_total{KIND="KIND_SHADOW",method="PREDICTIONS",modelName="1111",status="ok",} 0.0
// datacanvas_server_request_total{KIND="KIND_S",method="PREDICTIONS",modelName="1111",status="err",} 1.0
// datacanvas_server_request_total{KIND="KIND_M",method="PREDICTIONS",modelName="1111",status="ok",} 0.0
// datacanvas_server_request_total{KIND="KIND_S",method="PREDICTIONS",modelName="1111",status="ok",} 0.0

const (
	requestAcceptTotal = "datacanvas.server.request.accept.total"
	requestTotal       = "datacanvas.server.request.total"
	requestMax         = "datacanvas.server.request.seconds.max"
	requestMin         = "datacanvas.server.request.seconds.min"
	requestSum         = "datacanvas.server.request.seconds.sum"
	requestSeconds     = "datacanvas.server.request.seconds"
)

// OnHttpRequestHeaders Override types.DefaultHttpContext.
func (ctx *metricHttpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	if !endOfStream {
		return types.ActionContinue
	}
	// debug 打印header日志
	proxywasm.LogInfo("OnHttpRequestHeaders---------------------------------------------------")
	headers, _ := proxywasm.GetHttpRequestHeaders()
	for _, header := range headers {
		proxywasm.LogInfo("name = " + header[0] + ", value = " + header[1])
	}
	proxywasm.LogInfo("OnHttpRequestHeaders---------------------------------------------------")
	//// 添加当前时间
	//_ = proxywasm.AddHttpRequestHeader("APS-CURRENT-TIME", strconv.FormatInt(time.Now().UnixMilli(), 10))
	//if apiTest, _ := proxywasm.GetHttpRequestHeader("X-API-TEST"); apiTest != "" {
	//	return types.ActionContinue
	//}
	//if kind, _ := proxywasm.GetHttpRequestHeader("X-KIND-M"); kind != "" {
	//	return types.ActionContinue
	//}
	//// 模型id
	//modelId, _ := proxywasm.GetHttpRequestHeader("MODEL-ID")
	//if modelId == "" {
	//	return types.ActionContinue
	//}
	//// 模型name
	//modelName, _ := proxywasm.GetHttpRequestHeader("MODEL-NAME")
	//if modelName == "" {
	//	return types.ActionContinue
	//}
	//if host, _ := proxywasm.GetHttpRequestHeader("host"); !strings.HasSuffix(host, "shadow") {
	//	fqn := fmt.Sprintf("&modelName=%s&method=%s&", modelName, "PREDICTIONS")
	//	ctx.counter(requestAcceptTotal, fqn, 1)
	//}
	return types.ActionContinue
}

var shadowRegx = regexp.MustCompile("^.*-shadow:\\d+$|^.*-shadow$")

func (ctx *metricHttpContext) dispatchCallback(numHeaders, bodySize, numTrailers int) {
	proxywasm.LogInfo("dispatched httpbin")
}

// OnHttpResponseHeaders Override types.DefaultHttpContext.
func (ctx *metricHttpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("OnHttpResponseHeaders---------------------------------------------------")
	headers, _ := proxywasm.GetHttpResponseHeaders()
	for _, header := range headers {
		proxywasm.LogInfo("name = " + header[0] + ", value = " + header[1])
	}
	proxywasm.LogInfo("OnHttpResponseHeaders---------------------------------------------------")
	//modelId, _ := proxywasm.GetHttpResponseHeader("MODEL-ID")
	//if modelId == "" {
	//	return types.ActionContinue
	//}
	//proxywasm.LogInfo(modelId)
	//modelName, _ := proxywasm.GetHttpResponseHeader("MODEL-NAME")
	//if modelName == "" {
	//	return types.ActionContinue
	//}
	//proxywasm.LogInfo(modelName)
	//if _, err := proxywasm.DispatchHttpCall("httpbin", [][2]string{
	//	{":path", "/get"},
	//	{":method", "GET"},
	//	{"MODEL-ID", modelId},
	//	{"MODEL-NAME", modelName},
	//	{":authority", ""}},
	//	nil, nil, 50000, ctx.dispatchCallback); err != nil {
	//	proxywasm.LogInfo(err.Error())
	//}
	//proxywasm.LogInfo("DispatchHttpCall---------------------------------------------------")
	// Response的Headers中没有模型id和name直接返回
	//modelId, _ := proxywasm.GetHttpResponseHeader("MODEL-ID")
	//if modelId == "" {
	//	return types.ActionContinue
	//}
	//_ = proxywasm.RemoveHttpResponseHeader("MODEL-ID")
	//modelName, _ := proxywasm.GetHttpResponseHeader("MODEL-NAME")
	//if modelName == "" {
	//	return types.ActionContinue
	//}
	//_ = proxywasm.RemoveHttpResponseHeader("MODEL-NAME")
	//// 获取kind类型
	//var kind string
	//if headerValue, _ := proxywasm.GetHttpResponseHeader("X-KIND-M"); headerValue == "" {
	//	if host, _ := proxywasm.GetHttpResponseHeader("host"); !shadowRegx.MatchString(host) {
	//		kind = "KIND_S"
	//	} else {
	//		kind = "KIND_SHADOW"
	//	}
	//} else {
	//	kind = "KIND_M"
	//}
	//// Response的Headers中没有APS-CODE直接返回
	//apsCode, _ := proxywasm.GetHttpResponseHeader("APS-CODE")
	//if apsCode == "" {
	//	return types.ActionContinue
	//}
	//// 获取请求开始时间
	//currentTime, err := getCurrentTime()
	//if err != nil {
	//	return types.ActionContinue
	//}
	//_ = proxywasm.RemoveHttpResponseHeader("APS-CURRENT-TIME")
	//proxywasm.LogInfo("cost time: " + strconv.FormatInt(time.Now().UnixMilli()-currentTime, 10))
	//if apsCode == "0" {
	//	fqn := fmt.Sprintf("&modelName=%s&method=%s&status=%s&KIND=%s&", modelName, "PREDICTIONS", "ok", kind)
	//	ctx.counter(requestTotal, fqn, 1)
	//	ctx.histogram(requestSeconds, fqn, uint64(time.Now().UnixMilli()-currentTime))
	//	ctx.gauge(requestMin, fqn, time.Now().UnixMilli()-currentTime)
	//	ctx.gauge(requestMax, fqn, time.Now().UnixMilli()-currentTime)
	//	ctx.gauge(requestSum, fqn, time.Now().UnixMilli()-currentTime)
	//} else {
	//	fqn := fmt.Sprintf("&modelName=%s&method=%s&status=%s&KIND=%s&", modelName, "PREDICTIONS", "err", kind)
	//	ctx.counter(requestTotal, fqn, 1)
	//}
	return types.ActionContinue
}

func (ctx *metricHttpContext) histogram(id string, fqn string, value uint64) {
	histogram, ok := ctx.histograms[id]
	if !ok {
		histogram = proxywasm.DefineHistogramMetric(id + fqn)
		ctx.histograms[id+fqn] = histogram
	}
	proxywasm.LogInfo(id + fqn + " value = " + string(value))
	histogram.Record(value)
}

func (ctx *metricHttpContext) gauge(id string, fqn string, offset int64) {
	gauge, ok := ctx.gauges[id]
	if !ok {
		gauge = proxywasm.DefineGaugeMetric(id + fqn)
		ctx.gauges[id+fqn] = gauge
	}
	proxywasm.LogInfo(id + fqn + " offset = " + string(offset))
	gauge.Add(offset)
}

func (ctx *metricHttpContext) counter(id string, fqn string, offset uint64) {
	counter, ok := ctx.counters[id]
	if !ok {
		counter = proxywasm.DefineCounterMetric(id + fqn)
		ctx.counters[id+fqn] = counter
	}
	proxywasm.LogInfo(id + fqn + " offset")
	counter.Increment(offset)
}

func getCurrentTime() (int64, error) {
	if currentTimeStr, err := proxywasm.GetHttpResponseHeader("APS-CURRENT-TIME"); currentTimeStr != "" {
		currentTime, err := strconv.ParseInt(currentTimeStr, 10, 64)
		if err != nil {
			return 0, err
		}
		return currentTime, nil
	} else {
		return 0, err
	}
}
