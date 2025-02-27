package report

import (
	"backend-profitrack/helpers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Service interface {
	GetAllReportsService(ctx *gin.Context)
	GetDetailReportService(ctx *gin.Context)
	//ExportExcelService(ctx *gin.Context)
	DeleteDetailReportService(ctx *gin.Context)
	DeleteAllReportService(ctx *gin.Context)
}

type reportService struct {
	repository Repository
}

func NewReportService(repo Repository) Service {
	return &reportService{repository: repo}
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

//func (service *reportService) ExportExcelService(ctx *gin.Context) {
//	methodID, err := strconv.Atoi(ctx.Param("methodID"))
//	if err != nil {
//		response := map[string]string{"error": "Invalid method ID"}
//		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
//		return
//	}
//
//	period := ctx.PostForm("period")
//	if period == "" {
//		response := map[string]string{"error": "Period is required"}
//		helpers.ResponseJSON(ctx, http.StatusBadRequest, response)
//		return
//	}
//
//	reports, err := service.repository.GetAllReportRepository(methodID, period)
//	if err != nil {
//		helpers.ResponseJSON(ctx, http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
//		return
//	}
//
//	// Create a new Excel file
//	f := excelize.NewFile()
//	defer f.Close()
//
//	// Create a style for bold text
//	style, err := f.NewStyle(&excelize.Style{
//		Font: &excelize.Font{
//			Bold: true,
//		},
//	})
//	if err != nil {
//		// Handle error
//		helpers.ResponseJSON(ctx, http.StatusInternalServerError, gin.H{"error": "Gagal membuat style"})
//		return
//	}
//
//	// Create header
//	headers := []string{"Rank", "Nama Produk", "Skor Akhir", "Harga Beli", "Harga Jual", "Keuntungan", "Satuan", "Stok", "Stok Terjual"}
//	for i, header := range headers {
//		cell := string(rune('A'+i)) + "1"
//		f.SetCellValue("Sheet1", cell, header)
//		// Set the style for the header cell
//		f.SetCellStyle("Sheet1", cell, cell, style)
//	}
//
//	// Fill data
//	for i, report := range reports {
//		row := i + 2
//		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), i+1) // Rank
//		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), report.Product.Name)
//		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), report.FinalScore)
//		f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), report.Product.PurchaseCost)
//		f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), report.Product.PriceSale)
//		f.SetCellValue("Sheet1", fmt.Sprintf("F%d", row), report.Product.Profit)
//		f.SetCellValue("Sheet1", fmt.Sprintf("G%d", row), report.Product.Unit)
//		f.SetCellValue("Sheet1", fmt.Sprintf("H%d", row), report.Product.Stock)
//		f.SetCellValue("Sheet1", fmt.Sprintf("I%d", row), report.Product.Sold)
//	}
//
//	currentTime := time.Now()
//	fileName := fmt.Sprintf("data-produk-%s.xlsx", currentTime.Format("02-01-2006"))
//
//	// Set response header
//	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
//	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
//
//	if err = f.Write(ctx.Writer); err != nil {
//		helpers.ResponseJSON(ctx, http.StatusInternalServerError, gin.H{"error": "Gagal membuat file Excel"})
//		return
//	}
//}

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
