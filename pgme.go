package main

import (
    "bytes"
    "context"
    "encoding/csv"
    "fmt"
    "html/template"
    "net/http"
    "log"
    "os"
    "os/exec"
    "strings"
	"syscall"
	"os/signal"
	"path"
)


type PageVariables struct {
	PageTitle        string
	Metrics			 []string
	VersionInfo      map[string]string
}


var (
	// BuildTime is a time label of the moment when the binary was built
	BuildTime = "unset"
	// Commit is a last commit hash at the moment when the binary was built
	Commit = "unset"
	// Release is a semantic version of current build
	Release = "unset"
)

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}

	return value
}


// healthz is a liveness probe.
func healthz(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}


// name, index, temperature.gpu, utilization.gpu,
// utilization.memory, memory.total, memory.free, memory.used
func home(w http.ResponseWriter, r *http.Request) {

	metricList := []string {
		"temperature.gpu", "utilization.gpu",
		"utilization.memory", "memory.total", "memory.free", "memory.used"}

	verInfo := make(map[string]string)
	verInfo["Buildtime"] = BuildTime
	verInfo["Commit"] = Commit
	verInfo["Release"] = Release


	pv := PageVariables{
		PageTitle:   "Prometheus nVidia GPU Metrics Exporter",
		Metrics:     metricList,
		VersionInfo: verInfo,
	}

	filepath := path.Join(path.Dir("./template/home.html"), "home.html")
	template.ParseFiles()
	t, err := template.ParseFiles(filepath)
	if err != nil {
		log.Print("Template parsing error: ", err)
	}

	err = t.Execute(w, pv)
	if err != nil {
		log.Print("Template execution error: ", err)
	}

}


func metrics(response http.ResponseWriter, request *http.Request) {
    out, err := exec.Command(
        "nvidia-smi",
        "--query-gpu=name,index,temperature.gpu,utilization.gpu,utilization.memory,memory.total,memory.free,memory.used",
        "--format=csv,noheader,nounits").Output()

    if err != nil {
        log.Printf("ERROR: %s\n", err)
        return
    }

    csvReader := csv.NewReader(bytes.NewReader(out))
    csvReader.TrimLeadingSpace = true
    records, err := csvReader.ReadAll()

    if err != nil {
        log.Printf("%s\n", err)
        return
    }

    metricList := []string {
        "temperature.gpu", "utilization.gpu",
        "utilization.memory", "memory.total", "memory.free", "memory.used"}

    result := ""
    for _, row := range records {
        name := fmt.Sprintf("%s[%s]", row[0], row[1])
        for idx, value := range row[2:] {
            result = fmt.Sprintf("%s%s{gpu=\"%s\"} %s\n", result, metricList[idx], name, value)
        }
    }

    fmt.Fprintf(response, strings.Replace(result, ".", "_", -1))
}


func main() {
    log.Print("Starting the service...")
    port := getEnv("PORT", "9101");
    addr := ":"+port

    log.Print("- PORT set to "+ port +".  If  environment variable PORT is not set the default is 9101")


	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)


	srv := &http.Server{
		Addr: addr,
	}

	go func() {
		http.HandleFunc("/", home)
		http.HandleFunc("/healthz", healthz)
		http.HandleFunc("/metrics/", metrics)
		err := srv.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}

	}()

    log.Print("The service is listening on ", port)

	killSignal := <-interrupt
	switch killSignal {
	case os.Interrupt:
		log.Print("Got SIGINT...")
	case syscall.SIGTERM:
		log.Print("Got SIGTERM...")
	}

	log.Print("The service is shutting down...")
	srv.Shutdown(context.Background())
	log.Print("Done")


}
