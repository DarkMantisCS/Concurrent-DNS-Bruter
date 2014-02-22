package main

import (
  "flag"
  "fmt"
  "github.com/Southern/logger"
  "io/ioutil"
  "net"
  "regexp"
  "strings"
)

var (
  log       = logger.New()
  ValidWord = regexp.MustCompile(`^[a-zA-Z0-9_-].*$`)
)

func checkDomain(subdomain string) bool {
  if !ValidWord.MatchString(subdomain) {
    return false
  }

  return true
}

func checkDnsRequest(current string) (string, error) {
  res, err := net.LookupHost(current)

  if err != nil {
    return "", err
  }

  return fmt.Sprintf("%s", res), nil
}

func dictionaryAttack(server string, wordlist string, verbose bool) {
  words, err := ioutil.ReadFile(wordlist)

  if err != nil {
    log.Log("e", err.Error())
  }

  found := 0
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
        log.Log(fmt.Sprintf("%s found: \t %q", current, res))
      }
      tmp := fmt.Sprintf("%s - %s", current, res)
      output = append(output, tmp)
      found++
    }
  }

  log.
    Log(fmt.Sprintf("%d Total Found", found)).
    Log(fmt.Sprintf("%d Total Not Found", notFound)).
    Log(fmt.Sprintf("%+v", output))
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
  disallowedChars := []int{47, 46, 58, 59, 60, 61, 62, 63, 64, 91, 92, 93, 94, 95, 96}
  for i := 45; i <= 122; i++ {
    if stringInSlice(i, disallowedChars) {
      continue
    }

    for y := 45; y <= 122; y++ {
      if stringInSlice(y, disallowedChars) {
        continue
      }
      str := []string{string(i)}
      for x := 0; x <= length; x++ {
        j := x + y
        if j > 122 || stringInSlice(j, disallowedChars) {
          continue
        }
        str = append(str, string(j))
      }
      log.Log(fmt.Sprintf("%+v", str))
    }
  }
}

func main() {
  wordlist := flag.String("wordlist", "", "Wordlist if you want dictionary attack (new line delimited)")
  verbose := flag.Bool("verbose", false, "Set verbosity up (true|false)")
  length := flag.Int("length", 6, "Sets the max length of the brute forcer")

  flag.Parse()

  args := flag.Args()
  if len(args) == 0 {
    log.Log("c", "No host provided to search.")
  }

  for _, server := range args {
    // Use the wordlist over brute force
    if *wordlist != "" {
      dictionaryAttack(server, *wordlist, *verbose)
    } else {
      bruterAttack(server, *length, *verbose)
    }
  }
}