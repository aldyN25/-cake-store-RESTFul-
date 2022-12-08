package repo

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
)

func NewMockCake() (Cake, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	cake := NewCake(db, *zerolog.Ctx(context.Background()))
	return cake, mock
}

func Test_cake_Create(t *testing.T) {
	now := time.Now()
	ctx := context.Background()
	query := "INSERT INTO cake \\(title, description, image, rating, created_at\\) VALUES \\(\\?, \\?, \\?, \\?, \\?\\)"
	cake, mock := NewMockCake()
	defer cake.Close()

	input := CakeBaseModel{
		Title:       "test",
		Description: "test",
		Rating:      10,
		CreatedAt:   now,
	}

	type args struct {
		ctx   context.Context
		input CakeBaseModel
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		beforeFunc func()
	}{
		{
			name: "success",
			args: args{
				ctx:   ctx,
				input: input,
			},
			beforeFunc: func() {
				mock.ExpectPrepare(query).
					ExpectExec().
					WithArgs(input.Title, input.Description, input.Image, input.Rating, now).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "error prepare statement",
			args: args{
				ctx:   ctx,
				input: input,
			},
			beforeFunc: func() {
				mock.ExpectPrepare(query).WillReturnError(errors.New("foo"))
			},
			wantErr: true,
		},
		{
			name: "error",
			args: args{
				ctx:   ctx,
				input: input,
			},
			beforeFunc: func() {
				mock.ExpectPrepare(query).
					ExpectExec().
					WithArgs(input.Title, input.Description, input.Rating, now).
					WillReturnError(errors.New("foo"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			if err := cake.Create(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("cake.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cake_GetList(t *testing.T) {

	ctx := context.Background()
	now := time.Now()
	cake, mock := NewMockCake()
	defer cake.Close()

	query := "SELECT id, title, description, image, rating, created_at, updated_at FROM cake ORDER BY id asc LIMIT 10 OFFSET 0"

	output := CakeBaseModel{
		ID:          1,
		Title:       "test",
		Description: "test",
		Rating:      1,
		CreatedAt:   now,
	}

	type args struct {
		ctx    context.Context
		limit  int
		offset int
		search string
		sort   string
		sortBy string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput []CakeBaseModel
		beforeFunc func()
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				ctx:    ctx,
				limit:  10,
				offset: 0,
				search: "",
				sort:   "id",
				sortBy: "asc",
			},
			wantOutput: []CakeBaseModel{
				output,
			},
			beforeFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "image", "rating", "created_at", "updated_at"}).
					AddRow(output.ID, output.Title, output.Description, output.Image, output.Rating, output.CreatedAt, output.UpdatedAt)
				mock.ExpectQuery(query).WillReturnRows(rows)

			},
			wantErr: false,
		},
		{
			name: "error scan result",
			args: args{
				ctx:    ctx,
				limit:  10,
				offset: 0,
				search: "",
				sort:   "id",
				sortBy: "asc",
			},
			wantOutput: nil,
			beforeFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "description", "image", "rating", "created_at", "updated_at"}).
					AddRow(output.ID, output.Description, output.Image, output.Rating, output.CreatedAt, output.UpdatedAt)
				mock.ExpectQuery(query).WillReturnRows(rows)

			},
			wantErr: true,
		},
		{
			name: "error",
			args: args{
				ctx:    ctx,
				limit:  10,
				offset: 0,
				search: "",
				sort:   "id",
				sortBy: "asc",
			},
			wantOutput: nil,
			beforeFunc: func() {
				mock.ExpectQuery(query).WillReturnError(errors.New("foo"))

			},
			wantErr: true,
		},
		{
			name: "error with search query",
			args: args{
				ctx:    ctx,
				limit:  10,
				offset: 0,
				search: "chesscake",
				sort:   "id",
				sortBy: "asc",
			},
			wantOutput: nil,
			beforeFunc: func() {
				mock.ExpectQuery(query).WillReturnError(errors.New("foo"))

			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			gotOutput, err := cake.GetList(tt.args.ctx, tt.args.limit, tt.args.offset, tt.args.search, tt.args.sort, tt.args.sortBy)
			if (err != nil) != tt.wantErr {
				t.Errorf("cake.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("cake.GetList() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func Test_cake_GetDetail(t *testing.T) {
	ctx := context.Background()
	now := time.Now()
	cake, mock := NewMockCake()
	defer cake.Close()

	query := "SELECT id, title, description, image, rating, created_at, updated_at FROM cake WHERE id = ?"
	output := CakeBaseModel{
		ID:          1,
		Title:       "test",
		Description: "test",
		Rating:      1,
		CreatedAt:   now,
	}

	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name       string
		args       args
		wantOutput CakeBaseModel
		beforeFunc func()
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				id:  1,
			},
			wantOutput: output,
			beforeFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "image", "rating", "created_at", "updated_at"}).
					AddRow(output.ID, output.Title, output.Description, output.Image, output.Rating, output.CreatedAt, output.UpdatedAt)
				mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows)

			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				ctx: ctx,
				id:  1,
			},
			wantOutput: CakeBaseModel{},
			beforeFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description", "image", "rating", "created_at", "updated_at"}).
					AddRow(output.ID, output.Title, output.Description, output.Image, output.Rating, output.CreatedAt, output.UpdatedAt)
				mock.ExpectQuery(query).WithArgs(1).WillReturnRows(rows).WillReturnError(errors.New("foo"))

			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			gotOutput, err := cake.GetDetail(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("cake.GetDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("cake.GetDetail() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func Test_cake_Update(t *testing.T) {

	now := sql.NullTime{Time: time.Now(), Valid: true}
	ctx := context.Background()
	query := "UPDATE cake SET title = \\?, description = \\?, image = \\?, rating = \\?, updated_at = \\? WHERE id = 1"
	cake, mock := NewMockCake()
	defer cake.Close()

	input := CakeBaseModel{
		Title:       "test",
		Description: "test",
		Rating:      10,
		Image:       "test",
		UpdatedAt:   now,
		ID:          1,
	}

	type args struct {
		ctx   context.Context
		input CakeBaseModel
	}
	tests := []struct {
		name       string
		args       args
		beforeFunc func()
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				ctx:   ctx,
				input: input,
			},
			beforeFunc: func() {
				mock.ExpectPrepare(query).
					ExpectExec().
					WithArgs(input.Title, input.Description, input.Image, input.Rating, now).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				ctx:   ctx,
				input: input,
			},
			beforeFunc: func() {
				mock.ExpectPrepare(query).
					ExpectExec().
					WithArgs(input.Title, input.Description, input.Image, input.Rating, now).
					WillReturnResult(sqlmock.NewResult(0, 1)).WillReturnError(errors.New("foo"))
			},
			wantErr: true,
		},
		{
			name: "error prepare statement",
			args: args{
				ctx:   ctx,
				input: input,
			},
			beforeFunc: func() {
				mock.ExpectPrepare(query).
					WillReturnError(errors.New("foo"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			if err := cake.Update(tt.args.ctx, tt.args.input); (err != nil) != tt.wantErr {
				t.Errorf("cake.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cake_Delete(t *testing.T) {

	ctx := context.Background()
	id := 1
	query := "DELETE FROM cake WHERE id = \\?"
	cake, mock := NewMockCake()
	defer cake.Close()

	type args struct {
		ctx context.Context
		id  int
	}
	tests := []struct {
		name       string
		args       args
		beforeFunc func()
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				ctx: ctx,
				id:  id,
			},
			beforeFunc: func() {
				mock.ExpectPrepare(query).
					ExpectExec().
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				ctx: ctx,
				id:  id,
			},
			beforeFunc: func() {
				mock.ExpectPrepare(query).
					ExpectExec().
					WithArgs(id).
					WillReturnError(errors.New("foo"))
			},
			wantErr: true,
		},
		{
			name: "error prepare statement",
			args: args{
				ctx: ctx,
				id:  id,
			},
			beforeFunc: func() {
				mock.ExpectPrepare(query).
					WillReturnError(errors.New("foo"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			if err := cake.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("cake.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_cake_CountCake(t *testing.T) {

	ctx := context.Background()
	query := "SELECT count\\(id\\) FROM cake"
	cake, mock := NewMockCake()
	defer cake.Close()

	type args struct {
		ctx    context.Context
		search string
	}
	tests := []struct {
		name       string
		args       args
		wantCount  int
		beforeFunc func()
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				ctx:    ctx,
				search: "",
			},
			wantCount: 1,
			beforeFunc: func() {
				rows := sqlmock.NewRows([]string{"count"}).
					AddRow(1)

				mock.ExpectPrepare(query).
					ExpectQuery().
					WillReturnRows(rows)
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				ctx:    ctx,
				search: "",
			},
			wantCount: 0,
			beforeFunc: func() {
				mock.ExpectPrepare(query).
					ExpectQuery().
					WillReturnError(errors.New("foo"))
			},
			wantErr: true,
		},
		{
			name: "error prepare statement",
			args: args{
				ctx:    ctx,
				search: "",
			},
			wantCount: 0,
			beforeFunc: func() {
				mock.ExpectPrepare(query).WillReturnError(errors.New("foo"))
			},
			wantErr: true,
		},
		{
			name: "error with search param",
			args: args{
				ctx:    ctx,
				search: "chesscake",
			},
			wantCount: 0,
			beforeFunc: func() {
				mock.ExpectPrepare(query).
					ExpectQuery().
					WillReturnError(errors.New("foo"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.beforeFunc()
			c := cake
			gotCount, err := c.CountCake(tt.args.ctx, tt.args.search)
			if (err != nil) != tt.wantErr {
				t.Errorf("cake.CountCake() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCount != tt.wantCount {
				t.Errorf("cake.CountCake() = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}
