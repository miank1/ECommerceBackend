package repository

import (
	"ecommerce-backend/services/cartservice/internal/model"

	uuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type CartRepository struct {
	DB *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{DB: db}
}

// Find cart by user
func (r *CartRepository) GetByUserID(userID string) (*model.Cart, error) {
	var cart model.Cart
	if err := r.DB.Preload("Items").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}

// Create new cart
func (r *CartRepository) Create(cart *model.Cart) error {
	return r.DB.Create(cart).Error
}

// Save cart (with items)
func (r *CartRepository) Save(cart *model.Cart) error {
	return r.DB.Save(cart).Error
}

// GetItemByID finds a cart item by ID
func (r *CartRepository) GetItemByID(itemID string) (*model.CartItem, error) {
	var item model.CartItem
	if err := r.DB.First(&item, "id = ?", itemID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// UpdateItem updates a cart item
func (r *CartRepository) UpdateItem(item *model.CartItem) error {
	return r.DB.Save(item).Error
}

// DeleteItem removes a cart item by ID
func (r *CartRepository) DeleteItem(itemID string) error {
	return r.DB.Delete(&model.CartItem{}, "id = ?", itemID).Error
}

// DeleteCart deletes a cart and its items
func (r *CartRepository) DeleteCart(cartID string) error {
	return r.DB.Delete(&model.Cart{}, "id = ?", cartID).Error
}

// ClearCart deletes all items and the cart itself
func (r *CartRepository) ClearCart(cartID uuid.UUID) error {
	// Order of delete is important because of foreign key constraint
	if err := r.DB.Delete(&model.CartItem{}, "cart_id = ?", cartID).Error; err != nil {
		return err
	}
	return r.DB.Delete(&model.Cart{}, "id = ?", cartID).Error
}
