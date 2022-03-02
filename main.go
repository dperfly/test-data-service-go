package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type TxtFile struct {
	TotalNum int64
	Datas    []string
	curIndex *int64
}
type CsvFile struct {
	TotalNum int64
	Datas    [][]string
	curIndex *int64
}
type JsonFile struct {
	Datas map[string]interface{}
}
type FileDist struct {
	TxtFiles  map[string]TxtFile
	JsonFiles map[string]JsonFile
	CsvFiles  map[string]CsvFile
}

type CsvResult struct {
	Msg    string     `json:"msg"`
	Result [][]string `json:"result"`
}
type TxtResult struct {
	Msg    string   `json:"msg"`
	Result []string `json:"result"`
}
type JsonResult struct {
	Msg    string                 `json:"msg"`
	Result map[string]interface{} `json:"result"`
}

var fileDist = FileDist{
	TxtFiles:  make(map[string]TxtFile),
	JsonFiles: make(map[string]JsonFile),
	CsvFiles:  make(map[string]CsvFile),
}

func loadData(path string) {
	if path == "" {
		pwd, _ := os.Getwd()
		path = filepath.Join(pwd, "file")
	}
	filePathNames, _ := ioutil.ReadDir(path)

	for _, v := range filePathNames {

		AbsFilePath := filepath.Join(path, v.Name())
		log.Println("load file: ", AbsFilePath)
		fi, err := os.Open(AbsFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer fi.Close()

		if strings.HasSuffix(v.Name(), ".csv") {
			// csv
			ReadCsv := csv.NewReader(fi)
			readAll, err := ReadCsv.ReadAll()
			if err != nil {
				log.Fatal(err)
			}
			var curIndex int64
			curIndex++
			fileDist.CsvFiles[v.Name()] = CsvFile{
				Datas:    readAll,
				TotalNum: int64(len(readAll)),
				curIndex: &curIndex,
			}

		} else if strings.HasSuffix(v.Name(), ".json") {
			// json
			decoder := json.NewDecoder(fi)
			var jsonMap map[string]interface{}
			decoder.Decode(&jsonMap)
			fileDist.JsonFiles[v.Name()] = JsonFile{
				Datas: jsonMap,
			}

		} else {
			//txt
			br := bufio.NewReader(fi)
			linesData := make([]string, 0)
			for {
				a, _, c := br.ReadLine()
				if c == io.EOF {
					break
				}
				linesData = append(linesData, string(a))
			}
			var curIndex int64
			fileDist.TxtFiles[v.Name()] = TxtFile{
				Datas:    linesData,
				TotalNum: int64(len(linesData)),
				curIndex: &curIndex,
			}
		}
	}

}
func getTxt(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	fileName, ok := fileDist.TxtFiles[vars["title"]]

	if !ok {
		msg, _ := json.Marshal(TxtResult{
			Msg: "fileName not found",
		})
		writer.Write(msg)
		return
	}

	fileType := request.URL.Query().Get("type")
	var num int
	num, err := strconv.Atoi(request.URL.Query().Get("num"))
	if err != nil {
		num = 1
	}
	res := make([]string, 0)

	index := fileDist.TxtFiles[vars["title"]].curIndex
	switch fileType {
	case "random":
		rand.Seed(time.Now().Unix())
		for i := 0; i < num; i++ {
			randomNum := rand.Intn(int(fileName.TotalNum))
			res = append(res, fileName.Datas[randomNum])
		}
	default:
		if int64(num) > fileName.TotalNum {
			num = int(fileName.TotalNum)
		}
		if int64(num) <= 0 {
			num = 1
		}
		if atomic.LoadInt64(index)+int64(num) > fileName.TotalNum {
			atomic.SwapInt64(index, 0)

		}
		res = fileName.Datas[atomic.LoadInt64(index) : atomic.LoadInt64(index)+int64(num)]
		atomic.AddInt64(index, int64(num))

	}
	msg, _ := json.Marshal(TxtResult{
		Msg:    "",
		Result: res,
	})
	writer.Write(msg)

}
func getJson(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	fileName, ok := fileDist.JsonFiles[vars["title"]]
	if !ok {
		msg, _ := json.Marshal(JsonResult{
			Msg: "fileName not found",
		})
		writer.Write(msg)
		return
	}

	msg, _ := json.Marshal(JsonResult{
		Result: fileName.Datas,
	})

	writer.Write(msg)
}
func getCsv(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	csvFile, ok := fileDist.CsvFiles[vars["title"]]
	if !ok {
		msg, _ := json.Marshal(CsvResult{
			Msg: "csvFile not found",
		})
		writer.Write(msg)
		return
	}

	var num int
	num, err := strconv.Atoi(request.URL.Query().Get("num"))
	if err != nil || num <= 0 {
		num = 1
	}

	fileType := request.URL.Query().Get("type")
	res := make([][]string, 0)

	switch fileType {
	case "random":
		rand.Seed(time.Now().Unix())
		for i := 0; i < num; i++ {
			randomNum := rand.Intn(int(csvFile.TotalNum-1)) + 1
			res = append(res, csvFile.Datas[randomNum])
		}
	default:
		if int64(num) > csvFile.TotalNum-1 {
			num = int(csvFile.TotalNum) - 1
		}
		index := fileDist.CsvFiles[vars["title"]].curIndex
		if atomic.LoadInt64(index)+int64(num) > csvFile.TotalNum {
			atomic.SwapInt64(index, 1)
		}
		res = csvFile.Datas[atomic.LoadInt64(index) : atomic.LoadInt64(index)+int64(num)]
		atomic.AddInt64(index, int64(num))

	}
	msg, _ := json.Marshal(CsvResult{
		Msg:    "",
		Result: res,
	})
	writer.Write(msg)
}
func getData(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	if strings.HasSuffix(vars["title"], ".json") {
		getJson(writer, request)
	} else if strings.HasSuffix(vars["title"], ".csv") {
		getCsv(writer, request)
	} else {
		getTxt(writer, request)
	}
}

func main() {
	var ip string
	var port string
	var file string
	flag.StringVar(&ip, "ip", "", "ip")
	flag.StringVar(&port, "port", "8080", "port")
	flag.StringVar(&file, "file", "", "file")
	flag.Parse()
	loadData(file)
	addr := fmt.Sprint(ip, ":", port)
	r := mux.NewRouter()
	r.HandleFunc("/{title}", getData)
	log.Println("start server : ", "ip: ", ip, "port: ", port)
	err := http.ListenAndServe(addr, r)
	if err != nil {
		log.Fatal(err)
	}

}
