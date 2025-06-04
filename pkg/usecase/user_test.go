package usecase_test

import (
	"errors"
	"testing"

	mockrepository "ecommerce_clean_arch/pkg/Mock/MockRepository"
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/usecase"
	"ecommerce_clean_arch/pkg/utils/models"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_GetAllAddresses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mockrepository.NewMockUserRepository(ctrl)
	userUseCase := usecase.NewUserUseCase(mockUserRepo)

	testCases := map[string]struct {
		userID int
		stub   func(m *mockrepository.MockUserRepository)
		want   []domain.Address
		err    error
	}{
		"success": {
			userID: 1,
			stub: func(m *mockrepository.MockUserRepository) {
				m.EXPECT().GetAllAddresses(1).Return([]domain.Address{
					{
						HouseName: "House No 1",
						Street:    "panayur",
						City:      "ottapalam",
						District:  "Palakkad",
						State:     "kerala",
						Pin:       "679522",
					},
				}, nil)
			},
			want: []domain.Address{
				{
					HouseName: "House No 1",
					Street:    "panayur",
					City:      "ottapalam",
					District:  "Palakkad",
					State:     "kerala",
					Pin:       "679522",
				},
			},
			err: nil,
		},
		"failure": {
			userID: 1,
			stub: func(m *mockrepository.MockUserRepository) {
				m.EXPECT().GetAllAddresses(1).Return(nil, errors.New("error in getting addresses"))
			},
			want: []domain.Address{},
			err:  errors.New("error in getting addresses"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.stub(mockUserRepo)
			got, err := userUseCase.GetAllAddresses(tc.userID)

			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.err, err)
		})
	}
}
func Test_GetAllProducts(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mockrepository.NewMockUserRepository(ctrl)
	userUseCase := usecase.NewUserUseCase(mockProductRepo)

	testCases := map[string]struct {
		stub func(m *mockrepository.MockUserRepository)
		want []models.ProductResponse
		err  error
	}{
		"success": {
			stub: func(m *mockrepository.MockUserRepository) {
				m.EXPECT().GetProducts().Return([]models.ProductResponse{
					{
						ID:          1,
						Category_Id: 44,
						Name:        "nike",
						Quantity:    10,
						Stock:       5,
						Price:       3000,
						OfferPrice:  2000,
					},
				}, nil)
			},
			want: []models.ProductResponse{
				{
					ID:          1,
					Category_Id: 44,
					Name:        "nike",
					Quantity:    10,
					Stock:       5,
					Price:       3000,
					OfferPrice:  2000,
				},
			},
			err: nil,
		},
		"failure": {
			stub: func(m *mockrepository.MockUserRepository) {
				m.EXPECT().GetProducts().Return([]models.ProductResponse{}, errors.New("failed to fetch products"))
			},
			want: []models.ProductResponse{},
			err:  errors.New("failed to fetch products"),
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.stub(mockProductRepo)
			got, err := userUseCase.GetProducts()

			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.err, err)
		})
	}
}
func Test_GetAllCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockCategoryRepo := mockrepository.NewMockUserRepository(ctrl)
	userUseCase := usecase.NewUserUseCase(mockCategoryRepo)

	testCases := map[string]struct {
		stub func(m *mockrepository.MockUserRepository)
		want []domain.Category
		err  error
	}{
		"success": {
			stub: func(m *mockrepository.MockUserRepository) {
				m.EXPECT().ListCategory().Return([]domain.Category{
					{
						ID:               1,
						Category:         "Formal Shoes",
						Description:      "This is Formal Shoes",
						CategoryDiscount: 2,
					},
				}, nil)
			},
			want: []domain.Category{
				{
					ID:               1,
					Category:         "Formal Shoes",
					Description:      "This is Formal Shoes",
					CategoryDiscount: 2,
				},
			},
			err: nil,
		},
		"failure": {
			stub: func(m *mockrepository.MockUserRepository) {
				m.EXPECT().ListCategory().Return([]domain.Category{}, errors.New("failed to fetch categories"))
			},
			want: []domain.Category{},
			err:  errors.New("failed to fetch categories"),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.stub(mockCategoryRepo)
			got, err := userUseCase.ListCategory()

			assert.Equal(t, tc.want, got)
			assert.Equal(t, tc.err, err)
		})
	}
}
