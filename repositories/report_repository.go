package repositories

import (
	"belajar-go/models"
	"database/sql"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (repo *ReportRepository) GetReport(startDate, endDate time.Time) (*models.ReportResponse, error) {
	var report models.ReportResponse
	
	// Query untuk total revenue dan total transaksi
	queryRevenue := `
		SELECT 
			COALESCE(SUM(total_amount), 0) as total_revenue,
			COALESCE(COUNT(*), 0) as total_transaksi
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2
	`
	
	err := repo.db.QueryRow(queryRevenue, startDate, endDate).Scan(
		&report.TotalRevenue, 
		&report.TotalTransaksi,
	)
	if err != nil {
		return nil, err
	}
	
	// Query untuk produk terlaris
	queryBestProduct := `
		SELECT 
			p.name,
			COALESCE(SUM(td.quantity), 0) as qty_terjual
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.id, p.name
		ORDER BY qty_terjual DESC
		LIMIT 1
	`
	
	var bestProduct models.BestProduct
	err = repo.db.QueryRow(queryBestProduct, startDate, endDate).Scan(
		&bestProduct.Nama,
		&bestProduct.QtyTerjual,
	)
	
	if err == sql.ErrNoRows {
		// Jika tidak ada data, set produk terlaris ke nil
		report.ProdukTerlaris = nil
	} else if err != nil {
		return nil, err
	} else {
		report.ProdukTerlaris = &bestProduct
	}
	
	return &report, nil
}
