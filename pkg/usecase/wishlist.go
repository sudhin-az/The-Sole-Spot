package usecase

import (
	"ecommerce_clean_architecture/pkg/repository"
	"ecommerce_clean_architecture/pkg/utils/models"
	"errors"
)

type WishlistUseCase struct {
	wishlistRepo repository.WishlistRepository
}

func NewWishlistUseCase(repo repository.WishlistRepository) *WishlistUseCase {
	return &WishlistUseCase{
		wishlistRepo: repo,
	}
}

func (w *WishlistUseCase) AddToWishList(productID int, userID int) error {
	productExist, err := w.wishlistRepo.DoesProductExist(productID)
	if err != nil {
		return err
	}
	if !productExist {
		return errors.New("product does not exist")
	}
	exists, err := w.wishlistRepo.ProductExistInWishList(productID, userID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("product is already in the wishlist")
	}
	err = w.wishlistRepo.AddToWishList(userID, productID)
	if err != nil {
		return err
	}
	return nil
}
func (w *WishlistUseCase) GetWishList(userID int) ([]models.WishListResponse, error) {
	wishedProducts, err := w.wishlistRepo.GetWishList(userID)
	if err != nil {
		return []models.WishListResponse{}, err
	}
	return wishedProducts, nil
}

func (w *WishlistUseCase) RemoveFromWishList(productID int, userID int) error {

	exists, err := w.wishlistRepo.ProductExistInWishList(productID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("No product found in with this id")
	}
	err = w.wishlistRepo.RemoveFromWishList(userID, productID)
	if err != nil {
		return err
	}
	return nil
}
