package helper

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryUploader struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryUploader() (*CloudinaryUploader, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return nil, err
	}

	return &CloudinaryUploader{cld: cld}, nil
}

func (c *CloudinaryUploader) UploadResizedImage(
	ctx context.Context,
	filePath string,
	publicID string,
) (string, error) {

	result, err := c.cld.Upload.Upload(ctx, filePath, uploader.UploadParams{
		PublicID: publicID,
		Folder:   "missing-persons",
		Transformation: "c_fill,w_512,h_512,q_auto,f_auto",
	})

	if err != nil {
		return "", err
	}

	// 🔑 return URL
	return result.SecureURL, nil
}
