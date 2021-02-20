package sink

import (
	"time"

	v2 "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v2"
)

type requestFields struct {
	HeadersBytes uint64 `json:"headersBytes"`
	BodyBytes    uint64 `json:"bodyBytes"`
	RemoteIP     string `json:"remoteIP"`
}

type responseFields struct {
	Details                     string        `json:"details"`
	HeadersBytes                uint64        `json:"headersBytes"`
	BodyBytes                   uint64        `json:"bodyBytes"`
	TimeToLastRxByte            time.Duration `json:"timeToLastRxByte"`
	TimeToFirstUpstreamTxByte   time.Duration `json:"timeToFirstUpstreamTxByte"`
	TimeToLastUpstreamTxByte    time.Duration `json:"timeToLastUpstreamTxByte"`
	TimeToFirstUpstreamRxByte   time.Duration `json:"timeToFirstUpstreamRxByte"`
	TimeToLastUpstreamRxByte    time.Duration `json:"timeToLastUpstreamRxByte"`
	TimeToFirstDownstreamTxByte time.Duration `json:"timeToFirstDownstreamTxByte"`
	TimeToLastDownstreamTxByte  time.Duration `json:"timeToLastDownstreamTxByte"`
}

// AmbassadorLog is struct for JSON format to push ambassador access logs to Elasticsearch
type AmbassadorLog struct {
	Timestamp       string         `json:"timestamp"`
	TraceID         string         `json:"traceId"`
	Method          string         `json:"method"`
	Protocol        string         `json:"protocol"`
	Scheme          string         `json:"scheme"`
	StatusCode      uint32         `json:"statusCode"`
	Domain          string         `json:"domain"`
	Path            string         `json:"path"`
	UserAgent       string         `json:"userAgent"`
	Referer         string         `json:"referer"`
	ForwardedFor    string         `json:"forwardedFor"`
	UpstreamFailure string         `json:"upsteamFailure"`
	Time            time.Duration  `json:"time_ms"`
	Request         requestFields  `json:"request"`
	Response        responseFields `json:"response"`
}

func transform(bundle *v2.StreamAccessLogsMessage_HTTPAccessLogEntries) []AmbassadorLog {
	var ambassadorLog []AmbassadorLog
	for _, logEntry := range bundle.LogEntry {
		ambassadorLog = append(ambassadorLog, AmbassadorLog{
			Timestamp:       logEntry.CommonProperties.StartTime.AsTime().Format("2006-01-02T15:04:05-0700"),
			TraceID:         logEntry.Request.GetRequestId(),
			Method:          logEntry.Request.GetRequestMethod().String(),
			Protocol:        logEntry.GetProtocolVersion().String(),
			Scheme:          logEntry.Request.GetScheme(),
			StatusCode:      logEntry.Response.GetResponseCode().GetValue(),
			Domain:          logEntry.Request.GetAuthority(),
			Path:            logEntry.Request.GetPath(),
			UserAgent:       logEntry.Request.GetUserAgent(),
			Referer:         logEntry.Request.GetReferer(),
			ForwardedFor:    logEntry.Request.GetForwardedFor(),
			UpstreamFailure: logEntry.CommonProperties.GetUpstreamTransportFailureReason(),
			Time:            logEntry.CommonProperties.GetTimeToLastDownstreamTxByte().AsDuration() / time.Millisecond,
			Request: requestFields{
				HeadersBytes: logEntry.Request.GetRequestHeadersBytes(),
				BodyBytes:    logEntry.Request.GetRequestBodyBytes(),
				RemoteIP:     logEntry.CommonProperties.GetDownstreamRemoteAddress().GetSocketAddress().Address,
			},
			Response: responseFields{
				Details:                     logEntry.Response.GetResponseCodeDetails(),
				HeadersBytes:                logEntry.Response.GetResponseHeadersBytes(),
				BodyBytes:                   logEntry.Response.GetResponseBodyBytes(),
				TimeToLastRxByte:            logEntry.CommonProperties.GetTimeToLastRxByte().AsDuration() / time.Millisecond,
				TimeToFirstUpstreamTxByte:   logEntry.CommonProperties.GetTimeToFirstUpstreamTxByte().AsDuration() / time.Millisecond,
				TimeToLastUpstreamTxByte:    logEntry.CommonProperties.GetTimeToLastUpstreamTxByte().AsDuration() / time.Millisecond,
				TimeToFirstUpstreamRxByte:   logEntry.CommonProperties.GetTimeToFirstUpstreamRxByte().AsDuration() / time.Millisecond,
				TimeToLastUpstreamRxByte:    logEntry.CommonProperties.GetTimeToLastUpstreamRxByte().AsDuration() / time.Millisecond,
				TimeToFirstDownstreamTxByte: logEntry.CommonProperties.GetTimeToFirstDownstreamTxByte().AsDuration() / time.Millisecond,
				TimeToLastDownstreamTxByte:  logEntry.CommonProperties.GetTimeToLastDownstreamTxByte().AsDuration() / time.Millisecond,
			},
		})
	}
	return ambassadorLog
}
