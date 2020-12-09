package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"io"
	"os"
	"encoding/json"
)

// this is a comment
const (
	datapaq_path = "http://localhost:9998/datapaq"
)


func main() {
	log.Println("DP3 Remote calls")
	
	// GET CALLS
	getDatapaqVersion();
	getDatapaqStatus();
	//Get all container uid's, this can be used to call the scan method
	getContainerUids()
	//Get linear barcode
	getLinearBarcode()

	getContainerScanAsJson()

	// GET CALLS WITH IMAGE/PNG

	// // save last image with scale factor
	// scaleParam := make(map[string]string)
	// scaleParam["scaleFactor"] = "0.15"
	// getGETIMGRequest(datapaq_path +"/lastImage",scaleParam,"C:/temp/lastImg.png")

	// POST CALLS

	enableBarcodeScanner()

	disableBarcodeScanner()

	getLastSavedImage()

	shutDownDatapaqServer()
}

//Get datapaq version
func getDatapaqVersion(){
	var result map[string]interface{}
	version := getGETRequest(datapaq_path +"/version",nil)
	json.Unmarshal([]byte(version), &result)
	log.Println(result["version"])
}

//Get datapaq status
func getDatapaqStatus(){
	var result map[string]interface{}
	status := getGETRequest(datapaq_path +"/status",nil)
	json.Unmarshal([]byte(status), &result)
	log.Println(result["status"])
}

//Get all container uid's, this can be used to call the scan method
func getContainerUids(){
	var result []map[string]string
	getContainerUids := getGETRequest(datapaq_path +"/uids",nil)
	json.Unmarshal([]byte(getContainerUids), &result)
	for _, m := range result {
        for k, v := range m {
			log.Println(k +" = "+ v)
        }
    }
}

//Get linear barcode
func getLinearBarcode(){
	var result map[string]interface{}
	scanLinearBarcode := getGETRequest(datapaq_path +"/scanLinearBarcode",nil)
	json.Unmarshal([]byte(scanLinearBarcode), &result)
	if(result["error"] !=nil){
		log.Println(result["message"])	
	}else{
		log.Println(result["LinearBarcode"])
	}
}


//Scan a uid as JSON
func getContainerScanAsJson(){
	var result map[string]interface{}
	scanparam := make(map[string]string)
	scanparam["uid"] = "1"
	scan := getGETRequest(datapaq_path +"/scanAsJson",scanparam)
	json.Unmarshal([]byte(scan), &result)
	log.Println(result)
}


//enable barcode scanner
func enableBarcodeScanner(){
	var result map[string]interface{}
	enablescanner := getPOSTRequest(datapaq_path +"/enableBarcodeScanner",nil)
	json.Unmarshal([]byte(enablescanner), &result)
	log.Println(result["enableBarcodeScanner"])
}

//disable barcode scanner
func disableBarcodeScanner(){
	var result map[string]interface{}
	disablescanner := getPOSTRequest(datapaq_path +"/disableBarcodeScanner",nil)
	json.Unmarshal([]byte(disablescanner), &result)
	log.Println(result["disableBarcodeScanner"])
}

//shutdown datapaq webserver
func shutDownDatapaqServer(){
	var result map[string]interface{}
	shutdown := getPOSTRequest(datapaq_path +"/shutdown",nil)
	json.Unmarshal([]byte(shutdown), &result)
	log.Println(result["shutdown"])
}


// save last image specify the path
func getLastSavedImage(){
	var result map[string]interface{}
	saveImgMap := make(map[string]string)
	saveImgMap["path"] = "C:/temp/lastImg_1.png"
	saveImgMap["scaleFactor"] = "0.15"
	saveLastImage := getPOSTRequest(datapaq_path +"/saveLastImage",saveImgMap)
	json.Unmarshal([]byte(saveLastImage), &result)
	log.Println(result["saveLastImage"])
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