package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"io"
	"os"
)

// this is a comment
const (
	datapaq_path = "http://localhost:9998/datapaq"
)

func main() {
	log.Println("DP3 Remote calls")
	
	// GET CALLS

	//Get datapaq version
	version := getGETRequest(datapaq_path +"/version",nil)
	log.Println(version)

	// //Get all container uid's, this can be used to call the scan method
	// uids := getGETRequest(datapaq_path +"/uids",nil)
	// log.Println(uids)

	// //Get linear barcode
	// scanLinearBarcode := getGETRequest(datapaq_path +"/scanLinearBarcode",nil)
	// log.Println(scanLinearBarcode)

	// //Scan a uid as Text
	// scanparam := make(map[string]string)
	// scanparam["uid"] = "1"
	// scan := getGETRequest(datapaq_path +"/scanAsText",scanparam)
	// log.Println(scan)

	// // GET CALLS WITH IMAGE/PNG

	// // save last image with scale factor
	// scaleParam := make(map[string]string)
	// scaleParam["scaleFactor"] = "0.15"
	// getGETIMGRequest(datapaq_path +"/lastImage",scaleParam,"C:/temp/lastImg.png")

	// // POST CALLS

	// // enable barcode scanner
	// enablescanner := getPOSTRequest(datapaq_path +"/enableBarcodeScanner",nil)
	// log.Println(enablescanner)

	// // disable barcode scanner
	// disablescanner := getPOSTRequest(datapaq_path +"/disableBarcodeScanner",nil)
	// log.Println(disablescanner)

	// // save last image specify the path
	// saveImgMap := make(map[string]string)
	// saveImgMap["path"] = "C:/temp/lastImg_1.png"
	// saveImgMap["scaleFactor"] = "0.15"
	// saveLastImage := getPOSTRequest(datapaq_path +"/saveLastImage",saveImgMap)
	// log.Println(saveLastImage)

	// // shutdown webserver
	// shutdown := getPOSTRequest(datapaq_path +"/shutdown",nil)
	// log.Println(shutdown)
}

// Function for GET request
func getGETRequest(path string, parameters map[string]string ) string {
	u, _ := url.Parse(path)
	if parameters != nil && len(parameters) != 0 {
		q, _ := url.ParseQuery(u.RawQuery)
		for key, value := range parameters { // Order not specified 
			q.Add(key, value)
		}
		u.RawQuery = q.Encode()
	}

	resp, err := http.Get(u.String())

	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
    //Convert the body to type string
	sb := string(body)
	return sb
}

// Function for GET request when the return is of type IMG/PNG
func getGETIMGRequest(path string, parameters map[string]string,imgPath string ) {
	u, _ := url.Parse(path)
	if parameters != nil && len(parameters) != 0 {
		q, _ := url.ParseQuery(u.RawQuery)
		for key, value := range parameters { // Order not specified 
			q.Add(key, value)
		}
		u.RawQuery = q.Encode()
	}
	response, err := http.Get(u.String())
	defer response.Body.Close()

    //open a file for writing
    file, err := os.Create(imgPath)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // Use io.Copy to just dump the response body to the file. This supports huge files
    _, err = io.Copy(file, response.Body)
    if err != nil {
        log.Fatal(err)
    }
    log.Println("Success!")
}

// Function for POST request
func getPOSTRequest(path string, parameters map[string]string ) string {
	u, _ := url.Parse(path)
	if parameters != nil && len(parameters) != 0 {
		q, _ := url.ParseQuery(u.RawQuery)
		for key, value := range parameters { // Order not specified 
			q.Add(key, value)
		}
		u.RawQuery = q.Encode()
	}

	req, err := http.NewRequest("POST", u.String(),nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Add("Content-Type", "application/json")

	c := &http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	   log.Fatalln(err)
	}
    //Convert the body to type string
	sb := string(body)
	return sb
}