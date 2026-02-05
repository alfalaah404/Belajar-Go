package repositories

import (
	"belajar-go/models"
	"database/sql"
	"fmt"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// mengurangi N+1 query dengan batch insert
	if len(details) > 0 {
		query := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
		args := make([]interface{}, 0, len(details)*4)

		for i := range details {
			details[i].TransactionID = transactionID

			if i > 0 {
				query += ", "
			}

			pos := i * 4
			query += fmt.Sprintf("($%d, $%d, $%d, $%d)", pos+1, pos+2, pos+3, pos+4)

			args = append(args, transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		}

		_, err = tx.Exec(query, args...)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	res := &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}

	return res, nil
}
