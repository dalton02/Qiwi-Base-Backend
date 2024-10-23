package userController

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Função para obter os cookies do SIGAA
func getCookiesFromSigaa(login string, password string) (map[string]string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://sig.ufca.edu.br/sigaa/mobile/touch/login.jsf", nil)
	req.Header.Set("Cookie", "")
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	cookieHeader := resp.Header["Set-Cookie"]
	page := getInfoFromSigaa(login, password, string(cookieHeader[0]))
	return getEssentialHtml(page)
}

// Função para obter informações do SIGAA
func getInfoFromSigaa(login, password, cookie string) string {
	// Configurações do cabeçalho da requisição
	client := &http.Client{}
	data := url.Values{}
	data.Set("form-login", "form-login")
	data.Set("form-login:j_id_jsp_35904869_2", login)
	data.Set("form-login:j_id_jsp_35904869_3", password)
	data.Set("form-login:entrar", "Entrar")
	data.Set("javax.faces.ViewState", "j_id1")

	req, _ := http.NewRequest("POST", "https://sig.ufca.edu.br/sigaa/mobile/touch/login.jsf", strings.NewReader(data.Encode()))
	req.Header.Set("Host", "sig.ufca.edu.br")
	req.Header.Set("Referer", "https://sig.ufca.edu.br/sigaa/mobile/touch/login.jsf")
	req.Header.Set("Origin", "https://sig.ufca.edu.br")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("Cookie", cookie)

	resp, _ := client.Do(req)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}
func getEssentialHtml(page string) (map[string]string, error) {
	data := make(map[string]string)
	fullnameRegex := regexp.MustCompile(`<strong>(.*?)<\/strong>`)
	fullnameMatches := fullnameRegex.FindStringSubmatch(page)
	if len(fullnameMatches) > 0 {
		data["nome"] = fullnameMatches[1]
		data["nome"] = convertCharacters(data["nome"])
	}
	idRegex := regexp.MustCompile(`<span id="form-portal-discente:matricula">(.*?)<\/span>`)
	idMatches := idRegex.FindStringSubmatch(page)
	if len(idMatches) > 0 {
		data["codigo"] = idMatches[1]
	}
	courseRegex := regexp.MustCompile(`<\/strong><br\/>\s*(.*?)<br\/>`)
	courseMatches := courseRegex.FindStringSubmatch(page)
	if len(courseMatches) > 0 {
		data["curso"] = strings.TrimSpace(courseMatches[1])
		data["curso"] = convertCharacters(data["curso"])
	}
	if len(data) == 0 {
		return data, fmt.Errorf("Credenciais incorretas")
	}
	return data, nil
}

// Converte caracteres especiais para uma forma legível.
func convertCharacters(s string) string {
	s = strings.ReplaceAll(s, "&Ecirc;", "Ê")
	s = strings.ReplaceAll(s, "&Acirc;", "Â")
	s = strings.ReplaceAll(s, "&Ocirc;", "Ô")
	s = strings.ReplaceAll(s, "&Ucirc;", "Û")
	s = strings.ReplaceAll(s, "&Icirc;", "Î")
	s = strings.ReplaceAll(s, "&Ccedil;", "Ç")
	s = strings.ReplaceAll(s, "&Atilde;", "Ã")
	s = strings.ReplaceAll(s, "&Otilde;", "Õ")
	s = strings.ReplaceAll(s, "&Etilde;", "Ẽ")
	s = strings.ReplaceAll(s, "&Utilde;", "Ũ")
	s = strings.ReplaceAll(s, "&Itilde;", "Ĩ")
	return s
}

// // Função para obter os cookies do SIGAA.
// func getCookiesFromSigaa(login string, password string, response http.ResponseWriter) string {
// 	url := "https://sig.ufca.edu.br/sigaa/logar.do?dispatch=logOn"
// 	method := "POST"
// 	payload := &bytes.Buffer{}
// 	writer := multipart.NewWriter(payload)
// 	writer.WriteField("user.login", login)
// 	writer.WriteField("user.senha", password)
// 	writer.Close()
// 	client := &http.Client{}
// 	req, _ := http.NewRequest(method, url, payload)
// 	req.Header.Set("Content-Type", writer.FormDataContentType())
// 	req.Header.Set("User-Agent", "PostmanRuntime/7.42.0")
// 	req.Header.Set("Accept", "*/*")
// 	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
// 	req.Header.Set("Connection", "keep-alive")
// 	res, _ := client.Do(req)
// 	cookies := res.Cookies()
// 	cookieString := ""
// 	for _, cookie := range cookies {
// 		cookieString += fmt.Sprintf("%s=%s", cookie.Name, cookie.Value)
// 	}
// 	defer res.Body.Close()
// 	page := getPageFromSigaa(cookieString, response)
// 	return page
// }

// // Codigo perfeito com cookie de exemplo, negocio é procurar o cookie correto.......
// func getPageFromSigaa(cookies string, response http.ResponseWriter) string {
// 	url := "https://sig.ufca.edu.br/sigaa/portais/discente/discente.jsf"
// 	method := "GET"
// 	client := &http.Client{}
// 	req, _ := http.NewRequest(method, url, nil)
// 	req.Header.Add("Cookie", cookies)
// 	res, _ := client.Do(req)
// 	fmt.Println(res.StatusCode)
// 	body, _ := ioutil.ReadAll(res.Body)
// 	htmlContent := string(body)
// 	return htmlContent
// }
