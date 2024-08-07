package handler

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	"fliqt/config"
	"fliqt/internal/model"
	"fliqt/internal/service"
)

const (
	// 5MB
	maxFileSize = 5 * 1024 * 1024
	// pdf, image, docx, doc only
	allowedFileTypes = "application/pdf,image/jpeg,image/png,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/msword"
)

type FileHandler struct {
	cfg         *config.Config
	authService service.AuthServiceInterface
	s3Service   service.S3ServiceInterface
}

func NewFileHandler(
	cfg *config.Config,
	authService service.AuthServiceInterface,
	s3Service service.S3ServiceInterface,
) *FileHandler {
	return &FileHandler{
		cfg,
		authService,
		s3Service,
	}
}

type UploadFileRequest struct {
	ContentType string `json:"content_type" binding:"required"`
	FileName    string `json:"file_name" binding:"required"`
	FileSize    int64  `json:"file_size" binding:"required"`
}

type UploadFileResponse struct {
	ObjectKey string                 `json:"object_key"`
	URL       string                 `json:"url"`
	Metadata  map[string]interface{} `json:"metadata"`
}

func (h *FileHandler) GetUploadInfo(ctx *gin.Context) {
	var req UploadFileRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.Error(err)
		return
	}

	user, err := h.authService.CurrentUser(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	if req.FileSize > maxFileSize || strings.Contains(allowedFileTypes, req.ContentType) {
		ctx.Error(ErrBadRequest)
		return
	}

	ID := xid.New().String()
	objectKey := fmt.Sprintf("%s/%s", user.ID, ID)

	URL, err := h.s3Service.PresignUpload(ctx, h.cfg.S3Bucket, user.ID, objectKey, req.ContentType, req.FileSize)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, UploadFileResponse{
		ObjectKey: objectKey,
		URL:       URL,
		Metadata:  nil,
	})
}

func (h *FileHandler) GetDownloadInfo(ctx *gin.Context) {
	objectKey := ctx.Param("object_key")

	user, err := h.authService.CurrentUser(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	if user.Role == model.RoleCandidate && !strings.HasPrefix(objectKey, fmt.Sprintf("/%s", user.ID)) {
		ctx.Error(ErrForbidden)
		return
	}

	// Verify TOTP when user is HR or Interviewer
	if user.Role != model.RoleCandidate {
		passcode, ok := ctx.GetQuery("passcode")
		if !ok {
			ctx.Error(ErrBadRequest)
			return
		}

		if err := h.authService.VerifyTOTP(ctx, user.TotpSecret, passcode); err != nil {
			ctx.Error(err)
			return
		}
	}

	URL, err := h.s3Service.GetPresignDownloadURL(ctx, h.cfg.S3Bucket, objectKey)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(200, gin.H{
		"url": URL,
	})
}
