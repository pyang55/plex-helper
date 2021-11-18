package main

import (
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/c-bata/go-prompt"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	cp "github.com/stromland/cobra-prompt"
)

var ip string
var token string
var rating string
var minimal bool

type Directories struct {
	XMLName     xml.Name    `xml:"MediaContainer"`
	Directories []Directory `xml:"Directory"`
}

type Directory struct {
	Type   string `xml:"type,attr"`
	Key    string `xml:"key,attr"`
	Title  string `xml:"title,attr"`
	Rating string `xml:"audienceRating,attr"`
}

type Videos struct {
	XMLName xml.Name `xml:"MediaContainer"`
	Videos  []Video  `xml:"Video"`
}

type Video struct {
	Key    string  `xml:"key,attr"`
	Medias []Media `xml:"Media"`
	Rating string  `xml:"audienceRating,attr"`
	Title  string  `xml:"title,attr"`
}

type Media struct {
	Resolution string `xml:"videoResolution,attr"`
}

func GetHttpRequests(url string) io.ReadCloser {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// handle err
	}

	resp, err := client.Do(req)
	if err != nil {
		// handle err
	}
	//defer resp.Body.Close()

	return resp.Body
}

var qs = []*survey.Question{
	{
		Name: "prompt",
		Prompt: &survey.Select{
			Message: "Are you sure you like to proceed?",
			Options: []string{"yes", "no", "I would like to remove a few"},
		},
	},
}

var fa = []*survey.Question{
	{
		Name: "final",
		Prompt: &survey.Select{
			Message: "Final Answer?",
			Options: []string{"yes", "i suck"},
		},
	},
}

func DeleteHttpRequests(url string) io.ReadCloser {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		// handle err
	}

	resp, err := client.Do(req)
	if err != nil {
		// handle err
	}
	//defer resp.Body.Close()

	return resp.Body
}

func GetMovieSections(ip string, token string) []string {
	url := fmt.Sprintf("http://%s:32400/library/sections?type=movie&X-Plex-Token=%s", ip, token)
	resp := GetHttpRequests(url)
	byteValue, _ := ioutil.ReadAll(resp)
	// we initialize our Users array
	var dir Directories
	var sections []string
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &dir)
	for i := 0; i < len(dir.Directories); i++ {
		if dir.Directories[i].Type == "movie" {
			sections = append(sections, dir.Directories[i].Key)
		}
	}
	return sections
}

func GetShowSections(ip string, token string) []string {
	url := fmt.Sprintf("http://%s:32400/library/sections?type=show&X-Plex-Token=%s", ip, token)
	resp := GetHttpRequests(url)
	byteValue, _ := ioutil.ReadAll(resp)
	// we initialize our Users array
	var dir Directories
	var sections []string
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'users' which we defined above
	xml.Unmarshal(byteValue, &dir)
	for i := 0; i < len(dir.Directories); i++ {
		if dir.Directories[i].Type == "show" {
			sections = append(sections, dir.Directories[i].Key)
		}
	}
	return sections
}

func GetMoviesByRating(ip string, token string, rating string) ([]string, map[string]string) {
	var movies []string
	m := make(map[string]string)
	sections := GetMovieSections(ip, token)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Audience RT Rating", "Video Resolution"})

	for _, y := range sections {
		url := fmt.Sprintf("http://%s:32400/library/sections/%s/all/?audienceRating<=%s&X-Plex-Token=%s", ip, y, rating, token)
		resp := GetHttpRequests(url)
		byteValue, _ := ioutil.ReadAll(resp)
		var vid Videos

		xml.Unmarshal(byteValue, &vid)
		for i := 0; i < len(vid.Videos); i++ {
			movies = append(movies, vid.Videos[i].Title)
			m[vid.Videos[i].Title] = vid.Videos[i].Key
			t.AppendRow([]interface{}{vid.Videos[i].Title, vid.Videos[i].Rating, vid.Videos[i].Medias[0].Resolution})
		}
	}
	if minimal {
		for _, y := range movies {
			fmt.Println(y)
		}
	} else {
		t.SetStyle(table.StyleColoredBright)
		t.Render()
	}
	return movies, m
}

