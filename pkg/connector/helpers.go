package connector

import (
	"fmt"
	"time"

	v2 "github.com/conductorone/baton-sdk/pb/c1/connector/v2"
	"github.com/conductorone/baton-sdk/pkg/pagination"
)

const (
	resourcePageSize = 50
)

func parsePageToken(i string, resourceID *v2.ResourceId) (*pagination.Bag, string, error) {
	b := &pagination.Bag{}
	err := b.Unmarshal(i)
	if err != nil {
		return nil, "", err
	}

	if b.Current() == nil {
		b.Push(pagination.PageState{
			ResourceTypeID: resourceID.ResourceType,
			ResourceID:     resourceID.Resource,
		})
	}

	return b, b.PageToken(), nil
}

func handleNextPage(bag *pagination.Bag, nextCursor string) (string, error) {
	pageToken, err := bag.NextToken(nextCursor)
	if err != nil {
		return "", err
	}

	return pageToken, nil
}

func wrapError(err error, message string) error {
	return fmt.Errorf("miro-connector: %s: %w", message, err)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func parseTime(i string) (*time.Time, error) {
	if i == "" {
		return nil, nil
	}

	t, err := time.Parse(time.RFC3339, i)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
