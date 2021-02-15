package main

import (
	"WafLog/src/wafLog"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	r := gin.Default()
	//r.POST("/app/v1/creat", wafLog.CreateLog)
	//r.POST("/app/v1/update", wafLog.UpdateLog)
	r.GET("/log/v1/selectById", wafLog.RetrieveLog)
	r.GET("/log/v1/selectAll", wafLog.RetrieveAll)
	//r.POST("/app/v1/delete", wafLog.DeleteLog)

	go func() {
		r.Run(":9123") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	}()

	quit := make(chan os.Signal)
	defer wafLog.Close()
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")


}