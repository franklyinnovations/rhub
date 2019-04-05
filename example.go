package rhub

import (
	"fmt"
	"log"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	http.ServeFile(w, r, "home.html")
}

type User struct {
	Name string
}

func startServer() {
	hub := NewHub(1, "redis://test.iqidao.com:50002/0", "test-room-1")
	// hub.BeforeJoin(func(client IClient) error {
	// 	return nil
	// })
	// hub.Use(func(m *ClientHubMessage, next func()) {
	// 	fmt.Printf("use before %+v\n", m)
	// 	next()
	// 	fmt.Printf("user after %+v\n", m)
	// })
	// hub.Use(func(m *ClientHubMessage, next func()) {
	// 	fmt.Printf("use before2 %+v\n", m)
	// 	next()
	// 	fmt.Printf("user after2 %+v\n", m)
	// })
	hub.On("join", func(m *RedisHubMessage) {
		fmt.Println("join", *m)
	})
	hub.On("leave", func(m *RedisHubMessage) {
		fmt.Println("leave", *m)
	})
	hub.On("im", func(m *RedisHubMessage) {
		var str string
		fmt.Println("im ", string(*m.Data))
		err := m.Decode(&str)
		fmt.Println(str, err)
		// hub.SendRedis("im", "your say:"+str, nil)
		hub.SendWsAll("im", "your say:"+str)
		if "close" == str {
			fmt.Println("close client")
			hub.Close()
			// go func() { hub.UnregisterChan() <- m.Client }()
		}
	})

	go hub.Run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		clientProp := map[string]interface{}{"user": &User{Name: "jim"}}
		upgrader = DefaultUpgrader()
		ServeWs(hub, w, r, clientProp, upgrader)
	})
	fmt.Println("start server al 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
