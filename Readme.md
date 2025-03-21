# Sole-Spot E-commerce Project

## Project Overview
**Sole-Spot** is an e-commerce platform built using the **Clean Architecture** in **Go** (Golang). This project focuses on backend development, providing a structured and scalable architecture for managing products, users, orders, payments, and more.

## Tech Stack
- **Backend:** Golang (Gin Framework)
- **Database:** PostgreSQL
- **Authentication:** OTP-based, Single Sign-On Google
- **Payment Gateway:** Razorpay
- **Clean Architecture:** Modular and maintainable structure

---

## Features

### 1. Admin Side
#### a. Authentication
- Admin sign-in

#### b. User Management
- List users
- Block/unblock users

#### c. Category Management
- Add, edit, and soft delete categories

#### d. Product Management
- Add, edit, and soft delete products

#### e. Order Management
- List orders
- Change order status
- Cancel orders

#### f. Inventory/Stock Management
- Manage product stock levels

#### g. Offer & Coupon Management
- Product offer, category offer
- Create and delete coupons

#### h. Sales Report
- Generate reports (Daily, Weekly, Yearly, Custom date)
- Show discount and coupon deductions
- Overall sales count and order amount
- Download reports in **PDF or Excel**

#### i. Admin Dashboard
- Sales chart with filters (yearly, monthly, etc.)
- Best-selling products, categories, and brands (Top 10)
- Generate **Ledger Book (Optional)**

---

### 2. User Side
#### a. Authentication
- User signup and login with validation
- Signup using OTP with timer and resend functionality
- Single Sign-On Google

#### b. Product Listing & Details
- List all products with:
  - Ratings, price, discounts, and reviews
  - Stock availability and sold-out indicators
  - Related product recommendations
- Advanced search with sorting options:
  - Popularity
  - Price (Low to High, High to Low)
  - Average Ratings
  - Featured, New Arrivals
  - Alphabetical (A-Z, Z-A)

#### c. User Profile Management
- Show user details, addresses, and orders
- Edit profile, cancel orders, forgot password, change password
- Manage multiple addresses (Add, Edit, Delete)

#### d. Cart Management
- Add to cart, remove products, and list cart items
- Control quantity based on stock
- Limit max quantity per user
- Hide or show out-of-stock products based on filter

#### e. Checkout & Payment
- Select address for checkout (Multiple/saved addresses)
- **Cash on Delivery (COD) restrictions:** Orders above Rs. 1000 not allowed
- Apply delivery charges based on location (optional)
- Integrate online payments (Razorpay)
- Handle payment failures with status update and retry option

#### f. Order Management
- Order cancellation, order history, and status tracking
- Download invoice (PDF)

#### g. Wishlist & Wallet
- Add/remove products from wishlist
- Wallet system for canceled orders

---

## Installation & Setup
1. Clone the repository:
   ```sh
   git clone https://github.com/sudhin-az/The-Sole-Spot.git
   ```
2. Navigate to the project folder:
   ```sh
   cd sole-spot
   ```
3. Install dependencies:
   ```sh
   go mod tidy
   ```
4. Set up the `.env` file for database and SMTP configurations.
5. Run the application:
   ```sh
   go run main.go
   ```

## Database Schema
- **Users Table**: Stores user details
- **Products Table**: Stores product details
- **Catgories Table**: Stores category details
- **Reviews Table**: Stores product review details
- **Orders Table**: Stores order details
- **Cart Table**: Manages user carts
- **Coupons Table**: Stores available discount coupons

## API Endpoints
(To be documented separately)

## Future Enhancements
- Implement frontend using js
- Add additional payment gateways
- Implement AI-powered recommendations

## Contributing
Feel free to fork the repo and submit pull requests.

## License
This project is open-source under the MIT License.

