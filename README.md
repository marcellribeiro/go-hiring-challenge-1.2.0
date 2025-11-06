# Implementation Summary
### Implemented Features
I uploaded the .env file just for ease of testing. 
This is not a good practice for production code.

**API Endpoints:**
- `GET /categories` - List all categories with pagination and total count
- `POST /categories` - Create new category with validation (code and name required)
- `GET /catalog` - List products with category, pagination (offset/limit), and filters
- `GET /catalog/:code` - Get product details including category and variants
- Query parameters support:
  - `offset` (default: 0) and `limit` (default: 10, max: 100, min: 1)
  - `category` - Filter by category code
  - `priceLessThan` - Filter by maximum price

**Key Functionalities:**
- ✅ Products include category information in responses
- ✅ Offset-based pagination with configurable limits
- ✅ Total count returned for both products and categories
- ✅ Category filter applied to product listings
- ✅ Price filter (less than) applied to product listings
- ✅ Multiple filters can be combined
- ✅ Variants inherit product price when their price is 0
- ✅ Categories are persisted in the database
- ✅ Input validation for required fields

**Database Schema:**
- 3 initial categories: CLOTHING, SHOES, ACCESSORIES

**Test Coverage:**
- `app/api` - 100%
- `app/catalog` - 88.9%
- `app/categories` - 76.0%
- `app/product` - 80.0%

---

# Go Hiring Challenge

This repository contains a Go application for managing products and their prices, including functionalities for CRUD operations and seeding the database with initial data.

## Project Structure

1. **cmd/**: Contains the main application and seed command entry points.

   - `server/main.go`: The main application entry point, serves the REST API.
   - `seed/main.go`: Command to seed the database with initial product data.

2. **app/**: Contains the application logic.
3. **sql/**: Contains a very simple database migration scripts setup.
4. **models/**: Contains the data models and repositories used in the application.
5. `.env`: Environment variables file for configuration.

## Setup Code Repository

1. Create a github/bitbucket/gitlab repository and push all this code as-is.
2. Create a new branch, and provide a pull-request against the main branch with your changes. Instructions to follow.

## Application Setup

- Ensure you have Go installed on your machine.
- Ensure you have Docker installed on your machine.
- Important makefile targets:
  - `make tidy`: will install all dependencies.
  - `make docker-up`: will start the required infrastructure services via docker containers.
  - `make seed`: ⚠️ Will destroy and re-create the database tables.
  - `make test`: Will run the tests.
  - `make run`: Will start the application.
  - `make docker-down`: Will stop the docker containers.

Follow up for the assignemnt here: [ASSIGNMENT.md](ASSIGNMENT.md)
