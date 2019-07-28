package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func initLog() {
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})
	log.Info("logger initialised")
}

type RouterStatus struct {
	Success bool `json:"success"`
	Timeout bool `json:"timeout"`
	Data    struct {
		ApEnable24G         string `json:"ap_enable_24g"`
		ApSsid24G           string `json:"ap_ssid_24g"`
		ApMac24G            string `json:"ap_mac_24g"`
		ApSignal24G         string `json:"ap_signal_24g"`
		ApSecurity24G       string `json:"ap_security_24g"`
		ApChannel24G        int    `json:"ap_channel_24g"`
		Status              int    `json:"status"`
		ApSecurityStatus24G string `json:"ap_security_status_24g"`
		ApEnable5G          string `json:"ap_enable_5g"`
		ApSsid5G            string `json:"ap_ssid_5g"`
		ApMac5G             string `json:"ap_mac_5g"`
		ApSignal5G          string `json:"ap_signal_5g"`
		ApSecurity5G        string `json:"ap_security_5g"`
		ApChannel5G         int    `json:"ap_channel_5g"`
		ApSecurityStatus5G  string `json:"ap_security_status_5g"`
		Show2GFlag          string `json:"show2gFlag"`
		Show5GFlag          string `json:"show5gFlag"`
	} `json:"data"`
}

type Clients struct {
	Success bool `json:"success"`
	Timeout bool `json:"timeout"`
	Data    []struct {
		Mac      string `json:"mac"`
		Name     string `json:"name"`
		IP       string `json:"ip"`
		Ipaddr   string `json:"ipaddr"`
		Type     string `json:"type"`
		ConnType string `json:"conn_type"`
	} `json:"data"`
}

var clientList Clients
var routerStatus RouterStatus

var (
	ClientCountGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "client_count",
		Help: "Number of clients connected to tplink r450",
	})
)

var (
	FiveSignelGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "fiveg_signal_strength",
		Help: "5G signal strenth of upstream wifi",
	})
)

var pReg = prometheus.NewRegistry()

func register() {
	log.Info("Register Stats")
	pReg.MustRegister(FiveSignelGauge)
	pReg.MustRegister(ClientCountGauge)
	log.Info("Stats Registered")
}

func generateRouterStatus() {
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	body := strings.NewReader(`operation=read`)
	req, err := http.NewRequest("POST", "https://" + os.Getenv("TPIP") + "/data/router.status.json", body)
	if err != nil {
		log.Info("Failed to get router status")
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:66.0) Gecko/20100101 Firefox/66.0")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-GB,en;q=0.5")
	req.Header.Set("Referer", "https://" + os.Getenv("TPIP") + "/index")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "COOKIE=" + os.Getenv("ACCESSCOOKIE"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Info("Failed to get router status")
	}
	defer resp.Body.Close()
	log.Info(resp.Body)

	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	json.Unmarshal(ret, &routerStatus)

	log.Info("Router Status Success: ", routerStatus.Success)
	sigStrength, err := strconv.ParseFloat(routerStatus.Data.ApSignal5G, 64)
	if err != nil {
		log.Info("Error converting string to float")
	}
	log.Info("Router 5G Strength: ", sigStrength)
	FiveSignelGauge.Set(sigStrength)

}

func generateClientStatus() {
	// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

	body := strings.NewReader(`operation=read`)
	req, err := http.NewRequest("POST", "https://" + os.Getenv("TPIP") + "/data/device.all.json", body)
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Language", "en-GB,en;q=0.5")
	req.Header.Set("Referer", "https://" + os.Getenv("TPIP") + "/index")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Cookie", "COOKIE=" + os.Getenv("ACCESSCOOKIE"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()

	//clientList := clients{}
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	}
	json.Unmarshal(ret, &clientList)

	log.Info("Client List:", clientList.Data)
	log.Info("Client Count:", (len(clientList.Data)))
	log.Info("Client Success:", clientList.Success)
	ClientCountGauge.Set(float64(len(clientList.Data)))
}

func serve() {
	log.Info("starting http hander")
	handler := promhttp.HandlerFor(pReg, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	http.ListenAndServe(":8089", nil)
	log.Info("Serving on port 8089")
}

func init() {
	initLog()
	register()
}

func main() {
	go func() {
		serve()
	}()
	for {
		generateClientStatus()
		generateRouterStatus()
		time.Sleep(10 * time.Second)
	}
}
