package repository_test

import (
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/repository"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Test_CreateOrderItems(t *testing.T) {
	tests := []struct {
		name      string
		input     []domain.OrderItem
		setupMock func(mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name: "Succcessfully created OrderItems",
			input: []domain.OrderItem{
				{
					OrderID:    1,
					ProductID:  1,
					Quantity:   10,
					TotalPrice: 2000,
				},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "order_items" 
				("order_id","product_id","quantity","total_price") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs(1, 1, 10, 2000.0).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
				mock.ExpectCommit()
			},
			expectErr: false,
		},
		{
			name: "Failed to create OrderItems",
			input: []domain.OrderItem{
				{
					OrderID:    2,
					ProductID:  2,
					Quantity:   15,
					TotalPrice: 3000,
				},
			},
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "order_items" 
				("order_id","product_id","quantity","total_price") VALUES ($1,$2,$3,$4) RETURNING "id"`)).
					WithArgs(2, 2, 15, 3000.0).
					WillReturnError(gorm.ErrInvalidDB)
				mock.ExpectRollback()
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

			orderRepo := repository.NewOrderRepository(db)

			tt.setupMock(mock)

			tx := db.Begin()
			err = orderRepo.CreateOrderItems(tx, tt.input)

			if tt.expectErr {
				assert.Error(t, err)
				tx.Rollback()
			} else {
				assert.NoError(t, err)
				tx.Commit()
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_GetProductStock(t *testing.T) {
	tests := []struct {
		name      string
		productID int
		mockStock int
		setupMock func(mock sqlmock.Sqlmock, productID int, stock int)
		expectErr bool
	}{
		{
			name:      "Successfully retrieved the Stock from Product",
			productID: 1,
			mockStock: 10,
			setupMock: func(mock sqlmock.Sqlmock, productID, stock int) {
				mock.ExpectQuery(regexp.QuoteMeta(`select stock from products where id = $1`)).
					WithArgs(productID).
					WillReturnRows(sqlmock.NewRows([]string{"stock"}).
						AddRow(stock))
			},
			expectErr: false,
		},
		{
			name:      "Failed to retrieve the Stock from product",
			productID: 1,
			setupMock: func(mock sqlmock.Sqlmock, productID int, stock int) {
				mock.ExpectQuery(regexp.QuoteMeta(`select stock from products where id = $1`)).
					WithArgs(productID).
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

			productRepo := repository.NewOrderRepository(db)

			tt.setupMock(mock, tt.productID, tt.mockStock)

			stock, err := productRepo.GetProductStock(db, tt.productID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Equal(t, 0, stock)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockStock, stock)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
