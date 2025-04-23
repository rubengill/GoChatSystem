package main

import "regexp"

// Validate the nickname based on defined rules
func isValidNickname(nick string) bool {
    if len(nick) == 0 || len(nick) > 10 {
        return false
    }
    regex := regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{0,9}$`)
    return regex.MatchString(nick)
}
