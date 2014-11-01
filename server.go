package main

import (
    "code.google.com/p/go.crypto/bcrypt"
	"database/sql"
	"fmt"
	"github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/securecookie"
	"github.com/kr/s3/s3util"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// cookie handling

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type Signin struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type Register struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	Nickname string `form:"nickname"`
}


type Project struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Msgs struct {
	Id       int    `json:"id"`
	Metadata string `json:"message"`
	Filetype string `json:"filetype"`
}

type Groupdata struct {
	Id       int    `json:"id"`
	Metadata string `json:"message"`
	Filetype string `json:"filetype"`
	User_id  string `json:"user_id`
	Username string `json:"username`
}

type Groupnames struct {
	Id         int    `json:"groupid"`
	Group_name string `json:"groupname"`
}

type Msg struct {
	Message string `form:"message"`
}

type Wish struct {
	Id   string `form:"id"`
	Name string `form:"name"`
	Age  string `form:"age"`
}

// Chat top level
type Notify struct {
	sync.Mutex
	groups []*Group
}

// Room level
type Group struct {
	sync.Mutex
	name    string
	clients map[int]*Client
}

// Client stores all the channels available in the handler in a struct.
type Client struct {
	Name       string
	in         <-chan *Message
	out        chan<- *Message
	done       <-chan bool
	err        <-chan error
	disconnect chan<- int
}

// A simple Message struct
type Message struct {
	Type    string `json:"type"`
	From    string `json:"from"`
	GroupId string `json:"groupid"`
	Text    string `json:"text"`
}

// the chat
var notify *Notify

func (c *Notify) getGroup(name string) *Group {
	c.Lock()
	defer c.Unlock()

	for _, group := range c.groups {
		if group.name == name {
			return group
		}
	}

	r := &Group{sync.Mutex{}, name, make(map[int]*Client)}
	c.groups = append(c.groups, r)

	return r
}

// Add a client to a room
func (r *Group) appendClient(id int, client *Client) {
	r.Lock()
	_, ok := r.clients[id]
	if ok {
		//ctemp.disconnect<-1000
		fmt.Println("one client disconnected")
	}
	r.clients[id] = client
	fmt.Println("number of clients %v", len(r.clients))
	r.Unlock()
}

// Message all the other clients in the same room
func (r *Group) messageOtherClients(client *Client, msg *Message) {
	r.Lock()
	msg.From = client.Name
	//add all values to the database here
	for _, c := range r.clients {
		if c != client {
			c.out <- msg
			//c.disconnect <-1
		}
	}
	defer r.Unlock()
}

func newNotify() *Notify {
	return &Notify{sync.Mutex{}, make([]*Group, 0)}
}

func GetMessages(db *sql.DB, id string) []Msgs {

	msgresult, err := db.Query("select id_message,metadata,filetype from messages where user_id=" + id + " order by id_message desc")
	if err != nil {
		fmt.Println(err)
	}

	var (
		id_message int
		metadata   string
		filetype   string
	)

	p := make([]Msgs, 0)
	defer msgresult.Close()
	for msgresult.Next() {
		err := msgresult.Scan(&id_message, &metadata, &filetype)
		if err != nil {
			fmt.Println(err)
		} else {
			p = append(p, Msgs{id_message, metadata, filetype})
		}
	}
	return p
}

func GetGroups(db *sql.DB, id string) []Groupnames {

	result, err := db.Query("select group_name,group_id from groups where group_id IN(select group_id from usergroups where user_id=" + id + ") order by group_id")

	fmt.Println(" ")
	if err != nil {
		fmt.Println(err)
	}

	var (
		group_id   int
		group_name string
	)

	p := make([]Groupnames, 0)
	defer result.Close()
	for result.Next() {
		err := result.Scan(&group_name, &group_id)
		if err != nil {
			fmt.Println(err)
		} else {
			p = append(p, Groupnames{group_id, group_name})
		}
	}
	return p
}

