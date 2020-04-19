package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	UploadStats StatsTotals   `json:"uploadstats"`
	VersionInfo StatsVersions `json:"versioninfo"`
}

var statsMux sync.RWMutex
var collectedStats map[string]Stats = make(map[string]Stats, 0)

func serveAggregatedStats(w http.ResponseWriter, r *http.Request) {
	var versionInfo *StatsVersions = nil
	var aggregatedTotals StatsTotals
	uploaders := make(map[string]Stats, 0)
	statsMux.RLock()
	for name, stats := range collectedStats {
		uploaders[name] = stats
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

func collectStats() {
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
	pods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{
		LabelSelector: "app=siacdn-uploader",
	})
	if err != nil {
		log.Println("Could not list pods from kubernetes")
		return
	}
	var wg sync.WaitGroup
	for _, pod := range pods.Items {
		wg.Add(1)
		go func(name, ip string) {
			defer wg.Done()
			collectOne(name, ip)
		}(pod.Name, pod.Status.PodIP)
	}
	wg.Wait()
	return
}

func collectOne(name, ip string) {
	log.Println("About to collect from", name)
	var netClient = &http.Client{
		Timeout: time.Second * 600,
	}
	resp, err := netClient.Get(fmt.Sprintf("http://%s:8080/stats", ip))
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
	statsMux.Lock()
	collectedStats[name] = stats
	statsMux.Unlock()
	log.Println("Got", stats.UploadStats.NumFiles, "files on", name)
}

func main() {
	go collectStats()
	http.ListenAndServe("0.0.0.0:8080", http.HandlerFunc(serveAggregatedStats))
}
