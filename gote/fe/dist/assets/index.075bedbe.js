import{d as V,r as w,o as j,w as N,b as K,a as C,c as L,F as q,e as B,n as E,t as O,f as oe,g as I,m as U,h as ae,i as n,u as d,j as $,v as le,k as ie,l as X,p as ce,q as re,s as ue}from"./vendor.2f7dd286.js";const de=function(){const l=document.createElement("link").relList;if(l&&l.supports&&l.supports("modulepreload"))return;for(const o of document.querySelectorAll('link[rel="modulepreload"]'))r(o);new MutationObserver(o=>{for(const c of o)if(c.type==="childList")for(const m of c.addedNodes)m.tagName==="LINK"&&m.rel==="modulepreload"&&r(m)}).observe(document,{childList:!0,subtree:!0});function f(o){const c={};return o.integrity&&(c.integrity=o.integrity),o.referrerpolicy&&(c.referrerPolicy=o.referrerpolicy),o.crossorigin==="use-credentials"?c.credentials="include":o.crossorigin==="anonymous"?c.credentials="omit":c.credentials="same-origin",c}function r(o){if(o.ep)return;o.ep=!0;const c=f(o);fetch(o.href,c)}};de();var g=(e=>(e[e.ReadOnly=1]="ReadOnly",e[e.ReadWrite=2]="ReadWrite",e[e.Html=3]="Html",e))(g||{}),_=(e=>(e[e.All=0]="All",e[e.Tagged=1]="Tagged",e[e.Untagged=2]="Untagged",e))(_||{}),R=(e=>(e.MenuOpen="state.menuOpen",e.SelectedPath="state.selectedPath",e))(R||{});const fe={class:""},me=["onClick"],pe=V({props:{items:null},emits:["select"],setup(e,{emit:l}){const f=e,r=w("");let o=!0;j(()=>{c()}),N(f,()=>{c()});function c(){if(o&&f.items.length>0){o=!1;const v=K.get(R.SelectedPath);if(v!==""){for(let i of f.items)if(i.path===v){m(i);break}}}}function m(v){r.value=v.path,l("select",v),K.set(R.SelectedPath,r.value,{expires:365})}return(v,i)=>(C(),L("ul",fe,[(C(!0),L(q,null,B(e.items,(k,x)=>(C(),L("li",{key:x,onClick:P=>m(k),class:E({active:r.value==k.path})},O(k.title),11,me))),128))]))}});let W=ge();try{W=oe.exports.getapi()}catch{console.log("API is not ready!")}function ge(){return{listNotes(){return[{path:"clipboard.md",title:"clipboard",tags:[]},{path:"Go Keywords.md",title:"Go Keywords",tags:["go"]},{path:"midway.md",title:"midway",tags:["tool"]}]},getNote(e){switch(e){case"Go Keywords.md":return{content:`# Go Keywords

A minimal runnable go code that contain all key words:

\`\`\`go
package main

import _ "os"

func main() {
	goto a
a:
	const t = true
	defer main()
	go main()
	type X struct{}
	var (
		c chan X
		_ map[X]interface{}
	)
	for range c {
		break
		continue
	}
	if t {
	} else {
	}
	switch {
	case t:
		fallthrough
	default:
	}
	select {}
	return
}
\`\`\`

Ref to: https://go.dev/ref/spec#Keywords

\`\`\`
Keywords
The following keywords are reserved and may not be used as identifiers.

break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
\`\`\`
`};case"clipboard.md":return{content:"this is a clipboard"};default:return{content:""}}}}}window.api=W;require.config({paths:{vs:"https://cdn.jsdelivr.net/npm/monaco-editor@0.20.0/min/vs"}});let D;const ve=new Promise((e,l)=>{D=e});require(["vs/editor/editor.main"],function(){D(monaco)});let F=[];const he=()=>ve;async function _e(e){const l=await he();return F=l.languages.getLanguages(),console.log("create editor"),l.editor.create(e,{scrollBeyondLastLine:!1,automaticLayout:!0})}function G(e,l,f=""){if(l===""){e.setValue("");return}const r=ye(f);e.setModel(monaco.editor.createModel(l,r)),e.layout()}function ye(e){e=e.substr(e.lastIndexOf("/")+1);const l=e.lastIndexOf(".");if(l===-1)return"";const f=e.substr(l).toLowerCase();for(const r of F)if(r.extensions&&r.extensions.indexOf(f)!==-1)return r.id}const we={class:"nav border"},ke=n("div",{class:"nav-head"},null,-1),be={class:"content"},xe=n("span",{class:"material-icons icon"},"description",-1),Oe=n("span",{class:"text"},"All Notes",-1),Ce={class:"note-count"},Le=n("span",{class:"material-icons icon"},"local_offer",-1),Ee=n("span",{class:"text"},"Tags",-1),Pe={class:"note-count"},Re=re('<ul class="demo"><li class="item-2"><span class="text">typescript</span><span class="note-count">2</span></li><li class="item-2"><span class="text">vue</span><span class="note-count">1</span></li><li class="item-2"><span class="text">go</span><span class="note-count">1</span></li></ul>',1),Ae={class:""},Se=["onClick"],Te={class:"text"},Ne={class:"note-count"},Ke=n("span",{class:"material-icons icon"},"edit_off",-1),We=n("span",{class:"text"},"Untagged",-1),Ie={class:"note-count"},He={class:"note-list right-border"},Me={class:"head"},Ue=n("span",{class:"material-icons icon search-icon"},"search",-1),$e={class:"content"},Xe=n("ul",{class:"demo"},[n("li",null,"clipboard"),n("li",{class:"active"},"Go Keywords"),n("li",null,"midway"),n("li",null,"gote"),n("li",null,"Web Tech")],-1),Ge={class:"note"},Ve={class:"head"},je={key:0,class:"readonly"},qe=["innerHTML"],Be=n("pre",{class:"demo"},`# Go Keywords

A minimal runnable go code that contain all key words:

\`\`\`go
package main

import _ "os"

func main() {
	goto a
a:
	const t = true
	defer main()
	go main()
	type X struct{}
	var (
		c chan X
		_ map[X]interface{}
	)
	for range c {
		break
		continue
	}
	if t {
	} else {
	}
	switch {
	case t:
		fallthrough
	default:
	}
	select {}
	return
}
\`\`\`

Ref to: https://go.dev/ref/spec#Keywords

\`\`\`
Keywords
The following keywords are reserved and may not be used as identifiers.

break        default      func         interface    select
case         defer        go           map          struct
chan         else         goto         package      switch
const        fallthrough  if           range        type
continue     for          import       return       var
\`\`\``,-1),De=V({setup(e){const l=w(null),f=w(null),r=w(null),o=w(!1),c=w(!1),m=w([]),v=w(""),i=w(g.ReadOnly),k=w(""),x=w(_.All);let P;const J=I(()=>i.value!=g.Html?"":U.parse(v.value)),z=I(()=>{var y,p;let t=[];const s=x.value;switch(s){case _.All:t=m.value;break;case _.Tagged:for(let h of m.value)h.tags&&t.push(h);break;case _.Untagged:for(let h of m.value)h.tags||t.push(h);break;default:if(typeof s!="string"){console.log("invalid tag value: ",s),t=m.value;break}for(let h of m.value)for(let T of(y=h.tags)!=null?y:[])if(s===T){t.push(h);break}break}if(k.value.length<=1)return t;const a=[],u=[],b=k.value;for(let h of t)for(let T of(p=h.tags)!=null?p:[]){if(T===b){a.push(h);break}if(T.indexOf(b)!==-1){u.push(h);break}}return a.concat(u)}),A=I(()=>{let t=0,s=0,a=0;const u=new Map;for(let y of m.value)if(t++,y.tags==null||y.tags.length==0)a++;else{s++;for(let p of y.tags)u.set(p,(u.get(p)||0)+1)}const b=[];return u.forEach((y,p)=>{b.push([p,y])}),b.sort((y,p)=>y[0]==p[0]?0:y[0]<p[0]?-1:1),{total:t,tagged:s,untagged:a,tags:b}}),Y=()=>window.innerWidth<=600||navigator.userAgent.indexOf("Lumi")!=-1,Q=window.hljs;U.setOptions({highlight:function(t,s){return s===""?t:Q.highlight(t,{language:s}).value}}),j(async()=>{const t=K.get(R.MenuOpen)||JSON.stringify(!0);o.value=JSON.parse(t),c.value=!0,m.value=await W.listNotes(),i.value=g.Html,P=await _e(l.value),H(!0),document.addEventListener("keyup",M)}),ae(()=>{H(!1),document.removeEventListener("keyup",M)}),N([i,v],()=>{P!=null&&(i.value==g.ReadWrite?G(P,v.value,".md"):G(P,"",".md"))}),N([o,i],()=>{const t=document.body.style;o.value||i.value==g.ReadWrite?(t.removeProperty("height"),t.removeProperty("overflow")):(t.height="initial",t.overflow="auto")}),N(o,()=>{K.set(R.MenuOpen,JSON.stringify(o.value),{expires:365})});async function Z(t){Y()&&(o.value=!1),window.document.title=t.title;const s=(await W.getNote(t.path)).content;v.value!=s&&(f.value&&f.value.scrollTo(0,0),v.value=s)}function S(t){x.value=t,t!=_.All&&t!=_.Tagged&&(k.value="")}function ee(){i.value==g.ReadWrite?i.value=g.Html:i.value=g.ReadWrite}function te(t){const s=window.innerWidth,a=t.clientX;let u=0;a<=s/3?u=-1:a>=s*2/3&&(u=1),u!==0&&(t.preventDefault(),ne(u))}function ne(t){const s=window.innerHeight-40-14;window.scrollTo(window.scrollX,window.scrollY+t*s)}function se(t){let s=0;const a=()=>s=0,u=()=>s++,b=p=>{s<10&&!o.value&&i.value!=g.ReadWrite&&t(p)};return p=>{p?(document.addEventListener("mousedown",a),document.addEventListener("mousemove",u),document.addEventListener("mouseup",b)):(document.removeEventListener("mousedown",a),document.removeEventListener("mousemove",u),document.removeEventListener("mouseup",b))}}const H=se(te);function M(t){t.key=="Escape"&&(k.value="",r.value!=null&&r.value.focus(),t.preventDefault())}return(t,s)=>(C(),L("div",{class:E(["root",{close:!o.value,fullscreen:c.value}])},[n("span",{class:"material-icons menu-icon",onClick:s[0]||(s[0]=a=>o.value=!o.value)},O(o.value?"menu_open":"menu"),1),n("div",we,[ke,n("div",be,[n("div",{class:E(["item-1",{active:x.value==d(_).All}]),onClick:s[1]||(s[1]=a=>S(d(_).All))},[xe,Oe,n("span",Ce,O(d(A).total),1)],2),n("div",{class:E(["item-1",{active:x.value==d(_).Tagged}]),onClick:s[2]||(s[2]=a=>S(d(_).Tagged))},[Le,Ee,n("span",Pe,O(d(A).tagged),1)],2),Re,n("ul",Ae,[(C(!0),L(q,null,B(d(A).tags,(a,u)=>(C(),L("li",{class:E(["item-2",{active:x.value===a[0]}]),key:u,onClick:b=>S(a[0])},[n("span",Te,O(a[0]),1),n("span",Ne,O(a[1]),1)],10,Se))),128))]),n("div",{class:E(["item-1",{active:x.value==d(_).Untagged}]),onClick:s[3]||(s[3]=a=>S(d(_).Untagged))},[Ke,We,n("span",Ie,O(d(A).untagged),1)],2)])]),n("div",He,[n("div",Me,[$(n("input",{class:"search",type:"text",ref_key:"searchElement",ref:r,placeholder:"Search by tag...","onUpdate:modelValue":s[4]||(s[4]=a=>k.value=a),spellcheck:"false"},null,512),[[le,k.value]]),Ue]),n("div",$e,[ie(pe,{items:d(z),onSelect:Z},null,8,["items"]),Xe])]),n("div",Ge,[n("div",Ve,[n("span",{class:E(["material-icons icon",{active:i.value==d(g).ReadWrite}]),onClick:ee},"edit",2)]),n("div",{class:"content",ref_key:"contentElement",ref:f},[i.value==d(g).ReadOnly?(C(),L("pre",je,O(v.value),1)):X("",!0),$(n("div",{class:"readwrite",ref_key:"inputElement",ref:l},null,512),[[ce,i.value==d(g).ReadWrite]]),i.value==d(g).Html?(C(),L("div",{key:1,class:"html",innerHTML:d(J)},null,8,qe)):X("",!0),Be],512)])],2))}});ue(De).mount("#app");