func GetShowsByRating(ip string, token string, rating string) ([]string, map[string]string) {
	var shows []string
	m := make(map[string]string)
	sections := GetShowSections(ip, token)
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Audience RT Rating"})

	for _, y := range sections {
		url := fmt.Sprintf("http://%s:32400/library/sections/%s/all/?audienceRating<=%s&X-Plex-Token=%s", ip, y, rating, token)
		resp := GetHttpRequests(url)
		byteValue, _ := ioutil.ReadAll(resp)
		var vid Directories
		xml.Unmarshal(byteValue, &vid)
		for i := 0; i < len(vid.Directories); i++ {
			shows = append(shows, vid.Directories[i].Title)
			m[vid.Directories[i].Title] = vid.Directories[i].Key
			t.AppendRow([]interface{}{vid.Directories[i].Title, vid.Directories[i].Rating})
		}
	}
	if minimal {
		for _, y := range shows {
			fmt.Println(y)
		}
	} else {
		t.SetStyle(table.StyleColoredBright)
		t.Render()
	}
	return shows, m
}

var filter = &cobra.Command{
	Use:   "list",
	Short: "list function",
}

var delete = &cobra.Command{
	Use:   "delete",
	Short: "delete function",
}

var listmovies = &cobra.Command{
	Use:   "movies",
	Short: "filter movies by rating or version",
	Run: func(cmd *cobra.Command, args []string) {
		GetMoviesByRating(ip, token, rating)
		os.Exit(0)
	},
}

var listshows = &cobra.Command{
	Use:   "shows",
	Short: "filter shows by rating or version",
	Run: func(cmd *cobra.Command, args []string) {
		GetShowsByRating(ip, token, rating)
		os.Exit(0)
	},
}

var deletemovies = &cobra.Command{
	Use:   "movies",
	Short: "delete movies filtered by rating or version",
	Run: func(cmd *cobra.Command, args []string) {
		movies, movmap := GetMoviesByRating(ip, token, rating)
		Optionselector(movies, movmap)
		os.Exit(0)
	},
}

var deleteshows = &cobra.Command{
	Use:   "shows",
	Short: "delete shows filtered by rating or version",
	Run: func(cmd *cobra.Command, args []string) {
		shows, showmap := GetShowsByRating(ip, token, rating)
		Optionselector(shows, showmap)
		os.Exit(0)
	},
}

func Optionselector(item []string, maps map[string]string) {
	fmt.Println("These are the movies to be deleted")
	fmt.Println()

	answers := struct {
		Prompt string `survey:"prompt"` // or you can tag fields to match a specific name
	}{}

	finalanswer := struct {
		Prompt string `survey:"final"` // or you can tag fields to match a specific name
	}{}

	// perform the questions
	err := survey.Ask(qs, &answers)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	switch o := answers.Prompt; o {
	case "I would like to remove a few":
		mov := []string{}
		prompt := &survey.MultiSelect{
			Message: "What would you like to save",
			Options: item,
		}
		survey.AskOne(prompt, &mov, survey.WithPageSize(20))
		for _, y := range mov {
			item = Remove(item, y)
		}

		fmt.Println("These are the items you will be deleting")
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Name"})
		for _, mv := range item {
			t.AppendRow([]interface{}{mv})
		}
		t.SetStyle(table.StyleColoredBright)
		t.Render()
		err := survey.Ask(fa, &finalanswer)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if finalanswer.Prompt == "yes" {
			DeleteItem(item, maps)
		}
		if finalanswer.Prompt == "i suck" {
			Openbrowser("https://smartcdn.prod.postmedia.digital/nationalpost/wp-content/uploads/2019/06/flip-2.png?quality=90&strip=all&w=564&type=webp")
		}

	case "yes":
		DeleteItem(item, maps)
	case "no":
		fmt.Println("Goodbye")
		os.Exit(0)
	}
}

func HandleDynamicSuggestions(annotation string, _ prompt.Document) []prompt.Suggest {
	switch annotation {
	default:
		return []prompt.Suggest{}
	}
}

func Remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func Openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}

func DeleteItem(item []string, movmap map[string]string) {
	var wg sync.WaitGroup
	wg.Add(len(item))

	for _, y := range item {
		go func(y string) {
			defer wg.Done()
			res := strings.ReplaceAll(movmap[y], "/children", "")
			url := fmt.Sprintf("http://%s:32400%s?X-Plex-Token=%s", ip, res, token)
			fmt.Println(url)
			fmt.Println("Deleting ", y)
			DeleteHttpRequests(url)
		}(y)
	}
	wg.Wait()
}

func main() {
	var rootCmd = &cobra.Command{Use: "plex-helper"}
	listmovies.Flags().StringVarP(&ip, "ip-addr", "i", "127.0.0.1", "ip address of plex server")
	listmovies.Flags().StringVarP(&token, "token", "t", "", "Plex Token, more information https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token/")
	listmovies.Flags().StringVarP(&rating, "rating", "r", "", "Filter for movies of this rating or lower (based on Rotton Tomatoes)")
	listmovies.Flags().BoolVarP(&minimal, "minimal", "m", false, "Show simple output of list")
	listmovies.MarkFlagRequired("ip-addr")
	listmovies.MarkFlagRequired("token")
	listmovies.MarkFlagRequired("rating")

	listshows.Flags().StringVarP(&ip, "ip-addr", "i", "127.0.0.1", "ip address of plex server")
	listshows.Flags().StringVarP(&token, "token", "t", "", "Plex Token, more information https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token/")
	listshows.Flags().StringVarP(&rating, "rating", "r", "", "Filter for movies of this rating or lower (based on Rotton Tomatoes)")
	listshows.Flags().BoolVarP(&minimal, "minimal", "m", false, "Show simple output of list")

	listshows.MarkFlagRequired("ip-addr")
	listshows.MarkFlagRequired("token")
	listshows.MarkFlagRequired("rating")

	deletemovies.Flags().StringVarP(&ip, "ip-addr", "i", "127.0.0.1", "ip address of plex server")
	deletemovies.Flags().StringVarP(&token, "token", "t", "", "Plex Token, more information https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token/")
	deletemovies.Flags().StringVarP(&rating, "rating", "r", "", "Filter for movies of this rating or lower (based on Rotton Tomatoes)")
	deletemovies.MarkFlagRequired("ip-addr")
	deletemovies.MarkFlagRequired("token")
	deletemovies.MarkFlagRequired("rating")

	deleteshows.Flags().StringVarP(&ip, "ip-addr", "i", "127.0.0.1", "ip address of plex server")
	deleteshows.Flags().StringVarP(&token, "token", "t", "", "Plex Token, more information https://support.plex.tv/articles/204059436-finding-an-authentication-token-x-plex-token/")
	deleteshows.Flags().StringVarP(&rating, "rating", "r", "", "Filter for movies of this rating or lower (based on Rotton Tomatoes)")
	deleteshows.MarkFlagRequired("ip-addr")
	deleteshows.MarkFlagRequired("token")
	deleteshows.MarkFlagRequired("rating")

	rootCmd.AddCommand(filter)
	rootCmd.AddCommand(delete)
	filter.AddCommand(listmovies)
	filter.AddCommand(listshows)
	delete.AddCommand(deletemovies)
	delete.AddCommand(deleteshows)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.Execute()

	shell := &cp.CobraPrompt{
		RootCmd:                rootCmd,
		DynamicSuggestionsFunc: HandleDynamicSuggestions,
		ResetFlagsFlag:         false,
		GoPromptOptions: []prompt.Option{
			prompt.OptionTitle("plex-helper"),
			prompt.OptionPrefix("> plex-helper "),
			prompt.OptionMaxSuggestion(10),
			prompt.OptionCompletionOnDown(),
			prompt.OptionCloseOnControlC(),
			prompt.OptionSetExitCheckerOnInput(func(in string, breakline bool) bool {
				return breakline
			}),
		},
	}
	shell.Run()

}
