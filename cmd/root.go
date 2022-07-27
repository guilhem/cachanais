/*
Copyright © 2022 Guilhem Lettron

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cachanais",
	Short: "Populate cache of a website",
	Long: `cachanais populate cache by visiting all link of a website.
	You an set the real address if used localy.`,
	RunE: crawl,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	connectURL         string
	scrapingURL        string
	cookies            []string
	headers            []string
	filterQueryStrings bool
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cachanais.yaml)")

	rootCmd.Flags().StringVarP(&connectURL, "address", "a", "", "url to connect")

	rootCmd.Flags().StringVarP(&scrapingURL, "url", "u", "", "Url to scrape")
	rootCmd.MarkFlagRequired("url")

	rootCmd.Flags().StringSliceVar(&cookies, "cookies", []string{}, "Cookies to set in the form key:value")

	rootCmd.Flags().StringSliceVar(&headers, "headers", []string{}, "Headers to set in the form key:value")

	rootCmd.Flags().BoolVar(&filterQueryStrings, "filter-query-strings", false, "filter url with query strings")

	rootCmd.Example = "cachanais --url https://text.com --address http://localhost --cookies mycookie:sup --headers X-Cool:blop"

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func crawl(cmd *cobra.Command, args []string) error {
	u, err := url.Parse(scrapingURL)
	if err != nil {
		return err
	}

	var connectU *url.URL

	if connectURL == "" {
		*connectU = *u
	} else {
		connectU, err = url.Parse(connectURL)
		if err != nil {
			return err
		}
	}

	parsedHeaders := make(map[string]string, len(headers)+1)

	parsedHeaders["host"] = u.Host

	for _, header := range headers {
		h := strings.SplitN(header, ":", 2)
		if len(h) != 2 {
			return fmt.Errorf("problem with header '%s': %w", header, errors.New("Header malformed"))
		}
		parsedHeaders[h[0]] = h[1]
	}

	regexpQueryString := "[\\?&]([^&=]+)=([^&=]+)"

	c := colly.NewCollector(
		colly.AllowedDomains(u.Host, connectU.Host),
		colly.IgnoreRobotsTxt(),
		colly.Headers(parsedHeaders),
	)

	if filterQueryStrings {
		c.DisallowedURLFilters = []*regexp.Regexp{regexp.MustCompile(regexpQueryString)}
	}

	parsedCookies := make([]*http.Cookie, 0, len(cookies))

	for _, cookie := range cookies {
		c := strings.SplitN(cookie, ":", 2)
		if len(c) != 2 {
			return fmt.Errorf("problem with cookie '%s': %w", cookie, errors.New("cookie malformed"))
		}
		parsedCookies = append(parsedCookies, &http.Cookie{
			Name:  c[0],
			Value: c[1],
		})
	}

	if err := c.SetCookies(connectU.String(), parsedCookies); err != nil {
		return err
	}

	c.SetRequestTimeout(time.Minute)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if err := e.Request.Visit(link); err != nil {
			fmt.Println("error with link", link)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.URL.Host = connectU.Host
		r.URL.Scheme = connectU.Scheme
		fmt.Println("Visiting", r.URL)
	})

	if err := c.Visit(scrapingURL); err != nil {
		return err
	}

	return nil
}
