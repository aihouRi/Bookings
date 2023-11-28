package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/aihou/bookings/pkg/config"
	"github.com/aihou/bookings/pkg/models"
	"github.com/justinas/nosurf"
)

//var functions = template.FuncMap{}

var app *config.AppConfig

// Newtemplates sets the config for the template package
func Newtemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders templates using html/template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		//get the template cache from the app config
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	//这里是创建一个缓冲区，然后将模板渲染到缓冲区中
	buf := new(bytes.Buffer)

	td = AddDefaultData(td, r)

	//这里的Execute的作用是将模板渲染到缓冲区中
	_ = t.Execute(buf, td)

	//render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser", err)
	}
}

// createTemplateCache函数是创建一个模板缓存，这里的map的key是string，value是*template.Template
// 其目的为了提高性能，因为如果每次请求都要解析模板文件，那么性能就会很低，所以我们可以将模板文件
// 解析到模板中，然后将模板存入到map中，这样就可以提高性能了
func CreateTemplateCache() (map[string]*template.Template, error) {
	//initialize a new map to act as the cache
	myCache := map[string]*template.Template{}

	//get all of the files named *.page.tmpl.html from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	//range through all files ending with *.page.tmpl.html
	for _, page := range pages {
		//extract the file name
		//filepath.Base的作用是返回路径的最后一个元素，比如说
		//./templates/home.page.tmpl，返回的就是home.page.tmpl
		//注意这里的filepath.Base函数只拿出文件名，不拿出内容
		name := filepath.Base(page)

		//这里的template.New(name)是创建一个新的模板，然后用ParseFiles方法，
		//将模板文件解析到模板中，这里的name是模板的名字，也就是home.page.tmpl
		//这里的ts是一个指向模板的指针
		//ParseFiles的作用是将模板文件解析到模板中，比如说，我们有一个模板文件
		//home.page.tmpl，然后我们用ParseFiles方法将这个模板文件解析到模板中，
		//然后我们就可以用这个模板来渲染页面了
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			//ParseGlob的作用是将模板文件解析到模板中，这里的*是通配符，表示匹配所有的文件
			//比如说，我们有一个模板文件home.layout.tmpl.html，然后我们用ParseGlob方法将
			//这个模板文件解析到模板中，然后我们就可以用这个模板来渲染页面了
			//ts可以存入多个模板吗？可以的，ts是一个指向模板的指针，我们可以用ts来存入多个模板
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		//名称对应内容，比如说home.page.tmpl对应的是ts
		myCache[name] = ts
	}
	return myCache, nil
}
