package parse
import "strings"

func ExtractPort(rawURL string) string {
    text := rawURL

    // Remove scheme (http://, https://, etc.)
    schemePos := strings.Index(text, "://")
    if schemePos != -1 {
        text = text[schemePos+3:]
    }

    // Remove userinfo (user:pass@)
    atPos := strings.LastIndex(text, "@")
    if atPos != -1 {
        text = text[atPos+1:]
    }

    // Cut path, query or fragment
    endPos := strings.IndexAny(text, "/?#&")
    if endPos != -1 {
        text = text[:endPos]
    }

    // IPv6 literal: [address]:port
    if strings.HasPrefix(text, "[") {
        closingBracket := strings.Index(text, "]")
        if closingBracket == -1 {
            return ""
        }

        hasPort := false
        if len(text) > closingBracket+1 {
            if text[closingBracket+1] == ':' {
                hasPort = true
            }
        }

        if hasPort {
            return text[closingBracket+2:]
        }

        return ""
    }

    // IPv4 or hostname: host:port
    colonPos := strings.LastIndex(text, ":")
    if colonPos != -1 {
        return text[colonPos+1:]
    }

    return ""
}
