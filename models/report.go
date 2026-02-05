package models

type ReportResponse struct {
	TotalRevenue    int            `json:"total_revenue"`
	TotalTransaksi  int            `json:"total_transaksi"`
	ProdukTerlaris  *BestProduct   `json:"produk_terlaris"`
}

type BestProduct struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}
