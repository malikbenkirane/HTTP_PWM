package main



import(
	"fmt"
	"net/http"
	"log"
	"strconv"
	"html/template"
	"github.com/gorilla/securecookie"
	"./sqlite"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

//Cookies handling Functions 
func setSession(username string , w http.ResponseWriter){
	value := map[string]string{
		"name":username,
	}
	if encoded, err := cookieHandler.Encode("session",value); err == nil{
		cookie := &http.Cookie{
			Name :"session",
			Value : encoded,
			Path: "/",
			MaxAge : 3600,
		}
		http.SetCookie(w,cookie)
	}
}

func getUserName(r *http.Request) (username string){
	if cookie , err := r.Cookie("session"); err == nil{
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session",cookie.Value, &cookieValue); err == nil{
			username = cookieValue["name"]
		}
	}
	return username
}


func clearSession(r http.ResponseWriter){
	cookie := &http.Cookie{
		Name : "session",
		Value : "",
		Path : "/" ,
		MaxAge : -1,
	}
	http.SetCookie(r,cookie)
}

// End of cookies functions 


type USER struct{
	Username string
	Ok bool
	Updated bool
	Id string
	User string
	Password string
	Category string
}

type Userlist []USER

var(

	finallist Userlist
	p = USER{}
)



func home(w http.ResponseWriter, r *http.Request){
	fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
	if r.URL.Path != "/"{
		http.Error(w,"404 Page Not Found",http.StatusNotFound)
		fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
		return
	}
	user := getUserName(r)
	u := USER{Username: user}

	t,_ := template.ParseFiles("home.html")
	t.Execute(w,u)
}




//register Func
func registerhandle(w http.ResponseWriter, r *http.Request){
	fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
	if r.URL.Path != "/register/"{
		http.Error(w,"404 Page Not Found",http.StatusNotFound)
	}

	if r.Method == "GET"{
		if cookie,_ := r.Cookie("session");cookie != nil {
			fmt.Fprintf(w,"logout to register")
		}else{
			t,_ := template.ParseFiles("register.html")
			t.Execute(w,nil)
		}
	}else{
		r.ParseForm()
		user := r.FormValue("user")
		passw := r.FormValue("passw")
		email := r.FormValue("email")
		
		if reg := sqlite.Register(user,passw,email); reg != 1 {
			http.Redirect(w,r,"/login/",http.StatusFound)
		}
	}
}




func loginhandle(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/login/"{
		http.Error(w,"404 Page Not Found",http.StatusNotFound)
		fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
	}


	if r.Method == "POST"{
		fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
		r.ParseForm()
		user := r.FormValue("user")
		passw := r.FormValue("passw")

		if user != "" && passw != ""{
			islogin := sqlite.Login(user,passw)
			if islogin {
				setSession(user,w) // setting session for the user
				http.Redirect(w,r,"/",http.StatusFound) // redirect user to Home page
			}else{
				fmt.Fprintf(w,"username or password is incorrect!\n")
			}
		}else{
			fmt.Println(user,passw)
			fmt.Fprintf(w,"args missing!\n")
		}

	}else if r.Method == "GET"{
		fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
		if cookie ,_ := r.Cookie("session"); cookie != nil{
			fmt.Fprintf(w,"You are Already logged-in")
		}else{
			t,_ := template.ParseFiles("login.html")
			t.Execute(w,nil)
		}
		
	}
}


func logouthandle(w http.ResponseWriter, r *http.Request){
	fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
	clearSession(w)
	fmt.Fprintf(w,"You were logged out")
}



