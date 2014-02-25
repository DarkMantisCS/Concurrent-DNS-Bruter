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
  Log             = logger.New().Log
  ValidWord       = regexp.MustCompile(`^[a-zA-Z0-9_-].*$`)
  // DisallowedChars = []int{47, 46, 58, 59, 60, 61, 62, 63, 64, 91, 92, 93, 94, 95, 96}

  AllowedChars    = []rune{'a','b','c','d','e','f','g','h','i','j','k','l','m','n','o','p','q','r','s','t','u','v','w','x','y','z','0','1','2','3','4','5','6','7','8','9','-','_'}
)

/**
 * @todo This function really needs to be rethought out
 */
func bruterAttack(server string, length int, verbose bool) {
  chars := make([]rune, length)

  for c := 0; c < len(chars); c++ {
    for i := 0; i < len(chars); i++{
      if i == c {
        continue;
      }
      chars[c] = AllowedChars[i]
      for x := 0; x < len(chars); x++ {
        if x == i {
          continue;
        }
        chars[i] = AllowedChars[x]
        for y := 0; y < len(AllowedChars); y++ {
          if y == x {
            continue;
          }
          chars[x] = AllowedChars[y]
          fmt.Println(string(chars))
        }
      }
    }
    chars[c]++
  }


  // for c := 0; c < len(chars); c++ {
  //   for chars[c] <= 122 {
  //     if stringInSlice(int(chars[c]), DisallowedChars) {
  //       chars[c]++
  //       continue
  //     }
  //     for x := 0; x < len(chars); x++ {
  //       if x == c {
  //         continue
  //       }
  //       for y := 45; y <= 122; y++ {
  //         if stringInSlice(y, DisallowedChars) {
  //           continue
  //         }
  //         chars[x] = rune(y)
  //         // TODO: Handle bruter data
  //         Log(fmt.Sprintf("%s.%s", string(chars), server))
  //       }
  //     }
  //     chars[c]++
  //   }
  // }
}

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
    Log("e", err.Error())
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
        Log(fmt.Sprintf("%s found: \t %q", current, res))
      }
      tmp := fmt.Sprintf("%s - %s", current, res)
      output = append(output, tmp)
      found++
    }
  }

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


func main() {
  wordlist := flag.String("wordlist", "", "Wordlist if you want dictionary attack (new line delimited)")
  verbose := flag.Bool("verbose", false, "Set verbosity up (true|false)")
  length := flag.Int("length", 6, "Sets the max length of the brute forcer")

  flag.Parse()

  args := flag.Args()
  if len(args) == 0 {
    Log("c", "No host provided to search.")
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
