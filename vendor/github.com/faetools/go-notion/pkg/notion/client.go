package notion

import (
	"context"
	"fmt"
	"net/http"

	"github.com/faetools/client"
	"github.com/google/uuid"
)

const (
	versionHeader           = "Notion-Version"
	version                 = "2022-02-22"
	maxPageSize    PageSize = 100
	maxPageSizeInt          = 100
)

// NewDefaultClient returns a new client with the default options.
func NewDefaultClient(bearer string, opts ...client.Option) (*Client, error) {
	opts = append([]client.Option{
		client.WithBearer(bearer),
		client.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
			req.Header.Set(versionHeader, version)
			return nil
		}),
	}, opts...)

	return NewClient(opts...)
}

// GetAllBlocksWithChildren returns all blocks of a given page, including block children.
func (c Client) GetAllBlocksWithChildren(ctx context.Context, id Id) (Blocks, error) {
	blocks, err := c.GetAllBlocks(ctx, id)
	if err != nil {
		return nil, err
	}

	for _, b := range blocks {
		if err := c.getAllChildren(ctx, b); err != nil {
			return nil, err
		}
	}

	return blocks, nil
}

// GetAllBlocks returns all blocks of a given page or block.
func (c Client) GetAllBlocks(ctx context.Context, id Id) (Blocks, error) {
	blocks := Blocks{}

	var cursor *StartCursor
	for {
		resp, err := c.GetBlocks(ctx, id, &GetBlocksParams{
			PageSize:    maxPageSize,
			StartCursor: cursor,
		})
		if err != nil {
			return nil, fmt.Errorf("getting blocks for %s: %w", id, err)
		}

		switch resp.StatusCode() {
		case http.StatusOK: // ok
		case http.StatusBadRequest:
			return nil, resp.JSON400
		case http.StatusNotFound:
			return nil, resp.JSON404
		default:
			return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
		}

		blocks = append(blocks, resp.JSON200.Results...)

		if !resp.JSON200.HasMore {
			return blocks, nil
		}

		cursor = (*StartCursor)(&resp.JSON200.NextCursor)
	}
}

// ListAllUsers returns all users in the workspace.
func (c Client) ListAllUsers(ctx context.Context) (Users, error) {
	users := Users{}

	var cursor *StartCursor
	for {
		resp, err := c.ListUsers(ctx, &ListUsersParams{
			PageSize:    maxPageSize,
			StartCursor: cursor,
		})
		if err != nil {
			return nil, err
		}

		switch resp.StatusCode() {
		case http.StatusOK: // ok
		case http.StatusBadRequest:
			return nil, resp.JSON400
		case http.StatusNotFound:
			return nil, resp.JSON404
		case http.StatusTooManyRequests:
			return nil, resp.JSON429
		default:
			return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
		}

		users = append(users, resp.JSON200.Results...)

		if !resp.JSON200.HasMore {
			return users, nil
		}

		cursor = (*StartCursor)(&resp.JSON200.NextCursor)
	}
}

// GetDatabaseEntries return the database entries or an error.
func (c Client) GetDatabaseEntries(ctx context.Context, id Id, filter *Filter, sorts *Sorts) (Pages, error) {
	entries := Pages{}

	var cursor *UUID
	for {
		resp, err := c.QueryDatabase(ctx, id,
			QueryDatabaseJSONRequestBody{
				Filter:      filter,
				PageSize:    maxPageSizeInt,
				Sorts:       sorts,
				StartCursor: cursor,
			})
		if err != nil {
			return nil, err
		}

		switch resp.StatusCode() {
		case http.StatusOK: // ok
		case http.StatusBadRequest:
			return nil, resp.JSON400
		case http.StatusNotFound:
			return nil, resp.JSON404
		case http.StatusTooManyRequests:
			return nil, resp.JSON429
		default:
			return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
		}

		entries = append(entries, resp.JSON200.Results...)

		if !resp.JSON200.HasMore {
			return entries, nil
		}

		cursor = (*UUID)(&resp.JSON200.NextCursor)
	}
}

