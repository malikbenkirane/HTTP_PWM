package sqlite 



import (
	"strings"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"crypto/sha256"
	"encoding/base32"
	
)



const DB string = "./database.db"


func checkError(err error){
	if err != nil{
		panic(err)
	}
}


func CheckForUser(user string) int{
	db, err := sql.Open("sqlite3",DB)
	checkError(err)

	rows , err := db.Query("SELECT username FROM users")
	checkError(err)

	var username string
	var users []string

	for rows.Next(){
		err = rows.Scan(&username)
		checkError(err)
		if user == username{
			users = append(users,username)
		}
	}
	return len(users)
}



func Register(user string , passw string , email string) int{
	db, err := sql.Open("sqlite3",DB)
	checkError(err)
	
	var rslt int

	if isexist := CheckForUser(user); isexist == 1{
		rslt = 1
		db.Close()
	}else{
		//Hash the password then save it into the db 
		h := sha256.New()
		h.Write([]byte(passw))
		ph := h.Sum(nil)
		phx := strings.ToLower(base32.HexEncoding.EncodeToString(ph))
		stmt , err := db.Prepare("INSERT INTO users(username,password,email) VALUES (?,?,?)")
		stmt.Exec(user,phx,email)
		checkError(err)
		rslt = 0
		db.Close()
	}
	return rslt
}


func Login(user string , passw string) bool{
	db, err := sql.Open("sqlite3",DB)
	checkError(err)

	rows , err := db.Query("SELECT * FROM users")
	checkError(err)
 	
	var username string 
	var password string
	var email string
	var id int
	var result bool

	h := sha256.New()
	h.Write([]byte(passw))
	ph := h.Sum(nil)
	phx := strings.ToLower(base32.HexEncoding.EncodeToString(ph))

	
	for rows.Next(){
		err = rows.Scan(&id,&username,&password,&email)
		checkError(err)
		if username == user && password == phx{
			result = true
		}else{
			result = false
		}
	}

	rows.Close()
	return result
}


func GetStuff(uid int, category string) ([]string,[]string,[]string,[]string){
	db, err := sql.Open("sqlite3",DB)
	checkError(err)

	rows , err := db.Query("SELECT pwid ,user,passw,category,userid FROM passwords")
	checkError(err)

	var userid int
	var pwid string
	var user string
	var passw string
	var catg string
	var i[]string
	var u[]string
	var p[]string
	var c[]string

	if category == ""{
		for rows.Next(){
			err = rows.Scan(&pwid,&user,&passw,&catg,&userid)
			checkError(err)
			if uid == userid {
				i = append(i,pwid)
				u = append(u,user)
				p = append(p,passw)
				c = append(c,catg)
			}
		}
		return i,u,p,c
		rows.Close()
	}else if category != ""{
		for rows.Next(){
			err = rows.Scan(&pwid,&user,&passw,&catg,&userid)
			checkError(err)
			if uid == userid && category == catg{
				i = append(i,pwid)
				u = append(u,user)
				p = append(p,passw)
				c = append(c,catg)
			}
		}
		return i,u,p,c
		rows.Close()
	}
	return i,u,p,c
}


func Getuid(user string) int{
	db, err := sql.Open("sqlite3",DB)
	checkError(err)

	rows , err := db.Query("SELECT * FROM users")
	checkError(err)

	var uid int 
	var username string 
	var password string
	var email string 
	var id int 
	for rows.Next(){
		err = rows.Scan(&uid,&username,&password,&email)
		checkError(err)
		if username == user{
			id =  uid
		}
	}
	rows.Close()
	return id
}


func Update(id int,args ...string) []int{
	db, err := sql.Open("sqlite3",DB)
	checkError(err)

	defer db.Close()
	var f int
	var f0 int
	var f1 int
	var ff []int

	for i,_ := range args {
		if id != 0 && args[i] != "" && i == 0{
			stmt , err := db.Prepare("UPDATE passwords SET user=? WHERE pwid = ?")
			checkError(err)
			stmt.Exec(args[i],id)
			f = 1
			ff = append(ff,f)
			
		}else if id != 0 && args[i] != "" && i == 1{
			stmt , err := db.Prepare("UPDATE passwords SET passw=? WHERE pwid = ?")
			checkError(err)
			stmt.Exec(args[i],id)
			f0 = 1
			ff = append(ff,f0)
			
		}else if id != 0 && args[i] != "" && i == 2 {
			stmt , err := db.Prepare("UPDATE passwords SET category=? WHERE pwid = ?")
			checkError(err)
			stmt.Exec(args[i],id)
			f1 = 1
			ff = append(ff,f1)
			
		}else {
			f = 0
			ff = append(ff,f)
			
		}
	}
	return ff
}



func Save(user string, passwd string, category string, uid int)bool{
	db, err := sql.Open("sqlite3",DB)
	checkError(err)

	var ok bool 

	stmt , err := db.Prepare("INSERT INTO passwords(user,passw,category,userid) VALUES (?,?,?,?)")
	stmt.Exec(user,passwd,category,uid)

	if err != nil{
		ok = false
	}else{
		ok = true
	}
	db.Close()
	return ok
}

