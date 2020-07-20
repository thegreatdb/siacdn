package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// StatsTotals represents stats totals for a single Sia node
type StatsTotals struct {
	NumFiles  int   `json:"numfiles"`
	TotalSize int64 `json:"totalsize"`
}

// StatsVersions represents stats versions for a single Sia node
type StatsVersions struct {
	Version     string `json:"version"`
	GitRevision string `json:"gitrevision"`
}

// Stats represents overall stats about a single Sia node
type Stats struct {
	UploadStats      StatsTotals   `json:"uploadstats"`
	VersionInfo      StatsVersions `json:"versioninfo"`
	PerformanceStats interface{}   `json:"performancestats,omitempty"`
	Uploaders        interface{}   `json:"uploaders,omitempty"`
	Viewers          interface{}   `json:"viewers,omitempty"`
}

var statsMux sync.RWMutex
var collectedUploadStats map[string]Stats = make(map[string]Stats, 0)
var collectedViewStats map[string]Stats = make(map[string]Stats, 0)

// Alert represents a single alert from a single node
type Alert struct {
	Cause    string `json:"cause"`
	Msg      string `json:"msg"`
	Module   string `json:"module"`
	Severity string `json:"severity"`
}

// Alerts represents overall alerts about a single Sia node
type Alerts struct {
	Alerts         []Alert     `json:"alerts"`
	CriticalAlerts []Alert     `json:"criticalalerts"`
	ErrorAlerts    []Alert     `json:"erroralerts"`
	WarningAlerts  []Alert     `json:"warningalerts"`
	Uploaders      interface{} `json:"uploaders,omitempty"`
	Viewers        interface{} `json:"viewers,omitempty"`
}

// Count returns the total count of all alerts for this node
func (a *Alerts) Count() int {
	return len(a.Alerts) + len(a.CriticalAlerts) + len(a.ErrorAlerts) + len(a.WarningAlerts)
}

var alertsMux sync.RWMutex
var collectedUploadAlerts map[string]Alerts = make(map[string]Alerts, 0)
var collectedViewAlerts map[string]Alerts = make(map[string]Alerts, 0)

func serveAggregatedStats(w http.ResponseWriter, r *http.Request) {
	var versionInfo *StatsVersions = nil
	var aggregatedTotals StatsTotals
	uploaders := make(map[string]Stats, 0)
	viewers := make(map[string]Stats, 0)
	statsMux.RLock()
	for name, stats := range collectedUploadStats {
		uploaders[name] = stats
		if versionInfo == nil {
			versionInfo = &stats.VersionInfo
		}
		aggregatedTotals.NumFiles += stats.UploadStats.NumFiles
		aggregatedTotals.TotalSize += stats.UploadStats.TotalSize
	}
	for name, stats := range collectedViewStats {
		viewers[name] = stats
		if versionInfo == nil {
			versionInfo = &stats.VersionInfo
		}
		aggregatedTotals.NumFiles += stats.UploadStats.NumFiles
		aggregatedTotals.TotalSize += stats.UploadStats.TotalSize
	}
	statsMux.RUnlock()
	if versionInfo == nil {
		msg := fmt.Sprintf("Requested stats before it was collected")
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	resp := map[string]interface{}{
		"uploaders":   uploaders,
		"viewers":     viewers,
		"uploadstats": aggregatedTotals,
		"versioninfo": versionInfo,
	}
	encoded, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		msg := fmt.Sprintf("Could not encode JSON: %w", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
}

func serveAggregatedAlerts(w http.ResponseWriter, r *http.Request) {
	uploaders := make(map[string]Alerts, 0)
	viewers := make(map[string]Alerts, 0)
	alertsMux.RLock()
	count := 0
	for name, alerts := range collectedUploadAlerts {
		uploaders[name] = alerts
		count += alerts.Count()
	}
	for name, alerts := range collectedViewAlerts {
		viewers[name] = alerts
		count += alerts.Count()
	}
	alertsMux.RUnlock()
	resp := map[string]interface{}{
		"uploaders": uploaders,
		"viewers":   viewers,
		"count":     count,
	}
	encoded, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		msg := fmt.Sprintf("Could not encode JSON: %w", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(encoded)
}

func collectLoop() {
	first := true
	for {
		if first {
			first = false
		} else {
			time.Sleep(30 * time.Second)
		}
		collectAll()
	}
}

func collectAll() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Println("Could not configue kubernetes client")
		return
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Println("Could not set up kubernetes client from configuration")
		return
	}

	uploaderPods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{
		LabelSelector: "app=siacdn-uploader",
	})
	if err != nil {
		log.Println("Could not list uploader pods from kubernetes")
		return
	}

	viewerPods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{
		LabelSelector: "app=siacdn-viewer",
	})
	if err != nil {
		log.Println("Could not list uploader pods from kubernetes")
		return
	}

	var wg sync.WaitGroup
	for _, pod := range uploaderPods.Items {
		wg.Add(1)
		go func(name, ip string) {
			defer wg.Done()
			collectOneStats(name, ip, "stats", collectedUploadStats)
			collectOneAlerts(name, ip, "alerts", collectedUploadAlerts)
		}(sanitizePodName(pod.Name), pod.Status.PodIP)
	}
	for _, pod := range viewerPods.Items {
		wg.Add(1)
		go func(name, ip string) {
			defer wg.Done()
			collectOneStats(name, ip, "statsdown", collectedViewStats)
			collectOneAlerts(name, ip, "alertsdown", collectedViewAlerts)
		}(sanitizePodName(pod.Name), pod.Status.PodIP)
	}
	wg.Wait()
	return
}

