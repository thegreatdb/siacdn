package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

var collectedStats map[string]Stats = make(map[string]Stats, 0)

func serveAggregatedStats(w http.ResponseWriter, r *http.Request) {
	var versionInfo *StatsVersions = nil
	var aggregatedTotals StatsTotals
	uploaders := make(map[string]Stats, 0)
	for name, stats := range collectedStats {
		uploaders[name] = stats
		if versionInfo == nil {
			versionInfo = &stats.VersionInfo
		}
		aggregatedTotals.NumFiles += stats.UploadStats.NumFiles
		aggregatedTotals.TotalSize += stats.UploadStats.TotalSize
	}
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
		if err := collectStatsRun(); err != nil {
			log.Println("Got error collecting stats", err)
		}
	}
}

func collectStatsRun() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	pods, err := clientset.CoreV1().Pods("default").List(metav1.ListOptions{
		LabelSelector: "app=siacdn-uploader",
	})
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		log.Println("About to collect from", pod.Name)
		resp, err := http.Get(fmt.Sprintf("http://%s:8080/stats", pod.Status.PodIP))
		if err != nil {
			return err
		}
		dec := json.NewDecoder(resp.Body)
		//dec.DisallowUnknownFields()
		var stats Stats
		if err = dec.Decode(&stats); err != nil {
			return err
		}
		collectedStats[pod.Name] = stats
		log.Println("Got", stats.UploadStats.NumFiles, "files on", pod.Name)
	}
	return nil
}

func main() {
	go collectStats()
	http.ListenAndServe("0.0.0.0:8080", http.HandlerFunc(serveAggregatedStats))
}
