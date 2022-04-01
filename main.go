package main

import (
    "os"
    "fmt"
    "log"
    "bytes"
    "net/http"
    "encoding/json"
    "gopkg.in/ini.v1"
)

type Repo struct {
    Name string `json:"name"`
    Full_name string `json:"full_name"`
    Description string `json:"description"`
    Visibility string `json:"visibility"`
    Import_url string `json:"import_url"`
}

type Provider struct {
    Protocol string
    Host string
    Token string
    User string
    Pass string
    Endpoint string
}

func (p Provider) get_repos(repos *[]Repo) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", p.Protocol+p.Host+p.Endpoint, nil)
    req.Header.Add("Authorization", "token " + p.Token)
    resp, err := client.Do(req)
    if err != nil { log.Fatalln(err) }
    if resp.StatusCode != 200 { log.Fatalln("❌ KO " + resp.Status) }
    json.NewDecoder(resp.Body).Decode(&repos)
}

func (p Provider) put_repo(repo *Repo) {
    fmt.Printf("   Importing: %s \n", repo.Full_name)
    client := &http.Client{}
    payload := new(bytes.Buffer)
    json.NewEncoder(payload).Encode(repo)
    // fmt.Println(payload)
    req, err := http.NewRequest("POST", p.Protocol+p.Host+p.Endpoint, payload)
    req.Header.Set("Content-type", "application/json")
    req.Header.Add("Private-Token", p.Token)
    resp, err := client.Do(req)
    if err != nil { log.Fatalln(err) }
    if resp.StatusCode != 201 {
        fmt.Println("❌ KO " + resp.Status)
    } else {
        fmt.Printf("✅ OK")
    }
}

func main() {
    // Read the config file
    home, err := os.UserHomeDir()
    if err != nil { log.Fatalln(err) }
    config, err :=  ini.Load(home + "/.config/gogs2gitlab/gogs2gitlab.ini")
    if err != nil { log.Fatalln(err) }
    // Define providers
    gogs := Provider {
        Protocol: config.Section("").Key("gogs_proto").String(),
        Host: config.Section("").Key("gogs_host").String(),
        Token: config.Section("").Key("gogs_token").String(),
        User: config.Section("").Key("gogs_user").String(),
        Pass: config.Section("").Key("gogs_pass").String(),
        Endpoint: "/api/v1/user/repos",
    }
    gitlab := Provider {
        Protocol: config.Section("").Key("gitlab_proto").String(),
        Host: config.Section("").Key("gitlab_host").String(),
        Token: config.Section("").Key("gitlab_token").String(),
        User: config.Section("").Key("gitlab_user").String(),
        Pass: config.Section("").Key("gitlab_pass").String(),
        Endpoint: "/api/v4/projects",
    }
    // Start the game
    repos := []Repo{}
    gogs.get_repos(&repos)

    for _, repo := range repos {
        fmt.Printf("➡️  %s \n", repo.Full_name)
        repo.Visibility = "private"
        repo.Import_url = fmt.Sprintf("%s%s:%s@%s/%s", gogs.Protocol, gogs.User, gogs.Pass, gogs.Host, repo.Full_name)
        gitlab.put_repo(&repo)
    }
}
