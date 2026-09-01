package app

import (
	"net/http"
)

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(page))
}

const page = `<!doctype html><html><head><meta charset="utf-8"><title>JobPilot</title><style>body{font:16px system-ui;margin:0;background:#f6f8fb;color:#172033}header{padding:20px 8%;background:#14213d;color:#fff}main{max-width:1000px;margin:30px auto;padding:0 20px}.cards{display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:14px}.card,section{background:white;border-radius:10px;padding:18px;box-shadow:0 1px 4px #ccd}label{display:block;margin:12px 0 4px}input,textarea{width:100%;box-sizing:border-box;padding:8px}button{margin-top:16px;padding:10px 16px;background:#1f6feb;color:white;border:0;border-radius:6px}small{color:#657}</style></head><body><header><h1>JobPilot</h1><p>Safe job-application preparation — human approval required.</p></header><main id="root"></main><script type="module">import React from 'https://esm.sh/react@18';import{createRoot}from'https://esm.sh/react-dom@18/client';const e=React.createElement;function App(){const[p,setP]=React.useState(null),[m,setM]=React.useState({});React.useEffect(()=>{fetch('/api/candidate').then(x=>x.json()).then(setP);fetch('/api/dashboard').then(x=>x.json()).then(setM)},[]);if(!p)return e('main',null,'Loading…');const save=()=>fetch('/api/candidate',{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(p)}).then(x=>x.json()).then(setP);return e('main',null,e('div',{className:'cards'},...Object.entries(m).map(([k,v])=>e('div',{className:'card',key:k},e('small',null,k.replaceAll('_',' ')),e('h2',null,v)))),e('section',null,e('h2',null,'Candidate profile'),e('small',null,'Fields marked TODO must be verified before JobPilot can use them.'),['full_name','email','phone','location','linkedin_url','github_url','portfolio_url','years_experience','work_authorization','sponsorship_requirement'].map(k=>e('label',{key:k},k.replaceAll('_',' '),e('input',{value:p[k]||'',onChange:x=>setP({...p,[k]:x.target.value})}))),e('label',null,'Skills',e('textarea',{value:(p.skills||[]).join(', '),onChange:x=>setP({...p,skills:x.target.value.split(',').map(v=>v.trim()).filter(Boolean)})})),e('button',{onClick:save},'Save verified profile'),e('p',null,'Profile version: ',p.version)));}createRoot(document.getElementById('root')).render(e(App));</script></body></html>`