// GetNotionPage return the notion page or an error.
func (c Client) GetNotionPage(ctx context.Context, id Id) (*Page, error) {
	resp, err := c.GetPage(ctx, id)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

func ensureDatabaseIsValid(db *Database) {
	db.Object = "database"

	if db.Parent != nil && db.Parent.PageId != "" {
		db.Parent.Type = "page_id"
	}

	props := db.GetProperties()

	for _, prop := range props {
		if prop.Title != nil {
			return
		}
	}

	props["title"] = TitleProperty
	db.Properties = props
}

// GetNotionDatabase returns the notion database or an error.
func (c Client) GetNotionDatabase(ctx context.Context, id Id) (*Database, error) {
	resp, err := c.GetDatabase(ctx, id)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

// CreateNotionDatabase creates a notion database or returns an error.
func (c Client) CreateNotionDatabase(ctx context.Context, db Database) (*Database, error) {
	if db.Id == "" {
		// create a new UUID for the database
		db.Id = UUID(uuid.NewString())
	}

	ensureDatabaseIsValid(&db)

	resp, err := c.CreateDatabase(ctx, CreateDatabaseJSONRequestBody(db))
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

// UpdateNotionDatabase updates a notion database or returns an error.
func (c Client) UpdateNotionDatabase(ctx context.Context, db Database) (*Database, error) {
	// can't be present when updating
	db.Parent = nil
	db.CreatedTime = nil

	ensureDatabaseIsValid(&db)

	resp, err := c.UpdateDatabase(ctx, Id(db.Id), UpdateDatabaseJSONRequestBody(db))
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return resp.JSON200, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	// case http.StatusNotFound:
	// 	return nil, resp.JSON404
	// case http.StatusTooManyRequests:
	// 	return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

// CreateOrUpdateNotionDatabase either creates or updates a notion database or returns an error.
func (c Client) CreateOrUpdateNotionDatabase(ctx context.Context, db Database) (*Database, error) {
	if db.Id == "" {
		return c.CreateNotionDatabase(ctx, db)
	}

	existing, err := c.GetNotionDatabase(ctx, Id(db.Id))
	if err != nil {
		return nil, err
	}

	newProps, changed := existing.GetProperties().Merge(db.GetProperties())

	changed = changed || existing.Title.Raw() != db.Title.Raw()

	if !changed {
		return existing, nil
	}

	db.Properties = newProps

	return c.UpdateNotionDatabase(ctx, db)
}

// AppendNotionBlocks creates or returns an error.
func (c Client) AppendNotionBlocks(ctx context.Context, id Id, blocks Blocks) (*Blocks, error) {
	resp, err := c.AppendBlocks(ctx, id, AppendBlocksJSONRequestBody{
		Children: blocks,
	})
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode() {
	case http.StatusOK: // ok
		return &resp.JSON200.Results, nil
	case http.StatusBadRequest:
		return nil, resp.JSON400
	case http.StatusNotFound:
		return nil, resp.JSON404
	case http.StatusTooManyRequests:
		return nil, resp.JSON429
	default:
		return nil, fmt.Errorf("unknown error response: %v", string(resp.Body))
	}
}

// getAllChildren returns all block children of a given block.
func (c Client) getAllChildren(ctx context.Context, b Block) error {
	if !b.HasChildren {
		return nil
	}

	children, err := c.GetAllBlocks(ctx, Id(b.Id))
	if err != nil {
		return fmt.Errorf("getting children of %s: %w", b.Id, err)
	}

	// get nested children as well
	for _, child := range children {
		if err := c.getAllChildren(ctx, child); err != nil {
			return err
		}
	}

	switch b.Type {
	case BlockTypeBulletedListItem:
		b.BulletedListItem.Children = children
	default:
		return fmt.Errorf("unknown parent block of type %s", b.Type)
	}

	return nil
}
