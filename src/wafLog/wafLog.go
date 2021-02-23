package wafLog

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/exec"
)

//func CreateLog(c *gin.Context) {
//	var json LogDB
//	if err := c.ShouldBindJSON(&json); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	json.Id = uuid.NewV4()
//	json.Timestamp = time.Now().Unix()
//	fmt.Println(json)
//	Insert(json)
//	c.JSON(200, gin.H{
//		"message": "success",
//	})
//}
//
//func DeleteLog(c *gin.Context) {
//	var json LogDB
//	if err := c.ShouldBindJSON(&json); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	fmt.Println(json)
//	Delete(json.Id)
//	c.JSON(200, gin.H{
//		"message": "success",
//	})
//}
//
//func UpdateLog(c *gin.Context) {
//	//Update()
//	var json LogDB
//	if err := c.ShouldBindJSON(&json); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	fmt.Println(json)
//	Update(json)
//	c.JSON(200, gin.H{
//		"message": "success",
//	})
//}

const SUCCESS = 0
const FAIL = 1

type Logs struct {
	Id        string   `json:"id,omitempty"`
	Contents  []string `json:"contents,omitempty"`
	Length    int      `json:"length,omitempty"`
	Timestamp string   `json:"timestamp,omitempty" `
	Path      string   `json:"path,omitempty"`
}

type Ids struct {
	Ids []string `json:"ids,omitempty"`
}

func readFile(l *Logs) {
	file, err := os.Open(l.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineText := scanner.Text()
		l.Contents = append(l.Contents, lineText)
		l.Length += 1
	}
}

func readLogs(logDB []LogDB) ([]Logs, []error) {
	var errs []error
	var logs []Logs
	if len(logDB) > 0 {
		for _, l := range logDB {
			log := Logs{l.Id, nil, 0, l.Timestamp, l.LogPath}
			readFile(&log)
			logs = append(logs, log)
		}
	}
	return logs, errs
}

func RetrieveLog(c *gin.Context) {
	var ids Ids
	var logDB []LogDB
	if err := c.ShouldBindJSON(&ids); err != nil {
		log.Printf("Bind JSON ERROR: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	fmt.Println(ids)
	if len(ids.Ids) <= 0 {
		log.Println("Bind JSON ERROR: No id")
		c.JSON(http.StatusBadRequest, gin.H{"error": "No ids"})
		return
	}
	for _, i := range ids.Ids {
		l, err := Select(i)
		if err != nil {
			log.Fatalf("Select ERROR: %v\n", err)
		}
		logDB = append(logDB, l...)
	}
	logs, errs := readLogs(logDB)
	if len(errs) > 0 {
		for e := range errs {
			log.Fatalf("Read file ERROR: %v\n", e)
		}
	}
	c.JSON(http.StatusOK, logs)
}

func RetrieveAll(c *gin.Context) {
	logDB, err := SelectAll()
	if err != nil {
		log.Printf("Select ERROR: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}
	logs, errs := readLogs(logDB)
	if len(errs) > 0 {
		for e := range errs {
			log.Fatalf("Read file ERROR: %v\n", e)
		}
	}
	c.JSON(http.StatusOK, logs)
}

func getWafPid() string {
	cmd := exec.Command("/bin/sh", "-c", `ps -ef | grep -v "grep" | grep "waf.py" | awk '{print $2}'`)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return ""
	}
	fmt.Println(string(output))
	return string(output)
}

func wafStart() int {
	pid := getWafPid()
	if pid != "" {
		return FAIL
	}
	cmd := exec.Command("/bin/sh", "-c", "nohup /root/MicroWAF/start.sh &")
	fmt.Println(cmd.Args)
	if err := cmd.Run(); err != nil {
		log.Printf("CMD ERROR: %v\n", err)
	}
	return SUCCESS
}

func wafStop() int {
	pid := getWafPid()
	if pid == "" {
		return FAIL
	}
	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("kill -9 %s", pid))
	if err := cmd.Run(); err != nil {
		log.Printf("CMD ERROR: %v\n", err)
	}
	return SUCCESS
}

func WafStart(c *gin.Context) {
	if wafStart() == FAIL {
		c.JSON(http.StatusBadRequest, gin.H{"message": "waf is running"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Start OK"})
}

func WafStop(c *gin.Context) {
	if wafStop() == FAIL {
		c.JSON(http.StatusBadRequest, gin.H{"message": "waf is not running"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Stop OK"})
}

func WafRestart(c *gin.Context) {
	pid := getWafPid()
	if pid == "" {
		wafStart()
	} else {
		wafStop()
		wafStart()
	}

	c.JSON(http.StatusOK, gin.H{"message": "Restart OK"})
}
