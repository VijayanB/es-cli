package knn

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"odfe-cli/client"
	"odfe-cli/entity"
	gw "odfe-cli/gateway"
)

const (
	baseURL              = "_opendistro/_knn"
	statsURL             = baseURL + "/stats"
	nodeStatsURLTemplate = baseURL + "/%s/stats/%s"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen  -destination=mocks/mock_knn.go -package=mocks . Gateway

// Gateway interface to k-NN Plugin
type Gateway interface {
	GetStatistics(ctx context.Context, nodes string, names string) ([]byte, error)
}

type gateway struct {
	gw.HTTPGateway
}

// New creates new Gateway instance
func New(c *client.Client, p *entity.Profile) Gateway {
	return &gateway{*gw.NewHTTPGateway(c, p)}
}

//buildStatsURL to construct url for stats
func (g *gateway) buildStatsURL(nodes string, names string) (*url.URL, error) {
	endpoint, err := gw.GetValidEndpoint(g.Profile)
	if err != nil {
		return nil, err
	}
	path := statsURL
	if nodes != "" {
		path = fmt.Sprintf(nodeStatsURLTemplate, nodes, names)
	}
	endpoint.Path = path
	return endpoint, nil
}

/*GetStatistics provides information about the current status of the KNN Plugin.
GET /_opendistro/_knn/stats
{
    "_nodes" : {
        "total" : 1,
        "successful" : 1,
        "failed" : 0
    },
    "cluster_name" : "_run",
    "circuit_breaker_triggered" : false,
    "nodes" : {
        "HYMrXXsBSamUkcAjhjeN0w" : {
            "eviction_count" : 0,
            "miss_count" : 1,
            "graph_memory_usage" : 1,
            "graph_memory_usage_percentage" : 3.68,
            "graph_index_requests" : 7,
            "graph_index_errors" : 1,
            "knn_query_requests" : 4,
            "graph_query_requests" : 30,
            "graph_query_errors" : 15,
            "indices_in_cache" : {
                "myindex" : {
                    "graph_memory_usage" : 2,
                    "graph_memory_usage_percentage" : 3.68,
                    "graph_count" : 2
                }
            },
            "cache_capacity_reached" : false,
            "load_exception_count" : 0,
            "hit_count" : 0,
            "load_success_count" : 1,
            "total_load_time" : 2878745,
            "script_compilations" : 1,
            "script_compilation_errors" : 0,
            "script_query_requests" : 534,
            "script_query_errors" : 0
        }
    }
}
To filter stats query by nodeID and statName:
GET /_opendistro/_knn/nodeId1,nodeId2/stats/statName1,statName2
*/
func (g gateway) GetStatistics(ctx context.Context, nodes string, names string) ([]byte, error) {
	statsURL, err := g.buildStatsURL(nodes, names)
	if err != nil {
		return nil, err
	}
	detectorRequest, err := g.BuildRequest(ctx, http.MethodGet, "", statsURL.String(), gw.GetHeaders())
	if err != nil {
		return nil, err
	}
	response, err := g.Call(detectorRequest, http.StatusOK)
	if err != nil {
		return nil, err
	}
	return response, nil
}