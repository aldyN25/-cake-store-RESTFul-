package cake

import (
	"errors"
	"io"
	"mime/multipart"
	"net/url"
	"strconv"
	"time"
)

type CreateRequestJSON struct {
	Title       string    `json:"title"`
	Rating      float32   `json:"rating"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *CreateRequestJSON) Validate() error {

	if c.Title == "" {
		return errors.New("title cannot by empty")
	}

	if c.Description == "" {
		return errors.New("description cannot be empty")
	}

	if c.Rating <= 0 {
		return errors.New("rating cannot less than 1")
	}

	if c.Image == "" {
		return errors.New("image cannot be empty")
	}

	return nil

}

type CreateRequest struct {
	Title       string    `json:"title"`
	Rating      float32   `json:"rating"`
	Description string    `json:"description"`
	Image       io.Reader `json:"image,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *CreateRequest) ParseForm(v url.Values) error {

	if v.Get("title") == "" {
		return errors.New("title cannot by empty")
	}

	if v.Get("description") == "" {
		return errors.New("description cannot be empty")
	}

	if v.Get("rating") == "" {
		return errors.New("rating cannot be rmpty")
	}

	x, err := strconv.ParseFloat(v.Get("rating"), 32)
	if err != nil {
		return errors.New("rating is invalid")
	}

	if x < 0 {
		return errors.New("rating cannot less than 0")
	}

	c.Title = v.Get("title")
	c.Description = v.Get("description")
	c.Rating = float32(x)

	return nil
}

type UpdateRequest struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	Rating      float32        `json:"rating"`
	Description string         `json:"description"`
	Image       multipart.File `json:"image"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

func (c *UpdateRequest) ParseForm(v url.Values) (err error) {

	x := float64(0)
	if v.Get("rating") != "" {
		x, err = strconv.ParseFloat(v.Get("rating"), 32)
		if err != nil {
			return errors.New("rating is invalid")
		}
	}

	if x < 0 {
		return errors.New("rating cannot less than 0")
	}

	c.Title = v.Get("title")
	c.Description = v.Get("description")
	c.Rating = float32(x)

	return nil
}

type GetListRequest struct {
	Search string
	Sort   string
	SortBy string
}

func (c *GetListRequest) ParseQuery(v url.Values) {

	c.Search = v.Get("search")
	c.SortBy = v.Get("sort_by")
	c.Sort = v.Get("sort")

	if c.Sort == "" {
		c.Sort = "id,title"
	}

	if c.SortBy == "" {
		c.SortBy = "ASC"
	}
}
