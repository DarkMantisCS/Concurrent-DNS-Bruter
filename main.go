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

func checkDnsRequest(current string) (string, error){
    res, err := net.LookupHost(current)

    if err != nil {
        return "", err
    }

    return fmt.Sprintf("%s",res), nil
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

        res, err := checkDnsRequest(current)

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

func stringInSlice(a int, list []int) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}


/**
 * @todo This function really needs to be rethought out
 */
func bruterAttack(server string, length int, verbose bool) {
    disallowedChars := []int{47,46,58,59,60,61,62,63,64,91,92,93,94,95,96}
    for i := 45; i < 122; i++ {
        if stringInSlice(i, disallowedChars) {
            continue
        }

        for x := 0; x <= length; i++ {
            fmt.Println(string(i))
        }
    }
}

func main() {
    server      := flag.String("server",    "",     "The server to brute")
    wordlist    := flag.String("wordlist",  "",     "Wordlist if you want dictionary attack (new line delimited)")
    verbose     := flag.Bool("verbose",     false,  "Set verbosity up (true|false)")


    flag.Parse()

    bruterAttack(*server, 3, false);

    if *server == "" {
        log.Fatal("You must specify a server")
    }

    // Use the wordlist over brute force
    if *wordlist != "" {
        dictionaryAttack(*server, *wordlist, *verbose)
    }

    // Brute force method
}
