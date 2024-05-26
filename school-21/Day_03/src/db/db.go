package db

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

type Place struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Location struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	} `json:"location"`
}

type Store interface {
	GetPlaces(limit int, offset int) ([]Place, int, error)
	GetClosestPlaces(lat, lon float64, limit int) ([]Place, error)
}

type ElasticsearchStore struct {
	client *elasticsearch.Client
}

func NewElasticsearchStore(client *elasticsearch.Client) *ElasticsearchStore {
	return &ElasticsearchStore{client: client}
}

func (s *ElasticsearchStore) GetPlaces(limit int, offset int) ([]Place, int, error) {
	query := fmt.Sprintf(`{
		"from": %d,
		"size": %d,
		"query": {
			"match_all": {}
		}
	}`, offset, limit)

	res, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex("places"),
		s.client.Search.WithBody(strings.NewReader(query)),
		s.client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, 0, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return nil, 0, fmt.Errorf("error: %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, 0, fmt.Errorf("error parsing the response body: %s", err)
	}

	hits, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, 0, fmt.Errorf("unexpected response format")
	}
	totalHits := int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64))

	places := make([]Place, len(hits))
	for i, hit := range hits {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if !ok {
			return nil, 0, fmt.Errorf("unexpected source format")
		}
		location, ok := source["location"].(map[string]interface{})
		if !ok {
			return nil, 0, fmt.Errorf("unexpected location format")
		}
		places[i] = Place{
			ID:      source["id"].(string),
			Name:    source["name"].(string),
			Address: source["address"].(string),
			Phone:   source["phone"].(string),
			Location: struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			}{
				Lat: location["lat"].(float64),
				Lon: location["lon"].(float64),
			},
		}
	}

	return places, totalHits, nil
}

func (s *ElasticsearchStore) GetClosestPlaces(lat, lon float64, limit int) ([]Place, error) {
	query := fmt.Sprintf(`{
		"size": %d,
		"sort": [
			{
				"_geo_distance": {
					"location": {
						"lat": %f,
						"lon": %f
					},
					"order": "asc",
					"unit": "km",
					"mode": "min",
					"distance_type": "arc",
					"ignore_unmapped": true
				}
			}
		]
	}`, limit, lat, lon)

	res, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex("places"),
		s.client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("error parsing the response body: %s", err)
		} else {
			return nil, fmt.Errorf("error: %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("error parsing the response body: %s", err)
	}

	hits, ok := r["hits"].(map[string]interface{})["hits"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	places := make([]Place, len(hits))
	for i, hit := range hits {
		source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected source format")
		}
		location, ok := source["location"].(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unexpected location format")
		}
		places[i] = Place{
			ID:      source["id"].(string),
			Name:    source["name"].(string),
			Address: source["address"].(string),
			Phone:   source["phone"].(string),
			Location: struct {
				Lat float64 `json:"lat"`
				Lon float64 `json:"lon"`
			}{
				Lat: location["lat"].(float64),
				Lon: location["lon"].(float64),
			},
		}
	}

	return places, nil
}
