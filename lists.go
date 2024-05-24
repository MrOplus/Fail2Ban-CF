package main

import (
	"context"
	"errors"
	"github.com/cloudflare/cloudflare-go"
	"github.com/seculize/islazy/log"
)

func createList(ctx context.Context, api *cloudflare.API) (*cloudflare.List, error) {
	l, err := api.CreateList(ctx, &cloudflare.ResourceContainer{
		Identifier: config.AccountId,
	}, cloudflare.ListCreateParams{
		Name:        "blacklist",
		Kind:        cloudflare.IPListTypeIP,
		Description: "List of IP addresses to block",
	})
	if err != nil {
		return nil, err
	}
	return &l, nil
}
func getList(ctx context.Context, api *cloudflare.API) (*cloudflare.List, error) {
	lists, err := api.ListLists(ctx, &cloudflare.ResourceContainer{
		Identifier: config.AccountId,
	}, cloudflare.ListListsParams{})
	if err != nil {
		return nil, err
	}
	for _, list := range lists {
		if list.Name == "blacklist" {
			return &list, nil
		}
	}
	return nil, errors.New("list not found")
}

func addIP(ctx context.Context, api *cloudflare.API, list *cloudflare.List, ip, comment string) (*cloudflare.ListItemCreateResponse, error) {
	item, err := api.CreateListItemAsync(ctx, &cloudflare.ResourceContainer{
		Identifier: config.AccountId,
	}, cloudflare.ListCreateItemParams{
		ID: list.ID,
		Item: cloudflare.ListItemCreateRequest{
			IP:      &ip,
			Comment: comment,
		},
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func findIp(ctx context.Context, api *cloudflare.API, list *cloudflare.List, ip string) (*cloudflare.ListItem, error) {
	items, err := api.ListListItems(ctx, &cloudflare.ResourceContainer{
		Identifier: config.AccountId,
	}, cloudflare.ListListItemsParams{
		ID:     list.ID,
		Search: ip,
	})
	if err != nil {
		log.Fatal("Error listing items: ", err)
	}
	if len(items) == 0 {
		return nil, errors.New("item not found")
	}
	return &items[0], nil
}

func deleteIp(ctx context.Context, api *cloudflare.API, list *cloudflare.List, item *cloudflare.ListItem) error {
	resp, err := api.DeleteListItemsAsync(ctx, &cloudflare.ResourceContainer{
		Identifier: config.AccountId,
	}, cloudflare.ListDeleteItemsParams{
		ID: list.ID,
		Items: cloudflare.ListItemDeleteRequest{
			Items: []cloudflare.ListItemDeleteItemRequest{
				{
					ID: item.ID,
				},
			},
		},
	})
	if err != nil {
		return err
	}
	if resp.Success {
		return nil
	}
	return errors.New("error deleting item")
}
