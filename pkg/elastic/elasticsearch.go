package elastic

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esutil"
	l "github.com/krish512/ambassador_logservice_es/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

var logger = l.InitLogger()

//ES is instance of Elasticsearch Client connected to the server
var ES *elasticsearch.Client
var err error

//InitElasticsearch initiates connection to elasticsearch
func InitElasticsearch() {

	ElasticsearchEndpoints := []string{"http://localhost:9200"}

	if os.Getenv("ELASTICSEARCH_ENDPOINTS") != "" {
		endpoints := strings.Split(strings.TrimSpace(os.Getenv("ELASTICSEARCH_ENDPOINTS")), ",")
		ElasticsearchEndpoints = nil
		for _, endpoint := range endpoints {
			ElasticsearchEndpoints = append(ElasticsearchEndpoints, strings.TrimSpace(endpoint))
		}
	}
	logger.Debug("Used Elasticsearch Endpoints:", zap.Strings("endpoints", ElasticsearchEndpoints))

	// Test connection to elasticsearch server
	es, err := testConnection(ElasticsearchEndpoints)
	if err != nil {
		logger.Error("Failure encountered", zap.Error(err))
		os.Exit(1)
	}
	ES = es
}

// InitESBulkIndexer creates a BulkIndexer instance for stream to flush
func InitESBulkIndexer() esutil.BulkIndexer {

	ElasticsearchIndex := "ambassador"

	if os.Getenv("ELASTICSEARCH_INDEX") != "" {
		ElasticsearchIndex = os.Getenv("ELASTICSEARCH_INDEX")
	}

	ElasticsearchIndex = ElasticsearchIndex + "-" + time.Now().Format("2006.01.02")

	cfg := esutil.BulkIndexerConfig{
		NumWorkers:   2,
		Client:       ES,
		OnError:      onError,
		OnFlushStart: onFlushStart,
		OnFlushEnd:   onFlushEnd,
		Index:        ElasticsearchIndex,
	}

	es, err := esutil.NewBulkIndexer(cfg)
	if err != nil {
		logger.Error("Error creating Elasticsearch BulkIndexer:", zap.Error(err))
	}
	return es
}

func testConnection(endpoints []string) (*elasticsearch.Client, error) {

	var r map[string]interface{}

	cfg := elasticsearch.Config{
		Addresses:             endpoints,
		DiscoverNodesOnStart:  true,
		DiscoverNodesInterval: 300 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   2,
			ResponseHeaderTimeout: 30 * time.Second,
		},
	}

	if os.Getenv("ELASTICSEARCH_USERNAME") != "" {
		cfg.Username = os.Getenv("ELASTICSEARCH_USERNAME")
	}

	if os.Getenv("ELASTICSEARCH_PASSWORD") != "" {
		cfg.Password = os.Getenv("ELASTICSEARCH_PASSWORD")
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		logger.Error("Error creating Elasticsearch client:", zap.Error(err))
		return es, err
	}

	res, err := es.Info()
	if err != nil {
		logger.Error("Error getting response from Elasticsearch:", zap.Error(err))
		return es, err
	}

	defer res.Body.Close()
	// Check response status
	if res.IsError() {
		logger.Error("Response Error:", zap.String("res", res.String()))
		return es, errors.New(res.String())
	}
	// Deserialize the response into a map.
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		logger.Error("Error parsing the response body: %s", zap.Error(err))
		return es, err
	}

	// Print client and server version numbers.
	logger.Debug("Elasticsearch Client:", zap.String("elasticsearch.ClientVersion", elasticsearch.Version))
	logger.Debug("Elasticsearch Server:", zap.Any("elasticsearch.ServerVersion", r["version"].(map[string]interface{})["number"]))
	logger.Info("Connected to Elasticsearch successfully!")
	return es, nil
}

func onError(ctx context.Context, err error) {
	logger.Error("Error encountered while indexing", zap.Error(err))
}

func onFlushStart(ctx context.Context) context.Context {
	logger.Debug("Flush started")
	return ctx
}

func onFlushEnd(ctx context.Context) {
	logger.Debug("Flush ended")
}
