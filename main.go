package main

import (
    "fmt"
    "net"
    "flag"
    "io/ioutil"
    "log"
    "regexp"
    "strings"
)

// func checkByBrute(ns string) (bool, error) {
//     dns := net.Dial("tcp", )
// }

func checkDomain(subdomain string) bool {
    validWord := regexp.MustCompile(`^[a-zA-Z0-9_-].*$`)

    if !validWord.MatchString(subdomain) {
        return false
    }

    return true
}

func dictionaryAttack(server string, wordlist string, verbose bool) {
    words, err := ioutil.ReadFile(wordlist)

    if err != nil {
        panic(err)
    }

    found    := 0
    notFound := 0

    wordsInFile := strings.Split(string(words), "\n")

    var output []string

    for _, subdomain := range wordsInFile {

        subdomain = strings.TrimRight(subdomain, "\r")

        if !checkDomain(subdomain) {
            continue
        }

        current := fmt.Sprintf("%s.%s", subdomain, server)

        res, err := net.LookupHost(current)

        if err != nil {
            notFound++
        } else {

            if verbose {
                fmt.Printf("%s found: \t %q\r\n", current, res)
            }
            tmp := fmt.Sprintf("%s - %s", current, res)
            output = append(output, tmp)
            found++
        }
    }

    fmt.Printf("%d Total Found;\r\n%d Total NOT Found\r\n", found, notFound)
    fmt.Printf("%v\r\n", output)
}

func main() {
    server      := flag.String("server",    "",     "The server to brute")
    wordlist    := flag.String("wordlist",  "",     "Wordlist if you want dictionary attack (new line delimited)")
    verbose     := flag.Bool("verbose",     false,  "Set verbosity up (true|false)")


    flag.Parse()

    if *server == "" {
        log.Fatal("You must specify a server")
    }

    // Use the wordlist over brute force
    if *wordlist != "" {
        dictionaryAttack(*server, *wordlist, *verbose)
    }

    // Brute force method
}
