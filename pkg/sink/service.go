package sink

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/elastic/go-elasticsearch/v7/esutil"
	v2 "github.com/envoyproxy/go-control-plane/envoy/service/accesslog/v2"
	"github.com/golang/protobuf/jsonpb"
	"github.com/krish512/ambassador_logservice_es/pkg/elastic"
	l "github.com/krish512/ambassador_logservice_es/pkg/logger"
	"go.uber.org/zap"
)

var logger = l.InitLogger()

type server struct {
	marshaler jsonpb.Marshaler
}

var _ v2.AccessLogServiceServer = &server{}

// New AccessLogServiceServer
func New() v2.AccessLogServiceServer {
	return &server{}
}

func (s *server) StreamAccessLogs(stream v2.AccessLogService_StreamAccessLogsServer) error {
	logger.Info("Stream Received")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		logEntries := transform(in.GetHttpLogs())
		es := elastic.InitESBulkIndexer()
		for _, logEntry := range logEntries {
			var data []byte
			data, err = json.Marshal(logEntry)
			es.Add(context.Background(),
				esutil.BulkIndexerItem{
					Action:     "create",
					DocumentID: logEntry.TraceID,
					Body:       bytes.NewReader(data),
					OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
						logger.Debug("Success", zap.String("result", res.Result))
					},
					OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
						if err != nil {
							logger.Error("ERROR:", zap.Error(err))
						} else {
							logger.Error("ERROR:", zap.String("errorType", res.Error.Type), zap.String("errorReason", res.Error.Reason))
						}
					},
				})
			data = nil
		}
		es.Close(context.Background())
	}
}
