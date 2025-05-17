package usecase_test

import (
	"errors"
	"testing"

	mockrepository "ecommerce_clean_arch/pkg/Mock/MockRepository"
	"ecommerce_clean_arch/pkg/domain"
	"ecommerce_clean_arch/pkg/usecase"

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
						HouseName: "Athanikkal",
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
					HouseName: "Athanikkal",
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