func addElement(w http.ResponseWriter , r *http.Request){
	if r.URL.Path != "/add/"{
		http.Error(w,"404 Page Not Found",http.StatusNotFound)
		fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
		return
	}

	cookie , _ := r.Cookie("session")
	switch req := r.Method ;req {
		case "GET":
			if cookie == nil {
				fmt.Fprintf(w,"login to save passwords")
			}else{
				user := getUserName(r)
				p.Username = user
				t, _ := template.ParseFiles("add.html")
				t.Execute(w,p)
				p.Ok = false 
				
			}

		case "POST":
			fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
			r.ParseForm()
			user := r.Form.Get("user")
			passw := r.Form.Get("passw")
			category := r.Form.Get("category")
			username := getUserName(r)
			id := sqlite.Getuid(username)

			if ok := sqlite.Save(user,passw,category,id) ; ok{
				p.Ok = ok
				http.Redirect(w,r,"/add/",http.StatusFound)
			}
	}
}




func showhandle(w http.ResponseWriter , r *http.Request){
	fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)

	if r.URL.Path != "/show/"{
		http.Error(w,"404 Page Not Found", http.StatusNotFound)
		fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
	}
	cookie ,_ := r.Cookie("session")
	user := getUserName(r) // get username from session 

	switch req := r.Method ; req {
		case "GET":
			if cookie == nil{
				fmt.Fprintf(w,"login to show your stuff")
			}else{
				p.Username = user
				t,_ := template.ParseFiles("show.html")
				t.Execute(w,finallist)
			}
		
		case "POST":
			r.ParseForm()
			category := r.Form.Get("category")
			id := sqlite.Getuid(user)
			if category != "" {
				i,u,P,c := sqlite.GetStuff(id,category)
				finallist = finallist[:0]
				for x := 0 ; x<len(i);x++{
					p.Id = i[x]
					p.User = u[x]
					p.Password = P[x]
					p.Category = c[x]
					finallist = append(finallist,p)
				}
				
				http.Redirect(w,r,"/show/",http.StatusFound)
			}else{
				i,u,P,c := sqlite.GetStuff(id,category)
				finallist = finallist[:0]
				for x := 0 ; x<len(i);x++{
					p.Id = i[x]
					p.User = u[x]
					p.Password = P[x]
					p.Category = c[x]
					finallist = append(finallist,p)
				}
				
				http.Redirect(w,r,"/show/",http.StatusFound)
			}	
	}

}


func updatehandle(w http.ResponseWriter, r *http.Request){ 
	fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
	if r.URL.Path != "/update/"{
		http.Error(w,"404 Page Not Found", http.StatusNotFound)
		fmt.Printf("[%s] %s\n",r.Method,r.URL.Path)
	}
	cookie ,_ := r.Cookie("session")
	user := getUserName(r)
	switch req := r.Method ; req {
		case "GET":
			if cookie == nil{
				fmt.Fprintf(w,"login to update your stuff")
			}else{
				p.Username = user
				t,_ := template.ParseFiles("update.html")
				t.Execute(w,p)
				p.Updated = false 
			}
		case "POST":
			r.ParseForm()
			id,err := strconv.ParseInt(r.Form.Get("id")[0:], 10, 0)
			if err != nil{
				panic(err)
			}
			Nuser := r.Form.Get("user")
			Npassw := r.Form.Get("passw")
			Ncatg := r.Form.Get("catg")

			var ok bool = false
			if x := sqlite.Update(int(id),Nuser,Npassw,Ncatg) ; len(x) != 0 {

				for i:=0;i<len(x);i++{
					if x[i] != 0 && i == 0{
						ok = true
						p.Updated = ok
					}else if x[i] != 0 && i == 1{
						ok = true
						p.Updated = ok
					}else if x[i] != 0 && i == 2 {
						ok = true
						p.Updated = ok
					}
				}
				http.Redirect(w,r,"/update/",http.StatusFound)
			}
		}
}



func main(){
	
	http.HandleFunc("/",home) 
	http.HandleFunc("/add/",addElement)
	http.HandleFunc("/login/",loginhandle)
	http.HandleFunc("/register/",registerhandle)
	http.HandleFunc("/logout/",logouthandle)
	http.HandleFunc("/show/",showhandle)
	http.HandleFunc("/update/",updatehandle)

	fmt.Println("Running http server on http://localhost:8080 ")
	if err := http.ListenAndServe(":8080", nil) ; err != nil{
		log.Fatal(err)
	}
}

