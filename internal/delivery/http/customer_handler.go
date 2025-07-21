package http

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"xyz-multifinance/internal/model"
	"xyz-multifinance/internal/usecase"
	"xyz-multifinance/storage"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	usecase usecase.CustomerUsecase
}

func NewCustomerHandler(uc usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		usecase: uc,
	}
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	userIDStr := c.PostForm("user_id")
	nik := c.PostForm("nik")
	fullName := c.PostForm("full_name")
	legalName := c.PostForm("legal_name")
	placeBirth := c.PostForm("place_of_birth")
	dateBirthStr := c.PostForm("date_of_birth")
	salaryStr := c.PostForm("salary")

	if nik == "" || fullName == "" || userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields: user_id, nik, or full_name"})
		return
	}

	dateBirth, err := time.Parse("2006-01-02", dateBirthStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_of_birth format. Use YYYY-MM-DD"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	salary, err := strconv.ParseInt(salaryStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid salary"})
		return
	}

	var ktpPath, selfiePath string

	if ktpFileHeader, err := c.FormFile("ktp_photo"); err == nil {
		ktpFile, err := ktpFileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open ktp_photo"})
			return
		}
		defer ktpFile.Close()

		ktpPath, err = storage.SaveImage(ktpFile, ktpFileHeader, "ktp_"+nik)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if selfieFileHeader, err := c.FormFile("selfie_photo"); err == nil {
		selfieFile, err := selfieFileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open selfie_photo"})
			return
		}
		defer selfieFile.Close()

		selfiePath, err = storage.SaveImage(selfieFile, selfieFileHeader, "selfie_"+nik)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	customer := model.Customer{
		UserID:      uint(userID),
		NIK:         nik,
		FullName:    fullName,
		LegalName:   legalName,
		PlaceBirth:  placeBirth,
		DateBirth:   dateBirth,
		Salary:      salary,
		KTPPhoto:    ktpPath,
		SelfiePhoto: selfiePath,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.usecase.CreateCustomer(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Customer created",
		"customer": customer,
	})
}

func (h *CustomerHandler) GetCustomerByNIK(c *gin.Context) {
	nik := c.Param("nik")
	customer, err := h.usecase.GetCustomerByNIK(nik)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c *gin.Context) {
	nik := c.Param("nik")
	if nik == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "NIK is required"})
		return
	}

	oldCustomer, err := h.usecase.GetCustomerByNIK(nik)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	fullName := c.PostForm("full_name")
	legalName := c.PostForm("legal_name")
	placeOfBirth := c.PostForm("place_of_birth")
	dateBirthStr := c.PostForm("date_of_birth")
	salaryStr := c.PostForm("salary")

	updatedFields := make(map[string]interface{})

	if fullName != "" {
		updatedFields["full_name"] = fullName
	}
	if legalName != "" {
		updatedFields["legal_name"] = legalName
	}
	if placeOfBirth != "" {
		updatedFields["birth_place"] = placeOfBirth
	}
	if dateBirthStr != "" {
		dateBirth, err := time.Parse("2006-01-02", dateBirthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date_of_birth format. Use YYYY-MM-DD"})
			return
		}
		updatedFields["birth_date"] = dateBirth
	}
	if salaryStr != "" {
		salary, err := strconv.ParseInt(salaryStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid salary"})
			return
		}
		updatedFields["salary"] = salary
	}

	if ktpFileHeader, err := c.FormFile("ktp_photo"); err == nil {
		ktpFile, _ := ktpFileHeader.Open()
		defer ktpFile.Close()

		if oldCustomer.KTPPhoto != "" {
			_ = os.Remove(oldCustomer.KTPPhoto)
		}

		ktpPath, err := storage.SaveImage(ktpFile, ktpFileHeader, "ktp_"+nik)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		updatedFields["photo_ktp"] = ktpPath
	}

	if selfieFileHeader, err := c.FormFile("selfie_photo"); err == nil {
		selfieFile, _ := selfieFileHeader.Open()
		defer selfieFile.Close()

		if oldCustomer.SelfiePhoto != "" {
			_ = os.Remove(oldCustomer.SelfiePhoto)
		}

		selfiePath, err := storage.SaveImage(selfieFile, selfieFileHeader, "selfie_"+nik)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		updatedFields["photo_selfie"] = selfiePath
	}

	updatedFields["updated_at"] = time.Now()

	err = h.usecase.UpdateCustomer(nik, updatedFields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Customer updated"})
}

func (h *CustomerHandler) DeleteCustomer(c *gin.Context) {
	nik := c.Param("nik")
	err := h.usecase.DeleteCustomer(nik)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Customer deleted"})
}
