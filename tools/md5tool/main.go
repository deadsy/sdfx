package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

var exts = []string{"*.stl", "*.dxf", "*.png", "*.svg"}

func main() {
	var fUpdate bool
	var fDebug bool
	flag.BoolVar(&fUpdate, "update", false, "generate new MD5SUM file")
	flag.BoolVar(&fDebug, "debug", false, "enable debug messages")
	flag.Parse()

	if fDebug {
		log.SetLevel(log.DebugLevel)
	}

	if fUpdate {
		log.Debugf("Updating MD5SUM file")
		err := updateMD5sum()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	log.Debugf("Comparing files to MD5SUM")

	var issues []string

	hashes, err := readMD5sumFile()
	if err != nil {
		// log.Fatalf(err)
		issues = append(issues, fmt.Sprintf("Error: unable to read MD5SUM file: %s\n", err))
	}

	log.Debugf("hashes: %#v", hashes)

	fileList, err := getFileList()
	if err != nil {
		log.Fatal(err)
	}

	for fn := range hashes {
		if !find(&fileList, fn) {
			issues = append(issues, fmt.Sprintf("Error: file listed in MD5SUM but file not rendered: %s\n", fn))
		}
	}

	for _, fn := range fileList {
		var hval string
		var ok bool
		if hval, ok = hashes[fn]; !ok {
			issues = append(issues, fmt.Sprintf("Error: file in local filesystem but missing from MD5SUM file: %s\n", fn))
			continue
		}
		lval, err := getMD5sum(fn)
		if err != nil {
			log.Fatal(err) // TODO fix
		}
		if hval != lval {
			issues = append(issues, fmt.Sprintf("Error: file has changed: %s\n", fn))
		}
	}

	if len(issues) > 0 {
		for _, issue := range issues {
			fmt.Printf("%s", issue)
		}
		fmt.Printf("See %s/README.md for more information.\n\n", filepath.Dir(os.Args[0]))
		os.Exit(1)
	}
	fmt.Printf("No errors\n")
}

func find(slice *[]string, val string) bool {
	for _, item := range *slice {
		if item == val {
			return true
		}
	}
	return false
}

// getFileList retrieves all files matching the global exts variable
func getFileList() ([]string, error) {
	var matches []string

	for _, ext := range exts {
		match, err := filepath.Glob(ext)
		if err != nil {
			return nil, err
		}
		for _, m := range match {
			matches = append(matches, m)
		}
	}
	return matches, nil
}

func getMD5sum(fileName string) (string, error) {
	var md5String string
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashBytes := hash.Sum(nil)[:16]
	md5String = hex.EncodeToString(hashBytes)
	return md5String, nil
}

// updateMD5sum
func updateMD5sum() error {
	files, err := getFileList()
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Matches: %+v", files)

	f, err := os.Create("MD5SUM")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, fn := range files {
		md5sum, err := getMD5sum(fn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(f, "%s\t%s\n", md5sum, fn)
	}
	return nil
}

// readMD5sumFile reads the local MD5SUM file and populates and returns
// a map of filenames and their respective hashes
func readMD5sumFile() (map[string]string, error) {
	m := make(map[string]string)

	file, err := os.Open("MD5SUM")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		m[fields[1]] = fields[0]
	}

	log.Debugf("m: %+v", m)
	return m, nil
}
