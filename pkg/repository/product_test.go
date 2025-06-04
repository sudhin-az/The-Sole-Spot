package repository_test

import (
	"ecommerce_clean_arch/pkg/repository"
	"ecommerce_clean_arch/pkg/utils/models"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_AddProduct(t *testing.T) {
	tests := []struct {
		name      string
		input     models.AddProduct
		setupMock func(mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name: "Successfully added product",
			input: models.AddProduct{
				CategoryID: 1,
				Name:       "Adidas",
				Quantity:   20,
				Stock:      20,
				Price:      3000,
				OfferPrice: 2300,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(
					`INSERT INTO products (category_id, name, stock, quantity, price, offer_price) 
					VALUES ($1, $2, $3, $4, $5, $6) 
					RETURNING id, category_id, name, stock, quantity, price, offer_price`,
				)
				mock.ExpectQuery(query).
					WithArgs(1, "Adidas", 20, 20, 3000.0, 2300.0).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "category_id", "name", "stock", "quantity", "price", "offer_price",
					}).AddRow(1, 1, "Adidas", 20, 20, 3000.0, 2300.0))
			},
			expectErr: false,
		},
		{
			name: "Failed to add product - DB error",
			input: models.AddProduct{
				CategoryID: 2,
				Name:       "Puma",
				Quantity:   15,
				Stock:      15,
				Price:      2800,
				OfferPrice: 2400,
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(
					`INSERT INTO products (category_id, name, stock, quantity, price, offer_price) 
					VALUES ($1, $2, $3, $4, $5, $6) 
					RETURNING id, category_id, name, stock, quantity, price, offer_price`,
				)
				mock.ExpectQuery(query).
					WithArgs(2, "Puma", 15, 15, 2800.0, 2400.0).
					WillReturnError(gorm.ErrInvalidDB)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSQL, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockSQL.Close()

			db, err := gorm.Open(postgres.New(postgres.Config{
				Conn: mockSQL,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			assert.NoError(t, err)

			productRepo := repository.NewProductRepository(db)

			tt.setupMock(mock)

			_, err = productRepo.AddProduct(tt.input)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_UpdateProduct(t *testing.T) {
	tests := []struct {
		name      string
		input     models.ProductResponse
		productID int
		setupMock func(mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name: "Successfully updated product",
			input: models.ProductResponse{
				Category_Id: 1,
				Name:        "Convacs",
				Quantity:    20,
				Stock:       20,
				Price:       3000,
				OfferPrice:  2500,
			},
			productID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(
					`UPDATE products SET category_id = $1, name = $2, stock = $3, quantity = $4, price = $5, offer_price = $6 WHERE id = $7
				RETURNING id, category_id, name, stock, quantity, price, offer_price`,
				)
				mock.ExpectQuery(query).
					WithArgs(1, "Convacs", 20, 20, 3000.0, 2500.0, 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "category_id", "name", "stock", "quantity", "price", "offer_price",
					}).AddRow(1, 1, "Convacs", 20, 20, 3000.0, 2500.0))
			},
			expectErr: false,
		},
		{
			name: "Failed to Update Product - DB error",
			input: models.ProductResponse{
				Category_Id: 2,
				Name:        "Puma",
				Quantity:    15,
				Stock:       15,
				Price:       2800,
				OfferPrice:  2400,
			},
			productID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				query := regexp.QuoteMeta(
					`UPDATE products SET category_id = $1, name = $2, stock = $3, quantity = $4, price = $5, offer_price = $6 WHERE id = $7
				RETURNING id, category_id, name, stock, quantity, price, offer_price`,
				)
				mock.ExpectQuery(query).
					WithArgs(2, "Puma", 15, 15, 2800.0, 2400.0, 1).
					WillReturnError(gorm.ErrInvalidDB)
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSQL, mock, err := sqlmock.New()
			assert.NoError(t, err)
			defer mockSQL.Close()

			db, err := gorm.Open(postgres.New(postgres.Config{
				Conn: mockSQL,
			}), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})
			assert.NoError(t, err)

			productRepo := repository.NewProductRepository(db)

			tt.setupMock(mock)

			_, err = productRepo.UpdateProduct(tt.input, tt.productID)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
