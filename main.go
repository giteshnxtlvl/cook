package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

var m = make(map[interface{}]map[string][]string)
var params = make(map[string]string)
var pattern = []string{}
var version = "1.5"

// var verbose = false
var min int

const (
	blue   = "\u001b[38;5;14m"
	green  = "\u001b[38;5;46m"
	purple = "\u001b[38;5;207m"
	red    = "\u001b[38;5;196m"
	bold   = "\u001b[1m"
	white  = "\u001b[38;5;255m"
	reset  = "\u001b[0m"
)

var banner = `

                             
  ░            ░ ░      ░ ░  ░  ░   
  ░ ░        ░ ░ ░ ▒  ░ ░ ░ ▒  ░    
░░▒ ▒░    ░ ▒ ▒░   ░ ▒ ▒░ ░ ░ ▒ ░
░ ░▒ ▒  ░░ ▒░▒░▒░ ░ ▒░▒░▒░ ▒ ▒▒ ▓▒
 ▄████▄   ▒█████   ▒█████   ██ ▄█▀           A CUSTOMIZABLE WORDLIST
▒██▀ ▀█  ▒██▒  ██▒▒██▒  ██▒ ██▄█▒            AND PASSWORD GENERATOR
▒▓█    ▄ ▒██░  ██▒▒██░  ██▒▓███▄░ 
▒▓▓▄ ▄██▒▒██   ██░▒██   ██░▓██ █▄            dev by Gitesh Sharma 
 ▒▓███▀ ░░ ████▓▒░░ ████▓▒░▒██▒ █▄ ` + version + `      @giteshnxtlvl

`

var config = `
# Character set like crunch
charSet:
    sep   : [_- ] #common separators
    n     : [0123456789]
    A     : [ABCDEFGHIJKLMNOPQRSTUVWXYZ]
    a     : [abcdefghijklmnopqrstuvwxyz]
    aAn   : [abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789]
    An    : [ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789]
    an    : [abcdefghijklmnopqrstuvwxyz0123456789]
    aA    : [abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ]
    s     : ["!#$%&'()*+,-./:;<=>?@[\\]^_` + "`" + `{|}~&\""]
    all   : ["!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_` + "`" + `abcdefghijklmnopqrstuvwxyz{|}~\""]

# File to access from anywhere
files:
    raft_ext     : [E:\tools\wordlists\SecLists\Discovery\Web-Content\raft-large-extensions.txt]
    raft_dir     : [E:\tools\wordlists\SecLists\Discovery\Web-Content\raft-large-directories.txt]
    raft_files   : [E:\tools\wordlists\SecLists\Discovery\Web-Content\raft-large-files.txt]
    robot_1000   : [E:\tools\wordlists\SecLists\Discovery\Web-Content\RobotsDisallowed-Top1000.txt]

# Create your lists
lists:
    schemas      : [aim, callto, cvs, data, facetime, feed, file, ftp, git, gopher, gtalk, h323, hdl, http, https, imap, irc, irc6, ircs, itms, javascript, magnet, mailto, mms, msnim, news, nntp, prospero, rsync, rtsp, rtspu, sftp, shttp, sip, sips, skype, smb, snews, ssh, svn, svn, svn+ssh, telnet, tel, wais, ymsg]
    bypass       : ["%00", "%09", "%0A", "%0D", "%0D%0A"]

    admin_set    : [admin, root, su, superuser, administration]
    api          : [/v1/,/v2/,/v3/,/v4/,/v5/,/api/]
    pass_ends    : [123, "@123", "#123"]

    months       : [January, February, March, April, May, June, July, August, September, October, November, December]
    mons         : [Jan, Feb, Mar, Apr, May, Jun, Jul, Aug, Sep, Oct, Nov, Dec]

# Patterns
patterns:
    date : 
        - date(D,M,Y)
        - DMY  
        - MDY 
        - D/M/Y  
        - M/D/Y 
        - D-M-Y  
        - M-D-Y   
        - D.M.Y  
        - M.D.Y  
        - D.Y  
        - M.Y  
        - D.M

# Extension Set, . will added before using this
extensions:
    config  : [conf, confq, config]
    data    : [xml, json, yaml, yml]
    backup  : [bak, backup, backup1, backup2]
    exec    : [exe, msi, bin, command, sh, bat, crx]
    web     : [html, html5, htm, js, jsx, jsp, wasm, php, php3, php5, php7]
    iis     : [asax, ascx, asmx, aspx, exe, aspx.cs, ashx, axd, config, htm, jar, js, rdl, swf, txt, xls, xml, xsl, zpd, suo, sln]
    archive : [7z, a, apk, xapk, ar, bz2, cab, cpio, deb, dmg, egg, gz, iso, jar, lha, mar, pea, rar, rpm, s7z, shar, tar, tbz2, tgz, tlz, war, whl, xpi, zip, zipx, xz, pak, tar.gz, gz]
    code    : [c, cc, class, clj, cpp, cs, cxx, el, go, h, java, lua, m, m4, php, php3, php5, php7, pl, po, py, rb, rs, sh, swift, vb, vcxproj, xcodeproj, xml, diff, patch, js, jsx]

    #Rest
    css_type: [css, less, scss]
    sheet   : [ods, xls, xlsx, csv, ics vcf]
    slide   : [ppt, pptx, odp]
    font    : [eot, otf, ttf, woff, woff2]
    text    : [doc, docx, ebook, log, md, msg, odt, org, pages, pdf, rtf, rst, tex, txt, wpd, wps]
    audio   : [aac, aiff, ape, au, flac, gsm, it, m3u, m4a, mid, mod, mp3, mpa, pls, ra, s3m, sid, wav, wma, xm]
    book    : [mobi, epub, azw1, azw3, azw4, azw6, azw, cbr, cbz]
    video   : [3g2, 3gp, aaf, asf, avchd, avi, drc, flv, m2v, m4p, m4v, mkv, mng, mov, mp2, mp4, mpe, mpeg, mpg, mpv, mxf, nsv, ogg, ogv, ogm, qt, rm, rmvb, roq, srt, svi, vob, webm, wmv, yuv]
    image   : [3dm, 3ds, max, bmp, dds, gif, jpg, jpeg, png, psd, xcf, tga, thm, tif, tiff, yuv, ai, eps, ps, svg, dwg, dxf, gpx, kml, kmz, webp]
`