func collectOneStats(name, ip, pathPart string, statMap map[string]Stats) {
	log.Println("About to collect from", name)
	var netClient = &http.Client{
		Timeout: time.Second * 600,
	}
	resp, err := netClient.Get(fmt.Sprintf("http://%s:8080/"+pathPart, ip))
	if err != nil {
		log.Println("Could not collect stats from", name, err)
		//statsMux.Lock()
		//delete(collectedStats, name)
		//statsMux.Unlock()
		return
	}
	dec := json.NewDecoder(resp.Body)
	//dec.DisallowUnknownFields()
	var stats Stats
	if err = dec.Decode(&stats); err != nil {
		log.Println("Could not decode stats from", name, err)
		//statsMux.Lock()
		//delete(collectedStats, name)
		//statsMux.Unlock()
		return
	}
	if stats.Uploaders != nil {
		log.Println("Somehow got global stats for pod, bailing...")
		return
	}
	statsMux.Lock()
	statMap[name] = stats
	statsMux.Unlock()
	log.Println("Got", stats.UploadStats.NumFiles, "files on", name)
}

func collectOneAlerts(name, ip, pathPart string, alertsMap map[string]Alerts) {
	log.Println("About to collect from", name)
	var netClient = &http.Client{
		Timeout: time.Second * 600,
	}
	resp, err := netClient.Get(fmt.Sprintf("http://%s:8080/"+pathPart, ip))
	if err != nil {
		log.Println("Could not collect alerts from", name, err)
		//alertsMux.Lock()
		//delete(collectedAlerts, name)
		//alertsMux.Unlock()
		return
	}
	dec := json.NewDecoder(resp.Body)
	//dec.DisallowUnknownFields()
	var alerts Alerts
	if err = dec.Decode(&alerts); err != nil {
		log.Println("Could not decode alerts from", name, err)
		//alertsMux.Lock()
		//delete(collectedAlerts, name)
		//alertsMux.Unlock()
		return
	}
	if alerts.Uploaders != nil {
		log.Println("Somehow got global alerts for pod, bailing...")
		return
	}
	alertsMux.Lock()
	alertsMap[name] = alerts
	alertsMux.Unlock()
	log.Println("Got", alerts.Count(), "alerts on", name)
}

func sanitizePodName(name string) string {
	splitStr := strings.Split(name, "-")
	if len(splitStr) >= 5 {
		return strings.Join(splitStr[:len(splitStr)-2], "-")
	}
	return name
}

func main() {
	go collectLoop()
	mux := http.NewServeMux()
	mux.HandleFunc("/stats", http.HandlerFunc(serveAggregatedStats))
	mux.HandleFunc("/alerts", http.HandlerFunc(serveAggregatedAlerts))
	http.ListenAndServe("0.0.0.0:8080", mux)
}