func GetGroupData(db *sql.DB, id string) []Groupdata {

	result, err := db.Query("select g.groupdata_id,g.metadata,g.filetype,g.user_id,u.username from groupdata g,username u where g.group_id=" + id + " and u.id=g.user_id order by g.groupdata_id desc")

	fmt.Println(" ")
	if err != nil {
		fmt.Println(err)
	}

	var (
		id_groupdata int
		metadata     string
		filetype     string
		user_id      string
		username     string
	)

	p := make([]Groupdata, 0)
	defer result.Close()
	for result.Next() {
		err := result.Scan(&id_groupdata, &metadata, &filetype, &user_id, &username)
		if err != nil {
			fmt.Println(err)
		} else {
			p = append(p, Groupdata{id_groupdata, metadata, filetype, user_id, username})
		}
	}
	return p
}

func main() {

	s3util.DefaultConfig.AccessKey = "your access key"
	s3util.DefaultConfig.SecretKey = "your security key"

	m := martini.Classic()
	m.Use(render.Renderer())

	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("megh", store))

	db, err := sql.Open("mysql", "root:123456@/userdb")

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	m.Map(db)
	defer db.Close()

	m.Get("/", func(r render.Render) {

		r.HTML(200, "index", nil)

	})

	m.Get("/elements", func(r render.Render, request *http.Request) {

		userName := getUserName(request)

		if userName == "" {
			r.HTML(200, "login", "Session Expired .Please Sign in again")
		} else {
			r.HTML(200, "user", userName)
		}
	})

	m.Get("/login", func(r render.Render) {
		r.HTML(200, "login", nil)
	})

	m.Post("/login", binding.Form(Signin{}), func(signin Signin, r render.Render, response http.ResponseWriter, s sessions.Session) {

		var id string
		var hashedPassword string
		//err := db.QueryRow("select id from username where username='" + signin.Email + "' and password='" + hashedPassword + "'").Scan(&id)
		err := db.QueryRow("select id,password from username where username='" + signin.Email + "'").Scan(&id,&hashedPassword)
		
		err2 := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(signin.Password))
		
		if err != nil && err2 !=nil {

			r.HTML(200, "login", "Email or password is incorrect.Try again")
		} else {

			setSession(signin.Email, response)
			r.HTML(200, "user", signin.Email)
		}
	})

	m.Get("/register", func(r render.Render) {
		r.HTML(200, "register", nil)
	})

	m.Post("/register", binding.Form(Register{}), func(register Register, r render.Render, s sessions.Session) {

       hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(register.Password), 10)
		_, err := db.Query("insert into username(username,password,date_joined,nickname) values(?,?,?,?)", register.Email,hashedPassword,time.Now().UTC(),register.Nickname)
		fmt.Println(register.Email, hashedPassword, time.Now().UTC(),register.Nickname)
		if err != nil {
			fmt.Println(err)
			r.HTML(200, "register", "Registration failed .Looks Like there is already an account with the same email id")
		} else {
			// s.Set("userId", id)
			r.HTML(200, "login", "Registration Successfull.  You can login now")
		}
	})

	m.Get("/logout", func(r render.Render, response http.ResponseWriter) {
		// s.Delete("userId")
		clearSession(response)
		r.HTML(200, "login", nil)
	})

	m.Post("/messages", binding.Form(Msg{}), func(msg Msg, r render.Render, request *http.Request) {

		var iddb string

		userName := getUserName(request)

		errs := db.QueryRow("select id from username where username='" + userName + "'").Scan(&iddb)
		if errs == nil {
			_, err := db.Query("insert into messages(metadata,filetype,user_id,timestamp) values(?,'text',?,?)", msg.Message, iddb,time.Now().UTC())
			if err != nil {
				panic(err.Error())
				// proper error handling instead of panic in your app
			}
			r.HTML(200, "user", nil)
		} else {

			panic(errs.Error())
		}

	})

	m.Post("/urls", binding.Form(Msg{}), func(msg Msg, r render.Render, request *http.Request) {

		var iddb string

		userName := getUserName(request)

		errs := db.QueryRow("select id from username where username='" + userName + "'").Scan(&iddb)
		if errs == nil {
			_, err := db.Query("insert into messages(metadata,filetype,user_id) values(?,'url',?)", msg.Message, iddb)
			if err != nil {
				panic(err.Error())
				// proper error handling instead of panic in your app
			}
			r.HTML(200, "user", nil)
		} else {

			panic(errs.Error())
		}

	})

	m.Post("/groupdata/:id", binding.Form(Msg{}), func(msg Msg, params martini.Params, r render.Render, request *http.Request) {

		var iddb string
		var count string
		userName := getUserName(request)

		errs := db.QueryRow("select id from username where username='" + userName + "'").Scan(&iddb)
		db.QueryRow("select count(*) from usergroups where user_id='" + iddb + "' and group_id='" + params["id"] + "'").Scan(&count)

		if errs == nil && count == "1" {
			_, err := db.Query("insert into groupdata(metadata,filetype,user_id,group_id,timestamp) values(?,'text',?,?,?)", msg.Message, iddb, params["id"],time.Now().UTC())
			if err != nil {
				panic(err.Error())
				// proper error handling instead of panic in your app
			}
			r.JSON(200, map[string]interface{}{"status": "Success"})
		} else {

			panic(errs.Error())
		}

	})

	m.Get("/groups", func(r render.Render, request *http.Request) {

		userName := getUserName(request)

		if userName == "" {
			r.HTML(200, "login", nil)
		} else {
			r.HTML(200, "groups", userName)
		}

	})

	m.Get("/groupslist", func(r render.Render, params martini.Params, request *http.Request) {
		var iddb string

		userName := getUserName(request)

		db.QueryRow("select id from username where username='" + userName + "'").Scan(&iddb)

		if userName == "" || iddb == "" {

			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		} else {
			msgs := GetGroups(db, iddb)
			r.JSON(200, map[string]interface{}{"status": "Success", "data": msgs})
		}
	})

	m.Post("/creategroup", binding.Form(Msg{}), func(msg Msg, r render.Render, request *http.Request) {
		var iddb string
		var groupid string
		userName := getUserName(request)

		errs := db.QueryRow("select id from username where username='" + userName + "'").Scan(&iddb)
		if errs == nil {
			_, err1 := db.Query("insert into groups(group_name,group_admin_id,group_timestamp) values(?,?,?)", msg.Message, iddb,time.Now().UTC())

			if err1 != nil {
				fmt.Println("err1 error")
				panic(err1.Error())
				// proper error handling instead of panic in your app
			}
			db.QueryRow("select group_id from groups where group_admin_id='" + iddb + "' and group_name='" + msg.Message + "'").Scan(&groupid)

			fmt.Println(groupid, iddb)
			db.Query("insert into usergroups(group_id,user_id) values(?,?)", groupid, iddb)

			r.JSON(200, map[string]interface{}{"status": "Success"})
		} else {
			fmt.Println("errs error")
			panic(errs.Error())
		}

	})

	m.Get("/getgroupdatabyid/:id", func(r render.Render, params martini.Params, request *http.Request) {
		var iddb string
		var count string
		userName := getUserName(request)

		db.QueryRow("select id from username where username='" + userName + "'").Scan(&iddb)
		db.QueryRow("select count(*) from usergroups where user_id='" + iddb + "' and group_id='" + params["id"] + "'").Scan(&count)

		if userName == "" || count == "0" {

			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		} else if count == "1" {
			msgs := GetGroupData(db, params["id"])
			r.JSON(200, map[string]interface{}{"status": "Success", "data": msgs})
		}
	})

	m.Get("/getmessagesbyid/:id", func(r render.Render, params martini.Params, request *http.Request) {
		var iddb string

		userName := getUserName(request)

		db.QueryRow("select id from username where username='" + userName + "'").Scan(&iddb)

		if userName == "" || iddb != params["id"] {

			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		} else if iddb == params["id"] {
			msgs := GetMessages(db, params["id"])
			r.JSON(200, map[string]interface{}{"status": "Success", "data": msgs})
		}
	})

	m.Get("/getmessages", func(r render.Render, params martini.Params, request *http.Request) {
		var iddb string

		userName := getUserName(request)

		db.QueryRow("select id from username where username='" + userName + "'").Scan(&iddb)

		if userName == "" || iddb == "" {

			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		} else {
			msgs := GetMessages(db, iddb)
			r.JSON(200, map[string]interface{}{"status": "Success", "data": msgs})
		}
	})

	notify = newNotify()

	m.Get("/sockets/:clientname", sockets.JSON(Message{}), func(params martini.Params, request *http.Request, receiver <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {

		client := &Client{params["clientname"], receiver, sender, done, err, disconnect}
		var iddb string

		userName := getUserName(request)

		db.QueryRow("select id from username where username='" + userName + "'").Scan(&iddb)

		groups := GetGroups(db, iddb)

		fmt.Println("groups length %v", len(groups))
		for _, group := range groups {

			fmt.Println(group.Id)
			r := notify.getGroup(strconv.Itoa(group.Id))
			var id, _ = strconv.Atoi(iddb)
			r.appendClient(id, client)
		}

		// A single select can be used to do all the messaging
		for {
			select {
			case <-client.err:
				// Don't try to do this:
				// client.out <- &Message{"system", "system", "There has been an error with your connection"}
				// The socket connection is already long gone.
				// Use the error for statistics etc
			case msg := <-client.in:
				r := notify.getGroup(msg.GroupId)
				r.messageOtherClients(client, msg)
			case <-client.done:
				return 200, "OK"
			}
		}
	})
	m.Post("/uploadfile", uploadHandler)
	m.Post("/groupuploadfile/:id", groupUploadHandler)
	m.Get("/getfile", ReturnFile)
	m.Run()

}

func uploadHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, rend render.Render) {
	file, header, err := r.FormFile("file") // the FormFile function takes in the POST input id file
	defer file.Close()

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	l, _ := s3util.Create("https://s3-us-west-2.amazonaws.com/megh-uploads/"+header.Filename, nil, nil)
	_, err2 := io.Copy(l, file)
	l.Close()
	if err2 != nil {
		fmt.Fprintln(w, err)
	}
	// write the content from POST to the file

	var iddb string

	userName := getUserName(r)

	errs := db.QueryRow("select id from username where username='" + userName + "'").Scan(&iddb)
	if errs == nil {
		_, err := db.Query("insert into messages(metadata,filetype,user_id) values(?,'file',?)", header.Filename, iddb)
		if err != nil {
			panic(err.Error())
			// proper error handling instead of panic in your app
		}
		rend.HTML(200, "elements", nil)
	} else {

		panic(errs.Error())
	}

}
func groupUploadHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, rend render.Render, params martini.Params) {
	file, header, err := r.FormFile("file") // the FormFile function takes in the POST input id file
	defer file.Close()

	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	l, _ := s3util.Create("https://s3-us-west-2.amazonaws.com/megh-uploads/"+header.Filename, nil, nil)
	_, err2 := io.Copy(l, file)
	l.Close()
	if err2 != nil {
		fmt.Fprintln(w, err)
	}

	var iddb string
	var count string
	userName := getUserName(r)

	errs := db.QueryRow("select id from username where username='" + userName + "'").Scan(&iddb)
	db.QueryRow("select count(*) from usergroups where user_id='" + iddb + "' and group_id='" + params["id"] + "'").Scan(&count)

	if errs == nil && count == "1" {
		_, err := db.Query("insert into groupdata(metadata,filetype,user_id,group_id) values(?,'file',?,?)", header.Filename, iddb, params["id"])

		if err != nil {
			panic(err.Error())
			// proper error handling instead of panic in your app
		}
		rend.HTML(200, "groups", userName)
	} else {

		panic(errs.Error())
	}

}

func ReturnFile(writer http.ResponseWriter, req *http.Request) {

	name := req.URL.Query()["filename"][0]
	fmt.Println(name)
	writer.Header().Set("Content-type", "multipart")
	writer.Header().Set("Content-Disposition", "attachment;filename="+name)
	http.ServeFile(writer, req, "/clouduploads/"+name)

}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}
