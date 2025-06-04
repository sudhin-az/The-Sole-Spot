package usecase_test

import (
	mockrepository "ecommerce_clean_arch/pkg/Mock/MockRepository"
	"ecommerce_clean_arch/pkg/usecase"

	"ecommerce_clean_arch/pkg/utils/models"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_AddProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mockrepository.NewMockProductRepository(ctrl)
	ProductUseCase := usecase.NewProductUseCase(mockProductRepo)

	testCases := map[string]struct {
		input          models.AddProduct
		expectedOutput models.ProductResponse
		stub           func(mockProductRepo *mockrepository.MockProductRepository)
		expectedError  error
	}{
		"success": {
			input: models.AddProduct{
				ID:         1,
				CategoryID: 1,
				Name:       "Nike",
				Stock:      20,
				Quantity:   10,
				Price:      1500,
				OfferPrice: 1200,
			},
			expectedOutput: models.ProductResponse{
				ID:          1,
				Category_Id: 1,
				Name:        "Nike",
				Stock:       20,
				Quantity:    10,
				Price:       1500,
				OfferPrice:  1200,
			},
			stub: func(mockProductRepo *mockrepository.MockProductRepository) {
				mockProductRepo.EXPECT().AddProduct(models.AddProduct{
					ID:         1,
					CategoryID: 1,
					Name:       "Nike",
					Stock:      20,
					Quantity:   10,
					Price:      1500,
					OfferPrice: 1200,
				}).Return(models.ProductResponse{
					ID:          1,
					Category_Id: 1,
					Name:        "Nike",
					Stock:       20,
					Quantity:    10,
					Price:       1500,
					OfferPrice:  1200,
				}, nil)
			},
			expectedError: nil,
		},
		"negative price": {
			input: models.AddProduct{
				ID:         1,
				CategoryID: 1,
				Name:       "Nike",
				Stock:      20,
				Quantity:   10,
				Price:      -1500,
				OfferPrice: 1000,
			},
			expectedOutput: models.ProductResponse{},
			stub: func(mockProductRepo *mockrepository.MockProductRepository) {

			},
			expectedError: errors.New("invalid quantity or price"),
		},
		"negative quantity": {
			input: models.AddProduct{
				ID:         1,
				CategoryID: 1,
				Name:       "Nike",
				Stock:      20,
				Quantity:   -10,
				Price:      1500,
				OfferPrice: 1000,
			},
			expectedOutput: models.ProductResponse{},
			stub: func(mockProductRepo *mockrepository.MockProductRepository) {

			},
			expectedError: errors.New("invalid quantity or price"),
		},
		"database error": {
			input: models.AddProduct{
				ID:         1,
				CategoryID: 1,
				Name:       "Nike",
				Stock:      20,
				Quantity:   10,
				Price:      1500,
				OfferPrice: 1000,
			},
			expectedOutput: models.ProductResponse{},
			stub: func(mockProductRepo *mockrepository.MockProductRepository) {
				mockProductRepo.EXPECT().AddProduct(gomock.Any()).Return(models.ProductResponse{},
					errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.stub(mockProductRepo)

			result, err := ProductUseCase.AddProduct(tc.input)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedOutput, result)
		})
	}
}

func Test_UpdateProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductRepo := mockrepository.NewMockProductRepository(ctrl)
	productUseCase := usecase.NewProductUseCase(mockProductRepo)

	testCases := map[string]struct {
		input          models.ProductResponse
		productID      int
		expectedOutput models.ProductResponse
		stub           func(mockProductRepo *mockrepository.MockProductRepository)
		expectedError  error
	}{
		"success": {
			input: models.ProductResponse{
				ID:          1,
				Category_Id: 1,
				Name:        "Convacs",
				Stock:       15,
				Quantity:    10,
				Price:       3000,
				OfferPrice:  2500,
			},
			productID: 1,
			expectedOutput: models.ProductResponse{
				ID:          1,
				Category_Id: 1,
				Name:        "Convacs",
				Stock:       15,
				Quantity:    10,
				Price:       3000,
				OfferPrice:  2500,
			},
			stub: func(mockProductRepo *mockrepository.MockProductRepository) {
				mockProductRepo.EXPECT().UpdateProduct(models.ProductResponse{
					ID:          1,
					Category_Id: 1,
					Name:        "Convacs",
					Stock:       15,
					Quantity:    10,
					Price:       3000,
					OfferPrice:  2500,
				}, 1).Return(models.ProductResponse{
					ID:          1,
					Category_Id: 1,
					Name:        "Convacs",
					Stock:       15,
					Quantity:    10,
					Price:       3000,
					OfferPrice:  2500,
				}, nil)
			},
			expectedError: nil,
		},
		"negative price": {
			input: models.ProductResponse{
				ID:          1,
				Category_Id: 1,
				Name:        "Convacs",
				Stock:       15,
				Quantity:    10,
				Price:       -3000,
				OfferPrice:  2500,
			},
			productID:      1,
			expectedOutput: models.ProductResponse{},
			stub: func(mockProductRepo *mockrepository.MockProductRepository) {

			},
			expectedError: errors.New("invalid quantity or price"),
		},
		"negative quantity": {
			input: models.ProductResponse{
				ID:          1,
				Category_Id: 1,
				Name:        "Convacs",
				Stock:       15,
				Quantity:    -10,
				Price:       3000,
				OfferPrice:  2500,
			},
			productID:      1,
			expectedOutput: models.ProductResponse{},
			stub: func(mockProductRepo *mockrepository.MockProductRepository) {

			},
			expectedError: errors.New("invalid quantity or price"),
		},
		"database error": {
			input: models.ProductResponse{
				ID:          1,
				Category_Id: 1,
				Name:        "Convacs",
				Stock:       15,
				Quantity:    10,
				Price:       3000,
				OfferPrice:  2500,
			},
			productID:      1,
			expectedOutput: models.ProductResponse{},
			stub: func(mockProductRepo *mockrepository.MockProductRepository) {
				mockProductRepo.EXPECT().UpdateProduct(gomock.Any(), gomock.Any()).
					Return(models.ProductResponse{}, errors.New("database error"))
			},
			expectedError: errors.New("database error"),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.stub(mockProductRepo)

			result, err := productUseCase.UpdateProduct(tc.input, tc.productID)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectedOutput, result)
		})
	}
}
