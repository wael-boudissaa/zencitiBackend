package utils

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"log"
)

func Credentials() (*cloudinary.Cloudinary, context.Context) {
	// Use your actual Cloudinary URL
	cld, err := cloudinary.NewFromURL("cloudinary://526228195934882:rNFR_a5RYnnaUncC3FcEO4ZogVU@dlsrfos8m")
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

	// Force secure URL (HTTPS)
	cld.Config.URL.Secure = true

	ctx := context.Background()
	return cld, ctx
}

func UploadImage(cld *cloudinary.Cloudinary, ctx context.Context,image string) {
	// Upload an image from a remote URL
	resp, err := cld.Upload.Upload(ctx, image,
		uploader.UploadParams{
			PublicID:       "quickstart_butterfly",
			UniqueFilename: api.Bool(false),
			Overwrite:      api.Bool(true),
		})
	if err != nil {
		log.Printf("Upload failed: %v\n", err)
		return
	}

	fmt.Println("**** 1. Uploaded Image ****\nDelivery URL:", resp.SecureURL, "\n")
}

func GetAssetInfo(cld *cloudinary.Cloudinary, ctx context.Context) {
	// Get image asset info
	resp, err := cld.Admin.Asset(ctx, admin.AssetParams{PublicID: "quickstart_butterfly"})
	if err != nil {
		log.Printf("Failed to get asset info: %v\n", err)
		return
	}

	fmt.Println("**** 2. Asset Details ****\n", resp, "\n")

	var tags []string
	if resp.Width > 900 {
		tags = []string{"large"}
	} else {
		tags = []string{"small"}
	}

	updateResp, err := cld.Admin.UpdateAsset(ctx, admin.UpdateAssetParams{
		PublicID: "quickstart_butterfly",
		Tags:     tags,
	})
	if err != nil {
		log.Printf("Failed to update asset tags: %v\n", err)
		return
	}

	fmt.Println("**** 3. Updated Tags ****\nTags:", updateResp.Tags, "\n")
}

func TransformImage(cld *cloudinary.Cloudinary, ctx context.Context) {
	// Create transformation
	qsImg, err := cld.Image("quickstart_butterfly")
	if err != nil {
		log.Printf("Failed to create image object: %v\n", err)
		return
	}

	qsImg.Transformation = "r_max/e_sepia"

	transformedURL, err := qsImg.String()
	if err != nil {
		log.Printf("Failed to get transformed image URL: %v\n", err)
		return
	}

	fmt.Println("**** 4. Transformed Image ****\nTransformation URL:", transformedURL, "\n")
}
