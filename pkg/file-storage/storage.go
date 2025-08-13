package file_storage

import (
	"fmt"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var storageBasePath = "./storage"

func Store(ctx *fiber.Ctx, file *multipart.FileHeader, dir string, isPublic bool) (path string, err error) {
	if isPublic {
		storageBasePath += "/public/"
	} else {
		storageBasePath += "/private/"
	}

	fileName := fmt.Sprintf("%s-%s", strconv.Itoa(int(time.Now().Unix())), file.Filename)

	storagePath := fmt.Sprintf("%s%s", storageBasePath, dir)

	if _, err = os.Stat(storagePath); os.IsNotExist(err) {
		if err = os.MkdirAll(storagePath, os.ModePerm); err != nil {
			return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	fullPath := fmt.Sprintf("%s/%s", storagePath, fileName)

	if err = ctx.SaveFile(file, fullPath); err != nil {
		return "", fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return fmt.Sprintf("%s/%s", dir, fileName), nil
}

func Remove(path string, isPublic bool) (err error) {
	if isPublic {
		storageBasePath += "/public/"
	} else {
		storageBasePath += "/private/"
	}

	filePath := fmt.Sprintf("%s%s", storageBasePath, path)

	if err = os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fiber.NewError(fiber.StatusNotFound, "File not found")
		}

		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return nil
}

func GetStorageURL(ctx *fiber.Ctx) string {
	return strings.Join([]string{ctx.BaseURL(), "storage"}, "/")
}
