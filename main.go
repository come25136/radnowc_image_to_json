package main

import (
	"encoding/json"
	"image/color"
	"image/png"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go-playground/validator/v10"
)

// GetParams はjmaに投げるパラメーターでもある
type GetParams struct {
	Date int `validate:"required,min=202001010000"` // 過去データ参照する気無いので簡易
	X    int `validate:"required,min=19,max=44"`    // 数値範囲外=日本の観測範囲外
	Y    int `validate:"required,min=18,max=46"`    // 数値範囲外=日本の観測範囲外
	Z    int `validate:"required,len=6"`            // 1から6まである、XYを6基準で指定してあるので6固定
}

func radnowcHandler(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
	}

	date, _ := strconv.Atoi(r.QueryStringParameters["date"])
	x, _ := strconv.Atoi(r.QueryStringParameters["x"])
	y, _ := strconv.Atoi(r.QueryStringParameters["y"])
	z, _ := strconv.Atoi(r.QueryStringParameters["z"])

	params := &GetParams{
		Date: date,
		X:    x,
		Y:    y,
		Z:    z,
	}

	validate := validator.New()
	err := validate.Struct(params)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	url := "https://www.jma.go.jp/jp/realtimerad/highresorad_tile/HRKSNC/" + strconv.Itoa(date) + "/" + strconv.Itoa(date) + "/zoom" + strconv.Itoa(z) + "/" + strconv.Itoa(x) + "_" + strconv.Itoa(y) + ".png"

	radnowcRes, err := http.Get(url)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}
	defer radnowcRes.Body.Close()

	if radnowcRes.StatusCode == http.StatusNotFound {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: http.StatusNotFound,
		}, nil
	}

	img, err := png.Decode(radnowcRes.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	var radnowcData [256][256]int

	// 画像タイルが256x256
	for pY := 0; pY < 256; pY++ {
		for pX := 0; pX < 256; pX++ {
			color := color.RGBAModel.Convert(img.At(pX, pY)).(color.RGBA)

			value := 0

			if color.A == 255 {
				if color.R == 242 && color.G == 242 && color.B == 255 {
					value = 1
				} else if color.R == 160 && color.G == 210 && color.B == 255 {
					value = 5
				} else if color.R == 33 && color.G == 140 && color.B == 255 {
					value = 10
				} else if color.R == 0 && color.G == 65 && color.B == 255 {
					value = 20
				} else if color.R == 250 && color.G == 245 && color.B == 0 {
					value = 30
				} else if color.R == 255 && color.G == 153 && color.B == 0 {
					value = 50
				} else if color.R == 255 && color.G == 40 && color.B == 0 {
					value = 80
				} else if color.R == 180 && color.G == 0 && color.B == 104 {
					value = 100
				}
			}

			radnowcData[pY][pX] = value
		}
	}

	res, err := json.Marshal(radnowcData[:])

	if err != nil {
		return events.APIGatewayProxyResponse{
			Headers:    headers,
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Headers:    headers,
		StatusCode: http.StatusOK,
		Body:       string(res),
	}, nil
}

func main() {
	lambda.Start(radnowcHandler)
}