func valueInSlice(list []string, val string) bool {
	for _, l := range list {
		if l == val {
			return true
		}
	}
	return false
}

func findRegex(file, expresssion string) []string {
	founded := []string{}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return []string{file + ":" + expresssion}
	}

	r, err := regexp.Compile(expresssion)
	if err != nil {
		panic(err)
	}

	data := strings.ReplaceAll(string(content), "\r", "")
	extensionsList := r.FindAllString(data, -1)

	for _, found := range extensionsList {
		if !valueInSlice(founded, found) {
			founded = append(founded, found)
		}
	}

	return founded
}

func fileValues(file string) []string {
	content, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println(file)
		panic(err)
	}

	return strings.Split(strings.ReplaceAll(string(content), "\r", ""), "\n")
}

func cookConfig() {

	configFile := os.Getenv("COOK")
	content := []byte(config)

	if len(configFile) == 0 {

	} else if _, err := os.Stat(configFile); err == nil {
		// If file exists
		content, err = ioutil.ReadFile(configFile)
		if err != nil {
			fmt.Printf("Config File Reading Error: %v\n", err)
		}

		//If file is empty
		if len(content) == 0 {
			ioutil.WriteFile(configFile, []byte(config), 0644)
			content = []byte(config)
		}

	} else if os.IsNotExist(err) {
		err := ioutil.WriteFile(configFile, []byte(config), 0644)
		if err != nil {
			fmt.Printf("Config File Writing Error: %v\n", err)
		}
	}

	err := yaml.Unmarshal([]byte(content), &m)

	if err != nil {
		fmt.Printf("error: %v", err)
	}
}

func applyColumnCases(columnValues []string, columnNum int) {
	temp := []string{}

	// Using cases for columnValues
	if _, exists := columnCases[columnNum]; exists {

		//All cases
		if valueInSlice(columnCases[columnNum], "A") {
			for _, t := range final {
				for _, v := range columnValues {
					temp = append(temp, t+strings.ToUpper(v))
					temp = append(temp, t+strings.ToLower(v))
					temp = append(temp, t+strings.Title(v))
				}
			}
		} else {

			if valueInSlice(columnCases[columnNum], "U") {
				for _, t := range final {
					for _, v := range columnValues {
						temp = append(temp, t+strings.ToUpper(v))
					}
				}
			}

			if valueInSlice(columnCases[columnNum], "L") {
				for _, t := range final {
					for _, v := range columnValues {
						temp = append(temp, t+strings.ToLower(v))
					}
				}
			}

			if valueInSlice(columnCases[columnNum], "T") {
				for _, t := range final {
					for _, v := range columnValues {
						temp = append(temp, t+strings.Title(v))
					}
				}
			}
		}

	} else {
		for _, t := range final {
			for _, v := range columnValues {
				temp = append(temp, t+v)
			}
		}
	}

	final = temp
}

//Initializing with empty string, so loops will run for 1st column
var final = []string{""}

func main() {

	cookConfig()
	parseInput(os.Args[1:])

	for columnNum, param := range pattern {

		columnValues := []string{}

		for _, p := range strings.Split(param, ",") {

			val, success := parseRanges(p)
			if success {
				columnValues = append(columnValues, val...)
				continue
			}
			if val, exists := params[p]; exists {
				columnValues = append(columnValues, parseValue(val)...)
				continue
			}
			if val, exists := m["charSet"][p]; exists {
				chars := strings.Split(val[0], "")
				columnValues = append(columnValues, chars...)
				continue
			}
			if val, exists := m["files"][p]; exists {
				columnValues = append(columnValues, fileValues(val[0])...)
				continue
			}
			if val, exists := m["lists"][p]; exists {
				columnValues = append(columnValues, val...)
				continue
			}
			if val, exists := m["extensions"][p]; exists {
				for _, ext := range val {
					columnValues = append(columnValues, "."+ext)
				}
				continue
			}

			columnValues = append(columnValues, p)
		}

		applyColumnCases(columnValues, columnNum)

		if columnNum >= min {
			for _, v := range final {
				fmt.Println(v)
			}
		}
	}
}
