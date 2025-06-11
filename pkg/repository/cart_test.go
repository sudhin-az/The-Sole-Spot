package repository_test

import (
	"ecommerce_clean_arch/pkg/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_GetCartItem(t *testing.T) {
	tests := []struct {
		name      string
		userID    int
		productID int
		setupMock func(mock sqlmock.Sqlmock)
		expectErr bool
	}{
		{
			name:      "Successfully retrieved the cart item",
			userID:    1,
			productID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {

				mock.ExpectQuery(`SELECT \* FROM "carts" WHERE \(user_id = \$1 AND product_id = \$2\) AND "carts"\."deleted_at" IS NULL ORDER BY "carts"\."id" LIMIT \$3`).
					WithArgs(1, 1, 1).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "user_id", "product_id", "quantity", "price", "offer_price", "category_discount", "total_price", "created_at", "deleted_at",
					}).
						AddRow(1, 1, 1, 10, 2000, 1500, 2.0, 3000.0, time.Now(), nil))
			},
			expectErr: false,
		},
		{
			name:      "Query failed",
			userID:    1,
			productID: 1,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT \* FROM "carts" WHERE \(user_id = \$1 AND product_id = \$2\) AND "carts"\."deleted_at" IS NULL ORDER BY "carts"\."id" LIMIT \$3`).
					WithArgs(1, 1, 1).
					WillReturnError(gorm.ErrInvalidDB)
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			defer mockDB.Close()
			assert.NoError(t, err)

			db, err := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			assert.NoError(t, err)

			cartRepo := repository.NewCartRepository(db)
			tt.setupMock(mock)

			cartItems, err := cartRepo.GetCartItem(tt.userID, tt.productID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, cartItems)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cartItems)
				assert.Equal(t, tt.userID, cartItems.UserID)
				assert.Equal(t, tt.productID, cartItems.ProductID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
