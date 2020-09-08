package updater

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/micro/cli/v2"
	goruntime "github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/runtime"
)

var (
	// repository to detect changes from, e.g. micro/micro
	repository string
	// reference to detect changes from, e.g. latest
	reference string
	// latestCommit is the SHA of the latest commit
	latestCommit string

	// client to use for http requests
	client = new(http.Client)
	// updateFrequency is the time interval at which GitHub will be polled to check for changes
	updateFrequency = time.Minute * 5
	// mux is used to make the application thread safe, however the update frequence should be high
	// enough so that multiple gorountines aren't running at once
	mux = new(sync.Mutex)

	// Flags specific to the updater service
	Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "repository",
			Usage:   "Set the repository, e.g. micro/micro",
			EnvVars: []string{"MICRO_UPDATER_REPOSITORY"},
			Value:   "micro/micro",
		},
		&cli.StringFlag{
			Name:    "reference",
			Usage:   "Set the reference, e.g. latest",
			EnvVars: []string{"MICRO_UPDATER_REFERENCE"},
			Value:   "master",
		},
		&cli.StringFlag{
			Name:    "latest_commit",
			Usage:   "Set the latest commit SHA",
			EnvVars: []string{"MICRO_UPDATER_LATEST_COMMIT"},
			Value:   "",
		},
	}
)

// Run the updater service
func Run(cli *cli.Context) error {
	// create the service
	srv := service.New(
		service.Name("updater"),
		service.Version("latest"),
	)

	// load the configuration
	repository = cli.String("repository")
	reference = cli.String("reference")
	latestCommit = cli.String("latest_commit")
	fmt.Printf("Updater setup for %v:%v. Latest commit: '%v'\n", repository, reference, latestCommit)

	// updates periodically async
	t := time.NewTicker(updateFrequency)
	go func() {
		for {
			fmt.Println("Checking for updates")
			checkForUpdates()
			<-t.C
		}
	}()

	// run the service
	return srv.Run()
}

// checkForUpdates loads the latest commit and will update any services which have been changed
// since we last checked
func checkForUpdates() {
	mux.Lock()
	defer mux.Unlock()

	commit, err := getLatestCommit()
	if err != nil {
		fmt.Printf("Error getting latest commit: %v\n", err)
		// logger.Errorf("Error getting latest commit: %v", err)
		return
	}

	// this is the first time we loaded the commit, don't restart any services
	if len(latestCommit) == 0 {
		latestCommit = commit
		fmt.Printf("Latest commit has been initialized as %v\n", latestCommit)
		// logger.Debugf("Latest commit has been initialized as %v", latestCommit)
		return
	}

	// commit hasn't changed since last time we checked
	if latestCommit == commit {
		fmt.Printf("Latest commit is still %v\n", latestCommit)
		// logger.Debugf("Latest commit is still %v", latestCommit)
		return
	}

	// determine which files have changed since the service last changed
	files, err := getFilesChanged(latestCommit, commit)
	if err != nil {
		fmt.Printf("Error loading files changed since last commit: %v\n", err)
		// logger.Errorf("Error loading files changed since last commit: %v", err)
		return
	}

	// updateAll is a boolean indicating if all the serivces need to be updated, this would be the
	// case if a file impacting multiple services is amended, e.g. "go.mod".
	var updateAll bool

	// serviceNames is a map containing all the names of the services. Services reside at services/[name].
	// We are using a map to prevent duplicate values.
	var serviceNames map[string]bool

	for _, f := range files {
		// add the service name, e.g. "runtime" if the file is within a service/[name] directory, e.g.
		// service/runtime/server/server.go. If ths service does not match this pattern, the file could
		// apply to any service so we want to update them all
		if comps := strings.Split(f, "/"); len(comps) > 2 && comps[0] == "service" {
			serviceNames[string(comps[0])] = true
		} else {
			updateAll = true
			break
		}
	}

	// update all the services and then exit
	if updateAll {
		fmt.Printf("Updating all services\n")
		// logger.Debugf("Updating all services")

		srvs, err := runtime.Read(goruntime.ReadNamespace("default"))
		if err != nil {
			fmt.Printf("Error reading services from runtime: %v\n", err)
			// logger.Errorf("Error reading services from runtime: %v", err)
			return
		}
		for _, srv := range srvs {
			if len(srv.Name) == 0 || srv.Name == "updater" {
				fmt.Printf("Skipping service '%v'\n", srv.Name)
				continue
			}

			fmt.Printf("Updating service %v\n", srv.Name)
			// logger.Debugf("Updating service %v", srv.Name)

			if err := runtime.Update(srv); err != nil {
				fmt.Printf("Error updating %v service: %v\n", srv.Name, err)
				// logger.Errorf("Error updating %v service: %v", srv.Name, err)
			}
		}

		latestCommit = commit
		return
	}

	// update all the services which had a file changed
	for name := range serviceNames {
		srvs, err := runtime.Read(goruntime.ReadService(name), goruntime.ReadNamespace("default"))
		if err != nil {
			fmt.Printf("Error reading service: %v\n", err)
			// logger.Errorf("Error reading service: %v", err)
			continue
		} else if len(srvs) == 0 {
			fmt.Printf("Service %v not found\n", name)
			// logger.Debugf("Service %v not found", name)
			continue
		}

		fmt.Printf("Updating service %v\n", name)
		// logger.Debugf("Updating service %v", name)
		if err := runtime.Update(srvs[0]); err != nil {
			fmt.Printf("Error updating %v service: %v\n", srvs[0].Name, err)
			// logger.Errorf("Error updating %v service: %v", srvs[0].Name, err)
		}
	}
}

// getLatestCommit returns the latest commit SHA
func getLatestCommit() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%v/commits/%v", repository, reference)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	bytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	if rsp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error getting commits %v. Body: %v", rsp.Status, string(bytes))
	}
	rsp.Body.Close()

	var data struct {
		SHA string `json:"sha"`
	}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return "", err
	}

	return data.SHA, nil
}

// getFilesChanged returns the names of the files which have been changed between two commits
func getFilesChanged(shaOne, shaTwo string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%v/compare/%v...%v", repository, shaOne, shaTwo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	rsp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	rsp.Body.Close()

	var data struct {
		Files []struct {
			Filename string `json:"filename"`
		} `json:"files"`
	}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	filenames := make([]string, len(data.Files))
	for i, f := range data.Files {
		filenames[i] = f.Filename
	}
	return filenames, nil
}
