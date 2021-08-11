package ginjet

import (
	"github.com/CloudyKit/jet"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/fatih/structs"
	"net/http"
)

var htmlContentType = []string{"text/html; charset=utf-8"}

// JetRender is a custom Gin template renderer using Jet
type JetRender struct {
	Options  *RenderOptions
	Template *jet.Template
	Variables jet.VarMap
	Data     interface{}
	
	globals map[string]interface{}
}

// New creates a new JetRender instance with custom Options.
func New(options *RenderOptions) *JetRender {
	return &JetRender{
		Options: options,
	}
}

// Default creates a JetRender instance with default options.
func Default() *JetRender {
	return New(DefaultOptions())
}


func (r JetRender) Instance(name string, data interface{}) render.Render {
	set := jet.NewHTMLSet(r.Options.TemplateDir)
	
	//设置全局变量
	if r.globals!=nil{
		for key,value:=range r.globals{
			set.AddGlobal(key,value)
		}
	}
	
	t, err := set.GetTemplate(name)

	if err != nil {
		panic(err)
	}
	
	var v jet.VarMap
	if data != nil {
		vars, ok := data.(jet.VarMap)
		if ok == false {

			varMap, ok := data.(gin.H)

			if !ok {
				//varMap, err := data.(map[string]interface{})
				varMap = structs.Map(data) //不是gin.H类型就是map[string]interface{}类型
			}

			v = make(jet.VarMap)

			for key, value := range varMap {
				v.Set(key, value)
			}
		} else {
			v = vars
		}
	}

	//fmt.Println(v)
	
	return JetRender{
		Data:     data,
		Options:  r.Options,
		Template: t,
		Variables: v,
	}
}

func (r JetRender) AddGlobal(key string, value interface{}) {
	if r.globals == nil {
		r.globals = make(map[string]interface{})
	}
	r.globals[key] = value
}

func (r JetRender) AddGlobalFunc(key string, fn jet.Func) {
	r.AddGlobal(key, fn)
}

func (r JetRender) Render(w http.ResponseWriter) error {
	// Unless already set, write the Content-Type header.
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{r.Options.ContentType}
	}
	if err := r.Template.Execute(w, r.Variables, r.Data); err != nil {
		return err
	}
	return nil
}

func (r JetRender) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{r.Options.ContentType}
	}
}
