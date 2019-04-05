# rhub
distribute websocket = websocket + redis


### Example:
```go
func startServer() {
	hub := NewHub(1, conf.Redis, "test-room-1")
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
		ServeWs(hub, w, r, clientProp, DefaultUpgrader(), DefaultWsConfig())
	})
	fmt.Println("start server al 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```