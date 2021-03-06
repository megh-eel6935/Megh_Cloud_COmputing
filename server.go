package main

import (
	"code.google.com/p/go-uuid/uuid"
	"code.google.com/p/go.crypto/bcrypt"
	"encoding/json"
	"fmt"
	"github.com/alexjlockwood/gcm"
	"github.com/beatrichartz/martini-sockets"
	"github.com/go-martini/martini"
	"github.com/gorilla/securecookie"
	"github.com/kr/s3/s3util"
	"github.com/martini-contrib/binding"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/smugmug/godynamo/conf"
	"github.com/smugmug/godynamo/conf_file"
	ep "github.com/smugmug/godynamo/endpoint"
	delete_item "github.com/smugmug/godynamo/endpoints/delete_item"
	get "github.com/smugmug/godynamo/endpoints/get_item"
	put "github.com/smugmug/godynamo/endpoints/put_item"
	query "github.com/smugmug/godynamo/endpoints/query"
	update_item "github.com/smugmug/godynamo/endpoints/update_item"
	"github.com/smugmug/godynamo/types/attributevalue"
	"github.com/smugmug/godynamo/types/condition"
	"io"
	"math/rand"
	"net/http"
	"sort"
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

type MobileSignin struct {
	Email           string `form:"email"`
	Password        string `form:"password"`
	Registration_id string `form:"registration_id"`
}

type Pkey struct {
	Email     string `form:"email"`
	Publickey string `form:"publickey"`
	Groupname string `form:"groupname"`
	Groupid   string `form:"groupid"`
}

type Register struct {
	Email    string `form:"email"`
	Password string `form:"password"`
	Nickname string `form:"nickname"`
}

type Msgs struct {
	Id       int    `json:"id"`
	Metadata string `json:"message"`
	Filetype string `json:"filetype"`
}

type GroupnamesJson struct {
	Id          string `json:"group_id"`
	Username    string `json:"username"`
	Group_name  string `json:"group_name"`
	Group_admin string `json:"group_admin"`
	Timestamp   string `json:"timestamp"`
}

type GroupSocketsJson struct {
	Id           string `json:"group_id"`
	SocketServer string `json:"socketserver"`
}

type GroupusernamesJson struct {
	Id          string `json:"group_id"`
	Username    string `json:"username"`
	Group_admin string `json:"group_admin"`
	Timestamp   string `json:"timestamp"`
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
	clients map[string]*Client
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

type GroupdataJson struct {
	Id        string `json:"groupdata_id"`
	Content   string `json:"content"`
	Filetype  string `json:"content_type"`
	Username  string `json:"username"`
	Timestamp string `json:"timestamp"`
}

type UserDataJson struct {
	Content      string `json:"content"`
	Content_type string `json:"content_type"`
	Timestamp    string `json:"timestamp"`
}

type strd struct {
	S string
}

type numd struct {
	N string
}

type groupitems []struct {
	Timestamp    numd
	Username     strd
	Group_name   strd
	Group_id     strd
	Usergroup_id strd
	Group_admin  strd
}

type AndroidUserDetailsItem struct {
	Username        strd
	Timestamp       numd
	Publickey       strd
	Displayname     strd
	Password        strd
	Registration_id strd
}

type AndroidUserDetails struct {
	Item AndroidUserDetailsItem
}

type groupuseritems []struct {
	Group_id    strd
	Username    strd
	Timestamp   numd
	Group_admin numd
}

type groupsList struct {
	Count        int
	Items        groupitems
	ScannedCount int
}
type groupdetailitem struct {
	Username     strd
	Group_name   strd
	Timestamp    numd
	Group_id     strd
	Usergroup_id strd
	Group_admin  strd
}
type groupdetail struct {
	Item groupdetailitem
}

type groupusersList struct {
	Count        int
	Items        groupuseritems
	ScannedCount int
}

type groupdataitems []struct {
	Username     strd
	Content_type strd
	Timestamp    strd
	Group_id     strd
	Groupdata_id strd
	Content      strd
}
type groupdatasingleitem struct {
	Username     strd
	Content_type strd
	Timestamp    strd
	Group_id     strd
	Groupdata_id strd
	Content      strd
}

type groupdataitem struct {
	Item groupdatasingleitem
}

type groupdata struct {
	Count        int
	Items        groupdataitems
	ScannedCount int
}

type userdataitems []struct {
	Username     strd
	Content_type strd
	Timestamp    strd
	Content      strd
}

type userdata struct {
	Count        int
	Items        userdataitems
	ScannedCount int
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

	r := &Group{sync.Mutex{}, name, make(map[string]*Client)}
	c.groups = append(c.groups, r)

	return r
}

// Add a client to a room
func (r *Group) appendClient(username string, client *Client) {
	r.Lock()
	_, ok := r.clients[username]
	if ok {
		//ctemp.disconnect<-1000
		fmt.Println("one client disconnected")
	}
	r.clients[username] = client
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

func main() {

	conf_file.Read()
	if conf.Vals.Initialized == false {
		panic("the conf.Vals global conf struct has not been initialized")
	}

	s3util.DefaultConfig.AccessKey = ""
	s3util.DefaultConfig.SecretKey = ""

	m := martini.Classic()
	m.Use(render.Renderer())

	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("megh", store))

	m.Get("/", func(r render.Render) {

		r.HTML(200, "index", nil)

	})

	m.Get("/home", func(r render.Render, request *http.Request) {

		userName := getUserName(request)

		if userName == "" {
			r.HTML(200, "index", "Session Expired .Please Sign in again")
		} else {
			r.HTML(200, "user", userName)
		}
	})

	m.Get("/chromehome", func(r render.Render, request *http.Request) {

		userName := getUserName(request)

		if userName == "" {
			r.HTML(200, "chromeindex", "Session Expired .Please Sign in again")
		} else {
			r.HTML(200, "chromeuser", userName)
		}
	})

	m.Get("/chromefileupload", func(r render.Render, request *http.Request) {

		userName := getUserName(request)

		if userName == "" {
			r.JSON(200, map[string]interface{}{"status": "Seesion failed . Login again"})
		} else {
			r.HTML(200, "chromefileupload", "click to upload files")
		}
	})

	m.Post("/home", binding.Form(Signin{}), func(signin Signin, r render.Render, response http.ResponseWriter, s sessions.Session) {

		//err := db.QueryRow("select id,password from username where username='" + signin.Email + "'").Scan(&id,&hashedPassword)
		get1 := get.NewGetItem()
		get1.TableName = "users"
		fmt.Println("entir form is ", signin)
		fmt.Println("username is ", signin.Email)
		get1.Key["username"] = &attributevalue.AttributeValue{S: signin.Email}
		body, code, err := get1.EndpointReq()
		fmt.Println("body is %s", body)
		if len(body) <= 2 || err != nil || code != http.StatusOK {
			fmt.Printf("get failed %d %v %d \n", code, err, len(body))
			r.HTML(200, "index", "Username doesot exists in the database")
		} else {
			payload := []byte(body)
			var result map[string]interface{}
			if err := json.Unmarshal(payload, &result); err != nil {
				panic(err)
			}

			row := result["Item"].(map[string]interface{})
			password := row["password"].(map[string]interface{})
			publickey := row["publickey"].(map[string]interface{})

			err2 := bcrypt.CompareHashAndPassword([]byte(password["S"].(string)), []byte(signin.Password))

			if err2 != nil {

				r.HTML(200, "index", "Email or password is incorrect.Try again")
			} else {
				var servername []string
				servername = decideServer(signin.Email)
				setSession(signin.Email, publickey["S"].(string), response, servername)
				r.HTML(200, "user", signin.Email)
			}
		}
	})

	m.Post("/chromelogin", binding.Form(Signin{}), func(signin Signin, r render.Render, response http.ResponseWriter, s sessions.Session) {

		//err := db.QueryRow("select id,password from username where username='" + signin.Email + "'").Scan(&id,&hashedPassword)
		get1 := get.NewGetItem()
		get1.TableName = "users"
		fmt.Println("entir form is ", signin)
		fmt.Println("username is ", signin.Email)
		get1.Key["username"] = &attributevalue.AttributeValue{S: signin.Email}

		body, code, err := get1.EndpointReq()
		fmt.Println("body is %s", body)
		if len(body) <= 2 || err != nil || code != http.StatusOK {
			fmt.Printf("get failed %d %v %d \n", code, err, len(body))
			r.HTML(200, "index", "Username doesot exists in the database")
		} else {
			payload := []byte(body)
			var result map[string]interface{}
			if err := json.Unmarshal(payload, &result); err != nil {
				panic(err)
			}

			row := result["Item"].(map[string]interface{})
			password := row["password"].(map[string]interface{})
			publickey := row["publickey"].(map[string]interface{})

			err2 := bcrypt.CompareHashAndPassword([]byte(password["S"].(string)), []byte(signin.Password))

			if err2 != nil {

				r.HTML(200, "chromeindex", "Email or password is incorrect.Try again")
			} else {
				setSession(signin.Email, publickey["S"].(string), response, nil)
				r.HTML(200, "chromeuser", signin.Email)
			}
		}
	})

	m.Post("/mobilelogin", binding.Form(MobileSignin{}), func(signin MobileSignin, r render.Render, response http.ResponseWriter, s sessions.Session) {

		//err := db.QueryRow("select id,password from username where username='" + signin.Email + "'").Scan(&id,&hashedPassword)
		get1 := get.NewGetItem()
		get1.TableName = "users"

		fmt.Println("entir form is ", signin)
		fmt.Println("username is ", signin.Email)
		get1.Key["username"] = &attributevalue.AttributeValue{S: signin.Email}

		body, code, err := get1.EndpointReq()
		fmt.Println("body is %s", body)
		if len(body) <= 2 || err != nil || code != http.StatusOK {
			fmt.Printf("get failed %d %v %d \n", code, err, len(body))
			r.JSON(200, map[string]interface{}{"status": "Username doesot exists in the database"})
		} else {
			payload := []byte(body)
			var result map[string]interface{}
			if err := json.Unmarshal(payload, &result); err != nil {
				panic(err)
			}

			row := result["Item"].(map[string]interface{})
			password := row["password"].(map[string]interface{})
			publickey := row["publickey"].(map[string]interface{})
			displayname := row["displayname"].(map[string]interface{})
			err2 := bcrypt.CompareHashAndPassword([]byte(password["S"].(string)), []byte(signin.Password))

			if err2 != nil {
				r.JSON(200, map[string]interface{}{"status": "Email or password is incorrect.Try again"})

			} else {

				ctime := fmt.Sprintf("%v", time.Now().Unix())
				put1 := put.NewPutItem()
				put1.TableName = "users"
				put1.Item["username"] = &attributevalue.AttributeValue{S: signin.Email + "android"}
				put1.Item["password"] = &attributevalue.AttributeValue{S: password["S"].(string)}
				put1.Item["displayname"] = &attributevalue.AttributeValue{S: displayname["S"].(string)}
				put1.Item["timestamp"] = &attributevalue.AttributeValue{N: ctime}
				put1.Item["publickey"] = &attributevalue.AttributeValue{S: publickey["S"].(string)}
				put1.Item["registration_id"] = &attributevalue.AttributeValue{S: signin.Registration_id}
				body1, code1, err1 := put1.EndpointReq()

				if err1 != nil || code != http.StatusOK {

					fmt.Printf("put failed %d %v %s\n", code1, err1, body1)
					r.JSON(200, map[string]interface{}{"status": "Registration failed .Try again"})

				} else {

					setSession(signin.Email, publickey["S"].(string), response, nil)

					r.JSON(200, map[string]interface{}{"status": "Success"})
				}

			}
		}
	})

	m.Post("/adduser", binding.Form(Pkey{}), func(pkey Pkey, r render.Render, request *http.Request, response http.ResponseWriter, s sessions.Session) {
		ctime := fmt.Sprintf("%v", time.Now().Unix())
		username := getUserName(request)
		get1 := get.NewGetItem()
		get1.TableName = "users"
		get1.Key["username"] = &attributevalue.AttributeValue{S: pkey.Email}
		body, code, err := get1.EndpointReq()
		fmt.Println("body is %s", body)
		if len(body) <= 2 || err != nil || code != http.StatusOK {
			fmt.Printf("get failed %d %v %d \n", code, err, len(body))
			r.HTML(200, "login", "Username doesot exists in the database")
		} else {
			payload := []byte(body)
			var result map[string]interface{}
			if err := json.Unmarshal(payload, &result); err != nil {
				panic(err)
			}

			row := result["Item"].(map[string]interface{})
			upublickey := row["publickey"].(map[string]interface{})

			if pkey.Publickey == upublickey["S"].(string) {

				put2 := put.NewPutItem()
				put2.TableName = "usergroups"
				fmt.Println("value to be sinserted", pkey.Email)
				put2.Item["username"] = &attributevalue.AttributeValue{S: pkey.Email}
				put2.Item["group_id"] = &attributevalue.AttributeValue{S: pkey.Groupid}
				put2.Item["group_name"] = &attributevalue.AttributeValue{S: pkey.Groupname}
				put2.Item["group_admin"] = &attributevalue.AttributeValue{S: username}
				put2.Item["timestamp"] = &attributevalue.AttributeValue{N: ctime}

				body2, code2, err2 := put2.EndpointReq()
				if err2 != nil || code2 != http.StatusOK {
					fmt.Printf("put failed %d %v %s\n", code2, err2, body2)
					r.JSON(200, map[string]interface{}{"status": "failure"})
				} else {
					put3 := put.NewPutItem()
					put3.TableName = "groupusers"
					put3.Item["group_id"] = &attributevalue.AttributeValue{S: pkey.Groupid}
					put3.Item["username"] = &attributevalue.AttributeValue{S: pkey.Email}
					put3.Item["timestamp"] = &attributevalue.AttributeValue{N: ctime}
					put3.Item["group_admin"] = &attributevalue.AttributeValue{N: "0"}

					body3, code3, err3 := put3.EndpointReq()
					if err3 != nil || code3 != http.StatusOK {
						fmt.Printf("put failed 4%d %v %s\n", code3, err3, body3)
						r.JSON(200, map[string]interface{}{"status": "success"})
					} else {
						r.JSON(200, map[string]interface{}{"status": "failure"})
					}
				}
			}
		}
	})

	m.Get("/register", func(r render.Render) {
		r.HTML(200, "register", nil)
	})

	m.Post("/register", binding.Form(Register{}), func(register Register, r render.Render, s sessions.Session) {

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(register.Password), 10)
		// _, err := db.Query("insert into username(username,password,date_joined,nickname) values(?,?,?,?)", register.Email,hashedPassword,time.Now().UTC(),register.Nickname)
		get1 := get.NewGetItem()
		get1.TableName = "users"
		get1.Key["username"] = &attributevalue.AttributeValue{S: register.Email}

		body, code, err := get1.EndpointReq()
		fmt.Println("body is %s", body)
		if len(body) > 2 || err != nil || code != http.StatusOK {
			fmt.Printf("get failed %d %v %d \n", code, err, len(body))
			r.HTML(200, "register", "Registration failed .Looks Like there is already an account with the same email id")
		} else {

			ctime := fmt.Sprintf("%v", time.Now().Unix())
			put1 := put.NewPutItem()
			put1.TableName = "users"
			put1.Item["username"] = &attributevalue.AttributeValue{S: register.Email}
			put1.Item["password"] = &attributevalue.AttributeValue{S: string(hashedPassword)}
			put1.Item["displayname"] = &attributevalue.AttributeValue{S: register.Nickname}
			put1.Item["timestamp"] = &attributevalue.AttributeValue{N: ctime}
			put1.Item["publickey"] = &attributevalue.AttributeValue{S: strconv.Itoa(rand.Intn(100000))}
			body1, code1, err1 := put1.EndpointReq()

			if err1 != nil || code != http.StatusOK {
				fmt.Printf("put failed %d %v %s\n", code1, err1, body1)
				r.HTML(200, "register", "Registration failed .Try again")

			} else {
				r.HTML(200, "index", "Registration Successfull.  You can login now")
				fmt.Printf("%v\n%v\n,%v\n", body1, code1, err1)
			}

		}

	})

	m.Get("/logout", func(r render.Render, response http.ResponseWriter) {
		// s.Delete("userId")
		clearSession(response)
		r.HTML(200, "index", nil)
	})

	m.Get("/chromelogout", func(r render.Render, response http.ResponseWriter) {
		// s.Delete("userId")
		clearSession(response)
		r.HTML(200, "chromeindex", nil)
	})

	m.Post("/messages", binding.Form(Msg{}), func(msg Msg, r render.Render, request *http.Request) {

		username := getUserName(request)
		present := validateUsername(username)
		if present {

			if insertToUserData(username, msg.Message, "text") {
				r.JSON(200, map[string]interface{}{"status": "success"})

			} else {
				r.JSON(200, map[string]interface{}{"status": "failure"})
			}
		} else {

			r.JSON(200, map[string]interface{}{"Access denied": "Unauthorized request"})
		}

	})

	m.Post("/urls", binding.Form(Msg{}), func(msg Msg, r render.Render, request *http.Request) {

		username := getUserName(request)
		present := validateUsername(username)
		if present {
			if insertToUserData(username, msg.Message, "url") {
				r.JSON(200, map[string]interface{}{"status": "success"})
			} else {
				r.JSON(200, map[string]interface{}{"status": "failure"})
			}
		} else {

			r.JSON(200, map[string]interface{}{"Access denied": "Unauthorized request"})
		}

	})

	m.Post("/groupdata/:id", binding.Form(Msg{}), func(msg Msg, params martini.Params, r render.Render, request *http.Request) {

		username := getUserName(request)
		present := validateUsername(username)
		if present {
			if insertToGroupData(params["id"], username, msg.Message, "text") {
				r.JSON(200, map[string]interface{}{"status": "success"})
				go sendNotificationsToAndroid(params["id"], username, msg.Message, "text")
			} else {
				r.JSON(200, map[string]interface{}{"status": "failure"})
			}
		} else {

			r.JSON(200, map[string]interface{}{"Access denied": "Unauthorized request"})
		}

	})

	m.Get("/deletegroup/:groupid", func(params martini.Params, r render.Render, request *http.Request) {

		username := getUserName(request)
		present := validateUsername(username)

		if present {
			decide := deleteGroup(params["groupid"], username)

			if decide == 0 {
				r.JSON(200, map[string]interface{}{"status": "Not able to Delete .Server error"})
			} else if decide == 1 {
				r.JSON(200, map[string]interface{}{"status": "Success"})
			} else if decide == 2 {
				r.JSON(200, map[string]interface{}{"status": "Not Authorized to delete"})
			}
		} else {
			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		}

	})

	m.Post("/groupdataurl/:id", binding.Form(Msg{}), func(msg Msg, params martini.Params, r render.Render, request *http.Request) {

		username := getUserName(request)
		present := validateUsername(username)
		if present {
			r.HTML(200, "user", nil)
			if insertToGroupData(params["id"], username, msg.Message, "url") {
				r.JSON(200, map[string]interface{}{"status": "success"})
				go sendNotificationsToAndroid(params["id"], username, msg.Message, "url")
			} else {
				r.JSON(200, map[string]interface{}{"status": "failure"})
			}
		} else {

			r.JSON(200, map[string]interface{}{"Access denied": "Unauthorized request"})
		}

	})

	m.Get("/usersingroup/:id", func(params martini.Params, r render.Render, request *http.Request) {

		userName := getUserName(request)
		present := validateUsername(userName)
		if present {

			msgs, decide := getUsersInGroupList(params["id"])

			if decide {

				r.JSON(200, map[string]interface{}{"status": "Success", "data": msgs})

			} else {
				r.JSON(200, map[string]interface{}{"status": "Access denied"})
			}
		} else {
			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		}

	})

	m.Get("/deleteuserdata/:timestamp", func(params martini.Params, r render.Render, request *http.Request) {

		userName := getUserName(request)
		present := validateUsername(userName)
		if present {
			decide := deleteuserdata(userName, params["timestamp"])

			if decide {
				r.JSON(200, map[string]interface{}{"status": "Success"})
			} else {
				r.JSON(200, map[string]interface{}{"status": "Access denied"})
			}
		} else {
			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		}

	})

	m.Get("/deletegroupdata/:groupid/:timestamp", func(params martini.Params, r render.Render, request *http.Request) {

		userName := getUserName(request)
		present := validateUsername(userName)

		if present {
			decide := deletegroupdata(userName, params["timestamp"], params["groupid"])

			if decide == 0 {
				r.JSON(200, map[string]interface{}{"status": "Not able to Delete .Server error"})
			} else if decide == 1 {
				r.JSON(200, map[string]interface{}{"status": "Success"})
			} else if decide == 2 {
				r.JSON(200, map[string]interface{}{"status": "Not Authorized to delete"})
			}
		} else {
			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		}

	})

	m.Get("/deleteuserfromgroup/:groupid/:userid", func(params martini.Params, r render.Render, request *http.Request) {

		userName := getUserName(request)
		present := validateUsername(userName)

		if present {
			decide := deleteuserfromgroup(params["groupid"], params["userid"])
			if decide {
				r.JSON(200, map[string]interface{}{"status": "Success"})
			} else {
				r.JSON(200, map[string]interface{}{"status": "Not able to Delete .Server error"})
			}
		} else {
			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		}

	})

	m.Get("/removemefromgroup/:groupid", func(params martini.Params, r render.Render, request *http.Request) {

		userName := getUserName(request)
		present := validateUsername(userName)

		if present {
			decide := deleteuserfromgroup(params["groupid"], userName)
			if decide {
				r.JSON(200, map[string]interface{}{"status": "Success"})
			} else {
				r.JSON(200, map[string]interface{}{"status": "Not able to Delete .Server error"})
			}
		} else {
			r.JSON(200, map[string]interface{}{"status": "Access denied"})
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

		userName := getUserName(request)

		msgs, decide := getGroupsList(userName)
		fmt.Println(msgs)

		if decide {

			r.JSON(200, map[string]interface{}{"status": "Success", "data": msgs})

		} else {
			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		}
	})

	m.Post("/creategroup", binding.Form(Msg{}), func(msg Msg, r render.Render, request *http.Request) {
		username := getUserName(request)
		present := validateUsername(username)
		if present {
			r.HTML(200, "user", nil)
			if createGroup(username, msg.Message) {
				r.JSON(200, map[string]interface{}{"status": "success"})
			} else {
				r.JSON(200, map[string]interface{}{"status": "failure"})
			}
		} else {

			r.JSON(200, map[string]interface{}{"Access denied": "Unauthorized request"})
		}
	})

	m.Get("/getgroupdatabyid/:id", func(r render.Render, params martini.Params, request *http.Request) {

		userName := getUserName(request)
		msgs, decide := getGroupDataById(userName, params["id"])

		if decide {

			r.JSON(200, map[string]interface{}{"status": "Success", "data": msgs})

		} else {
			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		}
	})

	m.Get("/getmessages", func(r render.Render, params martini.Params, request *http.Request) {
		userName := getUserName(request)

		msgs, decide := getUserData(userName)

		if decide {

			r.JSON(200, map[string]interface{}{"status": "Success", "data": msgs})

		} else {
			r.JSON(200, map[string]interface{}{"status": "Access denied"})
		}
	})

	notify = newNotify()

	m.Get("/sockets/:clientname/:cookie", sockets.JSON(Message{}), func(params martini.Params, request *http.Request, receiver <-chan *Message, sender chan<- *Message, done <-chan bool, disconnect chan<- int, err <-chan error) (int, string) {

		client := &Client{params["clientname"], receiver, sender, done, err, disconnect}

		//userName :=getUserNameFromCookie(params["cookie"])

		groups, _ := getGroupsList(params["cookie"])

		fmt.Println("groups length %v", len(groups))
		for _, group := range groups {

			r := notify.getGroup(group.Id)

			r.appendClient(params["cookie"], client)
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
	m.Post("/uploadfilefromurl", uploadurlHandler)
	m.Post("/groupuploadfile/:id", groupUploadHandler)
	m.Post("/groupuploadfilefromurl/:id", groupUploadUrlHandler)
	m.Run()

}

func decideServer(username string) []string {

	msgs, _ := getGroupSocketsList(username)

	var groupmap = make(map[string]int)
	var count int
	fmt.Println("msg values", msgs)
	for _, elem := range msgs {
		if elem.SocketServer != "" {
			count = count + 1
			v, ok := groupmap[elem.SocketServer]
			if ok {
				groupmap[elem.SocketServer] = v + 1
			} else {
				groupmap[elem.SocketServer] = 1
			}

		}
	}

	var isempty bool = true
	var values []int
	var servernames []string
	for k, v := range groupmap {
		isempty = false
		values = append(values, v)
		servernames = append(servernames, k)
	}
	sort.Ints(values)
	var topservername string

	fmt.Println("values", groupmap, servernames, topservername)
	if !isempty {

		top := values[len(values)-1]
		for k, v := range groupmap {
			if v == top {
				topservername = k
			}
		}
	} else {

		topservername = getServername()

	}

	for _, elem := range msgs {
		fmt.Println(elem)
		if elem.SocketServer == "" {
			up1 := update_item.NewUpdateItem()
			new_attr_val := topservername
			up1.TableName = "groups"
			up1.Key["group_id"] = &attributevalue.AttributeValue{S: elem.Id}
			up1.AttributeUpdates = attributevalue.NewAttributeValueUpdateMap()
			up_avu := attributevalue.NewAttributeValueUpdate()
			up_avu.Action = update_item.ACTION_PUT
			up_avu.Value = &attributevalue.AttributeValue{S: new_attr_val}
			up1.AttributeUpdates["socketserver"] = up_avu
			up1.ReturnValues = update_item.RETVAL_ALL_NEW
			update_item_json, update_item_json_err := json.Marshal(up1)
			if update_item_json_err != nil {
				fmt.Printf("%v\n", update_item_json_err)

			}
			fmt.Printf("%s\n", string(update_item_json))

			body, code, err := up1.EndpointReq()
			if err != nil || code != http.StatusOK {
				fmt.Printf("update item failed %d %v %s\n", code, err, body)

			}
			fmt.Printf("%v\n%v\n,%v\n", body, code, err)
			fmt.Println("no scoket server")
		}
	}
	updateServerCount(count, topservername)
	return servernames
}

func getServername() string {

	get1 := get.NewGetItem()
	get1.TableName = "groupdata"
	get1.Key["group_id"] = &attributevalue.AttributeValue{S: "socketconnections"}
	get1.Key["timestamp"] = &attributevalue.AttributeValue{S: "1024"}
	body, code, err := get1.EndpointReq()
	fmt.Println("body is %s", body)
	if len(body) <= 2 || err != nil || code != http.StatusOK {
		fmt.Printf("get failed %d %v %d \n", code, err, len(body))
	} else {
		payload := []byte(body)
		var result map[string]interface{}
		if err := json.Unmarshal(payload, &result); err != nil {
			panic(err)
		}
		row := result["Item"].(map[string]interface{})
		fmt.Println(row)
		server1 := row["server1"].(map[string]interface{})
		server2 := row["server2"].(map[string]interface{})
		server1count, _ := strconv.Atoi(server1["S"].(string))
		server2count, _ := strconv.Atoi(server2["S"].(string))
		if server1count > server2count {
			return "ec2-54-149-41-188.us-west-2.compute.amazonaws.com"
		} else {
			return "ec2-54-149-42-11.us-west-2.compute.amazonaws.com"
		}
	}

	return "ec2-54-149-41-188.us-west-2.compute.amazonaws.com"
}

func updateServerCount(count int, topservername string) {

	var servername string
	var newservercount int
	if topservername == "ec2-54-149-41-188.us-west-2.compute.amazonaws.com" {
		servername = "server1"

	} else {
		servername = "server2"
	}

	get1 := get.NewGetItem()
	get1.TableName = "groupdata"
	get1.Key["group_id"] = &attributevalue.AttributeValue{S: "socketconnections"}
	get1.Key["timestamp"] = &attributevalue.AttributeValue{S: "1024"}
	body, code, err := get1.EndpointReq()
	fmt.Println("body is %s", body)
	if len(body) <= 2 || err != nil || code != http.StatusOK {
		fmt.Printf("get failed %d %v %d \n", code, err, len(body))
	} else {
		payload := []byte(body)
		var result map[string]interface{}
		if err := json.Unmarshal(payload, &result); err != nil {
			panic(err)
		}
		row := result["Item"].(map[string]interface{})
		server := row[servername].(map[string]interface{})

		servercount, _ := strconv.Atoi(server["S"].(string))
		newservercount = servercount + count

		hk := "socketconnections"
		rk := "1024"

		up1 := update_item.NewUpdateItem()

		up1.TableName = "groupdata"
		up1.Key["group_id"] = &attributevalue.AttributeValue{S: hk}
		up1.Key["timestamp"] = &attributevalue.AttributeValue{S: rk}

		up1.AttributeUpdates = attributevalue.NewAttributeValueUpdateMap()
		up_avu := attributevalue.NewAttributeValueUpdate()
		up_avu.Action = update_item.ACTION_PUT
		up_avu.Value = &attributevalue.AttributeValue{S: strconv.Itoa(newservercount)}
		up1.AttributeUpdates["server1"] = up_avu

		update_item_json, update_item_json_err := json.Marshal(up1)

		if update_item_json_err != nil {
			fmt.Printf("%v\n", update_item_json_err)

		}
		fmt.Printf("%s\n", string(update_item_json))

		body, code, err := up1.EndpointReq()
		if err != nil || code != http.StatusOK {
			fmt.Printf("update item failed %d %v %s\n", code, err, body)

		}

	}

}

func getGroupSocketsList(username string) ([]GroupSocketsJson, bool) {
	msgs, _ := getGroupsList(username)
	p := make([]GroupSocketsJson, 0)

	fmt.Println("msg values", msgs)
	for _, elem := range msgs {
		get1 := get.NewGetItem()
		get1.TableName = "groups"
		get1.Key["group_id"] = &attributevalue.AttributeValue{S: elem.Id}
		body, code, err := get1.EndpointReq()
		fmt.Println("body is %s", body)
		if len(body) <= 2 || err != nil || code != http.StatusOK {
			fmt.Printf("get failed %d %v %d \n", code, err, len(body))
		} else {
			payload := []byte(body)
			var result map[string]interface{}
			if err := json.Unmarshal(payload, &result); err != nil {
				panic(err)
			}
			row := result["Item"].(map[string]interface{})
			fmt.Println(row)
			socketserver := row["socketserver"].(map[string]interface{})

			p = append(p, GroupSocketsJson{elem.Id, socketserver["S"].(string)})

		}
	}
	return p, true
}
func sendNotificationsToAndroid(groupid string, username string, content string, content_type string) {

	fmt.Println("calling notification to android")

	usersingroup, _ := getUsersInGroupList(groupid)
	fmt.Println("got the follwing users from list ", usersingroup)
	var regIDs []string
	for _, elem := range usersingroup {

		// You can increase efficiency here in future

		fmt.Println("Username in forloop ", elem.Username)
		get1 := get.NewGetItem()
		get1.TableName = "users"
		get1.Key["username"] = &attributevalue.AttributeValue{S: elem.Username + "android"}
		body, code, err := get1.EndpointReq()

		if err != nil || code != http.StatusOK {
			fmt.Println("error at 1st get request")
			//return 0
		} else {
			var res AndroidUserDetails
			err3 := json.Unmarshal([]byte(body), &res)
			fmt.Println("Body for each get request", body)
			fmt.Println("res for each get request", res)
			if err3 != nil {
				//	return 0
			}

			if res.Item.Registration_id.S != "" {
				regIDs = append(regIDs, res.Item.Registration_id.S)
			}
			//regIDs := []string{"APA91bHdXL9Huqm23usTUR0w_xcdliT_jcNh_fA8zeSEa3f4KH7LLMjN_EORCMGH9VPimWyQ8azF2_1TFjbcq-NEtqW0FD6pLrjgfSA2mgHTX8SgSUw3-1m0MIjY8VidvKOdT-fIW_HMcj4SuI43r33pHewTvbvnqj_tlHkphKtokG5ymSUxxgk"}

		}
	}
	data := map[string]interface{}{"content_type": content_type, "content": content}
	msg := gcm.NewMessage(data, regIDs...)

	// Create a Sender to send the message.
	sender := &gcm.Sender{ApiKey: "AIzaSyCq83QZc-yl8pblJJOKtPRQ-rfbvrwg6ms"}

	// Send the message and receive the response after at most two retries.
	response, err := sender.Send(msg, 2)
	if err != nil {
		fmt.Println("Failed to send message:", err)
		return
	}

	fmt.Println(response)

}

func deleteuserdata(username string, timestamp string) bool {

	// DELETE AN ITEM
	del_item1 := delete_item.NewDeleteItem()
	del_item1.TableName = "userdata"
	del_item1.Key["username"] = &attributevalue.AttributeValue{S: username}
	del_item1.Key["timestamp"] = &attributevalue.AttributeValue{S: timestamp}

	body, code, err := del_item1.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("fail delete %d %v %s\n", code, err, body)

		return false
	}

	return true

}
func deletegroupdata(username string, timestamp string, groupid string) int {

	get1 := get.NewGetItem()
	get1.TableName = "groupdata"
	get1.Key["group_id"] = &attributevalue.AttributeValue{S: groupid}
	get1.Key["timestamp"] = &attributevalue.AttributeValue{S: timestamp}
	body, code, err := get1.EndpointReq()

	if err != nil || code != http.StatusOK {
		return 0
	} else {
		var res groupdataitem
		err3 := json.Unmarshal([]byte(body), &res)

		if err3 != nil {
			return 0
		}
		fmt.Println("usenmaesi ", res.Item.Username.S, " ", username)
		if res.Item.Username.S == username {
			// DELETE AN ITEM
			del_item1 := delete_item.NewDeleteItem()
			del_item1.TableName = "groupdata"
			del_item1.Key["group_id"] = &attributevalue.AttributeValue{S: groupid}
			del_item1.Key["timestamp"] = &attributevalue.AttributeValue{S: timestamp}
			body1, code1, err1 := del_item1.EndpointReq()

			if err1 != nil || code1 != http.StatusOK {
				fmt.Printf("fail delete %d %v %s\n", code1, err1, body1)

				return 0
			} else {
				return 1
			}

		} else {
			return 2
		}
	}
	return 0
}

func deleteGroup(groupid string, username string) int {

	get1 := get.NewGetItem()
	get1.TableName = "usergroups"
	get1.Key["group_id"] = &attributevalue.AttributeValue{S: groupid}
	get1.Key["username"] = &attributevalue.AttributeValue{S: username}
	body, code, err := get1.EndpointReq()

	if err != nil || code != http.StatusOK {
		fmt.Println("error at 1st get request")
		return 0
	} else {
		var res groupdetail
		err3 := json.Unmarshal([]byte(body), &res)

		if err3 != nil {
			return 0
		}

		if res.Item.Group_admin.S == username {

			usersingroup, _ := getUsersInGroupList(groupid)

			for _, elem := range usersingroup {

				// DELETE AN ITEM

				del_item1 := delete_item.NewDeleteItem()
				del_item1.TableName = "usergroups"
				del_item1.Key["group_id"] = &attributevalue.AttributeValue{S: groupid}

				del_item1.Key["username"] = &attributevalue.AttributeValue{S: elem.Username}

				body1, code1, err1 := del_item1.EndpointReq()

				if err1 != nil || code1 != http.StatusOK {
					fmt.Printf("fail delete %d %v %s\n", code1, err1, body1)
					return 0
				} else {

					fmt.Println("remvoing user", elem.Username)
					del_item1 := delete_item.NewDeleteItem()
					del_item1.TableName = "groupusers"
					del_item1.Key["group_id"] = &attributevalue.AttributeValue{S: groupid}

					del_item1.Key["username"] = &attributevalue.AttributeValue{S: elem.Username}

					body1, code1, err1 := del_item1.EndpointReq()

					if err1 != nil || code1 != http.StatusOK {
						fmt.Printf("fail delete %d %v %s\n", code1, err1, body1)
						return 0
					}
				}
			}

			return 1
		} else {
			return 2
		}
	}
	return 0
}

func deleteuserfromgroup(groupid string, userid string) bool {

	del_item1 := delete_item.NewDeleteItem()
	del_item1.TableName = "groupusers"
	del_item1.Key["group_id"] = &attributevalue.AttributeValue{S: groupid}
	del_item1.Key["username"] = &attributevalue.AttributeValue{S: userid}
	_, code, err := del_item1.EndpointReq()

	if err != nil || code != http.StatusOK {
		return false
	} else {

		// DELETE AN ITEM
		del_item2 := delete_item.NewDeleteItem()
		del_item2.TableName = "usergroups"
		del_item2.Key["group_id"] = &attributevalue.AttributeValue{S: groupid}
		del_item2.Key["username"] = &attributevalue.AttributeValue{S: userid}
		body1, code1, err1 := del_item2.EndpointReq()

		if err1 != nil || code1 != http.StatusOK {
			fmt.Printf("fail delete %d %v %s\n", code1, err1, body1)

			return false
		} else {
			return true
		}

	}
}
func getUserData(username string) ([]UserDataJson, bool) {

	q := query.NewQuery()
	q.TableName = "userdata"
	q.Select = ep.SELECT_ALL
	kc := condition.NewCondition()
	kc.AttributeValueList = make([]*attributevalue.AttributeValue, 1)
	kc.AttributeValueList[0] = &attributevalue.AttributeValue{S: username}
	kc.ComparisonOperator = query.OP_EQ
	var t bool
	z := &t
	*z = false
	q.ScanIndexForward = z
	q.Limit = 10000
	q.KeyConditions["username"] = kc

	body, code, err := q.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("query failed %d %v %s\n", code, err, body)
		return nil, false
	}

	var res userdata
	//str2:=body
	err2 := json.Unmarshal([]byte(body), &res)

	if err2 != nil {
		fmt.Println(err2)
	}

	p := make([]UserDataJson, 0)
	for _, elem := range res.Items {

		p = append(p, UserDataJson{elem.Content.S, elem.Content_type.S, elem.Timestamp.S})
	}
	return p, true

}

func getGroupDataById(username string, id string) ([]GroupdataJson, bool) {

	q := query.NewQuery()
	q.TableName = "groupdata"
	q.Select = ep.SELECT_ALL
	kc := condition.NewCondition()
	kc.AttributeValueList = make([]*attributevalue.AttributeValue, 1)
	kc.AttributeValueList[0] = &attributevalue.AttributeValue{S: id}
	kc.ComparisonOperator = query.OP_EQ
	var t bool
	z := &t
	*z = false
	q.ScanIndexForward = z
	q.Limit = 10000
	q.KeyConditions["group_id"] = kc
	k, _ := json.Marshal(q)
	fmt.Printf("JSON:%s\n", string(k))
	body, code, err := q.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("query failed %d %v %s\n", code, err, body)
		return nil, false
	}
	fmt.Printf("%v\n%v\n%v\n", body, code, err)

	var res groupdata
	//str2:=body
	err2 := json.Unmarshal([]byte(body), &res)

	if err2 != nil {
		fmt.Println(err2)
	}

	p := make([]GroupdataJson, 0)
	for _, elem := range res.Items {

		fmt.Println(elem.Username.S)
		p = append(p, GroupdataJson{elem.Groupdata_id.S, elem.Content.S, elem.Content_type.S, elem.Username.S, elem.Timestamp.S})
	}
	return p, true

}

func getUsersInGroupList(id string) ([]GroupusernamesJson, bool) {

	q := query.NewQuery()
	q.TableName = "groupusers"
	q.Select = ep.SELECT_ALL
	kc := condition.NewCondition()
	kc.AttributeValueList = make([]*attributevalue.AttributeValue, 1)
	kc.AttributeValueList[0] = &attributevalue.AttributeValue{S: id}
	kc.ComparisonOperator = query.OP_EQ
	q.Limit = 10000
	q.KeyConditions["group_id"] = kc
	k, _ := json.Marshal(q)
	fmt.Printf("JSON:%s\n", string(k))
	body, code, err := q.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("query failed %d %v %s\n", code, err, body)
		return nil, false
	}
	fmt.Printf("%v\n%v\n%v\n", body, code, err)

	var res groupusersList

	//str2:=body
	err2 := json.Unmarshal([]byte(body), &res)

	if err2 != nil {
		fmt.Println(err2)
	}

	arr := make([]GroupusernamesJson, 0)
	for _, elem := range res.Items {

		arr = append(arr, GroupusernamesJson{elem.Group_id.S, elem.Username.S, elem.Group_admin.N, elem.Timestamp.N})
	}

	return arr, true
}

func getGroupsList(name string) ([]GroupnamesJson, bool) {

	q := query.NewQuery()
	q.TableName = "usergroups"
	q.Select = ep.SELECT_ALL
	kc := condition.NewCondition()
	kc.AttributeValueList = make([]*attributevalue.AttributeValue, 1)
	kc.AttributeValueList[0] = &attributevalue.AttributeValue{S: name}
	kc.ComparisonOperator = query.OP_EQ
	q.Limit = 10000
	q.KeyConditions["username"] = kc
	k, _ := json.Marshal(q)
	fmt.Printf("JSON:%s\n", string(k))
	body, code, err := q.EndpointReq()
	if err != nil || code != http.StatusOK {
		fmt.Printf("query failed %d %v %s\n", code, err, body)
		return nil, false
	}
	fmt.Printf("%v\n%v\n%v\n", body, code, err)

	var res groupsList
	//str2:=body
	err2 := json.Unmarshal([]byte(body), &res)

	if err2 != nil {
		fmt.Println(err2)
	}
	p := make([]GroupnamesJson, 0)

	for _, elem := range res.Items {

		p = append(p, GroupnamesJson{elem.Group_id.S, elem.Username.S, elem.Group_name.S, elem.Group_admin.S, elem.Timestamp.N})
	}
	return p, true
}
func createGroup(username string, groupname string) bool {

	ctime := fmt.Sprintf("%v", time.Now().Unix())
	groupid := uuid.New()
	put1 := put.NewPutItem()
	put1.TableName = "groups"
	put1.Item["group_id"] = &attributevalue.AttributeValue{S: groupid}
	put1.Item["group_name"] = &attributevalue.AttributeValue{S: groupname}
	put1.Item["group_admin"] = &attributevalue.AttributeValue{S: username}
	put1.Item["timestamp"] = &attributevalue.AttributeValue{N: ctime}

	body1, code1, err1 := put1.EndpointReq()
	if err1 != nil || code1 != http.StatusOK {
		fmt.Printf("put failed %d %v %s\n", code1, err1, body1)
		return false
	} else {

		put2 := put.NewPutItem()
		put2.TableName = "usergroups"
		put2.Item["username"] = &attributevalue.AttributeValue{S: username}
		put2.Item["group_id"] = &attributevalue.AttributeValue{S: groupid}
		put2.Item["group_name"] = &attributevalue.AttributeValue{S: groupname}
		put2.Item["group_admin"] = &attributevalue.AttributeValue{S: username}
		put2.Item["timestamp"] = &attributevalue.AttributeValue{N: ctime}

		body2, code2, err2 := put2.EndpointReq()
		if err2 != nil || code2 != http.StatusOK {
			fmt.Printf("put failed %d %v %s\n", code2, err2, body2)
			return false
		} else {
			put3 := put.NewPutItem()
			put3.TableName = "groupusers"
			put3.Item["group_id"] = &attributevalue.AttributeValue{S: groupid}

			put3.Item["username"] = &attributevalue.AttributeValue{S: username}
			put3.Item["timestamp"] = &attributevalue.AttributeValue{N: ctime}
			put3.Item["group_admin"] = &attributevalue.AttributeValue{N: "1"}

			body3, code3, err3 := put3.EndpointReq()
			if err3 != nil || code3 != http.StatusOK {
				fmt.Printf("put failed 4%d %v %s\n", code3, err3, body3)
				return false
			} else {
				return true
			}
		}

	}
}
func validateUsername(name string) bool {
	get1 := get.NewGetItem()
	get1.TableName = "users"
	get1.Key["username"] = &attributevalue.AttributeValue{S: name}

	body, code, err := get1.EndpointReq()

	if len(body) <= 2 || err != nil || code != http.StatusOK {
		return false
	} else {
		return true
	}
}

func insertToUserData(username string, content string, ctype string) bool {

	ctime := fmt.Sprintf("%v", time.Now().UnixNano())
	put1 := put.NewPutItem()
	put1.TableName = "userdata"
	put1.Item["username"] = &attributevalue.AttributeValue{S: username}
	put1.Item["timestamp"] = &attributevalue.AttributeValue{S: ctime}
	put1.Item["content"] = &attributevalue.AttributeValue{S: content}
	put1.Item["content_type"] = &attributevalue.AttributeValue{S: ctype}

	body, code, err := put1.EndpointReq()

	if err != nil || code != http.StatusOK {
		fmt.Printf("put failed %d %v %s\n", code, err, body)
		return false
	} else {
		return true

	}

}

func insertToGroupData(groupid string, username string, content string, ctype string) bool {

	ctime := fmt.Sprintf("%v", time.Now().UnixNano())
	put1 := put.NewPutItem()
	put1.TableName = "groupdata"
	put1.Item["group_id"] = &attributevalue.AttributeValue{S: groupid}
	put1.Item["timestamp"] = &attributevalue.AttributeValue{S: ctime}
	put1.Item["username"] = &attributevalue.AttributeValue{S: username}
	put1.Item["content"] = &attributevalue.AttributeValue{S: content}
	put1.Item["content_type"] = &attributevalue.AttributeValue{S: ctype}

	body, code, err := put1.EndpointReq()

	if err != nil || code != http.StatusOK {
		fmt.Printf("put failed %d %v %s\n", code, err, body)
		return false
	} else {
		return true
	}

}

func uploadHandler(w http.ResponseWriter, r *http.Request, rend render.Render) {

	// write the content from POST to the file

	userName := getUserName(r)

	present := validateUsername(userName)

	if userName == "" {
		rend.JSON(200, map[string]interface{}{"status": "Unauthorized request"})
	}
	if present {
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
		if insertToUserData(userName, header.Filename, "file") {
			rend.JSON(200, map[string]interface{}{"status": "success"})
		} else {
			rend.JSON(200, map[string]interface{}{"status": "failure"})
		}
	} else {

		rend.JSON(200, map[string]interface{}{"Access denied": "Unauthorized request"})
	}

}

func uploadurlHandler(w http.ResponseWriter, r *http.Request, rend render.Render) {

	// write the content from POST to the fil
	userName := getUserName(r)

	present := validateUsername(userName)

	if userName == "" {
		rend.JSON(200, map[string]interface{}{"status": "Unauthorized request"})
	}
	if present {
		r.ParseForm()
		fmt.Println(r.Form)
		fmt.Println(r.FormValue("url"))

		resp, err := http.Get(r.FormValue("url"))
		defer resp.Body.Close()

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		ctime := fmt.Sprintf("%v", time.Now().Unix())
		filename := ctime + "" + r.FormValue("type")

		l, _ := s3util.Create("https://s3-us-west-2.amazonaws.com/megh-uploads/"+filename, nil, nil)
		_, err2 := io.Copy(l, resp.Body)
		l.Close()
		if err2 != nil {
			fmt.Fprintln(w, err)
		}
		if insertToUserData(userName, filename, "file") {
			rend.JSON(200, map[string]interface{}{"status": "success"})
		} else {
			rend.JSON(200, map[string]interface{}{"status": "failure"})
		}
	} else {

		rend.JSON(200, map[string]interface{}{"Access denied": "Unauthorized request"})
	}

}
func groupUploadHandler(w http.ResponseWriter, r *http.Request, rend render.Render, params martini.Params) {

	userName := getUserName(r)

	present := validateUsername(userName)

	if userName == "" {
		rend.JSON(200, map[string]interface{}{"status": "Unauthorized request"})
	}
	if present {
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

		if insertToGroupData(params["id"], userName, header.Filename, "file") {
			go sendNotificationsToAndroid(params["id"], userName, header.Filename, "file")
			rend.JSON(200, map[string]interface{}{"status": "success"})
		} else {
			rend.JSON(200, map[string]interface{}{"status": "failure"})
		}
	} else {

		rend.JSON(200, map[string]interface{}{"Access denied": "Unauthorized request"})
	}

}

func groupUploadUrlHandler(w http.ResponseWriter, r *http.Request, rend render.Render, params martini.Params) {

	// write the content from POST to the fil
	userName := getUserName(r)

	if userName == "" {
		rend.JSON(200, map[string]interface{}{"status": "Unauthorized request"})
	}
	present := validateUsername(userName)

	if present {
		r.ParseForm()
		fmt.Println(r.Form)
		fmt.Println(r.FormValue("url"))

		resp, err := http.Get(r.FormValue("url"))
		defer resp.Body.Close()

		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		ctime := fmt.Sprintf("%v", time.Now().Unix())
		filename := ctime + "" + r.FormValue("type")

		l, _ := s3util.Create("https://s3-us-west-2.amazonaws.com/megh-uploads/"+filename, nil, nil)
		_, err2 := io.Copy(l, resp.Body)
		l.Close()
		if err2 != nil {
			fmt.Fprintln(w, err)
		}
		if insertToGroupData(params["id"], userName, filename, "file") {
			go sendNotificationsToAndroid(params["id"], userName, filename, "file")
			rend.JSON(200, map[string]interface{}{"status": "success"})
		} else {
			rend.JSON(200, map[string]interface{}{"status": "failure"})
		}
	} else {

		rend.JSON(200, map[string]interface{}{"Access denied": "Unauthorized request"})
	}

}

func setSession(userName string, publickey string, response http.ResponseWriter, servernames []string) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		cookie2 := &http.Cookie{
			Name:  "publickey",
			Value: publickey,
			Path:  "/",
		}

		cookie3 := &http.Cookie{
			Name:  "username",
			Value: userName,
			Path:  "/",
		}
		if servernames != nil {
			var finalservername string
			for _, element := range servernames {
				finalservername = finalservername + ";" + element
			}
			fmt.Println("finalservername is ", finalservername, servernames)
			cookie4 := &http.Cookie{
				Name:  "servernames",
				Value: finalservername,
				Path:  "/",
			}
			http.SetCookie(response, cookie4)
		}
		http.SetCookie(response, cookie)
		http.SetCookie(response, cookie2)
		http.SetCookie(response, cookie3)

	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	cookie2 := &http.Cookie{
		Name:   "publickey",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	cookie3 := &http.Cookie{
		Name:   "username",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	cookie4 := &http.Cookie{
		Name:   "servernames",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
	http.SetCookie(response, cookie2)
	http.SetCookie(response, cookie3)
	http.SetCookie(response, cookie4)

}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		fmt.Println("cookies", cookie, cookie.Value)
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func getUserNameFromCookie(username string) (userName string) {

	cookieValue := make(map[string]string)
	var user string

	if err := cookieHandler.Decode("session", username, &cookieValue); err == nil {
		user = cookieValue["name"]
	}
	fmt.Println("username of the group", user)

	return user
}
