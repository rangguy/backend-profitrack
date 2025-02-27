package report

import (
	"backend-profitrack/helpers"
	"backend-profitrack/modules/method"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"net/http"
	"strconv"
	"time"
)

type Service interface {
	CountReportsService(ctx *gin.Context)
	GetAllReportsService(ctx *gin.Context)
	GetDetailReportService(ctx *gin.Context)
	ExportPDFService(ctx *gin.Context)
	DeleteDetailReportService(ctx *gin.Context)
	DeleteAllReportService(ctx *gin.Context)
}

type reportService struct {
	repository       Repository
	methodRepository method.Repository
}

func NewReportService(repo Repository, methodRepo method.Repository) Service {
	return &reportService{repository: repo, methodRepository: methodRepo}
}

func (service *reportService) CountReportsService(ctx *gin.Context) {
	result, err := service.repository.CountReportsRepository()
	if err != nil {
		response := map[string]string{"error": "gagal menghitung jumlah data laporan"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]int{"count": int(result)}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}

func (service *reportService) GetAllReportsService(ctx *gin.Context) {
	// Retrieve reports from the repository
	reports, err := service.repository.GetAllReportsRepository()
	if err != nil {
		response := map[string]string{"error": "Failed to retrieve reports"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if len(reports) == 0 {
		response := map[string]string{"error": "data laporan masih kosong"}
		helpers.ResponseJSON(ctx, http.StatusNotFound, response)
		return
	}

	// Respond with the retrieved reports
	helpers.ResponseJSON(ctx, http.StatusOK, reports)
}

func (service *reportService) GetDetailReportService(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := map[string]string{"error": "Invalid id"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	reports, err := service.repository.GetAllReportDetailRepository(id)
	if err != nil {
		response := map[string]string{"error": "gagal mendapatkan detail laporan"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	if len(reports) == 0 {
		response := map[string]string{"error": "data laporan masih kosong"}
		helpers.ResponseJSON(ctx, http.StatusNotFound, response)
		return
	}

	helpers.ResponseJSON(ctx, http.StatusOK, reports)
}

func (service *reportService) ExportPDFService(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := map[string]string{"error": "Invalid ID"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	reports, err := service.repository.GetAllReportDetailRepository(id)
	if err != nil {
		response := map[string]string{"error": "gagal mendapatkan detail laporan"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	getMethod, err := service.methodRepository.GetMethodByIdRepository(reports[0].MethodID)
	if err != nil {
		response := map[string]string{"error": "gagal mendapatkan nama metode"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	methodName := getMethod.Name

	// Create a new PDF document
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set margins (left, top, right)
	pdf.SetMargins(15, 15, 15)

	// Add title
	pdf.SetFont("Times", "B", 16)
	pdf.CellFormat(0, 10, fmt.Sprintf("Laporan Hasil Perhitungan %s", methodName), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Add descriptive paragraph
	pdf.SetFont("Times", "", 12)
	description := fmt.Sprintf(
		"Laporan ini menyajikan hasil perhitungan menggunakan sistem pendukung keputusan (SPK) dengan metode %s. "+
			"Penilaian didasarkan pada kriteria kinerja produk dari perspektif keuangan, termasuk Return On Investment, Net Profit Margin, dan Rasio Efisiensi. "+
			"Berikut ini adalah nilai akhir, peringkat, serta detail data masing-masing produk",
		methodName,
	)
	pdf.MultiCell(0, 8, description, "", "L", false)
	pdf.Ln(10)

	// Set headers with styling
	pdf.SetFont("Times", "B", 10)
	headers := []string{"Rank", "Nama Produk", "Skor Akhir", "Harga Beli", "Harga Jual", "Keuntungan", "Satuan", "Stok", "Stok Terjual"}
	colWidths := []float64{10, 30, 20, 20, 20, 25, 15, 15, 20} // Column widths in mm

	// Add header background
	pdf.SetFillColor(200, 200, 200) // Light gray background for header

	// Print headers
	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 10, header, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Reset fill color and set regular font for data
	pdf.SetFillColor(255, 255, 255)
	pdf.SetFont("Times", "", 10)

	// Add data rows
	for i, report := range reports {
		// Alternate row colors for better readability
		fill := i%2 == 1
		if fill {
			pdf.SetFillColor(240, 240, 240) // Very light gray for alternating rows
		} else {
			pdf.SetFillColor(255, 255, 255) // White
		}

		// Format numbers with appropriate separators
		purchaseCost := fmt.Sprintf(formatCurrency(float64(report.Product.PurchaseCost)))
		priceSale := fmt.Sprintf(formatCurrency(float64(report.Product.PriceSale)))
		profit := fmt.Sprintf(formatCurrency(float64(report.Product.Profit)))

		// Print data row
		pdf.CellFormat(colWidths[0], 8, fmt.Sprintf("%d", i+1), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[1], 8, report.Product.Name, "1", 0, "L", fill, 0, "")
		pdf.CellFormat(colWidths[2], 8, fmt.Sprintf("%f", report.FinalScore), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[3], 8, purchaseCost, "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[4], 8, priceSale, "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[5], 8, profit, "1", 0, "R", fill, 0, "")
		pdf.CellFormat(colWidths[6], 8, report.Product.Unit, "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[7], 8, fmt.Sprintf("%d", report.Product.Stock), "1", 0, "C", fill, 0, "")
		pdf.CellFormat(colWidths[8], 8, fmt.Sprintf("%d", report.Product.Sold), "1", 0, "C", fill, 0, "")
		pdf.Ln(-1)
	}

	// Add footer with date
	pdf.Ln(10)
	pdf.SetFont("Times", "I", 8)
	currentTime := time.Now()
	pdf.CellFormat(0, 10, fmt.Sprintf("Laporan dibuat pada: %s", reports[0].CreatedAt.Format("02-01-2006 15:04:05")), "", 0, "R", false, 0, "")

	// Set filename
	fileName := fmt.Sprintf("data-produk-%s.pdf", currentTime.Format("02-01-2006"))

	// Set headers
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))

	// Write PDF to response
	err = pdf.Output(ctx.Writer)
	if err != nil {
		response := map[string]string{"error": "gagal membuat file pdf"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}
}

// Helper function to format currency
func formatCurrency(amount float64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("Rp %.0f", amount)
}

func (service *reportService) DeleteDetailReportService(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := map[string]string{"error": "Invalid ID"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	_, err = service.repository.GetAllReportDetailRepository(id)
	if err != nil {
		response := map[string]string{"error": "gagal mendapatkan detail laporan"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	// menghapus berdasarkan report id
	err = service.repository.DeleteDetailReportRepository(id)
	if err != nil {
		response := map[string]string{"error": "gagal menghapus detail laporan"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "berhasil menghapus detail laporan"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)

}

func (service *reportService) DeleteAllReportService(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := map[string]string{"error": "Invalid ID"}
		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
		return
	}

	report, err := service.repository.GetReportByIDRepository(id)
	if err != nil {
		response := map[string]string{"error": "gagal mendapatkan laporan"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	err = service.repository.DeleteReportRepository(&report)
	if err != nil {
		response := map[string]string{"error": "gagal menghapus laporan"}
		helpers.ResponseJSON(ctx, http.StatusInternalServerError, response)
		return
	}

	response := map[string]string{"message": "berhasil menghapus laporan"}
	helpers.ResponseJSON(ctx, http.StatusOK, response)
}
