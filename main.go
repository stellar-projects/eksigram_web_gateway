package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.POST("/", routeWebPage)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		log.Println("Port env value not found, setting to defult")
		port = "5050"
	}
	e.Logger.Fatal(e.Start(strings.Join([]string{":", port}, "")))
}

func routeWebPage(c echo.Context) error {
	urlReq := new(ReqWebPage)
	err := c.Bind(urlReq)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("GET", urlReq.URL, nil)
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Safari/605.1.15")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	//req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		log.Println("Error get web: ", resp.Status)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("Error in defer read body: ", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	for key := range resp.Header {
		val := resp.Header[key]
		if len(val) > 0 {
			c.Response().Header().Add(key, val[0])
		}
	}
	return c.JSON(http.StatusOK, string(body))
}

type ReqWebPage struct {
	URL string `json:"url"`
}

//func startServer() {
//
//	e := echo.New()
//	e.IPExtractor = echo.ExtractIPFromRealIPHeader()
//	e.Use(middleware.CORS())
//	e.Use(middleware.Logger())
//	e.Use(middlewares.I18n())
//
//	if config.GetEnvironment() == config.RunEnvironmentDebug {
//		p := prometheus.NewPrometheus("echo", nil)
//		p.Use(e)
//	}
//
//	e.HTTPErrorHandler = common.CustomHttpErrorHandler
//	controllers.InitRoutes(e)
//	application.GetApplication().GetJwtBlackListHandler().InitializeJwtBlackList()
//	port := os.Getenv("API_PORT")
//	if port == "" {
//		port = "9090"
//	}
//	application.GetApplication().ApiServerSettings = &application.ApiServerSettings{
//		Port:    port,
//		BaseUrl: "https://overtune-api.com",
//	}
//
//	_ = media_uploader.NewContentAnalyzer().ProcessResult(3212)
//
//	enableHttps := os.Getenv("ENABLE_HTTPS") == "true"
//	if enableHttps {
//		e.Logger.Fatal(e.StartTLS(strings.Join([]string{":", port}, ""), "overtune.https.crt", "overtune.https.key"))
//	} else {
//		e.Logger.Fatal(e.Start(strings.Join([]string{":", port}, "")))
//	}
//}
