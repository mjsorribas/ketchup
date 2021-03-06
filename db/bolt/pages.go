package bolt

import (
	"github.com/boltdb/bolt"
	"github.com/golang/protobuf/proto"
	"github.com/satori/go.uuid"

	"github.com/ketchuphq/ketchup/proto/ketchup/api"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

const PAGE_BUCKET = "pages"

// GetPage returns a page from the database
func (m *Module) GetPage(uuid string) (*models.Page, error) {
	page := &models.Page{}
	err := m.Get(PAGE_BUCKET, uuid, page)
	if err != nil {
		return nil, err
	}
	return page, nil
}

// DeletePage deletes the page and (note!) also deletes
// related routes.
func (m *Module) DeletePage(page *models.Page) error {
	err := m.delete(PAGE_BUCKET, page)
	if err != nil {
		return err
	}
	if page.GetUuid() == "" {
		return nil
	}

	routes, err := m.ListRoutes(&api.ListRouteRequest{
		Options: &api.ListRouteRequest_ListRouteOptions{
			PageUuid: page.Uuid,
		},
	})
	if err != nil {
		return err
	}
	for _, route := range routes {
		err := m.DeleteRoute(route)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdatePage updates (creating if necessary) a new page.
// New pages and new content blocks will be assigned UUIDs
func (m *Module) UpdatePage(page *models.Page) error {
	if page.GetUuid() == "" {
		page.Uuid = proto.String(uuid.NewV4().String())
	}

	for _, c := range page.Contents {
		if c.GetUuid() == "" {
			c.Uuid = proto.String(uuid.NewV4().String())
		}
	}
	return m.Update(PAGE_BUCKET, page)
}

// ListPages lists all the pages stored in the DB (unsorted), without contents
func (m *Module) ListPages(opts *api.ListPageRequest) ([]*models.Page, error) {
	lst := []*models.Page{}
	filter := opts.GetOptions().GetFilter()

	err := m.Bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(PAGE_BUCKET))
		return b.ForEach(func(_, v []byte) error {
			pb := &models.Page{}
			err := proto.Unmarshal(v, pb)
			if err != nil {
				return err
			}
			if filter == api.ListPageRequest_draft && pb.GetPublishedAt() != 0 {
				return nil
			}
			if filter == api.ListPageRequest_published && pb.GetPublishedAt() == 0 {
				return nil
			}
			if pb.GetPublishedAt() == 0 {
				pb.PublishedAt = nil
			}
			for _, c := range pb.Contents {
				c.Value = nil
			}
			lst = append(lst, pb)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return lst, nil
}
