package lib

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func ScanDir(path string) []string{
	filesReturn := make([]string,0)
	files, err := ioutil.ReadDir(path)
	

    if err != nil {
        log.Fatal(err)
    }
	for _, f := range files {
		var tmpPath string
		tmpPath = fmt.Sprintf("%s/%s",path,f.Name());

		tmpFile,err := os.Open(tmpPath);

		if err !=nil{
			log.Fatal(err)
		}

		tmpFileInfo,_ := tmpFile.Stat()
		if tmpFileInfo.IsDir(){
			filesReturn = append(filesReturn,ScanDir(tmpPath)...)
		} else {
			filesReturn = append(filesReturn,path+"/"+f.Name())
		}
	}
	
	return filesReturn
}