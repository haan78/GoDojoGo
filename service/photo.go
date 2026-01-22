package service

import (
	"GoDojoGo/data"
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v5"
)

const (
	formField      = "file"
	maxUploadBytes = 10 << 20
	requestTimeout = 10 * time.Second
)

func SaveUserPhoto(c *echo.Context) error {

	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	userID, err := strconv.Atoi(c.Param("userId"))
	if err == nil && userID > 0 {
		fh, err := c.FormFile(formField)
		if err == nil {
			if fh.Header != nil {
				ct := fh.Header.Get("Content-Type")
				if ct != "" && strings.HasPrefix(ct, "image/") {
					f, err := fh.Open()
					if err == nil {
						defer f.Close()
						lr := io.LimitReader(f, maxUploadBytes+1)
						bArray, err := io.ReadAll(lr)
						if err == nil {
							if int64(len(bArray)) <= maxUploadBytes {
								err := data.SavePhotoData(ctx, int64(userID), bArray)
								if err == nil {
									return nil
								} else {
									return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
								}
							} else {
								return echo.NewHTTPError(http.StatusBadRequest, "maximum photo size is 10Mb")
							}
						} else {
							return echo.NewHTTPError(http.StatusBadRequest, err.Error())
						}
					} else {
						return echo.NewHTTPError(http.StatusBadRequest, err.Error())
					}
				} else {
					return echo.NewHTTPError(http.StatusBadRequest, "uploaded file must be an image")
				}
			} else {
				return echo.NewHTTPError(http.StatusBadRequest, "no file header")
			}
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, `missing form file field "file"`)
		}
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid userId")
	}
}

func GetUserPhoto(c *echo.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	userID, err := strconv.Atoi(c.Param("userId"))
	if err == nil && userID > 0 {
		img, err := data.GetUserPhoto(ctx, int64(userID))
		if err == nil {
			ct := http.DetectContentType(img)
			c.Blob(http.StatusOK, ct, img)
			return nil
		} else {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	} else {
		return err
	}
}
