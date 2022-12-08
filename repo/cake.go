package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type CakeBaseModel struct {
	ID          int          `db:"id"`
	Title       string       `db:"title"`
	Description string       `db:"decription"`
	Image       string       `db:"image"`
	Rating      float32      `db:"rating"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   sql.NullTime `db:"updated_at"`
}

type Cake interface {
	Create(ctx context.Context, input CakeBaseModel) error
	GetList(ctx context.Context, limit, offset int, search, sort, sortBy string) ([]CakeBaseModel, error)
	GetDetail(ctx context.Context, id int) (CakeBaseModel, error)
	Update(ctx context.Context, input CakeBaseModel) error
	CountCake(ctx context.Context, search string) (count int, err error)
	Delete(ctx context.Context, id int) error
	Close()
}

type cake struct {
	Log   zerolog.Logger
	MySQL *sql.DB
}

func NewCake(mysql *sql.DB, log zerolog.Logger) Cake {
	return &cake{
		MySQL: mysql,
		Log:   log,
	}
}

func (c *cake) Create(ctx context.Context, input CakeBaseModel) (err error) {

	query := "INSERT INTO cake (title, description, image, rating, created_at) VALUES (?, ?, ?, ?, ?)"

	stmt, err := c.MySQL.PrepareContext(ctx, query)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, input.Title, input.Description, input.Image, input.Rating, input.CreatedAt)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}

	return nil
}

func (c *cake) GetList(ctx context.Context, limit, offset int, search, sort, sortBy string) (output []CakeBaseModel, err error) {

	query := "SELECT id, title, description, image, rating, created_at, updated_at FROM cake"

	if search != "" {
		query += fmt.Sprintf(" WHERE title LIKE '%s%s%s' ", "%", search, "%")
	}

	query += fmt.Sprintf(" ORDER BY %s %s ", sort, sortBy)
	query += fmt.Sprintf(" LIMIT %d OFFSET %d ", limit, offset)

	rows, err := c.MySQL.QueryContext(ctx, query)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		cake := CakeBaseModel{}
		err = rows.Scan(&cake.ID, &cake.Title, &cake.Description, &cake.Image, &cake.Rating, &cake.CreatedAt, &cake.UpdatedAt)
		if err != nil {
			c.Log.Error().Msg(err.Error())
			return
		}
		output = append(output, cake)
	}
	return

}

func (c *cake) CountCake(ctx context.Context, search string) (count int, err error) {

	query := "SELECT count(id) FROM cake"

	if search != "" {
		query += fmt.Sprintf(" WHERE title LIKE '%s%s%s' ", "%", search, "%")
	}

	stmt, err := c.MySQL.Prepare(query)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}

	row := stmt.QueryRowContext(ctx)
	err = row.Scan(&count)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}

	return
}

func (c *cake) GetDetail(ctx context.Context, id int) (output CakeBaseModel, err error) {

	query := "SELECT id, title, description, image, rating, created_at, updated_at FROM cake WHERE id = ?"

	row := c.MySQL.QueryRowContext(ctx, query, id)
	err = row.Scan(&output.ID, &output.Title, &output.Description, &output.Image, &output.Rating, &output.CreatedAt, &output.UpdatedAt)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}

	return
}

func (c *cake) Update(ctx context.Context, input CakeBaseModel) (err error) {

	fields := []string{}
	values := []interface{}{}

	if input.Title != "" {
		fields = append(fields, "title = ?")
		values = append(values, input.Title)
	}

	if input.Description != "" {
		fields = append(fields, "description = ?")
		values = append(values, input.Description)
	}

	if input.Image != "" {
		fields = append(fields, "image = ?")
		values = append(values, input.Image)
	}

	if input.Rating != 0 {
		fields = append(fields, "rating = ?")
		values = append(values, input.Rating)
	}

	if input.UpdatedAt.Valid {
		fields = append(fields, "updated_at = ?")
		values = append(values, input.UpdatedAt)
	}

	query := fmt.Sprintf("UPDATE cake SET %s WHERE id = %d", strings.Join(fields, ", "), input.ID)

	stmt, err := c.MySQL.PrepareContext(ctx, query)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, values...)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}

	return
}

func (c *cake) Delete(ctx context.Context, id int) (err error) {

	query := "DELETE FROM cake WHERE id = ?"

	stmt, err := c.MySQL.PrepareContext(ctx, query)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		c.Log.Error().Msg(err.Error())
		return
	}

	return
}

func (r *cake) Close() {
	r.MySQL.Close()
}
