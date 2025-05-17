package repository_test

import (
	"regexp"
	"testing"
	"time"

	"ecommerce_clean_arch/pkg/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_GetAllAddresses(t *testing.T) {
	tests := []struct {
		name      string
		setupMock func(mock sqlmock.Sqlmock)
		expectErr bool
		expectLen int
	}{
		{
			name: "Successfully Retrieved all Addresses",
			setupMock: func(mock sqlmock.Sqlmock) {
				query := "SELECT * FROM addresses WHERE deleted_at IS NULL"

				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "user_id", "house_name", "street", "city", "district", "state", "pin", "created_at", "deleted_at",
					}).
						AddRow(101, 52, "HouseNo1", "vyttila", "kochi", "ekm", "ker", "366278", time.Now(), nil))
			},
			expectErr: false,
			expectLen: 1,
		},
		{
			name: "Query Failed",
			setupMock: func(mock sqlmock.Sqlmock) {
				query := "SELECT * FROM addresses WHERE deleted_at IS NULL"
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WillReturnError(gorm.ErrInvalidDB)
			},
			expectErr: true,
			expectLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			assert.NoError(t, err)

			db, err := gorm.Open(postgres.New(postgres.Config{
				Conn: mockDB,
			}), &gorm.Config{})
			assert.NoError(t, err)

			userRepo := repository.NewUserRepository(db)
			tt.setupMock(mock)

			addresses, err := userRepo.GetAllAddresses(1)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Len(t, addresses, 0)
			} else {
				assert.NoError(t, err)
				assert.Len(t, addresses, tt.expectLen)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
